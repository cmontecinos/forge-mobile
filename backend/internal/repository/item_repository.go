// Package repository provides data access patterns for {{.ProjectName}}.
package repository

import (
	"time"

	"github.com/{{.ProjectName}}/backend/internal/models"
	"github.com/{{.ProjectName}}/backend/internal/supabase"
)

// ItemRepository handles item data operations.
type ItemRepository struct {
	client *supabase.Client
}

// NewItemRepository creates a new item repository.
func NewItemRepository(client *supabase.Client) *ItemRepository {
	return &ItemRepository{client: client}
}

// Create inserts a new item for a user.
func (r *ItemRepository) Create(userID string, req models.CreateItemRequest, userToken string) (*models.Item, error) {
	item := map[string]interface{}{
		"user_id":     userID,
		"title":       req.Title,
		"description": req.Description,
		"completed":   false,
		"created_at":  time.Now().UTC(),
	}

	var result []models.Item
	err := r.client.InsertReturning("items", item, &result, userToken)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

// GetByID retrieves a single item by ID.
func (r *ItemRepository) GetByID(id string, userToken string) (*models.Item, error) {
	var item models.Item
	err := r.client.From("items").
		Eq("id", id).
		WithToken(userToken).
		Single().
		Execute(&item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetByUserID retrieves all items for a user.
func (r *ItemRepository) GetByUserID(userID string, userToken string) ([]models.Item, error) {
	var items []models.Item
	err := r.client.From("items").
		Eq("user_id", userID).
		Order("created_at", false).
		WithToken(userToken).
		Execute(&items)

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update modifies an existing item.
func (r *ItemRepository) Update(id string, req models.UpdateItemRequest, userToken string) (*models.Item, error) {
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}

	now := time.Now().UTC()
	updates["updated_at"] = now

	filters := []supabase.Filter{
		{Column: "id", Operator: supabase.OpEq, Value: id},
	}

	var result []models.Item
	err := r.client.UpdateReturning("items", updates, filters, &result, userToken)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

// Delete removes an item by ID.
func (r *ItemRepository) Delete(id string, userToken string) error {
	filters := []supabase.Filter{
		{Column: "id", Operator: supabase.OpEq, Value: id},
	}

	return r.client.Delete("items", filters, userToken)
}
