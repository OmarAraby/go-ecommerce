package api

import (
	"net/http"

	"github.com/OmarAraby/go-ecommerce/internal/api/health"
	"github.com/OmarAraby/go-ecommerce/internal/api/products"
)

func RegisterRoutes(mux *http.ServeMux, hh *health.Handler, ph *products.Handler) {
	mux.HandleFunc("GET /health", hh.Check)

	mux.HandleFunc("GET /products", ph.List)
	mux.HandleFunc("GET /products/{id}", ph.GetByID)
	mux.HandleFunc("POST /products", ph.Create)
	mux.HandleFunc("PUT /products/{id}", ph.Update)
	mux.HandleFunc("DELETE /products/{id}", ph.Delete)
}
