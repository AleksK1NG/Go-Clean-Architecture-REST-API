package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/dto"
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
	comment.AuthorID = user.UserID

	createdComment, err := u.commRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

// Update comment
func (u *useCase) Update(ctx context.Context, comment *dto.UpdateCommDTO) (*models.Comment, error) {
	comm, err := u.commRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String(), u.logger); err != nil {
		return nil, err
	}

	return u.commRepo.Update(ctx, comment)
}

// Delete comment
func (u *useCase) Delete(ctx context.Context, commentID uuid.UUID) error {
	comm, err := u.commRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String(), u.logger); err != nil {
		return err
	}
	return u.commRepo.Delete(ctx, commentID)
}

// GetByID comment
func (u *useCase) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	return u.commRepo.GetByID(ctx, commentID)
}

// GetAllByNewsID comments
func (u *useCase) GetAllByNewsID(ctx context.Context, query *dto.CommentsByNewsID) (*models.CommentsList, error) {
	return u.commRepo.GetAllByNewsID(ctx, query)
}
