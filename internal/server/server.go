package server

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	_ "github.com/AleksK1NG/api-mc/docs"
	"github.com/AleksK1NG/api-mc/pkg/db/aws"
	"github.com/AleksK1NG/api-mc/pkg/db/redis"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	certFile       = "ssl/server.crt"
	keyFile        = "ssl/server.pem"
	maxHeaderBytes = 1 << 20
)

// Server struct
type server struct {
	echo      *echo.Echo
	cfg       *config.Config
	db        *sqlx.DB
	redisPool redis.RedisPool
	awsClient aws.AWSClient
}

// New server constructor
func NewServer(cfg *config.Config, db *sqlx.DB, redisPool redis.RedisPool, awsS3Client aws.AWSClient) *server {
	return &server{echo: echo.New(), cfg: cfg, db: db, redisPool: redisPool, awsClient: awsS3Client}
}

func (s *server) Run() error {
	if s.cfg.Server.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout

		go func() {
			logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
			s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
			s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.Server.Port, certFile, keyFile); err != nil {
				logger.Fatalf("Error starting TLS server: ", err)
			}

		}()

		go func() {
			logger.Infof("Starting Debug server on PORT: %s", s.cfg.Server.PprofPort)
			if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
				logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}

	e := echo.New()
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := e.StartServer(server); err != nil {
			logger.Fatalf("Error starting server: ", err)
		}
	}()

	go func() {
		logger.Infof("Starting Debug server on PORT: %s", s.cfg.Server.PprofPort)
		if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
			logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
