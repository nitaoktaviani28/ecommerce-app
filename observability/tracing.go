package observability

import (
	"context"

	// OpenTelemetry API global
	"go.opentelemetry.io/otel"

	// Exporter OTLP berbasis HTTP untuk mengirim trace ke collector (Alloy / Tempo)
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	// Resource digunakan untuk mendefinisikan identitas service
	"go.opentelemetry.io/otel/sdk/resource"

	// SDK tracing OpenTelemetry
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	// Semantic conventions untuk atribut standar OpenTelemetry
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// initTracing menginisialisasi sistem tracing menggunakan OpenTelemetry.
// Fungsi ini menyiapkan exporter, identitas service, dan tracer provider,
// kemudian mendaftarkannya secara global agar dapat digunakan
// oleh seluruh bagian aplikasi.
func initTracing() error {
	// Membuat exporter OTLP HTTP untuk mengirim trace ke collector.
	// Endpoint diambil dari environment variable agar fleksibel
	// untuk berbagai environment (local, staging, production).
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(
			GetEnv(
				"OTEL_EXPORTER_OTLP_ENDPOINT",
				"http://alloy.monitoring.svc.cluster.local:4318",
			),
		),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	// Mendefinisikan resource OpenTelemetry sebagai identitas service.
	// service.name digunakan oleh Tempo dan Grafana untuk mengelompokkan trace.
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(
				GetEnv("OTEL_SERVICE_NAME", "ecommerce-backend"),
			),
		),
	)
	if err != nil {
		return err
	}

	// Membuat TracerProvider sebagai mesin utama tracing.
	// Batcher digunakan agar pengiriman trace lebih efisien.
	// Sampler diset AlwaysSample untuk keperluan observability dan pembelajaran.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Mendaftarkan TracerProvider secara global.
	// Setelah ini, seluruh tracer di aplikasi akan menggunakan provider ini.
	otel.SetTracerProvider(tp)

	return nil
}
