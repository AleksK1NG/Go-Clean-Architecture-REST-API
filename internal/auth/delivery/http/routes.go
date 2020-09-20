package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/labstack/echo"
)

// Map auth routes
func MapAuthRoutes(ag *echo.Group, h auth.Handlers, uc auth.UseCase, cfg *config.Config, l *logger.Logger) {
	ag.GET("/:user_id", h.GetUserByID())
	ag.POST("", h.Create())
}
