package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"latihan-fiber/internal/auth"
	"latihan-fiber/internal/config"
	"latihan-fiber/internal/model"
	"latihan-fiber/internal/validasi"
)

func Registrasi(c fiber.Ctx) error {
	var permintaan PermintaanRegistrasi
	err := c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(permintaan.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pengguna := model.Pengguna{
		Nama:      permintaan.Nama,
		Email:     permintaan.Email,
		KataSandi: string(hash),
		Peran:     "mahasiswa",
	}

	err = config.DB.Create(&pengguna).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Email sudah terdaftar")
	}
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Registrasi berhasil",
		"data":   pengguna,
	})
}

func Login(c fiber.Ctx) error {
	var permintaan PermintaanLogin
	err := c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	// Pesan sengaja sama untuk email salah maupun kata sandi salah.
	kredensialSalah := fiber.NewError(fiber.StatusUnauthorized, "Email atau kata sandi tidak sesuai")

	var pengguna model.Pengguna
	err = config.DB.Where("email = ?", permintaan.Email).First(&pengguna).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return kredensialSalah
	}
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pengguna.KataSandi), []byte(permintaan.KataSandi))
	if err != nil {
		return kredensialSalah
	}

	konf := config.Muat()
	token, err := auth.TerbitkanToken(pengguna.ID, pengguna.Peran, konf.JwtSecret, konf.JwtTtlMenit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Login berhasil",
		"data": fiber.Map{
			"pengguna": pengguna,
			"token":    token,
			"tipe":     "Bearer",
		},
	})
}

func Profil(c fiber.Ctx) error {
	var pengguna model.Pengguna
	idPengguna, _ := c.Locals("id_pengguna").(uint)
	err := config.DB.First(&pengguna, idPengguna).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "Pengguna tidak ditemukan")
	}
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   pengguna,
	})
}

// SegarkanToken menerbitkan token baru berdasarkan token lama yang masih berlaku (Tugas 1).
// Peran dibaca ulang dari basis data agar perubahan peran ikut berlaku.
// Token lama langsung dicabut (rotasi) agar token yang bocor tidak tetap berlaku.
func SegarkanToken(c fiber.Ctx) error {
	// Cabut token lama lebih dulu: bila pencabutan gagal, token baru tidak diterbitkan.
	if err := cabutTokenAktif(c); err != nil {
		return err
	}

	idPengguna, _ := c.Locals("id_pengguna").(uint)

	var pengguna model.Pengguna
	err := config.DB.First(&pengguna, idPengguna).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusUnauthorized, "Pengguna tidak ditemukan")
	}
	if err != nil {
		return err
	}

	konf := config.Muat()
	tokenBaru, err := auth.TerbitkanToken(pengguna.ID, pengguna.Peran, konf.JwtSecret, konf.JwtTtlMenit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Token berhasil disegarkan",
		"data": fiber.Map{
			"token": tokenBaru,
			"tipe":  "Bearer",
		},
	})
}

// cabutTokenAktif memasukkan sidik jari token pada request saat ini ke daftar
// hitam. Entri yang sudah kedaluwarsa ikut dibersihkan agar tabel tidak membengkak.
func cabutTokenAktif(c fiber.Ctx) error {
	ringkas, _ := c.Locals("token_hash").(string)
	kedaluwarsa, _ := c.Locals("token_kedaluwarsa").(time.Time)

	err := config.DB.Where("kedaluwarsa_pada < ?", time.Now()).Delete(&model.DaftarHitamToken{}).Error
	if err != nil {
		return err
	}

	entri := model.DaftarHitamToken{
		TokenHash:       ringkas,
		KedaluwarsaPada: kedaluwarsa,
	}
	err = config.DB.Create(&entri).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}
	return nil
}

// Logout mencabut token dengan memasukkan sidik jarinya ke daftar hitam (Tugas 3).
// Sidik jari dan waktu kedaluwarsa sudah disiapkan oleh middleware Autentikasi.
func Logout(c fiber.Ctx) error {
	if err := cabutTokenAktif(c); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Logout berhasil, token telah dicabut",
	})
}
