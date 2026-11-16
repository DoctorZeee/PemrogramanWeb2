package validasi

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var pemeriksa = bangunPemeriksa()

type PesanGalat struct {
	Field string `json:"field"`
	Pesan string `json:"pesan"`
}

// bangunPemeriksa membuat validator yang melaporkan nama field menurut tag
// json (mis. "kata_sandi"), bukan nama field Go (mis. "KataSandi"),
// sehingga cocok dengan body yang dikirim klien.
func bangunPemeriksa() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		nama := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if nama == "-" {
			return ""
		}
		return nama
	})
	return v
}

func Periksa(data any) []PesanGalat {
	err := pemeriksa.Struct(data)
	if err == nil {
		return nil
	}

	var galatField validator.ValidationErrors
	if !errors.As(err, &galatField) {
		// Bukan galat aturan, melainkan salah pakai (mis. data bukan struct).
		return []PesanGalat{{Field: "", Pesan: "Data tidak dapat divalidasi"}}
	}

	daftar := []PesanGalat{}
	for _, item := range galatField {
		daftar = append(daftar, PesanGalat{
			Field: item.Field(),
			Pesan: susunPesan(item),
		})
	}
	return daftar
}

func susunPesan(item validator.FieldError) string {
	switch item.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return "Nilai minimal adalah " + item.Param()
	case "max":
		return "Nilai maksimal adalah " + item.Param()
	case "gte":
		return "Nilai tidak boleh kurang dari " + item.Param()
	case "lte":
		return "Nilai tidak boleh lebih dari " + item.Param()
	case "nim":
		return "NIM harus diawali huruf H diikuti angka"
	default:
		return "Nilai tidak memenuhi aturan " + item.Tag()
	}
}

// DaftarkanValidasiKustom mendaftarkan fungsi validasi kustom pada pemeriksa.
func DaftarkanValidasiKustom(tag string, fn validator.Func) error {
	return pemeriksa.RegisterValidation(tag, fn)
}
