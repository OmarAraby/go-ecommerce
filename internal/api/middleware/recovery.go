package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
)

// Recovery catches any panic in downstream handlers, logs the stack trace,
// and returns 500 instead of crashing the whole server.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				response.InternalError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
