package session

import (
	"github.com/AleksK1NG/api-mc/internal/models"
	"time"
)

// Session use case
type UCSession interface {
	Create(session models.Session, expire time.Duration) (string, error)
	Delete(sessionID string) error
	GetByID(sessionID string) (models.Session, error)
}
