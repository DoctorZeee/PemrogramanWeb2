package router

import (
	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/config"
	"latihan-fiber/internal/handler"
	"latihan-fiber/internal/middleware"
)

func Daftarkan(app *fiber.App, konf config.Konfigurasi) {
	api := app.Group("/api")

	// /status dibiarkan publik untuk pengecekan kesehatan layanan.
	api.Get("/status", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"sukses": true,
			"pesan":  "API Fiber aktif",
		})
	})

	// Seluruh rute mahasiswa dilindungi middleware KunciApi.
	mahasiswa := api.Group("/mahasiswa")
	mahasiswa.Use(middleware.KunciApi(konf.ApiKey))

	// Rute statis didaftarkan sebelum /:id agar tidak tertelan parameter.
	mahasiswa.Get("/terhapus", handler.DaftarMahasiswaTerhapus)
	mahasiswa.Get("/", handler.DaftarMahasiswa)
	mahasiswa.Post("/", handler.BuatMahasiswa)
	mahasiswa.Get("/:id", handler.DetailMahasiswa)
	mahasiswa.Put("/:id", handler.PerbaruiMahasiswa)
	mahasiswa.Delete("/:id", handler.HapusMahasiswa)
	// Tugas 3: pemulihan data yang sudah di-soft delete.
	mahasiswa.Put("/:id/pulihkan", handler.PulihkanMahasiswa)

	// Tugas 1: resource matakuliah dengan CRUD lengkap.
	matakuliah := api.Group("/matakuliah")
	matakuliah.Use(middleware.KunciApi(konf.ApiKey))

	matakuliah.Get("/", handler.DaftarMataKuliah)
	matakuliah.Post("/", handler.BuatMataKuliah)
	matakuliah.Get("/:id", handler.DetailMataKuliah)
	matakuliah.Put("/:id", handler.PerbaruiMataKuliah)
	matakuliah.Delete("/:id", handler.HapusMataKuliah)

	// Tugas 2: daftar program studi beserta jumlah mahasiswa.
	programStudi := api.Group("/program-studi")
	programStudi.Use(middleware.KunciApi(konf.ApiKey))

	programStudi.Get("/", handler.DaftarProgramStudi)
}
