package auth

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
)

// User repo interface
type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
}
