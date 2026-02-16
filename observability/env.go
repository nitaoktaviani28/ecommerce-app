package observability

import "os"

// GetEnv digunakan untuk mengambil nilai environment variable.
// Jika environment variable tidak tersedia atau bernilai kosong,
// maka nilai default akan digunakan sebagai pengganti.
// Fungsi ini membantu memisahkan konfigurasi dari kode aplikasi
// sehingga aplikasi dapat dijalankan di berbagai environment
// (local, staging, maupun production) tanpa perubahan source code.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
