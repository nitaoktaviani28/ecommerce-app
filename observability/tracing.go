package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initTracing() error {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://alloy.monitoring.svc.cluster.local:4318")),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(GetEnv("OTEL_SERVICE_NAME", "ecommerce-backend")),
		),
	)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	return nil
}
