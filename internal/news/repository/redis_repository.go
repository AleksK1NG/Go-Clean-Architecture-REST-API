package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRedisRepo.GetNewsByIDCtx")
	defer span.Finish()

	newsBytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "newsRedisRepo.GetNewsByIDCtx.redisClient.Get")
	}
	newsBase := &models.NewsBase{}
	if err = json.Unmarshal(newsBytes, newsBase); err != nil {
		return nil, errors.Wrap(err, "newsRedisRepo.GetNewsByIDCtx.json.Unmarshal")
	}

	return newsBase, nil
}

// Cache news item
func (n *newsRedisRepo) SetNewsCtx(ctx context.Context, key string, seconds int, news *models.NewsBase) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRedisRepo.SetNewsCtx")
	defer span.Finish()

	newsBytes, err := json.Marshal(news)
	if err != nil {
		return errors.Wrap(err, "newsRedisRepo.SetNewsCtx.json.Marshal")
	}
	if err = n.redisClient.Set(ctx, key, newsBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "newsRedisRepo.SetNewsCtx.redisClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *newsRedisRepo) DeleteNewsCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRedisRepo.DeleteNewsCtx")
	defer span.Finish()

	if err := n.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "newsRedisRepo.DeleteNewsCtx.redisClient.Del")
	}
	return nil
}
