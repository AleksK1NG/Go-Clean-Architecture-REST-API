package auth

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// Auth Redis repository interface
type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
