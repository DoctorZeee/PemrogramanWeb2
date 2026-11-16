package validasi

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// reNim mencocokkan NIM yang diawali huruf H diikuti hanya angka.
var reNim = regexp.MustCompile(`^H\d+$`)

func init() {
	err := DaftarkanValidasiKustom("nim", func(fl validator.FieldLevel) bool {
		return reNim.MatchString(fl.Field().String())
	})
	if err != nil {
		panic(err)
	}
}
