package news

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// News use case
type UseCase interface {
	Create(ctx context.Context, news *models.News) (*models.News, error)
}
