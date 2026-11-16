package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"latihan-fiber/internal/model"
	"latihan-fiber/internal/repository"
	"latihan-fiber/internal/validasi"
)

// CRUD matakuliah (Tugas 1), mengikuti pola repository dan handler mahasiswa.
func DaftarMataKuliah(c fiber.Ctx) error {
	halaman, _ := strconv.Atoi(c.Query("halaman", "1"))
	perHalaman, _ := strconv.Atoi(c.Query("per_halaman", "10"))
	semester, _ := strconv.Atoi(c.Query("semester", "0"))

	if halaman < 1 {
		halaman = 1
	}
	if perHalaman < 1 || perHalaman > 100 {
		perHalaman = 10
	}

	filter := repository.FilterMataKuliah{
		Cari:       c.Query("cari", ""),
		Semester:   semester,
		Halaman:    halaman,
		PerHalaman: perHalaman,
	}

	daftar, total, err := repository.AmbilMataKuliah(filter)
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	totalHalaman := int(total) / perHalaman
	if int(total)%perHalaman != 0 {
		totalHalaman = totalHalaman + 1
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   daftar,
		"meta": fiber.Map{
			"halaman":       halaman,
			"per_halaman":   perHalaman,
			"total":         total,
			"total_halaman": totalHalaman,
		},
	})
}

func DetailMataKuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	matakuliah, err := repository.AmbilMataKuliahPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   matakuliah,
	})
}

func BuatMataKuliah(c fiber.Ctx) error {
	var permintaan PermintaanMataKuliah
	err := c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	matakuliah := model.MataKuliah{
		Kode:     permintaan.Kode,
		Nama:     permintaan.Nama,
		Sks:      permintaan.Sks,
		Semester: permintaan.Semester,
	}

	err = repository.SimpanMataKuliah(&matakuliah)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal menyimpan data, periksa keunikan kode")
	}
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan data")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil dibuat",
		"data":   matakuliah,
	})
}

func PerbaruiMataKuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	matakuliah, err := repository.AmbilMataKuliahPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	var permintaan PermintaanMataKuliah
	err = c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	matakuliah.Kode = permintaan.Kode
	matakuliah.Nama = permintaan.Nama
	matakuliah.Sks = permintaan.Sks
	matakuliah.Semester = permintaan.Semester

	err = repository.PerbaruiMataKuliah(&matakuliah)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal memperbarui data, periksa keunikan kode")
	}
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil diperbarui",
		"data":   matakuliah,
	})
}

func HapusMataKuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	matakuliah, err := repository.AmbilMataKuliahPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	err = repository.HapusMataKuliah(&matakuliah)
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil dihapus",
	})
}
