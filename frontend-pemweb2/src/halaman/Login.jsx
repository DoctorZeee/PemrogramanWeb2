import { useState } from 'react'
import { useNavigate, Navigate } from 'react-router-dom'
import { login } from '../api/layananAuth'
import { ambilPesanGalat } from '../api/petaData'

export default function Login() {
  const navigasi = useNavigate()
  const [email, setEmail] = useState('')
  const [kataSandi, setKataSandi] = useState('')
  const [pesanGalat, setPesanGalat] = useState('')
  const [sedangMemuat, setSedangMemuat] = useState(false)

  if (sessionStorage.getItem('token')) {
    return <Navigate to="/mahasiswa" replace />
  }

  async function tanganiKirim(peristiwa) {
    peristiwa.preventDefault()
    setPesanGalat('')
    setSedangMemuat(true)
    try {
      const hasil = await login(email, kataSandi)
      sessionStorage.setItem('token', hasil.token)
      if (hasil.pengguna) {
        sessionStorage.setItem('pengguna', JSON.stringify(hasil.pengguna))
      }
      navigasi('/mahasiswa')
    } catch (galat) {
      setPesanGalat(ambilPesanGalat(galat))
    } finally {
      setSedangMemuat(false)
    }
  }

  return (
    <div style={{ maxWidth: 360, margin: '80px auto' }}>
      <h1>Masuk</h1>
      {pesanGalat && <p style={{ color: 'crimson' }}>{pesanGalat}</p>}
      <form onSubmit={tanganiKirim}>
        <div>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(peristiwa) => setEmail(peristiwa.target.value)}
            required
            style={{ width: '100%' }}
          />
        </div>
        <div style={{ marginTop: 12 }}>
          <label htmlFor="kata-sandi">Kata Sandi</label>
          <input
            id="kata-sandi"
            type="password"
            value={kataSandi}
            onChange={(peristiwa) => setKataSandi(peristiwa.target.value)}
            required
            style={{ width: '100%' }}
          />
        </div>
        <button type="submit" disabled={sedangMemuat} style={{ marginTop: 16 }}>
          {sedangMemuat ? 'Memproses' : 'Masuk'}
        </button>
      </form>
    </div>
  )
}
