package server

import (
	authHttp "github.com/AleksK1NG/api-mc/internal/auth/delivery/http"
	authRepository "github.com/AleksK1NG/api-mc/internal/auth/repository"
	authUseCase "github.com/AleksK1NG/api-mc/internal/auth/usecase"
	commentsHttp "github.com/AleksK1NG/api-mc/internal/comments/delivery/http"
	commentsRepository "github.com/AleksK1NG/api-mc/internal/comments/repository"
	commentsUseCase "github.com/AleksK1NG/api-mc/internal/comments/usecase"
	metricsMiddleware "github.com/AleksK1NG/api-mc/internal/middleware"
	newsHttp "github.com/AleksK1NG/api-mc/internal/news/delivery/http"
	newsRepository "github.com/AleksK1NG/api-mc/internal/news/repository"
	newsUseCase "github.com/AleksK1NG/api-mc/internal/news/usecase"
	sessionRepository "github.com/AleksK1NG/api-mc/internal/session/repository"
	"github.com/AleksK1NG/api-mc/internal/session/usecase"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/metric"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
	"net/http"
)

const (
	loggerFormat = "${time_rfc3339}: ${method} ${uri}, STATUS: ${status} latency: ${latency_human}, bytes_in: ${bytes_in} \n"
)

// Map Server Handlers
func (s *server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.config.Metrics.Url, s.config.Metrics.ServiceName)
	if err != nil {
		s.logger.Error("CreateMetrics", zap.String("ERROR", err.Error()))
	}
	s.logger.Info(
		"Metrics available",
		zap.String("URL", s.config.Metrics.Url),
		zap.String("ServiceName", s.config.Metrics.ServiceName),
	)

	e.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.RequestID())
	e.Use(metricsMiddleware.MetricsMiddleware(metrics))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: loggerFormat,
	}))
	// e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(metricsMiddleware.DebugMiddleware(s.config.Server.Debug, s.logger))

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	newsGroup := v1.Group("/news")
	commGroup := v1.Group("/comments")

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.logger, s.db, s.redis, "api-auth")
	nRepo := newsRepository.NewNewsRepository(s.logger, s.db, s.redis, "api-news")
	cRepo := commentsRepository.NewCommentsRepository(s.logger, s.db, s.redis)
	sRepo := sessionRepository.NewSessionRepository(s.redis, s.logger, "api-session", s.config)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.logger, s.config, aRepo)
	newsUC := newsUseCase.NewNewsUseCase(s.logger, s.config, nRepo)
	commUC := commentsUseCase.NewCommentsUseCase(s.logger, s.config, cRepo)
	sessUC := usecase.NewSessionUseCase(sRepo, s.logger, s.config)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.config, authUC, sessUC, s.logger)
	newsHandlers := newsHttp.NewNewsHandlers(s.config, newsUC, s.logger)
	commHandlers := commentsHttp.NewCommentsHandlers(s.config, commUC, s.logger)

	{
		authHttp.MapAuthRoutes(authGroup, authHandlers, authUC, sessUC, s.config, s.logger)
		newsHttp.MapNewsRoutes(newsGroup, newsHandlers, authUC, sessUC, s.config, s.logger)
		commentsHttp.MapCommentsRoutes(commGroup, commHandlers, authUC, sessUC, s.config, s.logger)

		health.GET("", func(c echo.Context) error {
			s.logger.Info("Health check", zap.String("RequestID", utils.GetRequestID(c)))
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		})
	}

	return nil
}
