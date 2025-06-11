package processors

import (
	"fmt"
	"reflect"
	// Anda mungkin butuh library sanitasi jika diperlukan
	// "github.com/microcosm-cc/bluemonday"
)

type Meta struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
	ImageURL    string `json:"img"`
}

// getFieldStringValue adalah helper yang menggunakan reflection untuk mencari nilai string
// dari sebuah struct dengan mencoba beberapa nama field secara berurutan.
func getFieldStringValue(data any, fieldNames ...string) string {
	val := reflect.ValueOf(data)

	// Jika data adalah pointer, kita butuh elemen yang ditunjuknya
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Pastikan kita bekerja dengan struct
	if val.Kind() != reflect.Struct {
		return ""
	}

	// Coba setiap nama field yang diberikan
	for _, fieldName := range fieldNames {
		field := val.FieldByName(fieldName)
		// Periksa apakah field ada dan tipenya string
		if field.IsValid() && field.Kind() == reflect.String {
			value := field.String()
			if value != "" {
				return value // Kembalikan nilai pertama yang tidak kosong
			}
		}
	}
	return "" // Kembalikan string kosong jika tidak ada yang ditemukan
}

// GenericPostProcessor membuat meta tags untuk model APAPUN.
func GenericPostProcessor(data any, baseURL string) (any, error) {
	// Gunakan helper untuk mendapatkan nilai dengan fallback
	// Coba ambil dari "MetaTitle", jika tidak ada coba "Title", jika tidak ada coba "Name"
	title := getFieldStringValue(data, "MetaTitle", "Title", "Name")

	// Untuk deskripsi
	description := getFieldStringValue(data, "MetaDescription", "Description")

	// Untuk path gambar
	imgPath := getFieldStringValue(data, "MetaImg", "Image", "Thumbnail", "Banner")

	// Sanitasi deskripsi (opsional, tapi praktik yang baik)
	// p := bluemonday.StripTagsPolicy()
	// sanitizedDescription := p.Sanitize(description)
	// shortDescription := ... (potong hingga 160 karakter)

	// Bangun URL gambar lengkap
	imgURL := ""
	if imgPath != "" {
		// Logika untuk menggabungkan baseURL dengan imgPath
		imgURL = fmt.Sprintf("%s/%s", baseURL, imgPath)
	}

	// Generate keywords sederhana dari judul
	keywords := ""
	if title != "" {
		keywords = title // Anda bisa membuat logic yang lebih canggih di sini
	}

	meta := Meta{
		Title:       title,
		Keywords:    keywords,
		Description: description, // Gunakan sanitizedDescription jika ada
		ImageURL:    imgURL,
	}

	return meta, nil
}
