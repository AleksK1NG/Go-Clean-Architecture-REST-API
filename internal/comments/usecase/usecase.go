package usecase

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix      = "api-comments:"
	durationSeconds = 3600
)

// Comments UseCase
type commentsUC struct {
	cfg       *config.Config
	commRepo  comments.Repository
	redisRepo redis.RedisPool
}

// Comments UseCase constructor
func NewCommentsUseCase(cfg *config.Config, commRepo comments.Repository, redisRepo redis.RedisPool) comments.UseCase {
	return &commentsUC{cfg: cfg, commRepo: commRepo, redisRepo: redisRepo}
}

// Create comment
func (u *commentsUC) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	return u.commRepo.Create(ctx, comment)
}

// Update comment
func (u *commentsUC) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	comm, err := u.commRepo.GetByID(ctx, comment.CommentID)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String()); err != nil {
		return nil, errors.WithMessage(err, "commentsUC Update ValidateIsOwner")
	}

	if err := u.redisRepo.Delete(u.createKey(comment.CommentID.String())); err != nil {
		logger.Errorf("commentsUC Update redis delete: %s", err)
	}

	updatedComment, err := u.commRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	if err := u.redisRepo.Delete(u.createKey(comment.CommentID.String())); err != nil {
		logger.Errorf("commentsUC Update redis delete: %s", err)
	}

	return updatedComment, nil
}

// Delete comment
func (u *commentsUC) Delete(ctx context.Context, commentID uuid.UUID) error {
	comm, err := u.commRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err := utils.ValidateIsOwner(ctx, comm.AuthorID.String()); err != nil {
		return errors.WithMessage(err, "commentsUC Delete ValidateIsOwner")
	}

	if err := u.commRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	if err := u.redisRepo.Delete(u.createKey(commentID.String())); err != nil {
		logger.Errorf("commentsUC Delete redis delete: %s", err)
	}

	return nil
}

// GetByID comment
func (u *commentsUC) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentBase, error) {
	comment := &models.CommentBase{}
	if err := u.redisRepo.GetJSONContext(ctx, u.createKey(commentID.String()), comment); err == nil {
		return comment, nil
	}

	comment, err := u.commRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if err := u.redisRepo.SetexJSONContext(ctx, u.createKey(commentID.String()), durationSeconds, comment); err != nil {
		logger.Errorf("commentsUC GetByID redis set: %s", err)
	}

	return comment, nil
}

// GetAllByNewsID comments
func (u *commentsUC) GetAllByNewsID(ctx context.Context, newsID uuid.UUID, query *utils.PaginationQuery) (*models.CommentsList, error) {
	return u.commRepo.GetAllByNewsID(ctx, newsID, query)
}

func (u *commentsUC) createKey(commentID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, commentID)
}
