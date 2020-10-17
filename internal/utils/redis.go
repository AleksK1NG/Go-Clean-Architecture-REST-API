package utils

import (
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

// Redis SETEX value as json
func RedisMarshalJSON(ctx context.Context, pool *redis.Pool, key string, duration int, v interface{}) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	valueBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", key, duration, valueBytes)
	if err != nil {
		return err
	}

	return nil
}

// Redis GET value as json
func RedisUnmarshalJSON(ctx context.Context, pool *redis.Pool, key string, v interface{}) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	valueBytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(valueBytes, v); err != nil {
		return err
	}

	return nil
}

// Redis DEL value as json
func RedisDeleteKey(ctx context.Context, pool *redis.Pool, key string) error {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
