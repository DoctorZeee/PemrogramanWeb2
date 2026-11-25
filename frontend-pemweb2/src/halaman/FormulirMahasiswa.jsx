import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { simpanMahasiswa, ambilDaftarProgramStudi } from '../api/layananMahasiswa'
import { ambilPesanGalat, ambilGalatField } from '../api/petaData'

const TAHUN_INI = new Date().getFullYear()

export default function FormulirMahasiswa() {
  const navigasi = useNavigate()
  const [programStudi, setProgramStudi] = useState([])
  const [bentuk, setBentuk] = useState({
    program_studi_id: '',
    nim: '',
    nama: '',
    email: '',
    angkatan: TAHUN_INI, // DIUBAH: memakai konstanta
    ipk: '',
    aktif: true,
  })
  const [galatField, setGalatField] = useState([])
  const [pesanGalat, setPesanGalat] = useState('')
  const [sedangMemuat, setSedangMemuat] = useState(false)

  useEffect(() => {
    let dibatalkan = false
    ambilDaftarProgramStudi()
      .then((daftar) => {
        if (dibatalkan === false) {
          setProgramStudi(daftar)
          if (daftar.length > 0) {
            setBentuk((b) => ({ ...b, program_studi_id: String(daftar[0].id) }))
          }
        }
      })
      .catch(() => {
        if (dibatalkan === false) {
          setPesanGalat('Gagal memuat daftar program studi')
        }
      })
    return () => {
      dibatalkan = true
    }
  }, [])

  function ubah(kunci, nilai) {
    setBentuk((b) => ({ ...b, [kunci]: nilai }))
  }

  function pesanUntuk(field) {
    const temuan = galatField.find((item) => item.field === field)
    return temuan ? temuan.pesan : null
  }

  async function tanganiKirim(peristiwa) {
    peristiwa.preventDefault()
    setPesanGalat('')
    setGalatField([])
    setSedangMemuat(true)
    try {
      await simpanMahasiswa({
        program_studi_id: Number(bentuk.program_studi_id),
        nim: bentuk.nim,
        nama: bentuk.nama,
        email: bentuk.email,
        angkatan: Number(bentuk.angkatan),
        ipk: Number(bentuk.ipk),
        aktif: bentuk.aktif,
      })
      navigasi('/mahasiswa')
    } catch (galat) {
      setPesanGalat(ambilPesanGalat(galat))
      setGalatField(ambilGalatField(galat))
    } finally {
      setSedangMemuat(false)
    }
  }

  const gayaInput = { width: '100%' }

  return (
    <div style={{ maxWidth: 480, margin: '40px auto' }}>
      <h1>Tambah Mahasiswa</h1>
      <Link to="/mahasiswa">Kembali ke daftar</Link>
      {pesanGalat && <p style={{ color: 'crimson' }}>{pesanGalat}</p>}
      <form onSubmit={tanganiKirim} style={{ marginTop: 16 }}>
        <p>
          <label htmlFor="program_studi_id">Program Studi</label>
          <br />
          <select
            id="program_studi_id"
            value={bentuk.program_studi_id}
            onChange={(peristiwa) => ubah('program_studi_id', peristiwa.target.value)}
            required
            style={gayaInput}
          >
            {programStudi.map((prodi) => (
              <option key={prodi.id} value={prodi.id}>
                {prodi.kode} - {prodi.nama}
              </option>
            ))}
          </select>
          {pesanUntuk('program_studi_id') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('program_studi_id')}</span>
          )}
        </p>
        <p>
          <label htmlFor="nim">NIM</label>
          <br />
          <input
            id="nim"
            type="text"
            value={bentuk.nim}
            onChange={(peristiwa) => ubah('nim', peristiwa.target.value)}
            required
            style={gayaInput}
          />
          {pesanUntuk('nim') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('nim')}</span>
          )}
        </p>
        <p>
          <label htmlFor="nama">Nama</label>
          <br />
          <input
            id="nama"
            type="text"
            value={bentuk.nama}
            onChange={(peristiwa) => ubah('nama', peristiwa.target.value)}
            required
            style={gayaInput}
          />
          {pesanUntuk('nama') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('nama')}</span>
          )}
        </p>
        <p>
          <label htmlFor="email">Email</label>
          <br />
          <input
            id="email"
            type="email"
            value={bentuk.email}
            onChange={(peristiwa) => ubah('email', peristiwa.target.value)}
            required
            style={gayaInput}
          />
          {pesanUntuk('email') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('email')}</span>
          )}
        </p>
        <p>
          <label htmlFor="angkatan">Angkatan</label>
          <br />
          <input
            id="angkatan"
            type="number"
            min="2000"
            max="2100"
            value={bentuk.angkatan}
            onChange={(peristiwa) => ubah('angkatan', peristiwa.target.value)}
            required
            style={gayaInput}
          />
          {pesanUntuk('angkatan') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('angkatan')}</span>
          )}
        </p>
        <p>
          <label htmlFor="ipk">IPK (0 - 4)</label>
          <br />
          <input
            id="ipk"
            type="number"
            min="0"
            max="4"
            step="0.01"
            value={bentuk.ipk}
            onChange={(peristiwa) => ubah('ipk', peristiwa.target.value)}
            required
            style={gayaInput}
          />
          {pesanUntuk('ipk') && (
            <span style={{ color: 'crimson' }}>{pesanUntuk('ipk')}</span>
          )}
        </p>
        <p>
          <label htmlFor="aktif">
            <input
              id="aktif"
              type="checkbox"
              checked={bentuk.aktif}
              onChange={(peristiwa) => ubah('aktif', peristiwa.target.checked)}
            />{' '}
            Status aktif
          </label>
        </p>
        <button type="submit" disabled={sedangMemuat}>
          {sedangMemuat ? 'Menyimpan' : 'Simpan'}
        </button>
      </form>
    </div>
  )
}