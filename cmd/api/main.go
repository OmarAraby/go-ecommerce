package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/OmarAraby/go-ecommerce/config"
	internlapi "github.com/OmarAraby/go-ecommerce/internal/api"
	"github.com/OmarAraby/go-ecommerce/internal/api/health"
	apimw "github.com/OmarAraby/go-ecommerce/internal/api/middleware"
	"github.com/OmarAraby/go-ecommerce/internal/api/products"
	"github.com/OmarAraby/go-ecommerce/internal/application"
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
	productSvc     := application.NewProductService(productRepo)
	productHandler := products.NewHandler(productSvc)
	healthHandler  := health.NewHandler(pool)

	// Register routes
	mux := http.NewServeMux()
	internlapi.RegisterRoutes(mux, healthHandler, productHandler)

	// Apply middleware — order matters: Recovery is outermost, then Logging, then CORS
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
