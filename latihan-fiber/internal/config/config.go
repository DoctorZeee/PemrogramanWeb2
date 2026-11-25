package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Konfigurasi struct {
	NamaAplikasi       string
	Host               string
	Port               string
	Lingkungan         string
	ApiKey             string
	JwtSecret          string
	JwtTtlMenit        int
	CorsAllowedOrigins []string
	DbHost             string
	DbPort             string
	DbNama             string
	DbUser             string
	DbSandi            string
}

func Muat() Konfigurasi {
	godotenv.Load()

	ttl, err := strconv.Atoi(ambil("JWT_TTL_MENIT", "120"))
	if err != nil {
		ttl = 120
	}

	originCORS := strings.Split(ambil("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ",")
	asalDiizinkan := make([]string, 0, len(originCORS))
	for _, asal := range originCORS {
		asal = strings.TrimSpace(asal)
		if asal != "" {
			asalDiizinkan = append(asalDiizinkan, asal)
		}
	}

	return Konfigurasi{
		NamaAplikasi:       ambil("APP_NAME", "Fiber API"),
		Host:               ambil("APP_HOST", ""),
		Port:               ambil("APP_PORT", "3000"),
		Lingkungan:         ambil("APP_ENV", "production"),
		ApiKey:             ambil("API_KEY", ""),
		DbHost:             ambil("DB_HOST", "127.0.0.1"),
		DbPort:             ambil("DB_PORT", "3306"),
		DbNama:             ambil("DB_NAME", "pemweb2_fiber"),
		DbUser:             ambil("DB_USER", "root"),
		DbSandi:            ambil("DB_PASS", ""),
		JwtSecret:          ambil("JWT_SECRET", "rahasia-bawaan"),
		JwtTtlMenit:        ttl,
		CorsAllowedOrigins: asalDiizinkan,
	}
}

func ambil(kunci string, bawaan string) string {
	nilai := os.Getenv(kunci)

	if nilai == "" {
		return bawaan
	}

	return nilai
}
