// Package models defines data structures for {{.ProjectName}}.
package models

import "time"

// Item represents a basic item entity.
// This is an example model - modify or replace with your own models.
type Item struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// CreateItemRequest represents the request to create an item.
type CreateItemRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
}

// UpdateItemRequest represents the request to update an item.
type UpdateItemRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// ItemResponse represents an item in API responses.
type ItemResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// ToResponse converts an Item to ItemResponse.
func (i *Item) ToResponse() ItemResponse {
	return ItemResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Completed:   i.Completed,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
