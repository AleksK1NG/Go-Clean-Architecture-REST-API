package repository

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/go-redis/redis/v8"
	"time"
)

// News redis repository
type newsRedisRepo struct {
	redisClient *redis.Client
}

// News redis repository constructor
func NewNewsRedisRepo(redisClient *redis.Client) news.RedisRepository {
	return &newsRedisRepo{redisClient: redisClient}
}

// Get new by id
func (n *newsRedisRepo) GetNewsByIDCtx(ctx context.Context, key string) (*models.NewsBase, error) {
	newsBytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	newsBase := &models.NewsBase{}
	if err = json.Unmarshal(newsBytes, newsBase); err != nil {
		return nil, err
	}

	return newsBase, nil
}

// Cache news item
func (n *newsRedisRepo) SetNewsCtx(ctx context.Context, key string, seconds int, news *models.NewsBase) error {
	newsBytes, err := json.Marshal(news)
	if err != nil {
		return err
	}
	return n.redisClient.Set(ctx, key, newsBytes, time.Second*time.Duration(seconds)).Err()
}

// Delete new item from cache
func (n *newsRedisRepo) DeleteNewsCtx(ctx context.Context, key string) error {
	return n.redisClient.Del(ctx, key).Err()
}
