package repository

import (
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/jmoiron/sqlx"
)

// Auth Repository
type AuthRepository struct {
	l  *logger.Logger
	db *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(l *logger.Logger, db *sqlx.DB) *AuthRepository {
	return &AuthRepository{l: l, db: db}
}

// Create user
func (r *AuthRepository) Create() error {
	return nil
}
