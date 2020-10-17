package main

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/server"
	"github.com/AleksK1NG/api-mc/pkg/db/postgres"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"go.uber.org/zap"
	"log"
)

const (
	configPath = "./config/config-local"
)

func main() {
	log.Println("Starting api server")

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		l.Fatal("Init postgres", zap.String("error", err.Error()))
	}
	defer psqlDB.Close()

	if psqlDB != nil {
		l.Info("Postgres connected", zap.String("Status", fmt.Sprintf("%#v", psqlDB.Stats())))

	}

	redisPool, err := redis.NewRedisPool(cfg)
	if err != nil {
		l.Fatal("Init REDIS", zap.String("error", err.Error()))
	}
	if redisPool != nil {
		l.Info("Redis connected", zap.String("Status", fmt.Sprintf("%#v", redisPool.Stats())))
	}

	s := server.NewServer(cfg, l, psqlDB, redisPool)
	log.Fatal(s.Run())
}
