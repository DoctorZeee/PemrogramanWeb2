import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ambilDaftarMahasiswa, hapusMahasiswa } from '../api/layananMahasiswa'
import { ambilPesanGalat } from '../api/petaData'

export default function DaftarMahasiswa() {
  const [daftar, setDaftar] = useState([])
  const [meta, setMeta] = useState(null)
  const [halaman, setHalaman] = useState(1)
  const [cari, setCari] = useState('')
  const [sedangMemuat, setSedangMemuat] = useState(false)
  const [pesanGalat, setPesanGalat] = useState(false)
  const [pesanSukses, setPesanSukses] = useState('')
  const [sedangMenghapus, setSedangMenghapus] = useState(null)
  const [versi, setVersi] = useState(0)

  const peran = JSON.parse(sessionStorage.getItem('pengguna') || 'null')?.peran
  const bolehHapus = peran === 'admin'

  useEffect(() => {
    let dibatalkan = false

    async function muatData() {
      setSedangMemuat(true)
      setPesanGalat('')
      try {
        const hasil = await ambilDaftarMahasiswa({
          halaman,
          perHalaman: 10,
          cari,
        })
        if (dibatalkan === false) {
          setDaftar(hasil.data)
          setMeta(hasil.meta)
        }
      } catch (galat) {
        if (dibatalkan === false) {
          setPesanGalat(ambilPesanGalat(galat))
        }
      } finally {
        if (dibatalkan === false) {
          setSedangMemuat(false)
        }
      }
    }

    // Debounce 400 ms: permintaan hanya dikirim setelah pengguna berhenti mengetik.
    const penunda = setTimeout(muatData, 400)

    return () => {
      dibatalkan = true
      clearTimeout(penunda)
    }
  }, [halaman, cari, versi])

  async function tanganiHapus(mahasiswa) {
    const yakin = window.confirm(`Hapus mahasiswa ${mahasiswa.nama} (${mahasiswa.nim})?`)
    if (!yakin) return

    setSedangMenghapus(mahasiswa.id)
    setPesanGalat('')
    setPesanSukses('')
    try {
      await hapusMahasiswa(mahasiswa.id)
      setPesanSukses(`Data ${mahasiswa.nama} berhasil dihapus`)
      // Jika item terakhir pada halaman terakhir terhapus, mundur satu halaman.
      if (daftar.length === 1 && halaman > 1) {
        setHalaman(halaman - 1)
      } else {
        setVersi((v) => v + 1)
      }
    } catch (galat) {
      setPesanGalat(ambilPesanGalat(galat))
    } finally {
      setSedangMenghapus(null)
    }
  }

  function keluar() {
    sessionStorage.removeItem('token')
    sessionStorage.removeItem('pengguna')
    window.location.href = '/login'
  }

  return (
    <div style={{ maxWidth: 900, margin: '40px auto' }}>
      <h1>Daftar Mahasiswa</h1>
      <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
        <input
          type="search"
          placeholder="Cari nama atau NIM"
          value={cari}
          onChange={(peristiwa) => {
            setHalaman(1)
            setCari(peristiwa.target.value)
          }}
        />
        <Link to="/mahasiswa/baru">Tambah Mahasiswa</Link>
        <button onClick={keluar}>Keluar</button>
      </div>
      {sedangMemuat && <p>Memuat data</p>}
      {pesanGalat && <p style={{ color: 'crimson' }}>{pesanGalat}</p>}
      {pesanSukses && <p style={{ color: 'seagreen' }}>{pesanSukses}</p>}
      {sedangMemuat === false && daftar.length === 0 && pesanGalat === '' && (
        <p>Data tidak ditemukan</p>
      )}
      {daftar.length > 0 && (
        <table border="1" cellPadding="8" style={{ width: '100%', marginTop: 16, borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th>NIM</th>
              <th>Nama</th>
              <th>Angkatan</th>
              <th>IPK</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((mahasiswa) => (
              <tr key={mahasiswa.id}>
                <td>{mahasiswa.nim}</td>
                <td>{mahasiswa.nama}</td>
                <td>{mahasiswa.angkatan}</td>
                <td>{mahasiswa.ipk}</td>
                <td>
                  {bolehHapus && (
                    <button
                      onClick={() => tanganiHapus(mahasiswa)}
                      disabled={sedangMenghapus !== null}
                    >
                      {sedangMenghapus === mahasiswa.id ? 'Menghapus' : 'Hapus'}
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {meta && (
        <div style={{ marginTop: 16 }}>
          <button onClick={() => setHalaman(halaman - 1)} disabled={halaman <= 1}>
            Sebelumnya
          </button>
          <span style={{ margin: '0 12px' }}>
            Halaman {meta.halaman} dari {meta.total_halaman} (total {meta.total})
          </span>
          <button
            onClick={() => setHalaman(halaman + 1)}
            disabled={meta.halaman >= meta.total_halaman}
          >
            Berikutnya
          </button>
        </div>
      )}
    </div>
  )
}
