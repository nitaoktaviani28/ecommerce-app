package handlers

import (
	"html/template"
	"net/http"

	"ecommerce-app/repository"

	// OpenTelemetry API untuk tracing
	"go.opentelemetry.io/otel"
)

/*
	Home adalah HTTP handler untuk halaman utama aplikasi (/).
	Fungsi ini bertanggung jawab untuk:
	- Membuat span tracing level HTTP
	- Mengambil data produk dari database
	- Merender halaman HTML
*/
func Home(w http.ResponseWriter, r *http.Request) {
	// Membuat tracer untuk layer handler
	tracer := otel.Tracer("handlers")

	// Membuat span tracing untuk request ke halaman utama
	// Span ini akan menjadi parent bagi span repository dan query database
	ctx, span := tracer.Start(r.Context(), "home_handler")
	defer span.End()

	// Mengambil daftar produk dari database
	// Context tracing diteruskan agar query database
	// terhubung dengan trace request HTTP
	products, err := repository.GetProducts(ctx)
	if err != nil {
		http.Error(w, "Failed to get products", 500)
		return
	}

	// Parsing template HTML halaman utama
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}

	// Merender halaman HTML dengan data produk
	tmpl.Execute(w, products)
}
