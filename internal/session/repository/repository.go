package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

const (
	basePrefix = "api-session:"
)

// Session repository
type sessionRepo struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
}

// Session repository constructor
func NewSessionRepository(redisClient *redis.Client, cfg *config.Config) session.SessRepository {
	return &sessionRepo{redisClient: redisClient, basePrefix: basePrefix, cfg: cfg}
}

// Create session in redis
func (s *sessionRepo) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	session.SessionID = uuid.New().String()
	sessionKey := s.createKey(session.SessionID)

	sessBytes, err := json.Marshal(session)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo CreateSession json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.WithMessage(err, "sessionRepo CreateSession redisClient.Set")
	}
	return sessionKey, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.WithMessage(err, "sessionRepo GetSessionByID redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, sess); err != nil {
		return nil, errors.WithMessage(err, "sessionRepo GetSessionByID json.Unmarshal")
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return errors.WithMessage(err, "sessionRepo DeleteByID")
	}
	return nil
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
