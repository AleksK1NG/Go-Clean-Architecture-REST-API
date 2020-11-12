package server

import (
	"github.com/AleksK1NG/api-mc/docs"
	"strings"

	//_ "github.com/AleksK1NG/api-mc/docs"
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
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/metric"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	_ "net/http/pprof" //prof
)

// Map Server Handlers
func (s *server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.config.Metrics.URL, s.config.Metrics.ServiceName)
	if err != nil {
		logger.Errorf("CreateMetrics Error: %s", err.Error())
	}
	logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.config.Metrics.URL,
		s.config.Metrics.ServiceName,
	)

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	nRepo := newsRepository.NewNewsRepository(s.db)
	cRepo := commentsRepository.NewCommentsRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisPool, s.config)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.config, aRepo, s.redisPool)
	newsUC := newsUseCase.NewNewsUseCase(s.config, nRepo, s.redisPool)
	commUC := commentsUseCase.NewCommentsUseCase(s.config, cRepo, s.redisPool)
	sessUC := usecase.NewSessionUseCase(sRepo, s.config)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.config, authUC, sessUC)
	newsHandlers := newsHttp.NewNewsHandlers(s.config, newsUC)
	commHandlers := commentsHttp.NewCommentsHandlers(s.config, commUC)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.config, []string{"*"})

	e.Use(mw.RequestLoggerMiddleware)

	docs.SwaggerInfo.Title = "Swagger Example API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if s.config.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

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
	e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	// e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.config.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	newsGroup := v1.Group("/news")
	commGroup := v1.Group("/comments")

	{
		authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
		newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw)
		commentsHttp.MapCommentsRoutes(commGroup, commHandlers, mw)

		health.GET("", func(c echo.Context) error {
			logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		})
	}

	return nil
}
