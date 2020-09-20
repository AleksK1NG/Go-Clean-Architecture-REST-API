package server

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// type deliveries struct {
// 	auth auth.Delivery
// }
//
// type useCases struct {
// 	auth auth.UseCase
// }
//
// type repositories struct {
// 	auth auth.Repository
// }

// Server struct
type server struct {
	echo   *echo.Echo
	config *config.Config
	l      *logger.Logger
}

// New server constructor
func NewServer(config *config.Config, logger *logger.Logger) *server {
	e := echo.New()

	return &server{e, config, logger}
}

// Run server depends on config SSL option
func (s *server) Run() error {
	if s.config.Server.SSL {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		s.MapRoutes(s.echo)

		s.echo.Server.ReadTimeout = time.Second * s.config.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.config.Server.WriteTimeout

		go func() {
			s.l.Info("Server is listening on port: %v", zap.String("PORT", s.config.Server.Port))
			if err := s.echo.StartTLS(s.config.Server.Port, certFile, keyFile); err != nil {
				s.l.Fatal("error starting TLS server %v", zap.String("echo.StartTLS", err.Error()))
			}

		}()

		go func() {
			s.l.Info("Starting Debug server on port: %v", zap.String("PORT", s.config.Server.PprofPort))
			if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.l.Error("PPROF", zap.String("PPROF ListenAndServe", err.Error()))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.l.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)

	} else {
		e := echo.New()
		s.MapRoutes(e)

		server := &http.Server{
			Addr:           s.config.Server.Port,
			ReadTimeout:    time.Second * s.config.Server.ReadTimeout,
			WriteTimeout:   time.Second * s.config.Server.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			s.l.Info("Server is listening on port: %v", zap.String("PORT", s.config.Server.Port))
			if err := e.StartServer(server); err != nil {
				s.l.Fatal("error starting TLS server %v", zap.String("StartServer", err.Error()))
			}
		}()

		go func() {
			s.l.Info("Starting Debug server on port: %v", zap.String("PORT", s.config.Server.PprofPort))
			if err := http.ListenAndServe(s.config.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.l.Error("PPROF", zap.String("PPROF ListenAndServe", err.Error()))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.l.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}
}

func (s *server) ConfigLayers() {

}
