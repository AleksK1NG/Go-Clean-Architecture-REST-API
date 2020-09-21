package server

import (
	"fmt"
	authHttp "github.com/AleksK1NG/api-mc/internal/auth/delivery/http"
	authRepository "github.com/AleksK1NG/api-mc/internal/auth/repository"
	authUseCase "github.com/AleksK1NG/api-mc/internal/auth/usecase"
	"github.com/AleksK1NG/api-mc/internal/db/postgres"
	"github.com/AleksK1NG/api-mc/internal/db/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

const (
	loggerFormat = "${time_rfc3339}: ${method} ${uri}, STATUS: ${status} latency: ${latency_human}, bytes_in: ${bytes_in} \n"
)

// Map Server Routes
func (s *server) MapRoutes(e *echo.Echo) error {
	psqlDB, err := postgres.NewPsqlDB(s.config)
	if err != nil {
		s.l.Error("", zap.String("init psql", err.Error()))
		return err
	}
	s.l.Info("Postgres connected", zap.String("DB Status: %#v", fmt.Sprintf("%#v", psqlDB.Stats())))
	redisConn := redis.NewRedisClient(s.config)
	s.l.Info("Redis connected", zap.String("port", fmt.Sprintf("%s", s.config.Redis.RedisAddr)))

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
	aRepo := authRepository.NewAuthRepository(s.l, psqlDB)

	// Init useCases
	aUseCase := authUseCase.NewAuthUseCase(s.l, s.config, aRepo, redisConn)

	// Init handlers
	aHandlers := authHttp.NewAuthHandlers(s.config, aUseCase, s.l)

	{
		authHttp.MapAuthRoutes(auth, aHandlers, s.config, s.l)
		// auth_routes.MapAuthRoutes(auth, s.h, s.useCases, s.config, s.logger)
		// post_routes.MapPostRoutes(post, s.h, s.useCases, s.config, s.logger)
		// comment_routes.MapCommentRoutes(comment, s.h, s.useCases, s.config, s.logger)
		health.GET("", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"status": "OK"})
		})
	}

	return nil
}
