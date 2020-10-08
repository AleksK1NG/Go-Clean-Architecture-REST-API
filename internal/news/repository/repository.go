package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/utils"
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

// Delete news by id
func (r repository) Delete(ctx context.Context, newsID uuid.UUID) error {

	result, err := r.db.ExecContext(ctx, deleteNews, newsID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := r.redis.Delete(newsID.String()); err != nil {
		r.logger.Error("REDIS Delete", zap.String("ERROR", err.Error()))
	}

	return nil
}

// Get news
func (r repository) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, err
	}

	var newsList = make([]*models.News, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getNews, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.News{}
		if err := rows.StructScan(n); err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.NewsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		News:       newsList,
	}, nil
}
