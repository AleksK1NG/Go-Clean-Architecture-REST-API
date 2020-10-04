package auth

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/google/uuid"
)

// User repository interface
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*dto.UserWithToken, error)
	Login(ctx context.Context, loginDTO *dto.LoginDTO) (*dto.UserWithToken, error)
	Update(ctx context.Context, user *models.UserUpdate) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByName(ctx context.Context, query *dto.FindUserQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
}
