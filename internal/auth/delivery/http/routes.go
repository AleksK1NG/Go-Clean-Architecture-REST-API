package http

import (
	"github.com/labstack/echo/v4"

	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/middleware"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/find", h.FindByName())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:user_id", h.GetUserByID())
	// authGroup.Use(middleware.AuthJWTMiddleware(authUC, cfg))
	authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", h.GetMe())
	authGroup.GET("/token", h.GetCSRFToken())
	authGroup.POST("/:user_id/avatar", h.UploadAvatar(), mw.CSRF)
	authGroup.PUT("/:user_id", h.Update(), mw.OwnerOrAdminMiddleware(), mw.CSRF)
	authGroup.DELETE("/:user_id", h.Delete(), mw.CSRF, mw.RoleBasedAuthMiddleware([]string{"admin"}))
}
