package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// News useCase
type useCase struct {
	cfg      *config.Config
	newsRepo news.Repository
}

// News use case constructor
func NewNewsUseCase(cfg *config.Config, newsRepo news.Repository) news.UseCase {
	return &useCase{cfg: cfg, newsRepo: newsRepo}
}

// Create news
func (u *useCase) Create(ctx context.Context, news *models.News) (*models.News, error) {
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
func (u *useCase) Update(ctx context.Context, news *models.News) (*models.News, error) {
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

	return updatedUser, nil
}

// Get news by id
func (u *useCase) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	return u.newsRepo.GetNewsByID(ctx, newsID)
}

// Delete news
func (u *useCase) Delete(ctx context.Context, newsID uuid.UUID) error {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return errors.WithMessage(err, "newsUC Update ValidateIsOwner")
	}

	if err := u.newsRepo.Delete(ctx, newsID); err != nil {
		return err
	}

	return nil
}

// Get news
func (u *useCase) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	return u.newsRepo.GetNews(ctx, pq)
}

// Find nes by title
func (u *useCase) SearchByTitle(ctx context.Context, title string, query *utils.PaginationQuery) (*models.NewsList, error) {
	return u.newsRepo.SearchByTitle(ctx, title, query)
}
