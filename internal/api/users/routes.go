package users

import (
	"net/http"

	apimw "github.com/OmarAraby/go-ecommerce/internal/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, jwtSecret string) {
	requireAuth := apimw.RequireAuth(jwtSecret)

	mux.Handle("GET /users/me", requireAuth(http.HandlerFunc(h.GetProfile)))
	mux.Handle("PUT /users/me", requireAuth(http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("PUT /users/me/email", requireAuth(http.HandlerFunc(h.ChangeEmail)))
	mux.Handle("PUT /users/me/password", requireAuth(http.HandlerFunc(h.ChangePassword)))
}
