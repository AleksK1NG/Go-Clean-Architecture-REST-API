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

	logger.InitLogger(cfg)
	logger.Infof("LogLevel: %s, Mode: %s, SSL: %v", cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		logger.Fatal("Init postgres", zap.String("error", err.Error()))
	} else {
		logger.Info("Postgres connected", zap.String("Status", fmt.Sprintf("%#v", psqlDB.Stats())))
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	logger.Info("Redis connected", zap.String("Status", fmt.Sprintf("%#v", *redisClient.GetPool().PoolStats())))

	s := server.NewServer(cfg, psqlDB, redisClient)
	log.Fatal(s.Run())
}
