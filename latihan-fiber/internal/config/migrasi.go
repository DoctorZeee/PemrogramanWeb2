package config

import (
	"log"

	"latihan-fiber/internal/model"
)

func Migrasi() {
	err := DB.AutoMigrate(&model.ProgramStudi{}, &model.Mahasiswa{}, &model.MataKuliah{}, &model.Pengguna{}, &model.DaftarHitamToken{})
	if err != nil {
		log.Fatal("Gagal menjalankan migrasi: ", err)
	}

	log.Println("Migrasi selesai")
}

func DataAwal() {
	var jumlah int64
	DB.Model(&model.ProgramStudi{}).Count(&jumlah)
	if jumlah > 0 {
		return
	}

	daftar := []model.ProgramStudi{
		{Kode: "TK", Nama: "Teknik Komputer", Jenjang: "S1"},
		{Kode: "IF", Nama: "Teknik Informatika", Jenjang: "S1"},
		{Kode: "TE", Nama: "Teknik Elektro", Jenjang: "S1"},
	}
	err := DB.Create(&daftar).Error
	if err != nil {
		log.Fatal("Gagal membuat data awal program studi: ", err)
	}

	log.Println("Data awal program studi dibuat")
}
