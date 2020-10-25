package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
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
		return nil, err
	}

	news.AuthorID = user.UserID

	if err = utils.ValidateStruct(ctx, news); err != nil {
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
	newsByID, err := u.newsRepo.GetNewsByID(ctx, news.NewsID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return nil, err
	}

	updatedUser, err := u.newsRepo.Update(ctx, news)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// Get news by id
func (u *useCase) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*models.NewsBase, error) {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return nil, err
	}

	return newsByID, nil
}

// Delete news
func (u *useCase) Delete(ctx context.Context, newsID uuid.UUID) error {
	newsByID, err := u.newsRepo.GetNewsByID(ctx, newsID)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, newsByID.AuthorID.String()); err != nil {
		return err
	}

	if err := u.newsRepo.Delete(ctx, newsID); err != nil {
		return err
	}

	return nil
}

// Get news
func (u *useCase) GetNews(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	newsList, err := u.newsRepo.GetNews(ctx, pq)
	if err != nil {
		return nil, err
	}

	return newsList, nil
}

// Find nes by title
func (u *useCase) SearchByTitle(ctx context.Context, req *dto.FindNewsDTO) (*models.NewsList, error) {
	newsList, err := u.newsRepo.SearchByTitle(ctx, req)
	if err != nil {
		return nil, err
	}

	return newsList, nil
}
