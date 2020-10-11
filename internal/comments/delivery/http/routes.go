package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/middleware"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
)

// Map news routes
func MapCommentsRoutes(commGroup *echo.Group, h comments.Handlers, authUC auth.UseCase, sessUC session.UCSession, cfg *config.Config, log *logger.Logger) {
	commGroup.POST("/create", h.Create(), middleware.AuthSessionMiddleware(sessUC, authUC, cfg, log))
}
