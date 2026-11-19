import klien from './klien'
import { badanLogin, normalisasiLogin } from './petaData'

export async function login(email, kataSandi) {
  const respons = await klien.post('/auth/login', badanLogin(email, kataSandi))
  return normalisasiLogin(respons)
}

export async function ambilProfil() {
  const respons = await klien.get('/auth/profil')
  return respons.data.data
}

export async function logout() {
  const respons = await klien.post('/auth/logout')
  return respons.data
}
