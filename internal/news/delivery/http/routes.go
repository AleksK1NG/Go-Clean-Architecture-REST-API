package http

import (
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/labstack/echo/v4"
)

// Map news routes
func MapNewsRoutes(newsGroup *echo.Group, h news.Handlers, mw *middleware.MiddlewareManager) {
	newsGroup.POST("/create", h.Create(), mw.AuthSessionMiddleware)
	newsGroup.PUT("/:news_id", h.Update(), mw.AuthSessionMiddleware)
	newsGroup.DELETE("/:news_id", h.Delete(), mw.AuthSessionMiddleware)
	newsGroup.GET("/:news_id", h.GetByID())
	newsGroup.GET("/search", h.SearchByTitle())
	newsGroup.GET("", h.GetNews())
}
