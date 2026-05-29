package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	HTTPPort string
	DB       DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN builds the PostgreSQL connection string.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// Load reads configuration from environment variables, using defaults where safe.
// Sensitive values (DB_PASSWORD) are required and cause an error if missing.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"), // required
			Name:     getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
