package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, ph *ProductHandler) {
	mux.HandleFunc("GET /products", ph.List)
	mux.HandleFunc("GET /products/{id}", ph.GetByID)
	mux.HandleFunc("POST /products", ph.Create)
	mux.HandleFunc("PUT /products/{id}", ph.Update)
	mux.HandleFunc("DELETE /products/{id}", ph.Delete)
}
