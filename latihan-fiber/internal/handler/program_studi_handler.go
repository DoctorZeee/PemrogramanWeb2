package handler

import (
	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/repository"
)

// DaftarProgramStudi melayani GET /api/program-studi (Tugas 2):
// seluruh program studi beserta jumlah mahasiswanya.
func DaftarProgramStudi(c fiber.Ctx) error {
	daftar, err := repository.AmbilSemuaProgramStudi()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   daftar,
	})
}
