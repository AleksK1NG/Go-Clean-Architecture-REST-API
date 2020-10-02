package server

import (
	authHttp "github.com/AleksK1NG/api-mc/internal/auth/delivery/http"
	authRepository "github.com/AleksK1NG/api-mc/internal/auth/repository"
	authUseCase "github.com/AleksK1NG/api-mc/internal/auth/usecase"
	"github.com/AleksK1NG/api-mc/internal/utils"
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

	// echo.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.RequestID())
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
	// echo.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	auth := v1.Group("/auth")
	// post := v1.Group("/posts")
	// comment := v1.Group("/comments")

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.logger, s.db)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.logger, s.config, aRepo, s.redis)

	// Init handlers
	aHandlers := authHttp.NewAuthHandlers(s.config, authUC, s.logger)
	{
		authHttp.MapAuthRoutes(auth, aHandlers, authUC, s.config, s.logger)
		// auth_routes.MapAuthRoutes(auth, s.h, s.useCases, s.config, s.logger)
		// post_routes.MapPostRoutes(post, s.h, s.useCases, s.config, s.logger)
		// comment_routes.MapCommentRoutes(comment, s.h, s.useCases, s.config, s.logger)
		health.GET("", func(c echo.Context) error {
			s.logger.Info("Health check", zap.String("RequestID", utils.GetRequestID(c)))
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		})
	}

	return nil
}
