package repository

import (
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Auth Repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(logger *logger.Logger, db *sqlx.DB) *repository {
	return &repository{logger, db}
}

// Create user
func (r *repository) Create() error {
	r.logger.Info("Call auth repo")
	return nil
}
