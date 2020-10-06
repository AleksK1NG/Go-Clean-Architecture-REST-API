package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// News Repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
	redis  *redis.RedisClient
}

// News repository constructor
func NewNewsRepository(logger *logger.Logger, db *sqlx.DB, redis *redis.RedisClient) *repository {
	return &repository{logger, db, redis}
}

// Create news
func (r repository) Create(ctx context.Context, news *models.News) (*models.News, error) {
	var n models.News

	if err := r.db.QueryRowxContext(
		ctx,
		createUser,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).StructScan(&n); err != nil {
		r.logger.Error("QueryRowxContext", zap.String("ERROR", err.Error()))
		return nil, err
	}

	return &n, nil
}
