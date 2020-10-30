package dto

import (
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
)

// Find user query DTO
type CommentsByNewsID struct {
	NewsID uuid.UUID              `json:"news_id"`
	PQ     *utils.PaginationQuery `json:"pq"`
}
