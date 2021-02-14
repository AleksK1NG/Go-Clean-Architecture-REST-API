package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

// News handlers
type newsHandlers struct {
	cfg    *config.Config
	newsUC news.UseCase
	logger logger.Logger
}

// NewNewsHandlers News handlers constructor
func NewNewsHandlers(cfg *config.Config, newsUC news.UseCase, logger logger.Logger) news.Handlers {
	return &newsHandlers{cfg: cfg, newsUC: newsUC, logger: logger}
}

// Create godoc
// @Summary Create news
// @Description Create news handler
// @Tags News
// @Accept json
// @Produce json
// @Success 201 {object} models.News
// @Router /news/create [post]
func (h newsHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.Create")
		defer span.Finish()

		n := &models.News{}
		if err := c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdNews, err := h.newsUC.Create(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdNews)
	}
}

// Update godoc
// @Summary Update news
// @Description Update news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} models.News
// @Router /news/{id} [put]
func (h newsHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.Update")
		defer span.Finish()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.News{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.NewsID = newsUUID

		updatedNews, err := h.newsUC.Update(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedNews)
	}
}

// GetByID godoc
// @Summary Get by id news
// @Description Get by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} models.News
// @Router /news/{id} [get]
func (h newsHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.GetByID")
		defer span.Finish()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newsByID, err := h.newsUC.GetNewsByID(ctx, newsUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newsByID)
	}
}

// Delete godoc
// @Summary Delete news
// @Description Delete by id news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {string} string	"ok"
// @Router /news/{id} [delete]
func (h newsHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.Delete")
		defer span.Finish()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.newsUC.Delete(ctx, newsUUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetNews godoc
// @Summary Get all news
// @Description Get all news with pagination
// @Tags News
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.NewsList
// @Router /news [get]
func (h newsHandlers) GetNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.GetNews")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newsList, err := h.newsUC.GetNews(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newsList)
	}
}

// SearchByTitle godoc
// @Summary Search by title
// @Description Search news by title
// @Tags News
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.NewsList
// @Router /news/search [get]
func (h newsHandlers) SearchByTitle() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "newsHandlers.SearchByTitle")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newsList, err := h.newsUC.SearchByTitle(ctx, c.QueryParam("title"), pq)

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newsList)
	}
}
