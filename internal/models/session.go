package models

import "github.com/google/uuid"

// Session model
type Session struct {
	SessionID string    `json:"session_id" redis:"session_id"`
	UserID    uuid.UUID `json:"user_id" redis:"user_id"`
}
