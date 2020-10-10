package usecase

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"time"
)

// Session use case
type useCase struct {
	sessionRepo session.SessRepository
	logger      *logger.Logger
	cfg         *config.Config
}

// New session use case constructor
func NewSessionUseCase(sessionRepo session.SessRepository, logger *logger.Logger, cfg *config.Config) session.UCSession {
	return &useCase{sessionRepo: sessionRepo, logger: logger, cfg: cfg}
}

// Create new session
func (u *useCase) CreateSession(session models.Session, expire time.Duration) (string, error) {
	return u.sessionRepo.CreateSession(session, expire)
}

// Delete session by id
func (u *useCase) DeleteByID(sessionID string) error {
	return u.sessionRepo.DeleteByID(sessionID)
}

// get session by id
func (u *useCase) GetSessionByID(sessionID string) (*models.Session, error) {
	return u.sessionRepo.GetSessionByID(sessionID)
}
