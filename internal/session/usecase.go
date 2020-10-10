package session

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"time"
)

// Session use case
type UCSession interface {
	CreateSession(ctx context.Context, session models.Session, expire time.Duration) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
