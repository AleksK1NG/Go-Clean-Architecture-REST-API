package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth Repository
type repository struct {
	db          *sqlx.DB
	redisClient redis.RedisPool
	basePrefix  string
}

// Auth Repository constructor
func NewAuthRepository(db *sqlx.DB, redisClient redis.RedisPool) auth.Repository {
	return &repository{db: db, redisClient: redisClient, basePrefix: basePrefix}
}

// Create new user
func (r *repository) Register(ctx context.Context, user *models.User) (*models.User, error) {

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City,
		&user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(&u); err != nil {
		return nil, err
	}

	return u, nil

}

// Update existing user
func (r *repository) Update(ctx context.Context, user *models.User) (*models.User, error) {

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City, &user.Gender,
		&user.Postcode, &user.Birthday, &user.UserID,
	); err != nil {
		return nil, err
	}

	if err := r.redisClient.DeleteContext(ctx, r.generateUserKey(u.UserID.String())); err != nil {
		logger.Errorf("DeleteContext: %s", err.Error())
	}
	return u, nil
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

	if err := r.redisClient.DeleteContext(ctx, r.generateUserKey(userID.String())); err != nil {
		logger.Errorf("DeleteContext: %s", err.Error())
	}

	return nil
}

// Get user by id
func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	user := &models.User{}

	if err := r.redisClient.GetJSONContext(ctx, r.generateUserKey(userID.String()), user); err != nil {
		logger.Errorf("GetJSONContext: %s", err.Error())
	}

	if err := r.db.QueryRowxContext(ctx, getUserQuery, userID).StructScan(user); err != nil {
		return nil, err
	}

	if err := r.redisClient.SetexJSONContext(ctx, r.generateUserKey(userID.String()), cacheDuration, user); err != nil {
		logger.Errorf("GetJSONContext: %s", err.Error())
	}

	return user, nil
}

// Find users by name
func (r *repository) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, name); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryxContext(ctx, findUsers, name, query.GetOffset(), query.GetLimit())
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

	var users = make([]*models.User, 0, query.GetSize())
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
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
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
func (r *repository) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {

	foundUser := &models.User{}

	if err := r.redisClient.GetJSONContext(ctx, r.generateUserKey(user.Email), foundUser); err != nil {
		logger.Errorf("GetJSONContext: %s", err.Error())
	}

	if err := r.db.QueryRowxContext(ctx, findUserByEmail, user.Email).StructScan(foundUser); err != nil {
		return nil, err
	}

	if err := r.redisClient.SetexJSONContext(ctx, r.generateUserKey(user.Email), cacheDuration, foundUser); err != nil {
		logger.Errorf("GetJSONContext: %s", err.Error())
	}

	return foundUser, nil
}

func (r *repository) generateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, userID)
}
