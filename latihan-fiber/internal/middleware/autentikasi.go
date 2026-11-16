package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/auth"
	"latihan-fiber/internal/config"
	"latihan-fiber/internal/model"
)

func Autentikasi() fiber.Handler {
	konf := config.Muat()
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Header Authorization tidak ditemukan")
		}

		bagian := strings.Fields(header)
		if len(bagian) != 2 || !strings.EqualFold(bagian[0], "Bearer") {
			return fiber.NewError(fiber.StatusUnauthorized, "Format header Authorization tidak sesuai")
		}
		nilaiToken := bagian[1]

		// Verifikasi tanda tangan lebih dulu (murah, tanpa basis data),
		// baru periksa daftar hitam.
		klaim, err := auth.PeriksaToken(nilaiToken, konf.JwtSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
		}

		id, err := strconv.ParseUint(klaim.Subject, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Isi token tidak sesuai")
		}

		ringkas := auth.HashToken(nilaiToken)
		var dicabut int64
		err = config.DB.Model(&model.DaftarHitamToken{}).
			Where("token_hash = ?", ringkas).
			Count(&dicabut).Error
		if err != nil {
			// Galat non-fiber: dicatat di log, klien hanya menerima pesan umum.
			return err
		}
		if dicabut > 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "Token sudah dicabut")
		}

		c.Locals("id_pengguna", uint(id))
		c.Locals("peran", klaim.Peran)
		c.Locals("token_hash", ringkas)
		c.Locals("token_kedaluwarsa", klaim.ExpiresAt.Time)
		return c.Next()
	}
}

func HanyaPeran(peranDiizinkan ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		peran, ok := c.Locals("peran").(string)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Peran pengguna tidak diketahui")
		}

		for _, p := range peranDiizinkan {
			if p == peran {
				return c.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "Anda tidak berwenang mengakses sumber daya ini")
	}
}
