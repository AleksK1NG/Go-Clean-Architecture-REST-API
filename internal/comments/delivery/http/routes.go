package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/labstack/echo/v4"
)

// Map news routes
func MapCommentsRoutes(commGroup *echo.Group, h comments.Handlers, authUC auth.UseCase, sessUC session.UCSession, cfg *config.Config) {
	commGroup.POST("", h.Create(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	commGroup.DELETE("/:comment_id", h.Delete(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	commGroup.PUT("/:comment_id", h.Update(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	commGroup.GET("/:comment_id", h.GetByID())
	commGroup.GET("/byNewsId/:news_id", h.GetAllByNewsID())
}
