package session

import (
	"github.com/AleksK1NG/api-mc/internal/models"
	"time"
)

// Session use case
type UCSession interface {
	CreateSession(session models.Session, expire time.Duration) (string, error)
	GetSessionByID(sessionID string) (*models.Session, error)
	DeleteByID(sessionID string) error
}
