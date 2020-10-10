package repository

import (
	"encoding/json"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"time"
)

// Session repository
type SessionRepository struct {
	redis      *redis.RedisClient
	logger     *logger.Logger
	basePrefix string
}

// Session repository constructor
func NewSessionRepository(redis *redis.RedisClient, logger *logger.Logger, basePrefix string) *SessionRepository {
	return &SessionRepository{redis: redis, logger: logger, basePrefix: basePrefix}
}

func (s *SessionRepository) createKey(sessionId string) string {
	return s.basePrefix + sessionId
}

func (s *SessionRepository) convertToString(session *models.Session) (string, error) {
	sessionJSON, err := json.Marshal(session)

	if err != nil {
		return "", err
	}

	return string(sessionJSON), nil
}

func (s *SessionRepository) convertFromString(sessionString string) (*models.Session, error) {
	var storedSession models.Session

	if err := json.Unmarshal([]byte(sessionString), &storedSession); err != nil {
		return nil, err
	}

	return &storedSession, nil
}

func (s *SessionRepository) convertToBytes(session *models.Session) ([]byte, error) {
	sessionJSON, err := json.Marshal(session)

	if err != nil {
		return nil, err
	}

	return sessionJSON, nil
}

func (s *SessionRepository) convertFromBytes(sessionBytes []byte) (*models.Session, error) {
	var storedSession models.Session

	if err := json.Unmarshal(sessionBytes, &storedSession); err != nil {
		return nil, err
	}

	return &storedSession, nil
}

// Create session in redis
func (s *SessionRepository) CreateSession(session models.Session, expire time.Duration) (string, error) {
	sessionKey := s.createKey(session.ID)

	session.ID = uuid.New().String()

	if err := s.redis.SetEXJSON(sessionKey, int(expire.Seconds()), &session); err != nil {
		return "", err
	}

	return session.ID, nil
}

// Get session by id
func (s *SessionRepository) GetSessionByID(sessionId string) (*models.Session, error) {
	key := s.createKey(sessionId)

	storedSession := &models.Session{}

	if err := s.redis.GetIfExistsJSON(key, storedSession); err != nil {
		return nil, err
	}

	return storedSession, nil
}

// Delete session by id
func (s *SessionRepository) DeleteByID(sessionId string) error {
	if err := s.redis.Delete(s.createKey(sessionId)); err != nil {
		return err
	}

	return nil
}
