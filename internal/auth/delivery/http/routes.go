package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/labstack/echo/v4"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, authUC auth.UseCase, sessUC session.UCSession, cfg *config.Config) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:user_id", h.GetUserByID())
	//authGroup.Use(middleware.AuthJWTMiddleware(authUC, cfg))
	authGroup.Use(middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	authGroup.PUT("/:user_id", h.Update(), middleware.OwnerOrAdminMiddleware())
	authGroup.DELETE("/:user_id", h.Delete(), middleware.RoleBasedAuthMiddleware([]string{"admin"}))
	authGroup.GET("/me", h.GetMe())
}
