package observability

import (
	"github.com/grafana/pyroscope-go"
)

// initProfiling menginisialisasi fitur profiling aplikasi menggunakan Grafana Pyroscope.
// Profiling dilakukan secara runtime untuk mengumpulkan informasi performa seperti
// penggunaan CPU, alokasi objek, dan penggunaan memori.
// Konfigurasi diambil dari environment variable agar dapat disesuaikan
// dengan berbagai environment tanpa perubahan kode.
func initProfiling() error {
	_, err := pyroscope.Start(pyroscope.Config{
		// ApplicationName digunakan sebagai identitas service di Pyroscope.
		// Nilai ini disamakan dengan service name pada tracing agar
		// profiling dapat dikorelasikan dengan trace.
		ApplicationName: GetEnv("OTEL_SERVICE_NAME", "ecommerce-backend"),

		// ServerAddress menentukan endpoint Pyroscope
		// tempat data profiling dikirimkan.
		ServerAddress: GetEnv(
			"PYROSCOPE_ENDPOINT",
			"http://pyroscope-distributor.monitoring.svc.cluster.local:4040",
		),

		// ProfileTypes menentukan jenis profiling yang diaktifkan.
		// CPU: penggunaan CPU per fungsi
		// AllocObjects: jumlah alokasi objek
		// AllocSpace: penggunaan memori
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
		},
	})
	return err
}

// initMetrics disediakan sebagai hook untuk inisialisasi metrics tambahan.
// Pada implementasi saat ini, metrics aplikasi didaftarkan langsung
// pada layer handler.
func initMetrics() {
	// Inisialisasi metrics tambahan dapat ditambahkan di sini jika diperlukan
}
