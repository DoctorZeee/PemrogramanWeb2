package handler

import (
	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/validasi"
)

// tolakValidasi membalas 422 dengan format yang sama untuk semua endpoint.
func tolakValidasi(c fiber.Ctx, galat []validasi.PesanGalat) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"sukses": false,
		"pesan":  "Data yang dikirim tidak valid",
		"galat":  galat,
	})
}
