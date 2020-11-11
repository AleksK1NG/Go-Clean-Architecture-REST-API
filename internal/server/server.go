package server

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	_ "github.com/AleksK1NG/api-mc/docs"
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
	config    *config.Config
	db        *sqlx.DB
	redisPool redis.RedisPool
}

// New server constructor
func NewServer(config *config.Config, db *sqlx.DB, redisPool redis.RedisPool) *server {
	e := echo.New()

	return &server{e, config, db, redisPool}
}

func (s *server) Run() error {
	if s.config.Server.SSL {

		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.config.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.config.Server.WriteTimeout

		go func() {
			logger.Infof("Server is listening on PORT: %s", s.config.Server.Port)
			s.echo.Server.ReadTimeout = time.Second * s.config.Server.ReadTimeout
			s.echo.Server.WriteTimeout = time.Second * s.config.Server.WriteTimeout
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.config.Server.Port, certFile, keyFile); err != nil {
				logger.Fatalf("Error starting TLS server: ", err.Error())
			}

		}()

		go func() {
			logger.Infof("Starting Debug server on PORT: %s", s.config.Server.PprofPort)
			if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
				logger.Errorf("Error PPROF ListenAndServe: %s", err.Error())
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
		Addr:           s.config.Server.Port,
		ReadTimeout:    time.Second * s.config.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.config.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		logger.Infof("Server is listening on PORT: %s", s.config.Server.Port)
		if err := e.StartServer(server); err != nil {
			logger.Fatalf("Error starting server: ", err.Error())
		}
	}()

	go func() {
		logger.Infof("Starting Debug server on PORT: %s", s.config.Server.PprofPort)
		if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
			logger.Errorf("Error PPROF ListenAndServe: %s", err.Error())
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
