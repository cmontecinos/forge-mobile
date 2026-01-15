// Package server - route definitions for {{.ProjectName}}.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// bindRoutes registers all HTTP routes
func (s *Server) bindRoutes() {
	// Health check endpoint
	s.echo.GET("/health", s.healthCheck)

	// API v1 group
	api := s.echo.Group("/api/v1")

	// TODO: Add your routes here
	// Example:
	// api.GET("/users", s.getUsers)
	// api.POST("/users", s.createUser)

	_ = api // Placeholder to avoid unused variable error
}

// healthCheck returns the server health status
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
