package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
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
func (u *useCase) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	return u.sessionRepo.CreateSession(ctx, session, expire)
}

// Delete session by id
func (u *useCase) DeleteByID(ctx context.Context, sessionID string) error {
	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

// get session by id
func (u *useCase) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
