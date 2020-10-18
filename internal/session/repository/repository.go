package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

const (
	basePrefix = "api-session:"
)

// Session repository
type sessionRepository struct {
	redisPool  *redis.Pool
	logger     *logger.Logger
	basePrefix string
	cfg        *config.Config
}

// Session repository constructor
func NewSessionRepository(redisPool *redis.Pool, log *logger.Logger, cfg *config.Config) session.SessRepository {
	return &sessionRepository{redisPool: redisPool, logger: log, basePrefix: basePrefix, cfg: cfg}
}

func (s *sessionRepository) createKey(sessionId string) string {
	return s.basePrefix + sessionId
}

// Create session in redis
func (s *sessionRepository) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	session.SessionID = uuid.New().String()
	sessionKey := s.createKey(session.SessionID)
	if err := utils.RedisMarshalJSON(ctx, s.redisPool, sessionKey, expire, session); err != nil {
		return "", err
	}
	return sessionKey, nil
}

// Get session by id
func (s *sessionRepository) GetSessionByID(ctx context.Context, sessionId string) (*models.Session, error) {
	sess := &models.Session{}
	if err := utils.RedisUnmarshalJSON(ctx, s.redisPool, sessionId, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepository) DeleteByID(ctx context.Context, sessionId string) error {
	return utils.RedisDeleteKey(ctx, s.redisPool, sessionId)
}
