package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
)

// Auth useCase
type useCase struct {
	logger   *logger.Logger
	cfg      *config.Config
	authRepo auth.Repository
}

// Auth useCase constructor
func NewAuthUseCase(logger *logger.Logger, cfg *config.Config, authRepo auth.Repository) auth.UseCase {
	return &useCase{logger: logger, cfg: cfg, authRepo: authRepo}
}

// Create new user
func (u *useCase) Register(ctx context.Context, user *models.User) (*dto.UserWithToken, error) {
	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(err.Error())
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update existing user
func (u *useCase) Update(ctx context.Context, user *models.UserUpdate) (*models.User, error) {
	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	if err := user.PrepareUpdate(); err != nil {
		return nil, err
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	updatedUser.SanitizePassword()

	return updatedUser, nil
}

// Delete new user
func (u *useCase) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

// Get user by id
func (u *useCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.SanitizePassword()

	return user, nil
}

// Find users by name
func (u *useCase) FindByName(ctx context.Context, query *dto.FindUserQuery) (*models.UsersList, error) {
	users, err := u.authRepo.FindByName(ctx, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Get users with pagination
func (u *useCase) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	users, err := u.authRepo.GetUsers(ctx, pq)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Login user, returns user model with jwt token
func (u *useCase) Login(ctx context.Context, loginDTO *dto.LoginDTO) (*dto.UserWithToken, error) {
	user, err := u.authRepo.FindByEmail(ctx, loginDTO)
	if err != nil {
		return nil, err
	}

	if err = user.ComparePasswords(loginDTO.Password); err != nil {
		return nil, err
	}

	user.SanitizePassword()

	token, err := utils.GenerateJWTToken(user, u.cfg)
	if err != nil {
		return nil, err
	}

	return &dto.UserWithToken{
		User:  user,
		Token: token,
	}, nil
}
