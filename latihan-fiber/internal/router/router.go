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

	mahasiswa.Get("/", handler.DaftarMahasiswa)
	mahasiswa.Get("/nim/:nim", handler.CariMahasiswaByNim)
	mahasiswa.Get("/:id", handler.DetailMahasiswa)
}
