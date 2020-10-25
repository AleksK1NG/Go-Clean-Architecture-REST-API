package main

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/server"
	"github.com/AleksK1NG/api-mc/pkg/db/postgres"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
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
		logger.Fatalf("Postgresql init: %s", err.Error())
	} else {
		logger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	logger.Info("Redis connected")

	s := server.NewServer(cfg, psqlDB, redisClient)
	log.Fatal(s.Run())
}
