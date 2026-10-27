package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
)

// KunciApi memeriksa header X-Api-Key. Permintaan ditolak dengan status 401
// apabila header kosong atau nilainya tidak sama dengan kunci pada berkas .env.
// Jika kunci di .env belum diisi, semua permintaan ditolak (aman secara bawaan).
func KunciApi(kunciValid string) fiber.Handler {
	return func(c fiber.Ctx) error {
		kunciKiriman := c.Get("X-Api-Key")

		if kunciValid == "" ||
			subtle.ConstantTimeCompare([]byte(kunciKiriman), []byte(kunciValid)) != 1 {
			return fiber.NewError(fiber.StatusUnauthorized, "API key tidak valid atau tidak ada")
		}

		return c.Next()
	}
}
