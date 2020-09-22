package usecase

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/db/redis"
	"github.com/AleksK1NG/api-mc/internal/logger"
)

// Auth useCase
type UseCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	authRepo auth.Repository
	redis    *redis.RedisClient
}

// Auth useCase constructor
func NewAuthUseCase(l *logger.Logger, c *config.Config, ar auth.Repository, r *redis.RedisClient) *UseCase {
	return &UseCase{l, c, ar, r}
}

// Create new user
func (u *UseCase) Create() error {
	return u.authRepo.Create()
}
