package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const penerbit = "pemweb2-fiber-api"

type Klaim struct {
	Peran string `json:"peran"`
	jwt.RegisteredClaims
}

func TerbitkanToken(idPengguna uint, peran string, rahasia string, ttlMenit int) (string, error) {
	sekarang := time.Now()

	// jti acak membuat setiap token unik walau diterbitkan pada detik yang
	// sama, sehingga token hasil refresh tidak pernah identik dengan token lama.
	baut := make([]byte, 16)
	if _, err := rand.Read(baut); err != nil {
		return "", err
	}

	klaim := Klaim{
		Peran: peran,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        hex.EncodeToString(baut),
			Subject:   strconv.FormatUint(uint64(idPengguna), 10),
			IssuedAt:  jwt.NewNumericDate(sekarang),
			ExpiresAt: jwt.NewNumericDate(sekarang.Add(time.Duration(ttlMenit) * time.Minute)),
			Issuer:    penerbit,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, klaim)
	return token.SignedString([]byte(rahasia))
}

// PeriksaToken memverifikasi tanda tangan, algoritma, penerbit, dan masa berlaku.
// Opsi WithValidMethods menolak algoritma selain HS256 (termasuk "none"),
// sehingga pemeriksaan metode manual pada keyfunc tidak diperlukan lagi.
func PeriksaToken(nilaiToken string, rahasia string) (*Klaim, error) {
	token, err := jwt.ParseWithClaims(nilaiToken, &Klaim{}, func(t *jwt.Token) (any, error) {
		return []byte(rahasia), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(penerbit),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	klaim, ok := token.Claims.(*Klaim)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return klaim, nil
}

// HashToken menghasilkan sidik jari SHA-256 token. Daftar hitam menyimpan
// sidik jari ini, bukan token mentah, agar kebocoran basis data tidak
// membocorkan token yang masih berlaku.
func HashToken(nilaiToken string) string {
	ringkas := sha256.Sum256([]byte(nilaiToken))
	return hex.EncodeToString(ringkas[:])
}
