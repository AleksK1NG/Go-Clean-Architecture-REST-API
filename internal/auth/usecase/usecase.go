package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/db/redis"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// Auth useCase
type useCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	authRepo auth.Repository
	redis    *redis.RedisClient
}

// Auth useCase constructor
func NewAuthUseCase(l *logger.Logger, c *config.Config, ar auth.Repository, r *redis.RedisClient) *useCase {
	return &useCase{l, c, ar, r}
}

// Create new user
func (u *useCase) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareCreate(); err != nil {
		return nil, err
	}

	createdUser, err := u.authRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
