package server

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/db/redis"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	certFile       = "ssl/server.crt"
	keyFile        = "ssl/server.pem"
	maxHeaderBytes = 1 << 20
)

// Server struct
type server struct {
	echo   *echo.Echo
	config *config.Config
	logger *logger.Logger
	db     *sqlx.DB
	redis  *redis.RedisClient
}

// New server constructor
func NewServer(config *config.Config, logger *logger.Logger, db *sqlx.DB, redis *redis.RedisClient) *server {
	e := echo.New()

	return &server{e, config, logger, db, redis}
}

// Run server depends on config SSL option
func (s *server) Run() error {
	if s.config.Server.SSL {

		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.config.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.config.Server.WriteTimeout

		go func() {
			s.logger.Info("Server is listening", zap.String("PORT", s.config.Server.Port))
			if err := s.echo.StartTLS(s.config.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatal("error starting TLS server", zap.String("echo.StartTLS", err.Error()))
			}

		}()

		go func() {
			s.logger.Info("Starting Debug server", zap.String("PORT", s.config.Server.PprofPort))
			if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Error("PPROF", zap.String("PPROF ListenAndServe", err.Error()))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)

	} else {
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
			s.logger.Info("Server is listening", zap.String("PORT", s.config.Server.Port))
			if err := e.StartServer(server); err != nil {
				s.logger.Fatal("error starting TLS server", zap.String("StartServer", err.Error()))
			}
		}()

		go func() {
			s.logger.Info("Starting Debug server", zap.String("PORT", s.config.Server.PprofPort))
			if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Error("PPROF", zap.String("PPROF ListenAndServe", err.Error()))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}
}

func (s *server) ConfigLayers() {

}
