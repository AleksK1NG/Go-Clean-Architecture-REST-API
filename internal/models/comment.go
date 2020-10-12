package models

import (
	"github.com/google/uuid"
	"time"
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
