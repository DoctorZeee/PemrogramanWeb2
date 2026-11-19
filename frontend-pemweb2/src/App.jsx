import { BrowserRouter, Route, Routes, Navigate } from 'react-router-dom'
import Login from './halaman/Login'
import DaftarMahasiswa from './halaman/DaftarMahasiswa'
import FormulirMahasiswa from './halaman/FormulirMahasiswa'
import RuteTerlindungi from './komponen/RuteTerlindungi'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/mahasiswa" replace />} />
        <Route path="/login" element={<Login />} />
        <Route
          path="/mahasiswa"
          element={
            <RuteTerlindungi>
              <DaftarMahasiswa />
            </RuteTerlindungi>
          }
        />
        <Route
          path="/mahasiswa/baru"
          element={
            <RuteTerlindungi>
              <FormulirMahasiswa />
            </RuteTerlindungi>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}
