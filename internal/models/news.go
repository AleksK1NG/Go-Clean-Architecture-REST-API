package models

import (
	"github.com/google/uuid"
	"time"
)

// News base model
type News struct {
	ID        uuid.UUID `json:"news_id" db:"news_id" validate:"omitempty,uuid"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id" validate:"omitempty,uuid"`
	Title     string    `json:"title" db:"title" validate:"required,gte=10"`
	Content   string    `json:"content" db:"content" validate:"required,gte=20"`
	ImageURL  *string   `json:"image_url,omitempty" db:"image_url" validate:"omitempty,lte=512,url"`
	Category  *string   `json:"category,omitempty" db:"category" validate:"omitempty,lte=10"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
