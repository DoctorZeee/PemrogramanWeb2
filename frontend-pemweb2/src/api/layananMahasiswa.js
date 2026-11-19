import klien from './klien'
import { parameterDaftar, normalisasiDaftar } from './petaData'

export async function ambilDaftarMahasiswa(parameter) {
  const respons = await klien.get('/mahasiswa', {
    params: parameterDaftar(parameter),
  })
  return normalisasiDaftar(respons)
}

export async function ambilDetailMahasiswa(id) {
  const respons = await klien.get(`/mahasiswa/${id}`)
  return respons.data.data
}

export async function simpanMahasiswa(data) {
  const respons = await klien.post('/mahasiswa', data)
  return respons.data
}

export async function hapusMahasiswa(id) {
  const respons = await klien.delete(`/mahasiswa/${id}`)
  return respons.data
}

// Dipakai oleh formulir tambah (Tugas 1) untuk mengisi pilihan program studi.
export async function ambilDaftarProgramStudi() {
  const respons = await klien.get('/program-studi')
  const isi = respons.data && respons.data.data
  return Array.isArray(isi) ? isi : []
}
