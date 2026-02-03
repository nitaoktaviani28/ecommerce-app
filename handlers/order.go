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

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration"},
		[]string{"method", "endpoint"},
	)
	ordersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "orders_created_total", Help: "Total orders created"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, ordersCreatedTotal)
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	tracer := otel.Tracer("handlers")
	ctx, span := tracer.Start(r.Context(), "checkout_handler")
	defer span.End()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		recordMetrics(r, 405, start)
		return
	}

	// CPU simulation for profiling
	simulateCPUWork()

	productID, _ := strconv.Atoi(r.FormValue("product_id"))
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	product, err := repository.GetProduct(ctx, productID)
	if err != nil {
		http.Error(w, "Product not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	total := product.Price * float64(quantity)
	orderID, err := repository.CreateOrder(ctx, productID, quantity, total)
	if err != nil {
		http.Error(w, "Failed to create order", 500)
		recordMetrics(r, 500, start)
		return
	}

	ordersCreatedTotal.Inc()
	log.Printf("method=POST path=/checkout status=200 duration=%v order_id=%d", time.Since(start), orderID)

	recordMetrics(r, 200, start)
	http.Redirect(w, r, fmt.Sprintf("/success?order_id=%d", orderID), 303)
}

func Success(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	tracer := otel.Tracer("handlers")
	ctx, span := tracer.Start(r.Context(), "success_handler")
	defer span.End()

	orderID, _ := strconv.Atoi(r.URL.Query().Get("order_id"))
	order, err := repository.GetOrder(ctx, orderID)
	if err != nil {
		http.Error(w, "Order not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	product, err := repository.GetProduct(ctx, order.ProductID)
	if err != nil {
		http.Error(w, "Product not found", 404)
		recordMetrics(r, 404, start)
		return
	}

	tmpl, err := template.ParseFiles("templates/success.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		recordMetrics(r, 500, start)
		return
	}

	data := struct {
		Order   repository.Order
		Product repository.Product
	}{*order, *product}

	tmpl.Execute(w, data)
	recordMetrics(r, 200, start)
}

func recordMetrics(r *http.Request, status int, start time.Time) {
	duration := time.Since(start)
	httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
	log.Printf("method=%s path=%s status=%d duration=%v", r.Method, r.URL.Path, status, duration)
}

func simulateCPUWork() {
	for i := 0; i < 2000000; i++ {
		_ = i * i * i
	}
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	runtime.GC()
}
