package news

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/google/uuid"
)

// News Repository
type Repository interface {
	Create(ctx context.Context, news *models.News) (*models.News, error)
	Update(ctx context.Context, news *models.News) (*models.News, error)
	GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.News, error)
	Delete(ctx context.Context, newsID uuid.UUID) error
}
