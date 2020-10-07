package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
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
	var n models.News
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create news", zap.String("ReqID", utils.GetRequestID(c)))

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

// Update news item handler
func (h handlers) Update() echo.HandlerFunc {
	var n models.News
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update", zap.String("ReqID", utils.GetRequestID(c)))

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			h.log.Error(
				"Update uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		if err := c.Bind(&n); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}
		n.ID = newsUUID

		updatedNews, err := h.newsUC.Update(ctx, &n)
		if err != nil {
			h.log.Error(
				"newsUC.Update",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"Created news",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", updatedNews.ID.String()),
		)

		return c.JSON(http.StatusOK, updatedNews)
	}
}
