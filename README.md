# 🛒 E-Commerce App (Go) — LGTM Observability Ready

Contoh aplikasi **E-Commerce Backend** berbasis **Go (monolith)** dengan frontend sederhana (HTML templates)
yang sudah terintegrasi **Grafana LGTM Stack** (Loki, Grafana, Tempo, Mimir, Pyroscope).

Repository ini dibuat untuk menunjukkan **pola pemisahan yang jelas** antara:
- **Business Logic (Developer)**
- **Observability & Infrastruktur (DevOps)**

Tujuan utama:  
👉 business code tetap bersih  
👉 observability bisa ditambahkan **tanpa mengotori logic aplikasi**

---

## ✨ Fitur Utama

- Backend Go (`net/http`)
- Frontend HTML (server-side rendering)
- PostgreSQL sebagai database
- Flow sederhana: product → checkout → order
- Observability lengkap:
  - **Tracing** → Tempo
  - **Metrics** → Mimir (Prometheus-compatible)
  - **Logs** → Loki
  - **Profiling (CPU & Memory)** → Pyroscope

---

## 🧠 Arsitektur & Ownership

### Struktur Project

<img width="332" height="524" alt="image" src="https://github.com/user-attachments/assets/2be19953-21c5-4731-814e-08b6e6b41088" />


---

### Ownership Rule

**Developer bertanggung jawab atas:**
- `handlers/`
- `repository/`
- `templates/`
- seluruh business logic aplikasi

**DevOps bertanggung jawab atas:**
- `observability/`
- Dockerfile
- deployment
- integrasi LGTM Stack (Loki, Grafana, Tempo, Mimir, Pyroscope)

> Business code **tidak perlu tahu** tentang Tempo / Pyroscope  
> Observability cukup di-*inject* lewat satu function call

---

## 🔌 Cara Integrasi Observability (Paling Penting)

Integrasi observability di `main.go` **cukup satu baris**:

```go
import "ecommerce-app/observability"

func main() {
    observability.Init()

    // business code di bawah ini
}


