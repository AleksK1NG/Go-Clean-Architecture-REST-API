package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Comments repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
	redis  *redis.RedisClient
}

// Comments Repository constructor
func NewCommentsRepository(logger *logger.Logger, db *sqlx.DB, redis *redis.RedisClient) comments.Repository {
	return &repository{logger, db, redis}
}

// Create comment
func (r *repository) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	c := &models.Comment{}
	if err := r.db.QueryRowxContext(
		ctx,
		createComment,
		&comment.AuthorID,
		&comment.NewsID,
		&comment.Message,
	).StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Update comment
func (r *repository) Update(ctx context.Context, comment *dto.UpdateCommDTO) (*models.Comment, error) {

	comm := &models.Comment{}
	if err := r.db.QueryRowxContext(ctx, updateComment, comment.Message, comment.ID).StructScan(comm); err != nil {
		return nil, err
	}

	if err := r.redis.Delete(comm.ID.String()); err != nil {
		r.logger.Error("REDIS", zap.String("ERROR", err.Error()))
	}

	return comm, nil
}

// Delete comment
func (r *repository) Delete(ctx context.Context, commentID uuid.UUID) error {

	result, err := r.db.ExecContext(ctx, deleteComment, commentID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := r.redis.Delete(commentID.String()); err != nil {
		r.logger.Error("REDIS", zap.String("ERROR", err.Error()))
	}

	return nil
}

// GetByID comment
func (r *repository) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {

	comment := &models.Comment{}

	if err := r.redis.GetIfExistsJSON(commentID.String(), comment); err != nil {
		r.logger.Error("REDIS", zap.String("ERROR", err.Error()))
	} else {
		return comment, nil
	}

	if err := r.db.GetContext(ctx, comment, getCommentByID, commentID); err != nil {
		return nil, err
	}

	if err := r.redis.SetEXJSON(comment.ID.String(), 3600, comment); err != nil {
		r.logger.Error("REDIS", zap.String("ERROR", err.Error()))
	}

	return comment, nil
}

// GetAllByNewsID comments
func (r *repository) GetAllByNewsID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}
