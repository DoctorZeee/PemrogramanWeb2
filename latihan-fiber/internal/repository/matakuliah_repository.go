package repository

import (
	"latihan-fiber/internal/config"
	"latihan-fiber/internal/model"
)

type FilterMataKuliah struct {
	Cari       string
	Semester   int
	Halaman    int
	PerHalaman int
}

func AmbilMataKuliah(filter FilterMataKuliah) ([]model.MataKuliah, int64, error) {
	var daftar []model.MataKuliah
	var total int64

	kueri := config.DB.Model(&model.MataKuliah{})

	if filter.Cari != "" {
		pola := "%" + filter.Cari + "%"
		kueri = kueri.Where("nama LIKE ? OR kode LIKE ?", pola, pola)
	}
	if filter.Semester > 0 {
		kueri = kueri.Where("semester = ?", filter.Semester)
	}

	err := kueri.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	lompat := (filter.Halaman - 1) * filter.PerHalaman
	err = kueri.Order("semester ASC, nama ASC").
		Limit(filter.PerHalaman).
		Offset(lompat).
		Find(&daftar).Error
	if err != nil {
		return nil, 0, err
	}

	return daftar, total, nil
}

func AmbilMataKuliahPerId(id uint) (model.MataKuliah, error) {
	var matakuliah model.MataKuliah
	err := config.DB.First(&matakuliah, id).Error
	return matakuliah, err
}

func SimpanMataKuliah(matakuliah *model.MataKuliah) error {
	return config.DB.Create(matakuliah).Error
}

func PerbaruiMataKuliah(matakuliah *model.MataKuliah) error {
	return config.DB.Save(matakuliah).Error
}

func HapusMataKuliah(matakuliah *model.MataKuliah) error {
	return config.DB.Delete(matakuliah).Error
}
