package repository

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
)

var db *sql.DB

type Product struct {
	ID    int
	Name  string
	Price float64
}

type Order struct {
	ID        int
	ProductID int
	Quantity  int
	Total     float64
	CreatedAt time.Time
}

func Init() error {
	dsn := getEnv("DATABASE_DSN", "postgres://user:pass@postgres:5432/shop?sslmode=disable")
	
	var err error
	db, err = otelsql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return setupTables()
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func setupTables() error {
	tracer := otel.Tracer("repository")
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

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	// Insert sample products
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		products := []Product{
			{Name: "Gaming Laptop", Price: 15000000},
			{Name: "Wireless Mouse", Price: 300000},
			{Name: "Mechanical Keyboard", Price: 800000},
			{Name: "4K Monitor", Price: 3500000},
		}
		for _, p := range products {
			db.ExecContext(ctx, "INSERT INTO products (name, price) VALUES ($1, $2)", p.Name, p.Price)
		}
	}

	return nil
}

func GetProducts(ctx context.Context) ([]Product, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "get_products")
	defer span.End()

	rows, err := db.QueryContext(ctx, "SELECT id, name, price FROM products ORDER BY id")
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

func GetProduct(ctx context.Context, id int) (*Product, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "get_product")
	defer span.End()

	var p Product
	err := db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Price)
	return &p, err
}

func CreateOrder(ctx context.Context, productID, quantity int, total float64) (int, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "create_order")
	defer span.End()

	var id int
	err := db.QueryRowContext(ctx,
		"INSERT INTO orders (product_id, quantity, total) VALUES ($1,$2,$3) RETURNING id",
		productID, quantity, total).Scan(&id)
	return id, err
}

func GetOrder(ctx context.Context, id int) (*Order, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "get_order")
	defer span.End()

	var o Order
	err := db.QueryRowContext(ctx,
		"SELECT id, product_id, quantity, total, created_at FROM orders WHERE id=$1", id).
		Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt)
	return &o, err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
