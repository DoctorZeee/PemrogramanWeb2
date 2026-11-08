package repository

import (
	"latihan-fiber/internal/config"
	"latihan-fiber/internal/model"
)

// HasilProgramStudi membungkus program studi beserta jumlah mahasiswanya
// hasil agregasi kueri.
type HasilProgramStudi struct {
	model.ProgramStudi
	JumlahMahasiswa int64 `json:"jumlah_mahasiswa"`
}

// AmbilSemuaProgramStudi mengembalikan seluruh program studi beserta jumlah
// mahasiswa aktif pada setiap program studi (Tugas 2).
func AmbilSemuaProgramStudi() ([]HasilProgramStudi, error) {
	var daftar []HasilProgramStudi

	err := config.DB.Model(&model.ProgramStudi{}).
		Select("program_studis.*, COUNT(mahasiswas.id) AS jumlah_mahasiswa").
		Group("program_studis.id").
		Joins("LEFT JOIN mahasiswas ON mahasiswas.program_studi_id = program_studis.id AND mahasiswas.dihapus_pada IS NULL").
		Order("program_studis.kode ASC").
		Find(&daftar).Error

	return daftar, err
}
