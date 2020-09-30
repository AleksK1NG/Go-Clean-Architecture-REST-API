package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/auth/dto"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
func (r *repository) Register(ctx context.Context, user *models.User) (*models.User, error) {

	createUserQuery := `INSERT INTO users (first_name, last_name, email, password, role, about, avatar, phone_number, address,
	               		city, gender, postcode, birthday, created_at, updated_at, login_date)
						VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, $12, $13, now(), now(), now()) 
						RETURNING *`

	var u models.User
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role,
		&user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City, &user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(&u); err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return &u, nil
}

// Delete existing user
func (r *repository) Delete(ctx context.Context, userID uuid.UUID) error {
	deleteUserQuery := `DELETE FROM users WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
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

	return nil
}

// Get user by id
func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	getUserQuery := `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       				 address, city, gender, postcode, birthday, created_at, updated_at, login_date  
					 FROM users 
					 WHERE user_id = $1`

	var user models.User
	if err := r.db.GetContext(ctx, &user, getUserQuery, userID); err != nil {
		return nil, err
	}

	return &user, nil
}

// Find users by name
func (r *repository) FindByName(ctx context.Context, query *dto.FindUserQuery) (*models.UsersList, error) {
	getTotalCount := `SELECT COUNT(user_id) FROM users WHERE first_name ILIKE '%' || $1 || '%' or last_name ILIKE '%' || $1 || '%'`

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, query.Name); err != nil {
		return nil, err
	}

	findUsers := `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, address,
	              city, gender, postcode, birthday, created_at, updated_at, login_date 
				  FROM users 
				  WHERE first_name ILIKE '%' || $1 || '%' or last_name ILIKE '%' || $1 || '%'
				  ORDER BY first_name, last_name`

	rows, err := r.db.QueryxContext(ctx, findUsers, query.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.PQ.GetSize()),
		Page:       query.PQ.GetPage(),
		Size:       query.PQ.GetSize(),
		HasMore:    utils.GetHasMore(query.PQ.GetPage(), totalCount, query.PQ.GetSize()),
		Users:      users,
	}, nil
}

// Get users with pagination
func (r *repository) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	getTotal := `SELECT COUNT(user_id) FROM users`

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, err
	}

	getUsers := `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       			 address, city, gender, postcode, birthday, created_at, updated_at, login_date
				 FROM users 
				 ORDER BY COALESCE(NULLIF($1, ''), first_name) OFFSET $2 LIMIT $3`

	var users []*models.User
	if err := r.db.SelectContext(
		ctx,
		&users,
		getUsers,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, err
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}
