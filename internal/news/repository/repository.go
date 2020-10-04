package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// News Repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
}

// News repository constructor
func NewNewsRepository(logger *logger.Logger, db *sqlx.DB) *repository {
	return &repository{logger, db}
}

// Create news
func (r repository) Create(ctx context.Context, news *models.News) (*models.News, error) {
	return &models.News{}, nil
}
