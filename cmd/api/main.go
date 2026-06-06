package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/OmarAraby/go-ecommerce/config"
	internlapi "github.com/OmarAraby/go-ecommerce/internal/api"
	apiauth "github.com/OmarAraby/go-ecommerce/internal/api/auth"
	"github.com/OmarAraby/go-ecommerce/internal/api/health"
	apimw "github.com/OmarAraby/go-ecommerce/internal/api/middleware"
	"github.com/OmarAraby/go-ecommerce/internal/api/orders"
	"github.com/OmarAraby/go-ecommerce/internal/api/products"
	"github.com/OmarAraby/go-ecommerce/internal/api/users"
	orderapp "github.com/OmarAraby/go-ecommerce/internal/application/services/order"
	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	userapp "github.com/OmarAraby/go-ecommerce/internal/application/services/user"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/postgres"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/storage"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()
	log.Println("connected to postgres ✓")

	// Wire dependencies
	localStorage   := storage.NewLocalStorage("uploads", "/uploads")
	productRepo    := postgres.NewProductRepo(pool)
	productSvc     := productapp.NewService(productRepo, localStorage)
	productHandler := products.NewHandler(productSvc)

	userRepo       := postgres.NewUserRepo(pool)
	refreshRepo    := postgres.NewRefreshTokenRepo(pool)
	resetRepo      := postgres.NewPasswordResetRepo(pool)
	userSvc        := userapp.NewService(userRepo, refreshRepo, resetRepo, cfg.JWTSecret)
	authHandler    := apiauth.NewHandler(userSvc)
	userHandler    := users.NewHandler(userSvc)

	orderRepo      := postgres.NewOrderRepo(pool)
	orderSvc       := orderapp.NewService(orderRepo, productRepo)
	orderHandler   := orders.NewHandler(orderSvc)

	healthHandler  := health.NewHandler(pool)

	// Register routes
	mux := http.NewServeMux()
	internlapi.RegisterRoutes(mux, healthHandler, productHandler, authHandler, userHandler, orderHandler, cfg.JWTSecret)

	// Serve uploaded files statically
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Apply middleware — outermost first: Recovery → Logging → CORS → RateLimit
	var handler http.Handler = mux
	handler = apimw.RateLimit(100, 20)(handler) // 100 req/s, burst of 20 per IP
	handler = apimw.CORS(handler)
	handler = apimw.Logging(handler)
	handler = apimw.Recovery(handler)

	addr := ":" + cfg.HTTPPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}
