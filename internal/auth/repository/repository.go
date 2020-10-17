package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Auth Repository
type repository struct {
	logger    *logger.Logger
	db        *sqlx.DB
	redisPool *redis.Pool
	prefix    string
}

// Auth Repository constructor
func NewAuthRepository(logger *logger.Logger, db *sqlx.DB, redis *redis.Pool, prefix string) auth.Repository {
	return &repository{logger, db, redis, prefix}
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

	if err := r.deleteUser(ctx, r.generateUserKey(u.UserID.String())); err != nil {
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

	if err := r.deleteUser(ctx, r.generateUserKey(userID.String())); err != nil {
		r.logger.Error("Delete", zap.String("ERROR", err.Error()))
	}

	return nil
}

// Get user by id
func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	userJSON, err := r.getUserJSON(ctx, r.generateUserKey(userID.String()))
	if err != nil {
		r.logger.Error("getUserJSON", zap.String("ERROR", err.Error()))
	}
	if userJSON != nil {
		r.logger.Info("FROM REDIS")
		return userJSON, nil
	}

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, userID).StructScan(user); err != nil {
		return nil, err
	}

	if err := r.setexUserJSON(ctx, r.generateUserKey(userID.String()), 50, user); err != nil {
		r.logger.Error("setexUserJSON", zap.String("ERROR", err.Error()))
	}

	return user, nil
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

	userJSON, err := r.getUserJSON(ctx, loginDTO.Email)
	if err != nil {
		r.logger.Error("getUserJSON", zap.String("ERROR", err.Error()))
	}
	if userJSON != nil {
		return userJSON, nil
	}

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmail, loginDTO.Email).StructScan(user); err != nil {
		return nil, err
	}

	if err := r.setexUserJSON(ctx, r.generateUserKey(loginDTO.Email), 50, user); err != nil {
		r.logger.Error("setexUserJSON", zap.String("ERROR", err.Error()))
	}

	return user, nil
}

func (r *repository) setexUserJSON(ctx context.Context, key string, duration int, user *models.User) error {
	conn, err := r.redisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", key, duration, userBytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) getUserJSON(ctx context.Context, key string) (*models.User, error) {
	conn, err := r.redisPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	userBytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	if err := json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) deleteUser(ctx context.Context, key string) error {
	conn, err := r.redisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) generateUserKey(userID string) string {
	return r.prefix + userID
}
