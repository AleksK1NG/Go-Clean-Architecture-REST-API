package session

import (
	"github.com/AleksK1NG/api-mc/internal/models"
	"time"
)

// Session repository
type SessRepo interface {
	Create(session models.Session, expire time.Duration) (string, error)
	GetSessByID(sessionID string) (models.Session, error)
	DeleteByID(sessionID string) error
}
