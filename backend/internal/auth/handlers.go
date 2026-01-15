// Package auth provides authentication handlers for the API.
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/{{.ProjectName}}/backend/internal/supabase"
)

// Handler handles authentication requests.
type Handler struct {
	supabase *supabase.Client
}

// NewHandler creates a new auth handler.
func NewHandler(supabaseClient *supabase.Client) *Handler {
	return &Handler{
		supabase: supabaseClient,
	}
}

// Register handles user registration.
// POST /auth/register
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Email and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Password must be at least 6 characters",
		})
	}

	resp, err := h.supabase.SignUp(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "auth_error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		User: User{
			ID:        resp.User.ID,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		},
	})
}

// Login handles user login.
// POST /auth/login
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Email and password are required",
		})
	}

	resp, err := h.supabase.SignIn(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "auth_error",
			Message: "Invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		User: User{
			ID:        resp.User.ID,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		},
	})
}

// Refresh handles token refresh.
// POST /auth/refresh
func (h *Handler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
	}

	resp, err := h.supabase.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "auth_error",
			Message: "Invalid or expired refresh token",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		User: User{
			ID:        resp.User.ID,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
		},
	})
}

// Logout handles user logout.
// POST /auth/logout
func (h *Handler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Authorization header is required",
		})
	}

	// Extract token from "Bearer <token>"
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	if token == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid authorization header",
		})
	}

	// Best effort logout - ignore errors
	_ = h.supabase.SignOut(token)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
