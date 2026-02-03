# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary and templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# Expose port
EXPOSE 8080

# Set environment variables dengan default values
#ENV OTEL_SERVICE_NAME=backend
#ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://alloy.monitoring.svc.cluster.local:4318
#ENV OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
#ENV OTEL_TRACES_SAMPLER=always_on
#NV OTEL_METRICS_EXPORTER=none
#ENV DATABASE_DSN=postgres://user:password@postgres:5432/dbname?sslmode=disable
#NV PYROSCOPE_ENDPOINT=http://pyroscope-distributor.monitoring.svc.cluster.local:4040

# Run the application
CMD ["./main"]
