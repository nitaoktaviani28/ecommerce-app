<img width="762" height="695" alt="image" src="https://github.com/user-attachments/assets/06ff367a-861e-40e1-bcde-9aa6582bf30c" /># 🛒 E-Commerce App (Go) — LGTM Observability Ready

Contoh aplikasi **E-Commerce Backend** berbasis **Go (monolith)** dengan frontend sederhana (HTML template)  
yang sudah terintegrasi **Grafana LGTM Stack** (Logs, Grafana, Tempo, Mimir, Pyroscope).

Fokus utama repository ini adalah **pemisahan yang jelas antara business logic dan observability**  
agar mudah di-maintain, di-review, dan di-scale.

---

## ✨ Fitur Utama

- Backend Go (`net/http`)
- Frontend HTML (server-side rendering)
- PostgreSQL sebagai database
- Alur product & order sederhana
- Observability lengkap:
  - **Tracing** → Tempo
  - **Metrics** → Mimir (Prometheus-compatible)
  - **Logs** → Loki
  - **Profiling (CPU & Memory)** → Pyroscope

---

## 🧠 Arsitektur & Ownership

### Struktur Project

```text
ecommerce-app/
├── main.go                # business entrypoint
│
├── handlers/              # HTTP handlers (business)
│   ├── product.go
│   └── order.go
│
├── repository/            # database access layer (business)
│   └── postgres.go
│
├── templates/             # frontend (HTML templates)
│   ├── index.html
│   └── success.html
│
├── observability/         # DEVOPS OWNED MODULE
│   ├── init.go            # bootstrap observability
│   ├── tracing.go         # OpenTelemetry → Tempo
│   ├── profiling.go       # Pyroscope profiling
│   └── env.go             # environment helpers
│
├── go.mod
└── Dockerfile
...
