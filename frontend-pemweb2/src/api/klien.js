import axios from 'axios'

// Tugas 3: berpindah backend cukup dengan VITE_API_AKTIF pada .env.
const backend = import.meta.env.VITE_API_AKTIF === 'laravel' ? 'laravel' : 'fiber'
export const backendAktif = backend

const baseURL =
  backend === 'laravel'
    ? import.meta.env.VITE_API_LARAVEL
    : import.meta.env.VITE_API_FIBER

const klien = axios.create({
  baseURL,
  headers: { Accept: 'application/json' },
  timeout: 10000,
})

// Interceptor permintaan: menyisipkan token pada setiap panggilan.
klien.interceptors.request.use((konfigurasi) => {
  const token = sessionStorage.getItem('token')
  if (token) {
    konfigurasi.headers.Authorization = `Bearer ${token}`
  }
  return konfigurasi
})

// Interceptor respons: menangani token kedaluwarsa secara terpusat.
klien.interceptors.response.use(
  (respons) => respons,
  (galat) => {
    if (galat.response && galat.response.status === 401) {
      sessionStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(galat)
  }
)

export default klien
