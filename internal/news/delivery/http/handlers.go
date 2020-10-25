package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// News handlers
type handlers struct {
	cfg    *config.Config
	newsUC news.UseCase
}

// News handlers constructor
func NewNewsHandlers(cfg *config.Config, newsUC news.UseCase) news.Handlers {
	return &handlers{cfg, newsUC}
}

// Create news handler
func (h handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		n := &models.News{}
		if err := c.Bind(n); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		createdNews, err := h.newsUC.Create(ctx, n)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, createdNews)
	}
}

// Update news item handler
func (h handlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		n := &models.News{}
		if err = c.Bind(n); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}
		n.NewsID = newsUUID

		updatedNews, err := h.newsUC.Update(ctx, n)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, updatedNews)
	}
}

// Get news by id
func (h handlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		newsByID, err := h.newsUC.GetNewsByID(ctx, newsUUID)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, newsByID)
	}
}

// Delete news handler
func (h handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		newsUUID, err := uuid.Parse(c.Param("news_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		if err := h.newsUC.Delete(ctx, newsUUID); err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get all news with pagination
func (h handlers) GetNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		newsList, err := h.newsUC.GetNews(ctx, pq)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, newsList)
	}
}

// Search by title
func (h handlers) SearchByTitle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		newsList, err := h.newsUC.SearchByTitle(ctx, &dto.FindNewsDTO{
			Title: c.QueryParam("title"),
			PQ:    pq,
		})

		if err != nil {
			return utils.ErrResponseWithLog(c, err)
		}

		return c.JSON(http.StatusOK, newsList)
	}
}
