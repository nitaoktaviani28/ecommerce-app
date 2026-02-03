package main

import (
	"log"
	"net/http"

	"ecommerce-app/handlers"
	"ecommerce-app/observability"
	"ecommerce-app/repository"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// =========================
	// OBSERVABILITY (SINGLE INJECTION POINT)
	// =========================
	observability.Init()

	// =========================
	// DATABASE
	// =========================
	if err := repository.Init(); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}
	defer repository.Close()

	// =========================
	// HTTP ROUTES
	// =========================
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/checkout", handlers.Checkout)
	http.HandleFunc("/success", handlers.Success)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("🚀 E-commerce app starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
