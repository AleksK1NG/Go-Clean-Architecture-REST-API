package news

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// News redis repository
type RedisRepository interface {
	GetNewsByIDCtx(ctx context.Context, key string) (*models.NewsBase, error)
	SetNewsCtx(ctx context.Context, key string, seconds int, news *models.NewsBase) error
	DeleteNewsCtx(ctx context.Context, key string) error
}
