package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
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
func (s *sessionRepo) CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRepo.CreateSession")
	defer span.Finish()

	sess.SessionID = uuid.New().String()
	sessionKey := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRepo.GetSessionByID")
	defer span.Finish()

	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRep.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRepo.DeleteByID")
	defer span.Finish()

	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return errors.Wrap(err, "sessionRepo.DeleteByID")
	}
	return nil
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
