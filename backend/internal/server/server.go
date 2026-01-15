// Package server provides the HTTP server for {{.ProjectName}}.
package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/{{.ProjectName}}/backend/internal/config"
)

// Server wraps the Echo instance and configuration
type Server struct {
	echo   *echo.Echo
	config *config.Config
}

// New creates a new server instance with middleware configured
func New(cfg *config.Config) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &Server{
		echo:   e,
		config: cfg,
	}
}

// Start binds routes and starts the HTTP server
func (s *Server) Start() error {
	// Bind all routes
	s.bindRoutes()

	// Start server
	return s.echo.Start(":" + s.config.Port)
}
