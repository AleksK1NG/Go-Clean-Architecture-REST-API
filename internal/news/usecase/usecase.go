package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

const (
	basePrefix    = "api-news:"
	cacheDuration = 3600
)

// News UseCase
type newsUC struct {
	cfg       *config.Config
	newsRepo  news.Repository
	redisRepo news.RedisRepository
	logger    logger.Logger
}

// News UseCase constructor
func NewNewsUseCase(cfg *config.Config, newsRepo news.Repository, redisRepo news.RedisRepository, logger logger.Logger) news.UseCase {
	return &newsUC{cfg: cfg, newsRepo: newsRepo, redisRepo: redisRepo, logger: logger}
}

// Create news
func (u *newsUC) Create(ctx context.Context, news *models.News) (*models.News, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.Create")
	defer span.Finish()

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "newsUC.Create.GetUserFromCtx"))
	}

	news.AuthorID = user.UserID

	if err = utils.ValidateStruct(ctx, news); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "newsUC.Create.ValidateStruct"))
	}

	n, err := u.newsRepo.Create(ctx, news)
	if err != nil {
		return nil, err
	}

	return n, err
}

// Update news item
func (u *newsUC) Update(ctx context.Context, news *models.News) (*models.News, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.Update")
	defer span.Finish()

	newsByID, err := u.newsRepo.GetNewsByID(ctx, news.NewsID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), u.logger); err != nil {
		return nil, httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "newsUC.Update.ValidateIsOwner"))
	}

	updatedUser, err := u.newsRepo.Update(ctx, news)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(news.NewsID.String())); err != nil {
		u.logger.Errorf("newsUC.Update.DeleteNewsCtx: %v", err)
	}

	return updatedUser, nil
}

// Get news by id
func (u *newsUC) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.GetNewsByID")
	defer span.Finish()

	newsBase, err := u.redisRepo.GetNewsByIDCtx(ctx, u.getKeyWithPrefix(newsID.String()))
	if err != nil {
		u.logger.Errorf("newsUC.GetNewsByID.GetNewsByIDCtx: %v", err)
	}
	if newsBase != nil {
		return newsBase, nil
	}

	n, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetNewsCtx(ctx, u.getKeyWithPrefix(newsID.String()), cacheDuration, n); err != nil {
		u.logger.Errorf("newsUC.GetNewsByID.SetNewsCtx: %s", err)
	}

	return n, nil
}

// Delete news
func (u *newsUC) Delete(ctx context.Context, newsID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.Delete")
	defer span.Finish()

	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String(), u.logger); err != nil {
		return httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "newsUC.Delete.ValidateIsOwner"))
	}

	if err = u.newsRepo.Delete(ctx, newsID); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(newsID.String())); err != nil {
		u.logger.Errorf("newsUC.Delete.DeleteNewsCtx: %v", err)
	}

	return nil
}

// Get news
func (u *newsUC) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.GetNews")
	defer span.Finish()

	return u.newsRepo.GetNews(ctx, pq)
}

// Find nes by title
func (u *newsUC) SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsUC.SearchByTitle")
	defer span.Finish()

	return u.newsRepo.SearchByTitle(ctx, title, query)
}

func (u *newsUC) getKeyWithPrefix(newsID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, newsID)
}
