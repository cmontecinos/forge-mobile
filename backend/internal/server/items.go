// Package server provides item handlers for {{.ProjectName}}.
package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	custommw "github.com/{{.ProjectName}}/backend/internal/middleware"
	"github.com/{{.ProjectName}}/backend/internal/models"
	"github.com/{{.ProjectName}}/backend/internal/repository"
)

// ItemHandler handles item-related requests.
type ItemHandler struct {
	repo *repository.ItemRepository
}

// NewItemHandler creates a new item handler.
func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{repo: repo}
}

// getToken extracts the access token from the Authorization header.
func getToken(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}

// ListItems returns all items for the authenticated user.
// GET /api/v1/items
func (h *ItemHandler) ListItems(c echo.Context) error {
	userID := custommw.GetUserID(c)
	if userID == "" {
		return Unauthorized(c, "User not authenticated")
	}

	token := getToken(c)
	items, err := h.repo.GetByUserID(userID, token)
	if err != nil {
		return InternalError(c, "Failed to fetch items")
	}

	// Convert to response format
	response := make([]models.ItemResponse, len(items))
	for i, item := range items {
		response[i] = item.ToResponse()
	}

	return c.JSON(http.StatusOK, response)
}

// GetItem returns a single item by ID.
// GET /api/v1/items/:id
func (h *ItemHandler) GetItem(c echo.Context) error {
	userID := custommw.GetUserID(c)
	if userID == "" {
		return Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return BadRequest(c, "Item ID is required")
	}

	token := getToken(c)
	item, err := h.repo.GetByID(id, token)
	if err != nil {
		return NotFound(c, "Item not found")
	}

	// Verify ownership
	if item.UserID != userID {
		return Forbidden(c, "Access denied")
	}

	return c.JSON(http.StatusOK, item.ToResponse())
}

// CreateItem creates a new item.
// POST /api/v1/items
func (h *ItemHandler) CreateItem(c echo.Context) error {
	userID := custommw.GetUserID(c)
	if userID == "" {
		return Unauthorized(c, "User not authenticated")
	}

	var req models.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	// Validate
	errors := make(map[string]string)
	if req.Title == "" {
		errors["title"] = "Title is required"
	}
	if len(errors) > 0 {
		return ValidationError(c, errors)
	}

	token := getToken(c)
	item, err := h.repo.Create(userID, req, token)
	if err != nil {
		return InternalError(c, "Failed to create item")
	}

	return c.JSON(http.StatusCreated, item.ToResponse())
}

// UpdateItem updates an existing item.
// PATCH /api/v1/items/:id
func (h *ItemHandler) UpdateItem(c echo.Context) error {
	userID := custommw.GetUserID(c)
	if userID == "" {
		return Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return BadRequest(c, "Item ID is required")
	}

	token := getToken(c)

	// Verify item exists and belongs to user
	existing, err := h.repo.GetByID(id, token)
	if err != nil {
		return NotFound(c, "Item not found")
	}
	if existing.UserID != userID {
		return Forbidden(c, "Access denied")
	}

	var req models.UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "Invalid request body")
	}

	item, err := h.repo.Update(id, req, token)
	if err != nil {
		return InternalError(c, "Failed to update item")
	}

	return c.JSON(http.StatusOK, item.ToResponse())
}

// DeleteItem removes an item.
// DELETE /api/v1/items/:id
func (h *ItemHandler) DeleteItem(c echo.Context) error {
	userID := custommw.GetUserID(c)
	if userID == "" {
		return Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return BadRequest(c, "Item ID is required")
	}

	token := getToken(c)

	// Verify item exists and belongs to user
	existing, err := h.repo.GetByID(id, token)
	if err != nil {
		return NotFound(c, "Item not found")
	}
	if existing.UserID != userID {
		return Forbidden(c, "Access denied")
	}

	if err := h.repo.Delete(id, token); err != nil {
		return InternalError(c, "Failed to delete item")
	}

	return c.NoContent(http.StatusNoContent)
}
