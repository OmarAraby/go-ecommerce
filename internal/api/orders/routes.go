package orders

import (
	"net/http"

	apimw "github.com/OmarAraby/go-ecommerce/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, jwtSecret string) {
	requireAuth := apimw.RequireAuth(jwtSecret)

	mux.Handle("POST /orders", requireAuth(http.HandlerFunc(h.Create)))
	mux.Handle("GET /orders", requireAuth(http.HandlerFunc(h.List)))
	mux.Handle("GET /orders/{id}", requireAuth(http.HandlerFunc(h.GetByID)))
}
