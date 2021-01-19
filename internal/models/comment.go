package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment model
type Comment struct {
	CommentID uuid.UUID `json:"comment_id" db:"comment_id" validate:"omitempty,uuid"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id" validate:"required"`
	NewsID    uuid.UUID `json:"news_id" db:"news_id" validate:"required"`
	Message   string    `json:"message" db:"message" validate:"required,gte=10"`
	Likes     int64     `json:"likes" db:"likes" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Base Comment response
type CommentBase struct {
	CommentID uuid.UUID `json:"comment_id" db:"comment_id" validate:"omitempty,uuid"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id" validate:"required"`
	Author    string    `json:"author" db:"author" validate:"required"`
	AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
	Message   string    `json:"message" db:"message" validate:"required,gte=10"`
	Likes     int64     `json:"likes" db:"likes" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// All News response
type CommentsList struct {
	TotalCount int            `json:"total_count"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	HasMore    bool           `json:"has_more"`
	Comments   []*CommentBase `json:"comments"`
}
