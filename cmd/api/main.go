package main

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/server"
	"go.uber.org/zap"
	"log"
)

const (
	configPath = "./config/config-docker"
)

func main() {
	log.Println("Starting auth server")

	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		l.Fatal("fatal", zap.String("LoadConfig", err.Error()))
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		l.Fatal("fatal", zap.String("ParseConfig", err.Error()))
	}

	s := server.NewServer(cfg, l)
	log.Fatal(s.Run())
}
