package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func HubungkanBasisData(konf Konfigurasi) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		konf.DbUser,
		konf.DbSandi,
		konf.DbHost,
		konf.DbPort,
		konf.DbNama,
	)

	// Logger Info (seluruh kueri SQL tercetak) hanya pada lingkungan lokal,
	// karena log SQL memuat nilai sensitif seperti hash kata sandi.
	tingkatLog := gormLogger.Warn
	if konf.Lingkungan == "local" {
		tingkatLog = gormLogger.Info
	}

	// TranslateError membuat galat duplikat/kunci asing dapat dikenali lewat
	// errors.Is(err, gorm.ErrDuplicatedKey) tanpa bergantung pada driver.
	koneksi, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         gormLogger.Default.LogMode(tingkatLog),
		TranslateError: true,
	})
	if err != nil {
		log.Fatal("Gagal terhubung ke basis data: ", err)
	}

	DB = koneksi
	log.Println("Koneksi basis data berhasil")
}
