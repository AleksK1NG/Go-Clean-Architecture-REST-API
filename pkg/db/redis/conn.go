package redis

import (
	"encoding/json"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Redis db interface
type RedisPool interface {
	GetPool() *redis.Pool
	PingContext(ctx context.Context) error
	GetBytesContext(ctx context.Context, key string) ([]byte, error)
	SetBytes(key string, value []byte) error
	SetexBytes(key string, durationSec int, value []byte) error
	Exists(key string) (bool, error)
	Delete(key string) error
	DeleteContext(ctx context.Context, key string) error
	GetKeys(pattern string) ([]string, error)
	Incr(counterKey string) (int, error)
	SetEXJSON(key string, seconds int, value interface{}) error
	SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error
	GetJSON(key string, model interface{}) error
	GetJSONContext(ctx context.Context, key string, model interface{}) error
}

// Redis client
type RedisClient struct {
	config *config.Config
	pool   *redis.Pool
}

// Returns new redis client
func NewRedisClient(config *config.Config) *RedisClient {
	redisHost := config.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}
	pool := newPool(redisHost)

	cleanupHook(pool)
	return &RedisClient{config: config, pool: pool}
}

func cleanupHook(pool *redis.Pool) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
		os.Exit(0)
	}()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		MaxActive:   12000,
		//Wait:        true,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Get pool
func (r *RedisClient) GetPool() *redis.Pool {
	return r.pool
}

// Redis ping method
func (r *RedisClient) PingContext(ctx context.Context) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	ping, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	log.Printf("PING: %v", ping)
	return nil
}

// Redis ping method
func (r *RedisClient) Ping() error {
	conn := r.pool.Get()
	defer conn.Close()

	ping, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	log.Printf("PING: %v", ping)
	return nil
}

// Get by key string, return []byte
func (r *RedisClient) GetBytesContext(ctx context.Context, key string) ([]byte, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var data []byte
	data, err = redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

// Set by key string, return []byte
func (r *RedisClient) SetBytes(key string, value []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

// Setex by key string, return []byte
func (r *RedisClient) SetexBytes(key string, durationSec int, value []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, durationSec, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

// Exists by key string, return bool
func (r *RedisClient) Exists(key string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Delete by key string
func (r *RedisClient) Delete(key string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// Delete by key string
func (r *RedisClient) DeleteContext(ctx context.Context, key string) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	return err
}

// Get by keys string
func (r *RedisClient) GetKeys(pattern string) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	iter := 0
	var keys []string
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

// Incr by key string
func (r *RedisClient) Incr(counterKey string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}

// Set JSON value
func (r *RedisClient) SetEXJSON(key string, seconds int, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	_, err = redis.String(conn.Do("SETEX", key, seconds, bytes))
	if err != nil {
		return err
	}

	return nil
}

// Set JSON value
func (r *RedisClient) SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	_, err = redis.String(conn.Do("SETEX", key, seconds, bytes))
	if err != nil {
		return err
	}

	return nil
}

// Get JSON value
func (r *RedisClient) GetJSONContext(ctx context.Context, key string, model interface{}) error {
	ctx, _ = context.WithTimeout(context.Background(), time.Second*5)
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", key))
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
	conn := r.pool.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &model); err != nil {
		return err
	}

	return nil
}

// Get JSON value
func (r *RedisClient) GetIfExistsJSON(key string, model interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if !ok {
		return httpErrors.NotExists
	}

	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &model); err != nil {
		return err
	}

	return nil
}
