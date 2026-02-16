package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"ecommerce-app/repository"

	// Library Prometheus untuk pembuatan metrics aplikasi
	"github.com/prometheus/client_golang/prometheus"

	// OpenTelemetry API untuk tracing
	"go.opentelemetry.io/otel"
)

/*
	DEFINISI METRICS PROMETHEUS
	Metrics ini berada di level HTTP dan bisnis aplikasi.
*/

var (
	// Menghitung total request HTTP berdasarkan method, endpoint, dan status
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Mengukur durasi request HTTP
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration",
		},
		[]string{"method", "endpoint"},
	)

	// Menghitung total order yang berhasil dibuat
	ordersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total orders created",
		},
	)
)

/*
	init() akan dijalankan otomatis saat aplikasi start.
	Seluruh metrics didaftarkan agar dapat di-scrape oleh Prometheus
	melalui endpoint /metrics.
*/
func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		ordersCreatedTotal,
	)
}

/*
	Checkout adalah HTTP handler untuk endpoint POST /checkout.
	Fungsi ini:
	- Membuat span tracing level HTTP
	- Mengakses database (produk & order)
	- Mencatat metrics dan log
*/
func Checkout(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Membuat tracer untuk layer handler
	tracer := otel.Tracer("handlers")

	// Membuat span utama untuk request checkout
	// Span ini akan menjadi parent untuk span database
	ctx, span := tracer.Start(r.Context(), "checkout_handler")
	defer span.End()

	// Validasi method HTTP
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		recordMetrics(r, 405, start)
		return
	}

	/*
		Simulasi beban CPU dan memory.
		Fungsi ini sengaja dibuat untuk menghasilkan data profiling
	yang dapat dianalisis di Grafana Pyroscope.
	*/
	simulateCPUWork()

	// Mengambil parameter dari form request
	productID, _ := strconv.Atoi(r.FormValue("product_id"))
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	// Mengambil data produk dari database
	// Context tracing diteruskan agar query database
	// menjadi bagian dari trace end-to-end
	product, err := repository.GetProduct(ctx, productID)
	if err != nil {
		http.Error(w, "Product not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	// Menghitung total harga
	total := product.Price * float64(quantity)

	// Membuat order baru di database
	orderID, err := repository.CreateOrder(ctx, productID, quantity, total)
	if err != nil {
		http.Error(w, "Failed to create order", 500)
		recordMetrics(r, 500, start)
		return
	}

	// Menambah counter jumlah order yang berhasil dibuat
	ordersCreatedTotal.Inc()

	// Logging berbasis teks untuk korelasi dengan metrics dan trace
	log.Printf(
		"method=POST path=/checkout status=200 duration=%v order_id=%d",
		time.Since(start),
		orderID,
	)

	// Mencatat metrics HTTP
	recordMetrics(r, 200, start)

	// Redirect ke halaman success
	http.Redirect(w, r, fmt.Sprintf("/success?order_id=%d", orderID), 303)
}

/*
	Success adalah HTTP handler untuk endpoint GET /success.
	Fungsi ini:
	- Mengambil data order dan produk dari database
	- Menampilkan halaman HTML sukses
*/
func Success(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	tracer := otel.Tracer("handlers")

	// Membuat span tracing untuk handler success
	ctx, span := tracer.Start(r.Context(), "success_handler")
	defer span.End()

	// Mengambil parameter order_id dari query string
	orderID, _ := strconv.Atoi(r.URL.Query().Get("order_id"))

	// Mengambil data order dari database
	order, err := repository.GetOrder(ctx, orderID)
	if err != nil {
		http.Error(w, "Order not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	// Mengambil data produk berdasarkan ProductID pada order
	product, err := repository.GetProduct(ctx, order.ProductID)
	if err != nil {
		http.Error(w, "Product not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	// Parsing template HTML
	tmpl, err := template.ParseFiles("templates/success.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		recordMetrics(r, 500, start)
		return
	}

	// Data yang dikirim ke template HTML
	data := struct {
		Order   repository.Order
		Product repository.Product
	}{
		*order,
		*product,
	}

	// Render template
	tmpl.Execute(w, data)

	// Mencatat metrics HTTP
	recordMetrics(r, 200, start)
}

/*
	recordMetrics bertanggung jawab mencatat metrics HTTP
	secara terpusat agar konsisten di seluruh handler.
*/
func recordMetrics(r *http.Request, status int, start time.Time) {
	duration := time.Since(start)

	httpRequestsTotal.
		WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).
		Inc()

	httpRequestDuration.
		WithLabelValues(r.Method, r.URL.Path).
		Observe(duration.Seconds())

	log.Printf(
		"method=%s path=%s status=%d duration=%v",
		r.Method,
		r.URL.Path,
		status,
		duration,
	)
}

/*
	simulateCPUWork digunakan untuk menghasilkan beban CPU
	dan alokasi memory agar dapat dianalisis menggunakan
	Grafana Pyroscope.
*/
func simulateCPUWork() {
	for i := 0; i < 2000000; i++ {
		_ = i * i * i
	}

	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	// Memicu garbage collection agar terlihat di profiling
	runtime.GC()
}
