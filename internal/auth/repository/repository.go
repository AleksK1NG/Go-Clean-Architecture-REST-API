package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
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

// Create new user
func (r *repository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// createUserQuery := `INSERT INTO users (first_name, last_name, email, password, role, about, avatar, phone_number, address,
	//                		city, role, gender, postcode, birthday, created_at, updated_at, login_date)
	// 					VALUES (:first_name, :last_name, :email, :password, :role, :about, :avatar, :phone_number, :address,
	//     				:city, :role, :gender, :postcode, :birthday, now(), now(), now())`
	return user, nil

}
