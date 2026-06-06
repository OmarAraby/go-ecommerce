package api

import (
	"net/http"

	"github.com/OmarAraby/go-ecommerce/internal/api/auth"
	"github.com/OmarAraby/go-ecommerce/internal/api/health"
	"github.com/OmarAraby/go-ecommerce/internal/api/orders"
	"github.com/OmarAraby/go-ecommerce/internal/api/products"
	"github.com/OmarAraby/go-ecommerce/internal/api/users"
)

func RegisterRoutes(
	mux *http.ServeMux,
	hh *health.Handler,
	ph *products.Handler,
	ah *auth.Handler,
	uh *users.Handler,
	oh *orders.Handler,
	jwtSecret string,
) {
	health.RegisterRoutes(mux, hh)
	auth.RegisterRoutes(mux, ah)
	products.RegisterRoutes(mux, ph, jwtSecret)
	users.RegisterRoutes(mux, uh, jwtSecret)
	orders.RegisterRoutes(mux, oh, jwtSecret)
}
