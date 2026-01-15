// Package server provides the HTTP server for {{.ProjectName}}.
package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/{{.ProjectName}}/backend/internal/auth"
	"github.com/{{.ProjectName}}/backend/internal/config"
	custommw "github.com/{{.ProjectName}}/backend/internal/middleware"
	"github.com/{{.ProjectName}}/backend/internal/supabase"
)

// Server wraps the Echo instance and configuration
type Server struct {
	echo        *echo.Echo
	config      *config.Config
	supabase    *supabase.Client
	authHandler *auth.Handler
	jwtConfig   custommw.JWTConfig
}

// New creates a new server instance with middleware configured
func New(cfg *config.Config) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize Supabase client
	supabaseClient := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey)

	// Initialize auth handler
	authHandler := auth.NewHandler(supabaseClient)

	// JWT configuration for protected routes
	jwtConfig := custommw.JWTConfig{
		JWTSecret: cfg.SupabaseJWTSecret,
	}

	return &Server{
		echo:        e,
		config:      cfg,
		supabase:    supabaseClient,
		authHandler: authHandler,
		jwtConfig:   jwtConfig,
	}
}

// Start binds routes and starts the HTTP server
func (s *Server) Start() error {
	// Bind all routes
	s.bindRoutes()

	// Start server
	return s.echo.Start(":" + s.config.Port)
}
