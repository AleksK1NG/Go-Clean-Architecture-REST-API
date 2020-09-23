package repository

import (
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Auth Repository
type repository struct {
	l  *logger.Logger
	db *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(l *logger.Logger, db *sqlx.DB) *repository {
	return &repository{l: l, db: db}
}

// Create user
func (r *repository) Create() error {
	return nil
}
