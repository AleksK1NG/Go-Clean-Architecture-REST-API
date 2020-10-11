package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
)

// Comments useCase
type useCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	commRepo comments.Repository
}

// Auth useCase constructor
func NewCommentsUseCase(l *logger.Logger, c *config.Config, commRepo comments.Repository) comments.UseCase {
	return &useCase{l, c, commRepo}
}

// Create comment
func (u *useCase) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	comment.AuthorID = user.ID

	createdComment, err := u.commRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return createdComment, nil
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
