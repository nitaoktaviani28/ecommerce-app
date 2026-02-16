package observability

import (
	"log"
)

// Init berfungsi untuk menginisialisasi seluruh komponen observability.
// Fungsi ini menjadi satu-satunya titik masuk (single entry point)
// yang dipanggil oleh aplikasi utama untuk mengaktifkan tracing,
// profiling, dan metrics.
// Pendekatan ini memastikan logika bisnis tidak bergantung langsung
// pada detail implementasi observability.
func Init() {
	log.Println("Initializing observability...")

	// Menginisialisasi komponen tracing (OpenTelemetry / Tempo)
	// Jika inisialisasi gagal, aplikasi tetap berjalan
	// tanpa menghentikan proses utama.
	if err := initTracing(); err != nil {
		log.Printf("Tracing init failed: %v", err)
	}

	// Menginisialisasi komponen profiling (Grafana Pyroscope)
	// Profiling bersifat opsional dan tidak mempengaruhi
	// fungsionalitas utama aplikasi.
	if err := initProfiling(); err != nil {
		log.Printf("Profiling init failed: %v", err)
	}

	// Menginisialisasi metrics aplikasi.
	// Pada implementasi saat ini, sebagian besar metrics
	// didaftarkan langsung pada layer handler.
	initMetrics()

	log.Println("Observability initialized")
}
