package usecase

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
}

// News UseCase constructor
func NewNewsUseCase(cfg *config.Config, newsRepo news.Repository, redisRepo news.RedisRepository) news.UseCase {
	return &newsUC{cfg: cfg, newsRepo: newsRepo, redisRepo: redisRepo}
}

// Create news
func (u *newsUC) Create(ctx context.Context, news *models.News) (*models.News, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "newsUC Create GetUserFromCtx")
	}

	news.AuthorID = user.UserID

	if err = utils.ValidateStruct(ctx, news); err != nil {
		return nil, errors.WithMessage(err, "newsUC Create ValidateStruct")
	}

	n, err := u.newsRepo.Create(ctx, news)
	if err != nil {
		return nil, err
	}

	return n, err
}

// Update news item
func (u *newsUC) Update(ctx context.Context, news *models.News) (*models.News, error) {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, news.NewsID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return nil, errors.WithMessage(err, "newsUC Update ValidateIsOwner")
	}

	updatedUser, err := u.newsRepo.Update(ctx, news)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(news.NewsID.String())); err != nil {
		logger.Errorf("newsUC Update redis delete: %s", err)
	}

	return updatedUser, nil
}

// Get news by id
func (u *newsUC) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	newsBase, err := u.redisRepo.GetNewsByIDCtx(ctx, u.getKeyWithPrefix(newsID.String()))
	if err != nil {
		logger.Errorf("newsUC GetNewsByID redisRepo.GetNewsByIDCtx: %v", err)
	}
	if newsBase != nil {
		return newsBase, nil
	}

	n, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetNewsCtx(ctx, u.getKeyWithPrefix(newsID.String()), cacheDuration, n); err != nil {
		logger.Errorf("newsUC GetNewsByID redisRepo.SetNewsCtx: %s", err)
	}

	return n, nil
}

// Delete news
func (u *newsUC) Delete(ctx context.Context, newsID uuid.UUID) error {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return errors.WithMessage(err, "newsUC Update ValidateIsOwner")
	}

	if err = u.newsRepo.Delete(ctx, newsID); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteNewsCtx(ctx, u.getKeyWithPrefix(newsID.String())); err != nil {
		logger.Errorf("newsUC Delete redisRepo.DeleteNewsCtx: %v", err)
	}

	return nil
}

// Get news
func (u *newsUC) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	return u.newsRepo.GetNews(ctx, pq)
}

// Find nes by title
func (u *newsUC) SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error) {
	return u.newsRepo.SearchByTitle(ctx, title, query)
}

func (u *newsUC) getKeyWithPrefix(newsID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, newsID)
}
