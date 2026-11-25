package handler_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	"latihan-fiber/internal/config"
	"latihan-fiber/internal/handler"
	"latihan-fiber/internal/router"
)

func TestStatusEndpointMemerlukanApiKey(t *testing.T) {
	app := fiber.New()
	router.Daftarkan(app, config.Konfigurasi{ApiKey: "kunci-uji"})

	permintaan := httptest.NewRequest("GET", "/api/status", nil)
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status tanpa API key: diharapkan %d, diperoleh %d", fiber.StatusUnauthorized, respons.StatusCode)
	}
}

func TestStatusEndpointBerhasilDenganApiKey(t *testing.T) {
	app := fiber.New()
	router.Daftarkan(app, config.Konfigurasi{ApiKey: "kunci-uji"})

	permintaan := httptest.NewRequest("GET", "/api/status", nil)
	permintaan.Header.Set("X-Api-Key", "kunci-uji")
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusOK {
		t.Fatalf("status dengan API key: diharapkan %d, diperoleh %d", fiber.StatusOK, respons.StatusCode)
	}
	isi, err := io.ReadAll(respons.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"sukses":true`) {
		t.Fatalf("respons tidak memuat sukses=true: %s", isi)
	}
}

func TestDaftarMahasiswaMenolakPermintaanTanpaToken(t *testing.T) {
	app := fiber.New()
	router.Daftarkan(app, config.Konfigurasi{})

	permintaan := httptest.NewRequest("GET", "/api/mahasiswa", nil)
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status tanpa token: diharapkan %d, diperoleh %d", fiber.StatusUnauthorized, respons.StatusCode)
	}
}

func TestDetailMahasiswaIdBukanAngka(t *testing.T) {
	app := fiber.New()
	app.Get("/api/mahasiswa/:id", handler.DetailMahasiswa)

	permintaan := httptest.NewRequest("GET", "/api/mahasiswa/abc", nil)
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status ID mahasiswa nonangka: diharapkan %d, diperoleh %d", fiber.StatusBadRequest, respons.StatusCode)
	}
}

func TestDetailMataKuliahIdBukanAngka(t *testing.T) {
	app := fiber.New()
	app.Get("/api/matakuliah/:id", handler.DetailMataKuliah)

	permintaan := httptest.NewRequest("GET", "/api/matakuliah/abc", nil)
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status ID mata kuliah nonangka: diharapkan %d, diperoleh %d", fiber.StatusBadRequest, respons.StatusCode)
	}
}

func TestBuatMahasiswaMenolakBodyJsonRusak(t *testing.T) {
	app := fiber.New()
	app.Post("/api/mahasiswa", handler.BuatMahasiswa)

	permintaan := httptest.NewRequest("POST", "/api/mahasiswa", strings.NewReader(`{"nim":`))
	permintaan.Header.Set("Content-Type", "application/json")
	respons, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer respons.Body.Close()

	if respons.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status body JSON rusak: diharapkan %d, diperoleh %d", fiber.StatusBadRequest, respons.StatusCode)
	}
}
