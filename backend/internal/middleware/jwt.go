// Package middleware provides HTTP middleware for {{.ProjectName}}.
package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTConfig holds JWT middleware configuration.
type JWTConfig struct {
	// JWTSecret is the secret key used to validate Supabase JWT tokens.
	// This should be your Supabase JWT secret from the dashboard.
	JWTSecret string
}

// Claims represents the JWT claims from Supabase.
type Claims struct {
	Sub   string `json:"sub"` // User ID
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth creates a JWT authentication middleware.
// It validates the Authorization header contains a valid Supabase JWT token.
func JWTAuth(config JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Authorization header is required",
				})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Invalid authorization header format",
				})
			}

			tokenString := parts[1]

			// Parse and validate the token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Invalid or expired token",
				})
			}

			// Extract claims and set in context
			if claims, ok := token.Claims.(*Claims); ok {
				c.Set("user_id", claims.Sub)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
			}

			return next(c)
		}
	}
}

// GetUserID extracts the user ID from the request context.
// Returns empty string if not found.
func GetUserID(c echo.Context) string {
	if id, ok := c.Get("user_id").(string); ok {
		return id
	}
	return ""
}

// GetUserEmail extracts the user email from the request context.
// Returns empty string if not found.
func GetUserEmail(c echo.Context) string {
	if email, ok := c.Get("user_email").(string); ok {
		return email
	}
	return ""
}
