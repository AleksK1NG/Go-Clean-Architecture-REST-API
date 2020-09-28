package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/labstack/echo"
)

// Map auth routes
func MapAuthRoutes(ag *echo.Group, h auth.Handlers, authUC auth.UseCase, cfg *config.Config) {
	ag.POST("/create", h.Create())
	ag.GET("/find", h.FindByName())
	ag.GET("/all", h.GetUsers())
	ag.GET("/:user_id", h.GetUserByID())
	ag.Use(middleware.AuthJWTMiddleware(authUC, cfg))
	ag.PUT("/:user_id", h.Update())
	ag.DELETE("/:user_id", h.Delete())
}
