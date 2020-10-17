package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Session repository
type sessionRepository struct {
	redisPool  *redis.Pool
	logger     *logger.Logger
	basePrefix string
	cfg        *config.Config
}

// Session repository constructor
func NewSessionRepository(redisPool *redis.Pool, log *logger.Logger, prefix string, cfg *config.Config) session.SessRepository {
	return &sessionRepository{redisPool, log, prefix, cfg}
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

func (s *sessionRepository) setexSessionJSON(ctx context.Context, key string, expire int, session *models.Session) error {
	conn, err := s.redisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	sessBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	values, err := redis.String(conn.Do("SETEX", key, expire, sessBytes))
	if err != nil {
		return err
	}

	s.logger.Info("REDIS SET", zap.String("values", fmt.Sprintf("%#v", values)))
	return nil
}

func (s *sessionRepository) getSessionJSON(ctx context.Context, key string) (*models.Session, error) {
	conn, err := s.redisPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	values, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	sess := &models.Session{}
	if err := json.Unmarshal(values, sess); err != nil {
		return nil, err
	}

	s.logger.Info("REDIS GET", zap.String("session", fmt.Sprintf("%#v", sess)))
	return sess, nil
}

func (s *sessionRepository) deleteSession(ctx context.Context, sessionID string) error {
	conn, err := s.redisPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	values, err := redis.Values(conn.Do("DEL", sessionID))
	if err != nil {
		return err
	}

	s.logger.Info("REDIS DEL", zap.String("values", fmt.Sprintf("%#v", values)))
	return nil
}

// Create session in redis
func (s *sessionRepository) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			session.SessionID = uuid.New().String()
			sessionKey := s.createKey(session.SessionID)
			s.logger.Info("CreateSession ID", zap.String("SessionID", session.SessionID))
			s.logger.Info("CreateSession sessionKey", zap.String("sessionKey", sessionKey))
			if err := s.setexSessionJSON(ctx, sessionKey, expire, session); err != nil {
				return "", err
			}
			return sessionKey, nil
		}
	}
}

// Get session by id
func (s *sessionRepository) GetSessionByID(ctx context.Context, sessionId string) (*models.Session, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// key := s.createKey(sessionId)
			sess, err := s.getSessionJSON(ctx, sessionId)
			if err != nil {
				return nil, err
			}
			return sess, nil
		}
	}
}

// Delete session by id
func (s *sessionRepository) DeleteByID(ctx context.Context, sessionId string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// key := s.createKey(sessionId)
			if err := s.deleteSession(ctx, sessionId); err != nil {
				return err
			}
			return nil
		}
	}
}
