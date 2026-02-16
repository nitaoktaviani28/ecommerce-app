package repository

import (
	"context"
	"database/sql"
	"os"
	"time"

	// Driver PostgreSQL
	_ "github.com/lib/pq"

	// otelsql membungkus database/sql agar setiap operasi database
	// (query, exec, query row) otomatis menghasilkan span OpenTelemetry
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	// OpenTelemetry API untuk pembuatan tracer dan span secara manual
	"go.opentelemetry.io/otel"
)

// db merupakan connection pool global yang digunakan oleh seluruh aplikasi
var db *sql.DB

// Product merepresentasikan entitas produk pada database
type Product struct {
	ID    int
	Name  string
	Price float64
}

// Order merepresentasikan entitas pesanan pada database
type Order struct {
	ID        int
	ProductID int
	Quantity  int
	Total     float64
	CreatedAt time.Time
}

// Init menginisialisasi koneksi database.
// Penggunaan otelsql.Open memastikan seluruh query database
// terinstrumentasi dan dapat dikirim sebagai trace ke Tempo.
func Init() error {
	// DATABASE_DSN diambil dari environment variable
	// untuk mendukung deployment di berbagai environment
	dsn := getEnv("DATABASE_DSN", "postgres://user:pass@postgres:5432/shop?sslmode=disable")

	var err error

	// otelsql.Open menggantikan sql.Open biasa dan
	// mengaktifkan tracing otomatis pada seluruh operasi SQL
	db, err = otelsql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// Validasi koneksi database saat aplikasi dijalankan
	// Jika gagal, aplikasi akan dihentikan (fail fast)
	if err := db.Ping(); err != nil {
		return err
	}

	// Inisialisasi skema database
	// Seluruh query di tahap ini juga akan tercatat dalam trace
	return setupTables()
}

// Close menutup koneksi database secara graceful
// Dipanggil saat aplikasi dihentikan
func Close() {
	if db != nil {
		db.Close()
	}
}

// setupTables membuat tabel dan data awal.
// Span manual dibuat agar seluruh query inisialisasi
// tergabung dalam satu segmen trace.
func setupTables() error {
	tracer := otel.Tracer("repository")

	// Membuat span untuk proses inisialisasi database
	ctx, span := tracer.Start(context.Background(), "setup_tables")
	defer span.End()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			price DECIMAL(10,2)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			product_id INTEGER,
			quantity INTEGER,
			total DECIMAL(10,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// ExecContext digunakan agar query SQL membawa context tracing
	// sehingga dapat muncul sebagai span pada Tempo
	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	// QueryContext memastikan query SELECT juga tercatat dalam trace
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&count)

	// Menambahkan data awal jika tabel produk masih kosong
	if count == 0 {
		products := []Product{
			{Name: "Gaming Laptop", Price: 15000000},
			{Name: "Wireless Mouse", Price: 300000},
			{Name: "Mechanical Keyboard", Price: 800000},
			{Name: "4K Monitor", Price: 3500000},
		}
		for _, p := range products {
			db.ExecContext(
				ctx,
				"INSERT INTO products (name, price) VALUES ($1, $2)",
				p.Name,
				p.Price,
			)
		}
	}

	return nil
}

// GetProducts mengambil seluruh data produk.
// Context dari handler HTTP diteruskan agar trace HTTP
// terhubung dengan trace query database.
func GetProducts(ctx context.Context) ([]Product, error) {
	tracer := otel.Tracer("repository")

	// Span ini menjadi child dari span HTTP request
	ctx, span := tracer.Start(ctx, "get_products")
	defer span.End()

	rows, err := db.QueryContext(
		ctx,
		"SELECT id, name, price FROM products ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Price)
		products = append(products, p)
	}

	return products, nil
}

// GetProduct mengambil data produk berdasarkan ID.
// Query database pada fungsi ini otomatis ter-trace ke Tempo.
func GetProduct(ctx context.Context, id int) (*Product, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "get_product")
	defer span.End()

	var p Product
	err := db.QueryRowContext(
		ctx,
		"SELECT id, name, price FROM products WHERE id=$1",
		id,
	).Scan(&p.ID, &p.Name, &p.Price)

	return &p, err
}

// CreateOrder menyimpan data pesanan baru ke database.
// Proses INSERT dan latensinya akan tercatat dalam trace.
func CreateOrder(ctx context.Context, productID, quantity int, total float64) (int, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "create_order")
	defer span.End()

	var id int
	err := db.QueryRowContext(
		ctx,
		"INSERT INTO orders (product_id, quantity, total) VALUES ($1,$2,$3) RETURNING id",
		productID,
		quantity,
		total,
	).Scan(&id)

	return id, err
}

// GetOrder mengambil data pesanan berdasarkan ID.
// Context tracing diteruskan agar query database
// menjadi bagian dari trace end-to-end.
func GetOrder(ctx context.Context, id int) (*Order, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "get_order")
	defer span.End()

	var o Order
	err := db.QueryRowContext(
		ctx,
		"SELECT id, product_id, quantity, total, created_at FROM orders WHERE id=$1",
		id,
	).Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt)

	return &o, err
}

// getEnv mengambil nilai environment variable.
// Jika tidak tersedia, nilai default akan digunakan.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
