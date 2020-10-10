package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, authUC auth.UseCase, sessUC session.UCSession, cfg *config.Config, log *logger.Logger) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:user_id", h.GetUserByID())
	authGroup.Use(middleware.AuthJWTMiddleware(authUC, cfg, log))
	authGroup.PUT("/:user_id", h.Update(), middleware.OwnerOrAdminMiddleware(log))
	authGroup.DELETE("/:user_id", h.Delete(), middleware.RoleBasedAuthMiddleware([]string{"admin"}, log))
	authGroup.GET("/me", h.GetMe())
}
