# frontend-pemweb2 — Frontend React (Modul 9)

Frontend React + Vite yang mengonsumsi RESTful API backend Fiber (`../latihan-fiber`)
atau Laravel (`../latihan-laravel`) melalui satu konfigurasi.

## Menjalankan

```bash
npm install
npm run dev     # http://localhost:5173
```

Backend Fiber:

```bash
cd ../latihan-fiber
go run ./cmd/api    # http://localhost:3000
```

Login memakai akun dari Modul 8 (`POST /api/auth/registrasi`). Peran `admin` dapat
menghapus data; `dosen` hanya dapat menambah; `mahasiswa` hanya melihat.

## Berpindah backend (Tugas 3)

Ubah satu baris pada `.env` lalu restart `npm run dev`:

```
VITE_API_AKTIF=fiber     # atau "laravel"
```

- `VITE_API_FIBER` — baseURL Fiber, `VITE_API_LARAVEL` — baseURL Laravel.
- Perbedaan nama field dan bentuk meta diseragamkan di `src/api/petaData.js`
  (lapisan adapter), sehingga halaman tidak menyentuh kontrak backend mana pun.

## Struktur

```
src/
  api/
    klien.js              axios + interceptor (token 401, redirect /login)
    petaData.js           adapter kontrak Laravel <-> Fiber
    layananAuth.js        login, profil, logout
    layananMahasiswa.js   daftar, detail, simpan, hapus, program studi
  halaman/
    Login.jsx             halaman masuk (Langkah 7)
    DaftarMahasiswa.jsx   daftar + pencarian debounce 400ms + pagination + hapus (Tugas 2)
    FormulirMahasiswa.jsx tambah data + pesan galat validasi server (Tugas 1)
  komponen/
    RuteTerlindungi.jsx   pengaman rute berbasis token (Langkah 9)
  App.jsx                 routing (Langkah 9)
```

## Perbandingan kontrak API kedua backend (Tugas 4)

| Aspek | Laravel (Modul 5) | Fiber (Modul 8) |
|---|---|---|
| Login | `POST /api/login`, body `{email, password}` | `POST /api/auth/login`, body `{email, kata_sandi}` |
| Isi respons sukses | `{...resource}` + `token` pada level atas | `{sukses, pesan, data:{pengguna, token, tipe}}` |
| Profil | `GET /api/user` | `GET /api/auth/profil` |
| Logout | `POST /api/logout` | `POST /api/auth/logout` |
| Query daftar | `page`, `per_page`, `q` | `halaman`, `per_halaman`, `cari`, `angkatan` |
| Meta pagination | `{current_page, last_page, per_page, total}` | `{halaman, per_halaman, total, total_halaman}` |
| Nama field mahasiswa | snake_case Laravel (`program_studi_id`, `nim`, ...) | snake_case Go (`program_studi_id`, `nim`, ...) |
| Validasi 422 | `{message, errors:{field:[pesan]}}` | `{sukses, pesan, galat:[{field, pesan}]}` |
| Pesan galat umum | `message` | `pesan` |
| Bentuk pembungkus | resource langsung / `{data}` (paginasi) | selalu `{sukses, pesan, data, meta}` |
| Port lokal | `127.0.0.1:8000` | `localhost:3000` |

### Usulan kontrak API tunggal

Agar frontend tidak perlu adapter, kedua backend sebaiknya menyepakati:

```jsonc
// Sukses
{ "sukses": true, "pesan": "Login berhasil", "data": { "...": "" }, "meta": null }

// Galat validasi (422)
{
  "sukses": false,
  "pesan": "Data yang dikirim tidak valid",
  "galat": [ { "field": "nim", "pesan": "NIM harus diawali huruf H diikuti angka" } ]
}

// Meta pagination
{ "halaman": 1, "per_halaman": 10, "total": 42, "total_halaman": 5 }
```

Syarat lain: endpoint `/auth/login|logout|profil`, query `halaman/per_halaman/cari`,
body `kata_sandi`, dan semua field snake_case yang sama di kedua backend.

## Catatan keamanan

- Token disimpan pada `sessionStorage`: hilang saat tab ditutup, tapi terekspos ke
  JavaScript (XSS). Cookie `HttpOnly` tidak terbaca JS namun perlu proteksi CSRF.
  Lihat Pertanyaan Pembahasan no. 2.
- Variabel `VITE_*` terkompilasi ke bundel — jangan menyimpan rahasia di situ.
