// Package server - route definitions for {{.ProjectName}}.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	custommw "github.com/{{.ProjectName}}/backend/internal/middleware"
)

// bindRoutes registers all HTTP routes
func (s *Server) bindRoutes() {
	// Health check endpoint
	s.echo.GET("/health", s.healthCheck)

	// Auth routes (public)
	auth := s.echo.Group("/auth")
	auth.POST("/register", s.authHandler.Register)
	auth.POST("/login", s.authHandler.Login)
	auth.POST("/refresh", s.authHandler.Refresh)
	auth.POST("/logout", s.authHandler.Logout)

	// API v1 group (protected routes with JWT auth)
	api := s.echo.Group("/api/v1", custommw.JWTAuth(s.jwtConfig))

	// Item routes
	api.GET("/items", s.itemHandler.ListItems)
	api.GET("/items/:id", s.itemHandler.GetItem)
	api.POST("/items", s.itemHandler.CreateItem)
	api.PATCH("/items/:id", s.itemHandler.UpdateItem)
	api.DELETE("/items/:id", s.itemHandler.DeleteItem)

	// TODO: Add more protected routes here
	// Access user in handlers with: custommw.GetUserID(c), custommw.GetUserEmail(c)
}

// healthCheck returns the server health status
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
