package repository

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Auth Repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
}

// Auth Repository constructor
func NewAuthRepository(logger *logger.Logger, db *sqlx.DB) auth.Repository {
	return &repository{logger, db}
}

// Create new user
func (r *repository) Create(ctx context.Context, user *models.User) (*models.User, error) {

	createUserQuery := `INSERT INTO users (first_name, last_name, email, password, role, about, avatar, phone_number, address,
	               		city, gender, postcode, birthday, created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, $12, $13, now(), now(), now()) RETURNING *`

	var u models.User
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role,
		&user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City, &user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(&u); err != nil {
		r.logger.Error("QueryRowxContext", zap.String("ERROR", err.Error()))
		return nil, err
	}

	r.logger.Info("USER", zap.String("USER", fmt.Sprintf("%#v", u)))

	return &u, nil

}

// Update existing user
func (r *repository) Update(ctx context.Context, user *models.UserUpdate) (*models.User, error) {
	updateUserQuery := `UPDATE users 
						SET first_name = COALESCE(NULLIF($1, ''), first_name),
						    last_name = COALESCE(NULLIF($2, ''), last_name),
						    email = COALESCE(NULLIF($3, ''), email),
						    role = COALESCE(NULLIF($4, ''), role),
						    about = COALESCE(NULLIF($5, ''), about),
						    avatar = COALESCE(NULLIF($6, ''), avatar),
						    phone_number = COALESCE(NULLIF($7, ''), phone_number),
						    address = COALESCE(NULLIF($8, ''), address),
						    city = COALESCE(NULLIF($9, ''), city),
						    gender = COALESCE(NULLIF($10, ''), gender),
						    postcode = COALESCE(NULLIF($11, 0), postcode),
						    birthday = COALESCE(NULLIF($12, '')::date, birthday),
						    updated_at = now()
						WHERE user_id = $13
						RETURNING *
						`

	var u models.User
	if err := r.db.GetContext(ctx, &u, updateUserQuery, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.About, &user.Avatar, &user.PhoneNumber,
		&user.Address, &user.City, &user.Gender, &user.Postcode, &user.Birthday, &user.ID,
	); err != nil {
		r.logger.Error("Get", zap.String("ERROR", err.Error()))
		return nil, err
	}
	r.logger.Info("USER", zap.String("USER", fmt.Sprintf("%#v", u)))
	return &u, nil
}
