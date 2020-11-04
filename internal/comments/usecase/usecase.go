package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Comments useCase
type useCase struct {
	cfg      *config.Config
	commRepo comments.Repository
}

// Auth useCase constructor
func NewCommentsUseCase(cfg *config.Config, commRepo comments.Repository) comments.UseCase {
	return &useCase{cfg: cfg, commRepo: commRepo}
}

// Create comment
func (u *useCase) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	return u.commRepo.Create(ctx, comment)
}

// Update comment
func (u *useCase) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	comm, err := u.commRepo.GetByID(ctx, comment.CommentID)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String()); err != nil {
		return nil, errors.WithMessage(err, "commentsUC Update ValidateIsOwner")
	}

	return u.commRepo.Update(ctx, comment)
}

// Delete comment
func (u *useCase) Delete(ctx context.Context, commentID uuid.UUID) error {
	comm, err := u.commRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String()); err != nil {
		return errors.WithMessage(err, "commentsUC Delete ValidateIsOwner")
	}
	return u.commRepo.Delete(ctx, commentID)
}

// GetByID comment
func (u *useCase) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error) {
	return u.commRepo.GetByID(ctx, commentID)
}

// GetAllByNewsID comments
func (u *useCase) GetAllByNewsID(ctx context.Context, newsID uuid.UUID, query *utils.PaginationQuery) (*models.CommentsList, error) {
	return u.commRepo.GetAllByNewsID(ctx, newsID, query)
}
