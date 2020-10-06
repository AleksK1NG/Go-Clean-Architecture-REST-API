package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
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
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create news", zap.String("ReqID", utils.GetRequestID(c)))

		var n models.News
		if err := c.Bind(&n); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		createdNews, err := h.newsUC.Create(ctx, &n)
		if err != nil {
			h.log.Error(
				"News create",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"Created news",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", createdNews.ID.String()),
		)

		return c.JSON(http.StatusOK, createdNews)
	}
}
