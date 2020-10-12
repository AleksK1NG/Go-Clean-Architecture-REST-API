package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/dto"
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
func NewNewsHandlers(cfg *config.Config, newsUC news.UseCase, log *logger.Logger) news.Handlers {
	return &handlers{cfg, newsUC, log}
}

// Create news handler
func (h handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create news", zap.String("ReqID", utils.GetRequestID(c)))

		n := &models.News{}
		if err := c.Bind(n); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		createdNews, err := h.newsUC.Create(ctx, n)
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
			zap.String("ID", createdNews.NewsID.String()),
		)

		return c.JSON(http.StatusOK, createdNews)
	}
}

// Update news item handler
func (h handlers) Update() echo.HandlerFunc {
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

		n := &models.News{}
		if err = c.Bind(n); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}
		n.NewsID = newsUUID

		updatedNews, err := h.newsUC.Update(ctx, n)
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
			zap.String("ID", updatedNews.NewsID.String()),
		)

		return c.JSON(http.StatusOK, updatedNews)
	}
}

// Get news by id
func (h handlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetByID", zap.String("ReqID", utils.GetRequestID(c)))

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			h.log.Error(
				"Update uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		newsByID, err := h.newsUC.GetNewsByID(ctx, newsUUID)
		if err != nil {
			h.log.Error(
				"newsUC.GetNewsByID",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"GetByID",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", newsByID.UserID.String()),
		)

		return c.JSON(http.StatusOK, newsByID)
	}
}

// Delete news handler
func (h handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetByID", zap.String("ReqID", utils.GetRequestID(c)))

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			h.log.Error(
				"Update uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		if err := h.newsUC.Delete(ctx, newsUUID); err != nil {
			h.log.Error(
				"newsUC.GetNewsByID",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"GetByID",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", newsUUID.String()),
		)

		return c.NoContent(http.StatusOK)
	}
}

// Get all news with pagination
func (h handlers) GetNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetByID", zap.String("ReqID", utils.GetRequestID(c)))

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			h.log.Error(
				"GetPaginationFromCtx",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		newsList, err := h.newsUC.GetNews(ctx, pq)
		if err != nil {
			h.log.Error(
				"newsUC.GetNewsByID",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"GetByID",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.Int("Length", len(newsList.News)),
		)

		return c.JSON(http.StatusOK, newsList)
	}
}

// Search by title
func (h handlers) SearchByTitle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetByID", zap.String("ReqID", utils.GetRequestID(c)))

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			h.log.Error(
				"GetPaginationFromCtx",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		newsList, err := h.newsUC.SearchByTitle(ctx, &dto.FindNewsDTO{
			Title: c.QueryParam("title"),
			PQ:    pq,
		})

		if err != nil {
			h.log.Error(
				"newsUC.GetNewsByID",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"GetByID",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.Int("Length", len(newsList.News)),
		)

		return c.JSON(http.StatusOK, newsList)
	}
}
