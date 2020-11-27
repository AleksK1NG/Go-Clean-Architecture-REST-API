package main

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/server"
	"github.com/AleksK1NG/api-mc/pkg/db/aws"
	"github.com/AleksK1NG/api-mc/pkg/db/postgres"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"log"
	"os"
)

// @title Go Example REST API
// @version 1.0
// @description Example Golang REST API
// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru
// @BasePath /api/v1
func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("LogLevel: %s, Mode: %s, SSL: %v", cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	awsClient, err := aws.NewAWSClient(cfg.AWS.Endpoint, cfg.AWS.MinioAccessKey, cfg.AWS.MinioSecretKey, cfg.AWS.UseSSL)
	if err != nil {
		appLogger.Errorf("AWS Client init: %s", err)
	}
	appLogger.Info("AWS S3 connected")

	s := server.NewServer(cfg, psqlDB, redisClient, awsClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
