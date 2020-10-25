package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/labstack/echo/v4"
)

// Map news routes
func MapNewsRoutes(newsGroup *echo.Group, h news.Handlers, authUC auth.UseCase, sessUC session.UCSession, cfg *config.Config) {
	newsGroup.POST("/create", h.Create(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	newsGroup.PUT("/:news_id", h.Update(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	newsGroup.DELETE("/:news_id", h.Delete(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg))
	newsGroup.GET("/:news_id", h.GetByID())
	newsGroup.GET("/search", h.SearchByTitle())
	newsGroup.GET("", h.GetNews())
}
