package router

import (
	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/config"
	"latihan-fiber/internal/handler"
	"latihan-fiber/internal/middleware"
)

func Daftarkan(app *fiber.App, konf config.Konfigurasi) {
	api := app.Group("/api")

	// Kunci API (Tugas 3 Modul 6) dipasang pada /api/status saja, agar rute
	// auth dan CRUD dapat diuji persis seperti contoh curl pada modul.
	api.Get("/status", middleware.KunciApi(konf.ApiKey), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"sukses": true,
			"pesan":  "API Fiber aktif",
		})
	})

	// Rute autentikasi (publik).
	api.Post("/auth/registrasi", handler.Registrasi)
	api.Post("/auth/login", handler.Login)

	// Fiber menjalankan handler sesuai urutan penulisan, jadi middleware
	// peran HARUS ditulis sebelum handler akhir.
	adminDosen := middleware.HanyaPeran("admin", "dosen")
	admin := middleware.HanyaPeran("admin")

	// Autentikasi JWT dipasang per rute, bukan lewat grup kosong, agar URL
	// yang tidak terdaftar tetap dijawab 404 oleh penangan papan bawah,
	// bukan 401 oleh middleware.
	autentikasi := middleware.Autentikasi()

	// Rute autentikasi terlindungi JWT.
	api.Get("/auth/profil", autentikasi, handler.Profil)
	api.Post("/auth/segarkan", autentikasi, handler.SegarkanToken)
	api.Post("/auth/logout", autentikasi, handler.Logout)

	// Mahasiswa: semua pengguna terautentikasi boleh melihat.
	mahasiswa := api.Group("/mahasiswa")
	mahasiswa.Get("/terhapus", autentikasi, admin, handler.DaftarMahasiswaTerhapus)
	mahasiswa.Get("/", autentikasi, handler.DaftarMahasiswa)
	mahasiswa.Get("/:id", autentikasi, handler.DetailMahasiswa)
	mahasiswa.Post("/", autentikasi, adminDosen, handler.BuatMahasiswa)
	mahasiswa.Put("/:id", autentikasi, adminDosen, handler.PerbaruiMahasiswa)
	mahasiswa.Delete("/:id", autentikasi, admin, handler.HapusMahasiswa)
	mahasiswa.Put("/:id/pulihkan", autentikasi, admin, handler.PulihkanMahasiswa)

	// Matakuliah: melihat untuk semua, mengubah hanya admin/dosen.
	matakuliah := api.Group("/matakuliah")
	matakuliah.Get("/", autentikasi, handler.DaftarMataKuliah)
	matakuliah.Get("/:id", autentikasi, handler.DetailMataKuliah)
	matakuliah.Post("/", autentikasi, adminDosen, handler.BuatMataKuliah)
	matakuliah.Put("/:id", autentikasi, adminDosen, handler.PerbaruiMataKuliah)
	matakuliah.Delete("/:id", autentikasi, admin, handler.HapusMataKuliah)

	// Program studi.
	programStudi := api.Group("/program-studi")
	programStudi.Get("/", autentikasi, handler.DaftarProgramStudi)
}
