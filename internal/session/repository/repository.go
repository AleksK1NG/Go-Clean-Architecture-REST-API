package repository

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix = "api-session:"
)

// Session repository
type sessionRepo struct {
	redisPool  redis.RedisPool
	basePrefix string
	cfg        *config.Config
}

// Session repository constructor
func NewSessionRepository(redisPool redis.RedisPool, cfg *config.Config) session.SessRepository {
	return &sessionRepo{redisPool: redisPool, basePrefix: basePrefix, cfg: cfg}
}

// Create session in redis
func (s *sessionRepo) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	session.SessionID = uuid.New().String()
	sessionKey := s.createKey(session.SessionID)

	if err := s.redisPool.SetexJSONContext(ctx, sessionKey, expire, session); err != nil {
		return "", errors.WithMessage(err, "sessRepo CreateSession SetexJSONContext")
	}

	return sessionKey, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sess := &models.Session{}
	if err := s.redisPool.GetJSONContext(ctx, sessionID, sess); err != nil {
		return nil, errors.WithMessage(err, "sessRepo GetSessionByID GetJSONContext")
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	return s.redisPool.DeleteContext(ctx, sessionID)
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
