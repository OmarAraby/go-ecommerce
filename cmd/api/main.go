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
	"github.com/OmarAraby/go-ecommerce/internal/api/products"
	"github.com/OmarAraby/go-ecommerce/internal/api/users"
	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
	userapp "github.com/OmarAraby/go-ecommerce/internal/application/services/user"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/postgres"
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
	productRepo    := postgres.NewProductRepo(pool)
	productSvc     := productapp.NewService(productRepo)
	productHandler := products.NewHandler(productSvc)

	userRepo       := postgres.NewUserRepo(pool)
	refreshRepo    := postgres.NewRefreshTokenRepo(pool)
	resetRepo      := postgres.NewPasswordResetRepo(pool)
	userSvc        := userapp.NewService(userRepo, refreshRepo, resetRepo, cfg.JWTSecret)
	authHandler    := apiauth.NewHandler(userSvc)
	userHandler    := users.NewHandler(userSvc)

	healthHandler  := health.NewHandler(pool)

	// Register routes
	mux := http.NewServeMux()
	internlapi.RegisterRoutes(mux, healthHandler, productHandler, authHandler, userHandler, cfg.JWTSecret)

	// Apply middleware — Recovery outermost, then Logging, then CORS
	var handler http.Handler = mux
	handler = apimw.CORS(handler)
	handler = apimw.Logging(handler)
	handler = apimw.Recovery(handler)

	addr := ":" + cfg.HTTPPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}
