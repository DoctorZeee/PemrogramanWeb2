# Deployment Modul 10

Folder ini berisi contoh konfigurasi untuk menyiapkan proyek pada server Linux. Berkas di sini adalah template: ganti domain, jalur proyek, akun layanan, socket PHP-FPM, sertifikat TLS, dan kredensial sesuai server. Jangan salin placeholder sebagai secret produksi.

## 1. Laravel 13

1. Siapkan PHP 8.4, ekstensi yang dibutuhkan Laravel, Composer, Nginx, PHP-FPM, dan MariaDB/MySQL.
2. Salin `latihan-laravel/.env.production.example` menjadi `.env` di server. Isi `APP_KEY`, URL, dan kredensial database. Pastikan `APP_ENV=production` dan `APP_DEBUG=false`.
3. Buat APP_KEY bila aplikasi belum memilikinya: `php artisan key:generate` (jalankan dari folder proyek, setelah `.env` dibuat).
4. Pasang dependensi dan optimalkan aplikasi:

   ```bash
   composer install --no-dev --optimize-autoloader --no-interaction
   php artisan optimize:clear
   php artisan config:cache
   php artisan route:cache
   php artisan view:cache
   php artisan migrate --force
   ```

5. Arahkan document root Nginx ke folder `public`, bukan akar proyek. Pastikan PHP-FPM berjalan dan socket pada konfigurasi Nginx sesuai sistem operasi.

> `php artisan serve` ditujukan untuk pengembangan, bukan server produksi. Contoh Nginx di `nginx/pemweb2.conf` menggunakan PHP-FPM untuk Laravel.

## 2. Fiber v3

1. Buat akun layanan tanpa hak login dan direktori proyek (contoh ini menggunakan Fedora/RHEL):

   ```bash
   sudo useradd --system --home-dir /var/www/pemweb2 --shell /sbin/nologin pemweb2
   sudo install -d -o pemweb2 -g pemweb2 /var/www/pemweb2
   ```

2. Salin `latihan-fiber/.env.production.example` ke lokasi rahasia server, misalnya `/etc/pemweb2/fiber.env`. Isi API key, JWT secret acak minimal 32 karakter, serta akses database. Pertahankan `APP_HOST=127.0.0.1` agar port Fiber tidak terbuka langsung ke jaringan publik. Amankan file environment dengan kepemilikan root dan izin `600`; systemd membacanya saat memulai layanan.
3. Letakkan source proyek pada `/var/www/pemweb2/latihan-fiber`, lalu bangun binary Linux dari folder proyek menggunakan toolchain Go yang sesuai:

   ```bash
   cd /var/www/pemweb2/latihan-fiber
   ./scripts/build-linux.sh
   sudo chown pemweb2:pemweb2 bin/api
   ```

4. Pasang unit `systemd/pemweb2-fiber.service` ke `/etc/systemd/system/`, sesuaikan jalur, user/group, file environment, dan nama layanan database.
5. Aktifkan layanan:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now pemweb2-fiber
   sudo systemctl status pemweb2-fiber
   journalctl -u pemweb2-fiber -f
   ```

## 3. Nginx dan frontend

- Contoh `nginx/pemweb2.conf` melayani Laravel pada `https://api.contoh.ac.id/api` dan Fiber pada `https://api.contoh.ac.id/fiber/api`. Prefix `/fiber/` dilepas oleh Nginx sebelum request diteruskan ke Fiber.
- Contoh ini memakai socket PHP-FPM Fedora `/run/php-fpm/www.sock`. Pada distribusi lain, sesuaikan ke socket PHP-FPM yang aktif.
- Siapkan sertifikat TLS terlebih dahulu (misalnya melalui Let's Encrypt), lalu aktifkan site dan jalankan `nginx -t` sebelum reload Nginx.
- Salin `frontend-pemweb2/.env.production.example` menjadi `.env.production`, sesuaikan URL API, lalu jalankan:

  ```bash
  npm ci
  npm run lint
  npm run build
  ```

  Publikasikan isi `frontend-pemweb2/dist` pada document root domain frontend. Semua variabel `VITE_*` masuk ke bundle publik; jangan letakkan password, API key, atau JWT secret di dalamnya.

## 4. CI dan pengujian

Workflow `.github/workflows/uji.yml` berjalan pada push dan pull request ke `main`:

- Fiber: `go vet ./...` dan `go test ./... -cover`.
- Laravel: `php artisan test`; konfigurasi PHPUnit menggunakan SQLite in-memory sehingga tidak menyentuh database pengembangan.
- Frontend: `npm run lint` dan `npm run build`.

Jalankan tes lokal dari direktori masing-masing proyek sebelum push.

## 5. Checklist sebelum rilis

- [ ] File `.env` dan kredensial tidak masuk Git.
- [ ] `APP_DEBUG=false`; pesan error publik tidak memuat stack trace atau rahasia.
- [ ] Endpoint perubahan data dilindungi autentikasi dan otorisasi.
- [ ] Validasi diterapkan pada input pengguna.
- [ ] Rate limit aktif pada endpoint login.
- [ ] CORS hanya mengizinkan origin frontend yang benar-benar dipakai.
- [ ] Akses publik menggunakan HTTPS; port backend hanya terbuka untuk localhost/reverse proxy.
- [ ] Dokumentasi `openapi.yaml` sesuai implementasi terbaru.
- [ ] Semua tes otomatis lulus dan build frontend berhasil.
- [ ] Backup database dan prosedur pemulihan sudah diuji.

Contoh pencadangan kedua database MySQL (jalankan di server dan simpan hasil di lokasi aman di luar web root):

```bash
mysqldump --single-transaction --routines --triggers -u pemweb2 -p pemweb2_laravel > pemweb2_laravel-$(date +%F).sql
mysqldump --single-transaction --routines --triggers -u pemweb2 -p pemweb2_fiber > pemweb2_fiber-$(date +%F).sql
```

Jangan menaruh hasil backup, file environment, sertifikat privat, atau kredensial di repositori Git.
