package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a PostgreSQL connection pool and verifies connectivity with a Ping.
// Returns a closed pool and an error if connection setup fails.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// 1. parse the DSN ==> data source name --> our connection string for the database
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	// 2. Pool sizing — tune for your workload
	cfg.MaxConns = 25                      // max number of connections in the pool
	cfg.MinConns = 2                       // min number of connections in the pool
	cfg.MaxConnLifetime = time.Hour        // max lifetime of a connection
	cfg.MaxConnIdleTime = 30 * time.Minute // idle time before connection is closed

	// 3. create the pool
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	// 4. Verify we can actually reach the database before returning the pool.
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // close the pool if ping fails
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}
