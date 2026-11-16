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

func DaftarMahasiswa(c fiber.Ctx) error {
	halaman, _ := strconv.Atoi(c.Query("halaman", "1"))
	perHalaman, _ := strconv.Atoi(c.Query("per_halaman", "10"))
	angkatan, _ := strconv.Atoi(c.Query("angkatan", "0"))

	if halaman < 1 {
		halaman = 1
	}
	if perHalaman < 1 || perHalaman > 100 {
		perHalaman = 10
	}

	filter := repository.FilterMahasiswa{
		Cari:       c.Query("cari", ""),
		Angkatan:   angkatan,
		Halaman:    halaman,
		PerHalaman: perHalaman,
	}

	daftar, total, err := repository.Ambil(filter)
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

func DetailMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	mahasiswa, err := repository.AmbilPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   mahasiswa,
	})
}

func BuatMahasiswa(c fiber.Ctx) error {
	var permintaan PermintaanMahasiswa
	err := c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	aktif := true
	if permintaan.Aktif != nil {
		aktif = *permintaan.Aktif
	}

	mahasiswa := model.Mahasiswa{
		ProgramStudiID: permintaan.ProgramStudiID,
		Nim:            permintaan.Nim,
		Nama:           permintaan.Nama,
		Email:          permintaan.Email,
		Angkatan:       permintaan.Angkatan,
		Ipk:            permintaan.Ipk,
		Aktif:          aktif,
	}

	err = repository.Simpan(&mahasiswa)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal menyimpan data, periksa keunikan NIM dan email")
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Program studi tidak ditemukan")
	}
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan data")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dibuat",
		"data":   mahasiswa,
	})
}

// PerbaruiMahasiswa melayani PUT (penggantian penuh): body divalidasi dengan
// aturan yang sama seperti pembuatan.
func PerbaruiMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	var permintaan PermintaanMahasiswa
	err = c.Bind().Body(&permintaan)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); galat != nil {
		return tolakValidasi(c, galat)
	}

	mahasiswa, err := repository.AmbilPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	mahasiswa.ProgramStudiID = permintaan.ProgramStudiID
	mahasiswa.Nim = permintaan.Nim
	mahasiswa.Nama = permintaan.Nama
	mahasiswa.Email = permintaan.Email
	mahasiswa.Angkatan = permintaan.Angkatan
	mahasiswa.Ipk = permintaan.Ipk
	if permintaan.Aktif != nil {
		mahasiswa.Aktif = *permintaan.Aktif
	}

	err = repository.Perbarui(&mahasiswa)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal memperbarui data, periksa keunikan NIM dan email")
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Program studi tidak ditemukan")
	}
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui data")
	}

	// Muat ulang agar relasi program_studi pada respons sesuai ID yang baru.
	mahasiswa, err = repository.AmbilPerId(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil diperbarui",
		"data":   mahasiswa,
	})
}

func HapusMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	mahasiswa, err := repository.AmbilPerId(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa tidak ditemukan")
		}
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	err = repository.Hapus(&mahasiswa)
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dihapus",
	})
}

// DaftarMahasiswaTerhapus melayani GET /api/mahasiswa/terhapus (Tugas 3):
// menampilkan data yang sudah di-soft delete.
func DaftarMahasiswaTerhapus(c fiber.Ctx) error {
	halaman, _ := strconv.Atoi(c.Query("halaman", "1"))
	perHalaman, _ := strconv.Atoi(c.Query("per_halaman", "10"))
	angkatan, _ := strconv.Atoi(c.Query("angkatan", "0"))

	if halaman < 1 {
		halaman = 1
	}
	if perHalaman < 1 || perHalaman > 100 {
		perHalaman = 10
	}

	filter := repository.FilterMahasiswa{
		Cari:       c.Query("cari", ""),
		Angkatan:   angkatan,
		Halaman:    halaman,
		PerHalaman: perHalaman,
	}

	daftar, total, err := repository.AmbilTerhapus(filter)
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

// PulihkanMahasiswa melayani PUT /api/mahasiswa/:id/pulihkan (Tugas 3):
// mengembalikan data yang sudah di-soft delete.
func PulihkanMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka")
	}

	terpulihkan, err := repository.Pulihkan(uint(id))
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memulihkan data")
	}
	if terpulihkan == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Data mahasiswa terhapus tidak ditemukan")
	}

	mahasiswa, err := repository.AmbilPerId(uint(id))
	if err != nil {
		log.Printf("Galat %s %s: %v", c.Method(), c.Path(), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dipulihkan",
		"data":   mahasiswa,
	})
}
