package utils

import (
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
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

func RedisSet(ctx context.Context, pool *redis.Pool, key string, duration int, v interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	do, err := conn.Do("SETEX", key, duration, string(bytes))
	if err != nil {
		return err
	}
	log.Printf("SET REPLY: %#v", do)
	return nil
}

func RedisGet(ctx context.Context, pool *redis.Pool, key string, v interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, v); err != nil {
		return err
	}

	log.Printf("GET REPLY: %#v", v)
	return nil
}
