package usecase

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/aws"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo redis.RedisPool
	awsClient aws.AWSClient
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo redis.RedisPool, awsClient aws.AWSClient) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, awsClient: awsClient}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "authUC Register PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.WithMessage(err, "authUC Register GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update existing user
func (u *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "authUC Register PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err := u.redisRepo.Delete(u.generateUserKey(user.UserID.String())); err != nil {
		logger.Errorf("AuthUC Update redis delete: %s", err)
	}

	return updatedUser, nil
}

// Delete new user
func (u *authUC) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := u.redisRepo.Delete(u.generateUserKey(userID.String())); err != nil {
		logger.Errorf("AuthUC Delete redis delete: %s", err)
	}

	return nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := u.redisRepo.GetJSONContext(ctx, u.generateUserKey(userID.String()), user); err == nil {
		return user, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetexJSONContext(ctx, u.generateUserKey(userID.String()), cacheDuration, user); err != nil {
		logger.Errorf("AuthUC GetByID redis set: %s", err)
	}

	user.SanitizePassword()

	return user, nil
}

// Find users by name
func (u *authUC) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	return u.authRepo.FindByName(ctx, name, query)
}

// Get users with pagination
func (u *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	return u.authRepo.GetUsers(ctx, pq)
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "authUC GetUsers ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.WithMessage(err, "authUC GetUsers GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Upload user avatar
func (u *authUC) UploadAvatar(ctx context.Context, file aws.UploadInput) error {
	uploadInfo, err := u.awsClient.FileUpload(ctx, file)
	if err != nil {
		return err
	}

	logger.Infof("UploadAvatar: %#v", uploadInfo)
	return nil
}

func (u *authUC) generateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
