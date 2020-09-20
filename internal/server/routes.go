package server

import (
	authHttp "github.com/AleksK1NG/api-mc/internal/auth/delivery/http"
	authRepository "github.com/AleksK1NG/api-mc/internal/auth/repository"
	authUseCase "github.com/AleksK1NG/api-mc/internal/auth/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	loggerFormat = "${time_rfc3339}: ${method} ${uri}, STATUS: ${status} latency: ${latency_human}, bytes_in: ${bytes_in} \n"
)

// Map Server Routes
func (s *server) MapRoutes(e *echo.Echo) {
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
	authRepo := authRepository.NewAuthRepository(s.l, nil)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.l, s.config, authRepo)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.config, authUC, s.l)

	{
		authHttp.MapAuthRoutes(auth, authHandlers, authUC, s.config, s.l)
		// auth_routes.MapAuthRoutes(auth, s.h, s.useCases, s.config, s.logger)
		// post_routes.MapPostRoutes(post, s.h, s.useCases, s.config, s.logger)
		// comment_routes.MapCommentRoutes(comment, s.h, s.useCases, s.config, s.logger)
		health.GET("", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"status": "OK"})
		})
	}
}
