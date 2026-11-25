# Modul 10 - Dokumentasi, Pengujian, dan Deployment API

Implementasi Modul 10 pada proyek ini disesuaikan dengan endpoint yang telah tersedia pada Laravel 13, Fiber v3, dan frontend React/Vite.

## Berkas yang ditambahkan atau diperbarui

| Bagian modul | Implementasi |
|---|---|
| Langkah 1 - OpenAPI | `openapi.yaml` di akar repositori, mencakup endpoint Laravel dan Fiber yang benar-benar terdaftar. |
| Langkah 2 - Pengujian Laravel | `latihan-laravel/tests/Feature/MahasiswaApiTest.php`; PHPUnit memakai SQLite in-memory dari `phpunit.xml`. |
| Perbaikan agar pengujian bermakna | `UpdateMahasiswaRequest::authorize()` diubah menjadi `true`; autentikasi/kemampuan tetap diperiksa middleware pada `routes/api.php`. |
| Langkah 3 - Pengujian Fiber | `latihan-fiber/internal/handler/mahasiswa_handler_test.go`, mencakup API key status, autentikasi daftar mahasiswa, ID nonangka, dan JSON rusak. |
| Langkah 4 - Konfigurasi Laravel | `latihan-laravel/.env.production.example`, `latihan-laravel/config/cors.php`, dan petunjuk optimasi di `deploy/README.md`. |
| Langkah 5 - Binary Fiber | `latihan-fiber/scripts/build-linux.sh`. |
| Langkah 6 - Systemd | `deploy/systemd/pemweb2-fiber.service`. |
| Langkah 7 - Reverse proxy | `deploy/nginx/pemweb2.conf` (Laravel melalui PHP-FPM, Fiber melalui reverse proxy). |
| Langkah 8 - Build frontend | `frontend-pemweb2/.env.production.example`; perintah build dijelaskan pada `deploy/README.md`. |
| Langkah 9 - CI | `.github/workflows/uji.yml` menjalankan tes Go, tes Laravel, lint, dan build frontend. |
| Langkah 10 - Checklist | Bagian checklist rilis dan contoh pencadangan database di `deploy/README.md`. |

## Menjalankan pemeriksaan di komputer lokal

### OpenAPI

Buka `openapi.yaml` melalui Swagger Editor atau ekstensi OpenAPI di editor kode. Dua server pengembangan yang tercantum adalah Laravel `http://127.0.0.1:8000/api` dan Fiber `http://127.0.0.1:3000/api`.

### Fiber

```bash
cd latihan-fiber
go vet ./...
go test ./... -cover
```

### Laravel

```bash
cd latihan-laravel
composer install
php artisan test
```

Pengujian PHPUnit menggunakan database SQLite `:memory:` sehingga tidak mengubah database pengembangan.

### Frontend

```bash
cd frontend-pemweb2
npm ci
npm run lint
npm run build
```

## Catatan implementasi dan batasan

- Nama dan bentuk endpoint di OpenAPI mengikuti implementasi aktual. Sebagai contoh, Laravel memakai `/auth/register` dan field `password`, sedangkan Fiber memakai `/auth/registrasi` dan `kata_sandi`.
- Endpoint mata kuliah, program studi, soft delete/pemulihan, dan status saat ini hanya didaftarkan pada Fiber; dokumentasi menandainya sebagai Fiber-only.
- Laravel menghapus mahasiswa secara permanen, sedangkan Fiber memakai soft delete dan menyediakan endpoint pemulihan.
- Berkas deployment merupakan template untuk server Linux, bukan deployment yang sudah dijalankan. Sesuaikan domain, path, sertifikat, socket PHP-FPM, service account, dan environment file pada server tujuan.
- Hasil `go test` dan `php artisan test` harus dijalankan kembali pada mesin pengembangan atau CI dengan toolchain dan dependensi proyek yang tersedia.
