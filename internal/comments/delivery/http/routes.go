package http

import (
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/labstack/echo/v4"
)

// Map news routes
func MapCommentsRoutes(commGroup *echo.Group, h comments.Handlers, mw *middleware.MiddlewareManager) {
	commGroup.POST("", h.Create(), mw.AuthSessionMiddleware)
	commGroup.DELETE("/:comment_id", h.Delete(), mw.AuthSessionMiddleware)
	commGroup.PUT("/:comment_id", h.Update(), mw.AuthSessionMiddleware)
	commGroup.GET("/:comment_id", h.GetByID())
	commGroup.GET("/byNewsId/:news_id", h.GetAllByNewsID())
}
