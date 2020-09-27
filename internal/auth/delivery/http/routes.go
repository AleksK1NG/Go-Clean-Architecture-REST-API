package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/labstack/echo"
)

// Map auth routes
func MapAuthRoutes(ag *echo.Group, h auth.Handlers, cfg *config.Config, l *logger.Logger) {
	ag.POST("/create", h.Create())
	ag.PUT("/:user_id", h.Update())
	ag.GET("/:user_id", h.GetUserByID())
}
