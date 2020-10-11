package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	panic("implement me")
}

// Update comment
func (r *repository) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}

// Delete comment
func (r *repository) Delete(ctx context.Context, commentID uuid.UUID) error {
	panic("implement me")
}

// GetByID comment
func (r *repository) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}

// GetAllByNewsID comments
func (r *repository) GetAllByNewsID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}
