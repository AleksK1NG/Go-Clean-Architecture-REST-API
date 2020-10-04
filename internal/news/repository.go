package news

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// News Repository
type Repository interface {
	Create(ctx context.Context, news *models.News) (*models.News, error)
}
