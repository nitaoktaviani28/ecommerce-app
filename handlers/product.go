package handlers

import (
	"html/template"
	"net/http"

	"ecommerce-app/repository"

	"go.opentelemetry.io/otel"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handlers")
	ctx, span := tracer.Start(r.Context(), "home_handler")
	defer span.End()

	products, err := repository.GetProducts(ctx)
	if err != nil {
		http.Error(w, "Failed to get products", 500)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}

	tmpl.Execute(w, products)
}
