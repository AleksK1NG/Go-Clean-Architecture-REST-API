//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// Auth repository interface
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, file models.UploadInput) (*models.User, error)
}
