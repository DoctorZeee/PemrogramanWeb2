package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/model"
	"latihan-fiber/internal/repository"
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
		return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
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

	matakuliah := model.MataKuliah{
		Kode:     permintaan.Kode,
		Nama:     permintaan.Nama,
		Sks:      permintaan.Sks,
		Semester: permintaan.Semester,
	}

	err = repository.SimpanMataKuliah(&matakuliah)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal menyimpan data, periksa keunikan kode")
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
		return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
	}

	var permintaan PermintaanMataKuliah
	err = c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if permintaan.Kode != "" {
		matakuliah.Kode = permintaan.Kode
	}
	if permintaan.Nama != "" {
		matakuliah.Nama = permintaan.Nama
	}
	if permintaan.Sks > 0 {
		matakuliah.Sks = permintaan.Sks
	}
	if permintaan.Semester > 0 {
		matakuliah.Semester = permintaan.Semester
	}

	err = repository.PerbaruiMataKuliah(&matakuliah)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal memperbarui data")
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
		return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
	}

	err = repository.HapusMataKuliah(&matakuliah)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil dihapus",
	})
}
