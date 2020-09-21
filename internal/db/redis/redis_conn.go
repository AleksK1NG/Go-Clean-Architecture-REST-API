package redis

import (
	"encoding/json"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Pool *redis.Pool
)

// Redis client
type RedisClient struct {
	config *config.Config
}

// Returns new redis client
func NewRedisClient(config *config.Config) *RedisClient {
	redisHost := config.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)

	cleanupHook()
	return &RedisClient{config: config}
}

// func init() {
// 	redisHost := os.Getenv("REDIS_HOST")
// 	if redisHost == "" {
// 		redisHost = ":6379"
// 	}
// 	Pool = newPool(redisHost)
// 	cleanupHook()
// }

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

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

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

// Redis ping method
func (r *RedisClient) Ping() error {
	conn := Pool.Get()
	defer conn.Close()

	ping, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	log.Printf("PING: %v", ping)
	return nil
}

// Get by key string, return []byte
func (r *RedisClient) GetBytes(key string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

// Get by key string, return []byte
func (r *RedisClient) GetMarshalToJSON(key string, value interface{}) error {
	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return fmt.Errorf("error getting key %s: %v", key, err)
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return nil
}

// Set by key string, return []byte
func (r *RedisClient) SetBytes(key string, value []byte) error {
	conn := Pool.Get()
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

// Set by key string, return []byte, time int seconds
func (r *RedisClient) SetExBytes(key string, seconds int, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, seconds, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

// Set by key string, return []byte, time int seconds
func (r *RedisClient) SetExBytesJSON(key string, seconds int, value interface{}) error {
	conn := Pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", key, seconds, bytes)
	if err != nil {
		return fmt.Errorf("error setting key %s to %s: %v", key, key, err)
	}

	return err
}

// Exists by key string, return bool
func (r *RedisClient) Exists(key string) (bool, error) {
	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Delete by key string
func (r *RedisClient) Delete(key string) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// Get by keys string
func (r *RedisClient) GetKeys(pattern string) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
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
	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}
