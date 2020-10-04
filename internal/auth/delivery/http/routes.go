package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/labstack/echo"
)

// Map auth routes
func MapAuthRoutes(ag *echo.Group, h auth.Handlers, authUC auth.UseCase, cfg *config.Config, logger *logger.Logger) {
	ag.POST("/register", h.Register())
	ag.POST("/login", h.Login())
	ag.POST("/logout", h.Logout())
	ag.GET("/find", h.FindByName())
	ag.GET("/all", h.GetUsers())
	ag.GET("/:user_id", h.GetUserByID())
	ag.Use(middleware.AuthJWTMiddleware(authUC, cfg, logger))
	ag.PUT("/:user_id", h.Update(), middleware.OwnerOrAdminMiddleware(logger))
	ag.DELETE("/:user_id", h.Delete(), middleware.RoleBasedAuthMiddleware([]string{"admin"}, logger))
	ag.GET("/me", h.GetMe())
}
