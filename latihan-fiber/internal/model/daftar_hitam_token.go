package model

import "time"

// DaftarHitamToken menyimpan sidik jari (SHA-256) token yang sudah dicabut
// (logout) sebelum kedaluwarsa.
type DaftarHitamToken struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TokenHash       string    `gorm:"size:64;uniqueIndex;not null" json:"token_hash"`
	KedaluwarsaPada time.Time `gorm:"index" json:"kedaluwarsa_pada"`
	CreatedAt       time.Time `json:"dibuat_pada"`
}
