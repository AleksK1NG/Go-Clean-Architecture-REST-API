package repository

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
)

// Session repository
type sessionRepository struct {
	redis      *redis.RedisClient
	logger     *logger.Logger
	basePrefix string
	cfg        *config.Config
}

// Session repository constructor
func NewSessionRepository(redis *redis.RedisClient, log *logger.Logger, prefix string, cfg *config.Config) session.SessRepository {
	return &sessionRepository{redis, log, prefix, cfg}
}

func (s *sessionRepository) createKey(sessionId string) string {
	return s.basePrefix + sessionId
}

func (s *sessionRepository) convertToString(session *models.Session) (string, error) {
	sessionJSON, err := json.Marshal(session)

	if err != nil {
		return "", err
	}

	return string(sessionJSON), nil
}

func (s *sessionRepository) convertFromString(sessionString string) (*models.Session, error) {
	var storedSession models.Session

	if err := json.Unmarshal([]byte(sessionString), &storedSession); err != nil {
		return nil, err
	}

	return &storedSession, nil
}

func (s *sessionRepository) convertToBytes(session *models.Session) ([]byte, error) {
	sessionJSON, err := json.Marshal(session)

	if err != nil {
		return nil, err
	}

	return sessionJSON, nil
}

func (s *sessionRepository) convertFromBytes(sessionBytes []byte) (*models.Session, error) {
	var storedSession models.Session

	if err := json.Unmarshal(sessionBytes, &storedSession); err != nil {
		return nil, err
	}

	return &storedSession, nil
}

// Create session in redis
func (s *sessionRepository) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		session.ID = uuid.New().String()
		sessionKey := s.createKey(session.ID)
		if err := s.redis.SetEXJSON(sessionKey, expire, &session); err != nil {
			return "", err
		}

		return session.ID, nil
	}
}

// Get session by id
func (s *sessionRepository) GetSessionByID(ctx context.Context, sessionId string) (*models.Session, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		key := s.createKey(sessionId)
		storedSession := &models.Session{}
		if err := s.redis.GetIfExistsJSON(key, storedSession); err != nil {
			return nil, err
		}

		return storedSession, nil
	}
}

// Delete session by id
func (s *sessionRepository) DeleteByID(ctx context.Context, sessionId string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.redis.Delete(s.createKey(sessionId)); err != nil {
			return err
		}

		return nil
	}
}
