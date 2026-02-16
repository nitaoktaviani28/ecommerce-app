package main

import (
	"log"
	"net/http"

	"ecommerce-app/handlers"
	"ecommerce-app/observability"
	"ecommerce-app/repository"

	// Handler Prometheus untuk expose endpoint /metrics
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
	main merupakan entry point aplikasi.
	Seluruh proses inisialisasi aplikasi dimulai dari fungsi ini,
	mulai dari observability, database, hingga HTTP server.
*/
func main() {

	// =========================
	// INISIALISASI OBSERVABILITY
	// =========================
	// Fungsi Init merupakan single entry point untuk mengaktifkan
	// tracing, profiling, dan metrics.
	// Pendekatan ini menjaga agar logika bisnis tidak bergantung
	// langsung pada detail implementasi observability.
	observability.Init()

	// =========================
	// INISIALISASI DATABASE
	// =========================
	// Membuka koneksi database dan melakukan validasi koneksi.
	// Jika inisialisasi gagal, aplikasi akan dihentikan (fail fast).
	if err := repository.Init(); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	// Menutup koneksi database secara graceful saat aplikasi berhenti.
	defer repository.Close()

	// =========================
	// REGISTRASI HTTP ROUTE
	// =========================
	// Endpoint aplikasi utama
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/checkout", handlers.Checkout)
	http.HandleFunc("/success", handlers.Success)

	// Endpoint metrics untuk Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// =========================
	// MENJALANKAN HTTP SERVER
	// =========================
	// Aplikasi akan mendengarkan request pada port 8080.
	// Jika server gagal dijalankan, aplikasi akan berhenti.
	log.Println("E-commerce app starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
