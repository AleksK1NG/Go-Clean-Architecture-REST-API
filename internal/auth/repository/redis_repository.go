package repository

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

// Auth redis repository
type authRedisRepo struct {
	redisClient *redis.Client
}

// Auth redis repository constructor
func NewAuthRedisRepo(redisClient *redis.Client) auth.RedisRepository {
	return &authRedisRepo{redisClient: redisClient}
}

// Get user by id
func (a *authRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	userBytes, err := a.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "authRedisRepo GetByIDCtx redisClient.Get")
	}
	user := &models.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.Wrap(err, "authRedisRepo GetByIDCtx json.Unmarshal")
	}
	return user, nil
}

// Cache user with duration in seconds
func (a *authRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo SetUserCtx json.Unmarshal")
	}
	if err = a.redisClient.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo SetUserCtx redisClient.Set")
	}
	return nil
}

// Delete user by key
func (a *authRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	if err := a.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo DeleteUserCtx redisClient.Del")
	}
	return nil
}
