package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
)

// Comments useCase
type useCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	authRepo auth.Repository
}

// Auth useCase constructor
func NewCommentsUseCase(l *logger.Logger, c *config.Config, ar auth.Repository) comments.UseCase {
	return &useCase{l, c, ar}
}

// Create comment
func (u *useCase) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}

// Update comment
func (u *useCase) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}

// Delete comment
func (u *useCase) Delete(ctx context.Context, commentID uuid.UUID) error {
	panic("implement me")
}

// GetByID comment
func (u *useCase) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}

// GetAllByNewsID comments
func (u *useCase) GetAllByNewsID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}
