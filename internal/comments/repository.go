package comments

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/google/uuid"
)

// Comments repository interface
type Repository interface {
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	Update(ctx context.Context, comment *dto.UpdateCommDTO) (*models.Comment, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
	GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	GetAllByNewsID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
}
