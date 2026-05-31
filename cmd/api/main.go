package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/OmarAraby/go-ecommerce/config"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/postgres"
)

func main() {
	// Load .env if present (dev only). In prod, env vars come from the orchestrator.
	_ = godotenv.Load()

	// 1. Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// 2. Connect to PostgreSQL
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()
	log.Println("connected to postgres ✓")

	// 3. Wire HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(pool))

	addr := ":" + cfg.HTTPPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// healthHandler returns a handler that reports app + database health.
// The pool is injected via closure — this is Go's idiomatic DI for handlers.
func healthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bound DB ping with a short timeout from the request context.
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbOK := pool.Ping(ctx) == nil

		status := http.StatusOK
		if !dbOK {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json") // 1.header
		w.WriteHeader(status)                              // 2.status
		json.NewEncoder(w).Encode(map[string]any{          // 3.body
			"status": "ok",
			"db":     dbOK,
		})
	}
}
