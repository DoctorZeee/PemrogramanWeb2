// Tugas 3: lapisan abstraksi perbedaan kontrak API Laravel dan Fiber.
// Semua layanan berbicara dengan bentuk kontrak "netral" di sini.
import { backendAktif } from './klien'

// --- Query daftar mahasiswa ---------------------------------------------
// Fiber   : halaman, per_halaman, cari, angkatan
// Laravel : page, per_page, q (per Modul 5)
export function parameterDaftar({ halaman, perHalaman, cari, angkatan }) {
  if (backendAktif === 'laravel') {
    return { page: halaman, per_page: perHalaman, q: cari || undefined }
  }
  return {
    halaman,
    per_halaman: perHalaman,
    cari: cari || undefined,
    angkatan: angkatan || undefined,
  }
}

// --- Meta pagination ------------------------------------------------------
// Fiber   : { halaman, per_halaman, total, total_halaman }
// Laravel : { current_page, last_page, per_page, total }
// Hasil selalu: { halaman, per_halaman, total, total_halaman }
export function normalisasiMeta(respons) {
  const m = respons.data && respons.data.meta
  if (!m) {
    return { halaman: 1, per_halaman: 0, total: 0, total_halaman: 1 }
  }
  if (backendAktif === 'laravel') {
    return {
      halaman: m.current_page,
      per_halaman: m.per_page,
      total: m.total,
      total_halaman: m.last_page,
    }
  }
  return {
    halaman: m.halaman,
    per_halaman: m.per_halaman,
    total: m.total,
    total_halaman: m.total_halaman,
  }
}

export function normalisasiDaftar(respons) {
  const isi = respons.data && respons.data.data
  return { data: Array.isArray(isi) ? isi : [], meta: normalisasiMeta(respons) }
}

// --- Login ----------------------------------------------------------------
// Fiber   : { email, kata_sandi } -> data.token, data.pengguna
// Laravel : { email, password }   -> data.token (Modul 5)
export function badanLogin(email, kataSandi) {
  if (backendAktif === 'laravel') {
    return { email, password: kataSandi }
  }
  return { email, kata_sandi: kataSandi }
}

export function normalisasiLogin(respons) {
  const data = (respons.data && respons.data.data) || {}
  return {
    token: data.token,
    pengguna: data.pengguna || data.user || null,
  }
}

// --- Pesan galat ----------------------------------------------------------
// Fiber: { sukses, pesan, galat?: [{field, pesan}] } (422)
// Laravel: { message, errors?: {field: [..]} } (422)
export function ambilPesanGalat(galat) {
  const respons = galat && galat.response
  if (!respons) {
    return 'Tidak dapat menghubungi server'
  }
  return (
    (respons.data && (respons.data.pesan || respons.data.message)) ||
    `Terjadi galat (HTTP ${respons.status})`
  )
}

export function ambilGalatField(galat) {
  const respons = galat && galat.response
  if (!respons || respons.status !== 422 || !respons.data) {
    return []
  }
  if (Array.isArray(respons.data.galat)) {
    // Fiber: [{field, pesan}]
    return respons.data.galat.map((item) => ({
      field: item.field,
      pesan: item.pesan,
    }))
  }
  if (respons.data.errors) {
    // Laravel: {field: [pesan, ...]}
    return Object.entries(respons.data.errors).map(([field, daftar]) => ({
      field,
      pesan: Array.isArray(daftar) ? daftar[0] : String(daftar),
    }))
  }
  return []
}
