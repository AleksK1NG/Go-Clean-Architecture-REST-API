package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-news:"
	cacheDuration = 3600
)

// News Repository
type repository struct {
	db         *sqlx.DB
	redisPool  redis.RedisPool
	basePrefix string
}

// News repository constructor
func NewNewsRepository(db *sqlx.DB, redisPool redis.RedisPool) news.Repository {
	return &repository{db: db, redisPool: redisPool, basePrefix: basePrefix}
}

// Create news
func (r repository) Create(ctx context.Context, news *models.News) (*models.News, error) {
	var n models.News

	if err := r.db.QueryRowxContext(
		ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).StructScan(&n); err != nil {
		return nil, errors.WithMessage(err, "newsRepo Create QueryRowxContext")
	}

	return &n, nil
}

// Update news item
func (r repository) Update(ctx context.Context, news *models.News) (*models.News, error) {

	var n models.News
	if err := r.db.QueryRowxContext(
		ctx,
		updateNews,
		&news.Title,
		&news.Content,
		&news.ImageURL,
		&news.Category,
	).StructScan(&n); err != nil {
		return nil, errors.WithMessage(err, "newsRepo Update QueryRowxContext")
	}

	return &n, nil
}

// Get single news by id
func (r repository) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {

	n := &models.NewsBase{}

	if err := r.redisPool.GetJSONContext(ctx, r.getKeyWithPrefix(newsID.String()), n); err == nil {
		return n, nil
	}

	if err := r.db.GetContext(ctx, n, getNewsByID, newsID); err != nil {
		return nil, errors.WithMessage(err, "newsRepo GetNewsByID GetContext")
	}

	if err := r.redisPool.SetexJSONContext(ctx, r.getKeyWithPrefix(newsID.String()), cacheDuration, n); err != nil {
		logger.Errorf("SetexJSONContext Error: %s", err.Error())
	}

	return n, nil
}

// Delete news by id
func (r repository) Delete(ctx context.Context, newsID uuid.UUID) error {

	result, err := r.db.ExecContext(ctx, deleteNews, newsID)
	if err != nil {
		return errors.WithMessage(err, "newsRepo Delete ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "newsRepo Delete RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.WithMessage(sql.ErrNoRows, "newsRepo Deleteno no rowsAffected")
	}

	return nil
}

// Get news
func (r repository) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.WithMessage(err, "newsRepo GetNews GetContext")
	}

	if totalCount == 0 {
		return &models.NewsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			News:       make([]*models.News, 0),
		}, nil
	}

	var newsList = make([]*models.News, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getNews, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.WithMessage(err, "newsRepo GetNews QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.News{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.WithMessage(err, "newsRepo GetNews StructScan")
		}
		newsList = append(newsList, n)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, "newsRepo GetNews rows.Err")
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

// Find news by title
func (r repository) SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByTitleCount, title); err != nil {
		return nil, errors.WithMessage(err, "newsRepo SearchByTitle GetContext")
	}
	if totalCount == 0 {
		return &models.NewsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			News:       make([]*models.News, 0),
		}, nil
	}

	var newsList = make([]*models.News, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, findByTitle, title, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.WithMessage(err, "newsRepo SearchByTitle QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.News{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.WithMessage(err, "newsRepo SearchByTitle StructScan")
		}
		newsList = append(newsList, n)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, "newsRepo SearchByTitle rows.Err")
	}

	return &models.NewsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		News:       newsList,
	}, nil
}

func (r *repository) getKeyWithPrefix(newsID string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, newsID)
}
