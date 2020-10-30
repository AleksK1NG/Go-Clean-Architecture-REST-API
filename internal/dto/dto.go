package dto

import (
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
)

// Find user query DTO
type FindNewsDTO struct {
	Title string `json:"title" validate:"required"`
	PQ    *utils.PaginationQuery
}

// Update Comment DTO
type UpdateCommDTO struct {
	ID      uuid.UUID `json:"comment_id" db:"comment_id" validate:"omitempty,uuid"`
	Message string    `json:"message" db:"password" validate:"required,gte=0"`
}

// Find user query DTO
type CommentsByNewsID struct {
	NewsID uuid.UUID              `json:"news_id"`
	PQ     *utils.PaginationQuery `json:"pq"`
}
