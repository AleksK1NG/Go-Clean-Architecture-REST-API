package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
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
func NewNewsRepository(logger *logger.Logger, db *sqlx.DB, redis *redis.RedisClient) news.Repository {
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

// Update news item
func (r repository) Update(ctx context.Context, news *models.News) (*models.News, error) {

	var n models.News
	if err := r.db.QueryRowxContext(
		ctx,
		updateUser,
		&news.Title,
		&news.Content,
		&news.ImageURL,
		&news.Category,
	).StructScan(&n); err != nil {
		return nil, err
	}

	if err := r.redis.Delete(n.ID.String()); err != nil {
		r.logger.Error("redis.Delete", zap.String("ERROR", err.Error()))
	}

	return &n, nil
}

// Get single news by id
func (r repository) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.News, error) {
	var n models.News
	if err := r.redis.GetIfExistsJSON(newsID.String(), &n); err != nil {
		if err != errors.NotExists {
			r.logger.Error("REDIS GetIfExistsJSON", zap.String("ERROR", err.Error()))
		}
	} else {
		return &n, nil
	}

	if err := r.db.GetContext(ctx, &n, getNewsByID, newsID); err != nil {
		return nil, err
	}

	if err := r.redis.SetEXJSON(n.ID.String(), 50, &n); err != nil {
		r.logger.Error("SetEXJSON", zap.String("ERROR", err.Error()))
	}

	return &n, nil
}
