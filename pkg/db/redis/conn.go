package redis

import (
	"encoding/json"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Redis db interface
type RedisPool interface {
	GetClient() *redis.Client
	Ping()
	PingContext(ctx context.Context)
	GetBytesContext(ctx context.Context, key string) ([]byte, error)
	SetBytes(key string, value []byte) error
	SetexBytes(key string, duration int, value []byte) error
	Exists(key string) (int64, error)
	Delete(key string) error
	DeleteContext(ctx context.Context, key string) error
	IncrCtx(ctx context.Context, counterKey string) (int64, error)
	SetexJSON(key string, seconds int, value interface{}) error
	SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error
	GetJSON(key string, model interface{}) error
	GetJSONContext(ctx context.Context, key string, model interface{}) error
}

// Redis client
type RedisClient struct {
	config *config.Config
	client *redis.Client
}

// Returns new redis client
func NewRedisClient(config *config.Config) *RedisClient {
	redisHost := config.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: 200,
		PoolSize:     12000,
		PoolTimeout:  240 * time.Second,
		Password:     "", // no password set
		DB:           0,  // use default DB
	})

	cleanupHook(client)
	return &RedisClient{config: config, client: client}
}

func cleanupHook(client *redis.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		if err := client.Close(); err != nil {
			logger.Errorf("RedisClient Close: %s", err.Error())
		}
		os.Exit(0)
	}()
}

// Get Redis Client
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Redis ping method
func (r *RedisClient) PingContext(ctx context.Context) {
	s := r.client.Ping(ctx).String()
	logger.Infof("PING: %s", s)
}

// Redis ping method
func (r *RedisClient) Ping() {
	s := r.client.Ping(context.Background()).String()
	logger.Infof("PING: %s", s)
}

// Get by key string, return []byte
func (r *RedisClient) GetBytesContext(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// Set by key string, return []byte
func (r *RedisClient) SetBytes(key string, value []byte) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

// Setex by key string, return []byte
func (r *RedisClient) SetexBytes(key string, duration int, value []byte) error {
	seconds := time.Second * time.Duration(duration)
	return r.client.Set(context.Background(), key, value, seconds).Err()
}

// Exists by key string, return bool
func (r *RedisClient) Exists(key string) (int64, error) {
	return r.client.Exists(context.Background(), key).Result()
}

// Delete by key string
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// Delete by key string
func (r *RedisClient) DeleteContext(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Incr by key string
func (r *RedisClient) IncrCtx(ctx context.Context, counterKey string) (int64, error) {
	return r.client.Incr(ctx, counterKey).Result()
}

// Set JSON value
func (r *RedisClient) SetexJSON(key string, seconds int, value interface{}) error {
	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	duration := time.Second * time.Duration(seconds)
	return r.client.Set(context.Background(), key, bytes, duration).Err()
}

// Set JSON value
func (r *RedisClient) SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error {
	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	duration := time.Second * time.Duration(seconds)
	return r.client.Set(ctx, key, bytes, duration).Err()
}

// Get JSON value
func (r *RedisClient) GetJSONContext(ctx context.Context, key string, model interface{}) error {
	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &model); err != nil {
		return err
	}
	return nil
}

// Get JSON value
func (r *RedisClient) GetJSON(key string, model interface{}) error {
	bytes, err := r.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &model); err != nil {
		return err
	}
	return nil
}
