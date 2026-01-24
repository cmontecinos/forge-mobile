// Package server provides error response helpers for {{.ProjectName}}.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// BadRequest returns a 400 Bad Request response.
func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}

// ValidationError returns a 400 Bad Request with field-level errors.
func ValidationError(c echo.Context, errors map[string]string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "validation_error",
		Message: "Validation failed",
		Details: errors,
	})
}

// Unauthorized returns a 401 Unauthorized response.
func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    "unauthorized",
		Message: message,
	})
}

// Forbidden returns a 403 Forbidden response.
func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    "forbidden",
		Message: message,
	})
}

// NotFound returns a 404 Not Found response.
func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "not_found",
		Message: message,
	})
}

// InternalError returns a 500 Internal Server Error response.
func InternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "internal_error",
		Message: message,
	})
}
