package redis

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Create new redis pool
func NewRedisPool(cfg *config.Config) (*redis.Pool, error) {
	redisHost := cfg.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}
	pool := newPool(redisHost)

	cleanupHook()
	return pool, nil
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     60,
		MaxActive:   250,
		Wait:        true,
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
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := Pool.Close(); err != nil {
			return
		}
		os.Exit(0)
	}()
}
