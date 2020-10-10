package models

import "github.com/google/uuid"

// Session model
type Session struct {
	ID     string    `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}
