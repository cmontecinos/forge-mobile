// Package config handles environment configuration for {{.ProjectName}}.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port        string
	SupabaseURL string
	SupabaseKey string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		SupabaseURL: getEnv("SUPABASE_URL", ""),
		SupabaseKey: getEnv("SUPABASE_KEY", ""),
	}

	return cfg, nil
}

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
