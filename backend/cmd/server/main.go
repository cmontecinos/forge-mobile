// Package main is the entry point for the {{.ProjectName}} server.
package main

import (
	"log"

	"github.com/{{.ProjectName}}/backend/internal/config"
	"github.com/{{.ProjectName}}/backend/internal/server"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
