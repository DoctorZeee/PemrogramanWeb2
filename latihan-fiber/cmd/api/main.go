package main

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"latihan-fiber/internal/config"
	"latihan-fiber/internal/middleware"
	"latihan-fiber/internal/router"
)

func main() {
	konf := config.Muat()
	periksaRahasia(konf)
	config.HubungkanBasisData(konf)
	config.Migrasi()
	config.DataAwal()
	app := fiber.New(fiber.Config{
		AppName:      konf.NamaAplikasi,
		ErrorHandler: penanganGalat,
	})
	app.Use(recoverer.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(middleware.WaktuProses())
	router.Daftarkan(app, konf)
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Endpoint tidak tersedia",
		})
	})
	log.Fatal(app.Listen(":" + konf.Port))
}
func penanganGalat(c fiber.Ctx, err error) error {
	kode := fiber.StatusInternalServerError
	pesan := "Terjadi kesalahan pada server"
	var galatFiber *fiber.Error
	if errors.As(err, &galatFiber) {
		kode = galatFiber.Code
		pesan = galatFiber.Message
	}

	if kode == fiber.StatusInternalServerError {
		log.Printf("Galat internal: %s %s: %v", c.Method(), c.Path(), err)
	}

	return c.Status(kode).JSON(fiber.Map{
		"sukses": false,
		"pesan":  pesan,
	})
}

// periksaRahasia menolak JWT_SECRET yang lemah di luar lingkungan lokal.
func periksaRahasia(konf config.Konfigurasi) {
	lemah := len(konf.JwtSecret) < 32 ||
		konf.JwtSecret == "rahasia-bawaan" ||
		strings.HasPrefix(konf.JwtSecret, "ubah-nilai")
	if !lemah {
		return
	}
	if konf.Lingkungan == "local" {
		log.Println("PERINGATAN: JWT_SECRET lemah, ganti sebelum dipakai di luar lokal")
		return
	}
	log.Fatal("JWT_SECRET terlalu lemah: gunakan nilai acak minimal 32 karakter")
}
