package redis

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/gomodule/redigo/redis"
	"log"
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

	cleanupHook(pool)
	return pool, nil
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   300,
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

func cleanupHook(pool *redis.Pool) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if err := pool.Close(); err != nil {
			log.Printf("POOL CLOSE ERROR: %s", err.Error())
			return
		}
		os.Exit(0)
	}()
}
