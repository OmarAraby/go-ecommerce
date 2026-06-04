package products

import (
	"net/http"

	apimw "github.com/OmarAraby/go-ecommerce/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, jwtSecret string) {
	requireAuth := apimw.RequireAuth(jwtSecret)

	// Public
	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("GET /products/{id}", h.GetByID)

	// Protected
	mux.Handle("POST /products", requireAuth(http.HandlerFunc(h.Create)))
	mux.Handle("PUT /products/{id}", requireAuth(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /products/{id}", requireAuth(http.HandlerFunc(h.Delete)))
	mux.Handle("POST /products/{id}/images", requireAuth(http.HandlerFunc(h.UploadImage)))
	mux.Handle("DELETE /products/{id}/images/{imageId}", requireAuth(http.HandlerFunc(h.DeleteImage)))
	mux.Handle("PUT /products/{id}/images/{imageId}/main", requireAuth(http.HandlerFunc(h.SetMainImage)))
}
