package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// News Repository
type newsRepo struct {
	db *sqlx.DB
}

// News repository constructor
func NewNewsRepository(db *sqlx.DB) news.Repository {
	return &newsRepo{db: db}
}

// Create news
func (r *newsRepo) Create(ctx context.Context, news *models.News) (*models.News, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Create")
	defer span.Finish()

	var n models.News
	if err := r.db.QueryRowxContext(
		ctx,
		createNews,
		&news.AuthorID,
		&news.Title,
		&news.Content,
		&news.Category,
	).StructScan(&n); err != nil {
		return nil, errors.Wrap(err, "newsRepo.Create.QueryRowxContext")
	}

	return &n, nil
}

// Update news item
func (r *newsRepo) Update(ctx context.Context, news *models.News) (*models.News, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Update")
	defer span.Finish()

	var n models.News
	if err := r.db.QueryRowxContext(
		ctx,
		updateNews,
		&news.Title,
		&news.Content,
		&news.ImageURL,
		&news.Category,
		&news.NewsID,
	).StructScan(&n); err != nil {
		return nil, errors.Wrap(err, "newsRepo.Update.QueryRowxContext")
	}

	return &n, nil
}

// Get single news by id
func (r *newsRepo) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.GetNewsByID")
	defer span.Finish()

	n := &models.NewsBase{}
	if err := r.db.GetContext(ctx, n, getNewsByID, newsID); err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetNewsByID.GetContext")
	}

	return n, nil
}

// Delete news by id
func (r *newsRepo) Delete(ctx context.Context, newsID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteNews, newsID)
	if err != nil {
		return errors.Wrap(err, "newsRepo.Delete.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "newsRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "newsRepo.Delete.rowsAffected")
	}

	return nil
}

// Get news
func (r *newsRepo) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.GetNews")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetNews.GetContext.totalCount")
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
		return nil, errors.Wrap(err, "newsRepo.GetNews.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.News{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "newsRepo.GetNews.StructScan")
		}
		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetNews.rows.Err")
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
func (r *newsRepo) SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.SearchByTitle")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByTitleCount, title); err != nil {
		return nil, errors.Wrap(err, "newsRepo.SearchByTitle.GetContext")
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
		return nil, errors.Wrap(err, "newsRepo.SearchByTitle.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.News{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "newsRepo.SearchByTitle.StructScan")
		}
		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "newsRepo.SearchByTitle.rows.Err")
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
