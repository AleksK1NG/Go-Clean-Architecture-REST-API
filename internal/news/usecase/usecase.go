package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/logger"
)

// News useCase
type useCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	newsRepo news.Repository
}

// News use case constructor
func NewNewsUseCase(logger *logger.Logger, cfg *config.Config, newsRepo news.Repository) news.UseCase {
	return &useCase{logger, cfg, newsRepo}
}

// Create news
func (u *useCase) Create(ctx context.Context, news *models.News) (*models.News, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	news.AuthorID = user.ID

	if err := utils.ValidateStruct(ctx, news); err != nil {
		return nil, err
	}

	n, err := u.newsRepo.Create(ctx, news)
	if err != nil {
		return nil, err
	}

	return n, err
}

// Update news item
func (u *useCase) Update(ctx context.Context, news *models.News) (*models.News, error) {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, news.ID)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return nil, err
	}

	updatedUser, err := u.newsRepo.Update(ctx, news)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
