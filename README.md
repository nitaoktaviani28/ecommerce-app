# 🛒 E-Commerce App (Go) — LGTM Observability Ready

Contoh aplikasi **E-Commerce Backend** berbasis **Go (monolith)** dengan frontend sederhana (HTML template)  
yang sudah terintegrasi **Grafana LGTM Stack** (Logs, Grafana, Tempo, Mimir, Pyroscope).

Fokus utama repository ini adalah **pemisahan yang jelas antara business logic dan observability**  
agar mudah di-maintain, di-review, dan di-scale.

---

## ✨ Fitur Utama

- Backend Go (net/http)
- Frontend HTML (server-side rendering)
- PostgreSQL sebagai database
- Order & product flow sederhana
- Observability lengkap:
  - **Tracing** → Tempo
  - **Metrics** → Mimir (Prometheus-compatible)
  - **Logs** → Loki
  - **Profiling (CPU & Memory)** → Pyroscope

---

## 🧠 Arsitektur & Ownership

Struktur project:

ecommerce-app/
├── main.go # business entrypoint
├── handlers/ # HTTP handlers (business)
│ ├── product.go
│ └── order.go
├── repository/ # database access (business)
│ └── postgres.go
├── templates/ # frontend (HTML)
│ ├── index.html
│ └── success.html
├── observability/ # DEVOPS OWNED MODULE
│ ├── init.go
│ ├── tracing.go
│ ├── profiling.go
│ └── env.go
├── go.mod
└── Dockerfile

### Ownership Rule
- **Developer**: `handlers/`, `repository/`, `templates/`, business logic
- **DevOps**: `observability/`, Dockerfile, deployment, LGTM integration

> Business code **tidak perlu tahu** tentang Tempo / Pyroscope  
> Observability cukup di-*inject* lewat satu function call.

---

## 🔌 Cara Integrasi Observability (Paling Penting)

Di `main.go`, integrasi observability **cukup satu baris**:

```go
import "ecommerce-app/observability"

func main() {
    observability.Init()

    // business code di bawah ini
}
