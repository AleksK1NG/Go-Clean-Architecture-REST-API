package repository

import (
	"context"
	"database/sql"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Auth Repository
type repository struct {
	logger *logger.Logger
	db     *sqlx.DB
	redis  *redis.RedisClient
}

// Auth Repository constructor
func NewAuthRepository(logger *logger.Logger, db *sqlx.DB, redis *redis.RedisClient) auth.Repository {
	return &repository{logger, db, redis}
}

// Create new user
func (r *repository) Register(ctx context.Context, user *models.User) (*models.User, error) {

	var u models.User
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City,
		&user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(&u); err != nil {
		return nil, err
	}

	return &u, nil

}

// Update existing user
func (r *repository) Update(ctx context.Context, user *models.UserUpdate) (*models.User, error) {

	var u models.User
	if err := r.db.GetContext(ctx, &u, updateUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City, &user.Gender,
		&user.Postcode, &user.Birthday, &user.ID,
	); err != nil {
		return nil, err
	}

	if err := r.redis.Delete(u.UserID.String()); err != nil {
		r.logger.Error("Delete", zap.String("ERROR", err.Error()))
	}

	return &u, nil
}

// Delete existing user
func (r *repository) Delete(ctx context.Context, userID uuid.UUID) error {

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

	if err := r.redis.Delete(userID.String()); err != nil {
		r.logger.Error("Delete", zap.String("ERROR", err.Error()))
	}

	return nil
}

// Get user by id
func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.redis.GetIfExistsJSON(userID.String(), &user); err != nil {
		if err != httpErrors.NotExists {
			r.logger.Error("GetIfExistsJSON", zap.String("ERROR", err.Error()))
		}
	} else {
		return &user, nil
	}

	if err := r.db.GetContext(ctx, &user, getUserQuery, userID); err != nil {
		return nil, err
	}

	if err := r.redis.SetEXJSON(userID.String(), 3600, &user); err != nil {
		r.logger.Error("SetEXJSON", zap.String("ERROR", err.Error()))
	}

	return &user, nil
}

// Find users by name
func (r *repository) FindByName(ctx context.Context, query *dto.FindUserQuery) (*models.UsersList, error) {

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, query.Name); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryxContext(ctx, findUsers, query.Name, query.PQ.GetOffset(), query.PQ.GetLimit())
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			}
		}
	}()

	var users = make([]*models.User, 0, query.PQ.GetSize())
	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
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

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, err
	}

	var users = make([]*models.User, 0, pq.GetSize())
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

// Find user by email
func (r *repository) FindByEmail(ctx context.Context, loginDTO *dto.LoginDTO) (*models.User, error) {
	var user models.User

	if err := r.redis.GetIfExistsJSON(loginDTO.Email, &user); err != nil {
		if err != httpErrors.NotExists {
			r.logger.Error("GetIfExistsJSON", zap.String("ERROR", err.Error()))
		}
	} else {
		return &user, nil
	}

	if err := r.db.GetContext(ctx, &user, findUserByEmail, loginDTO.Email); err != nil {
		return nil, err
	}

	if err := r.redis.SetEXJSON(loginDTO.Email, 3600, &user); err != nil {
		r.logger.Error("SetEXJSON", zap.String("ERROR", err.Error()))
	}

	return &user, nil
}
