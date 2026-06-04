package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/auth"
)

// RequireAuth validates the Bearer token and injects claims into the request context.
// Protected handlers retrieve claims with auth.ContextKey.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Unauthorized(w, "missing or malformed authorization header")
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims, err := auth.ValidateAccessToken(tokenStr, secret)
			if err != nil {
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			// Inject claims so handlers can read user ID, role, etc.
			ctx := context.WithValue(r.Context(), auth.ContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
