package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
)

// News handlers
type handlers struct {
	cfg    *config.Config
	newsUC news.UseCase
	log    *logger.Logger
}

// News handlers constructor
func NewNewsHandlers(cfg *config.Config, newsUC news.UseCase, log *logger.Logger) *handlers {
	return &handlers{cfg, newsUC, log}
}

// Create news handler
func (h handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "Ok")
	}
}
