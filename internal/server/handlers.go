package server

import (
	authHttp "github.com/AleksK1NG/api-mc/internal/auth/delivery/http"
	authRepository "github.com/AleksK1NG/api-mc/internal/auth/repository"
	authUseCase "github.com/AleksK1NG/api-mc/internal/auth/usecase"
	commentsHttp "github.com/AleksK1NG/api-mc/internal/comments/delivery/http"
	commentsRepository "github.com/AleksK1NG/api-mc/internal/comments/repository"
	commentsUseCase "github.com/AleksK1NG/api-mc/internal/comments/usecase"
	apiMiddlewares "github.com/AleksK1NG/api-mc/internal/middleware"
	newsHttp "github.com/AleksK1NG/api-mc/internal/news/delivery/http"
	newsRepository "github.com/AleksK1NG/api-mc/internal/news/repository"
	newsUseCase "github.com/AleksK1NG/api-mc/internal/news/usecase"
	sessionRepository "github.com/AleksK1NG/api-mc/internal/session/repository"
	"github.com/AleksK1NG/api-mc/internal/session/usecase"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/metric"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	_ "net/http/pprof"
)

// Map Server Handlers
func (s *server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.config.Metrics.Url, s.config.Metrics.ServiceName)
	if err != nil {
		logger.Errorf("CreateMetrics Error: %s", err.Error())
	}
	logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.config.Metrics.Url,
		s.config.Metrics.ServiceName,
	)

	if s.config.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Use(apiMiddlewares.RequestLoggerMiddleware())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(apiMiddlewares.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	// e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.config.Server.Debug {
		e.Use(apiMiddlewares.DebugMiddleware(s.config.Server.Debug))
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	newsGroup := v1.Group("/news")
	commGroup := v1.Group("/comments")

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db, s.redisPool)
	nRepo := newsRepository.NewNewsRepository(s.db, s.redisPool)
	cRepo := commentsRepository.NewCommentsRepository(s.db, s.redisPool)
	sRepo := sessionRepository.NewSessionRepository(s.redisPool, s.config)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.config, aRepo)
	newsUC := newsUseCase.NewNewsUseCase(s.config, nRepo)
	commUC := commentsUseCase.NewCommentsUseCase(s.config, cRepo)
	sessUC := usecase.NewSessionUseCase(sRepo, s.config)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.config, authUC, sessUC)
	newsHandlers := newsHttp.NewNewsHandlers(s.config, newsUC)
	commHandlers := commentsHttp.NewCommentsHandlers(s.config, commUC)

	{
		authHttp.MapAuthRoutes(authGroup, authHandlers, authUC, sessUC, s.config)
		newsHttp.MapNewsRoutes(newsGroup, newsHandlers, authUC, sessUC, s.config)
		commentsHttp.MapCommentsRoutes(commGroup, commHandlers, authUC, sessUC, s.config)

		health.GET("", func(c echo.Context) error {
			logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		})
	}

	return nil
}
