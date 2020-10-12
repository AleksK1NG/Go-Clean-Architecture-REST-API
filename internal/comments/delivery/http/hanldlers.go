package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

// Comments handlers
type handlers struct {
	cfg   *config.Config
	comUC comments.UseCase
	log   *logger.Logger
}

// Comments handlers constructor
func NewCommentsHandlers(cfg *config.Config, comUC comments.UseCase, log *logger.Logger) comments.Handlers {
	return &handlers{cfg: cfg, comUC: comUC, log: log}
}

// Create comment
func (h *handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create", zap.String("ReqID", utils.GetRequestID(c)))

		comment := &models.Comment{}
		if err := c.Bind(comment); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		createdComment, err := h.comUC.Create(ctx, comment)
		if err != nil {
			h.log.Error(
				"comUC.Create",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"createdComment",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", createdComment.ID.String()),
		)

		return c.JSON(http.StatusCreated, createdComment)
	}
}

// Update comment
func (h *handlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update", zap.String("ReqID", utils.GetRequestID(c)))

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			h.log.Error(
				"uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		comm := &dto.UpdateCommDTO{}
		if err := c.Bind(comm); err != nil {
			h.log.Error(
				"c.Bind",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}
		comm.ID = commID

		updatedComment, err := h.comUC.Update(ctx, comm)
		if err != nil {
			h.log.Error(
				"comUC.Update",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info(
			"updatedComment",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", updatedComment.ID.String()),
		)

		return c.JSON(http.StatusOK, updatedComment)
	}
}

// Delete comment
func (h *handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update", zap.String("ReqID", utils.GetRequestID(c)))

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			h.log.Error(
				"uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		if err := h.comUC.Delete(ctx, commID); err != nil {
			h.log.Error(
				"comUC.Delete",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetByID comment
func (h *handlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update", zap.String("ReqID", utils.GetRequestID(c)))

		commID, err := uuid.Parse(c.Param("comment_id"))
		if err != nil {
			h.log.Error(
				"uuid.Parse",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		comment, err := h.comUC.GetByID(ctx, commID)
		if err != nil {
			h.log.Error(
				"comUC.GetByID",
				zap.String("ReqID", utils.GetRequestID(c)),
				zap.String("Error:", err.Error()),
			)
			return c.JSON(errors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// GetAllByNewsID comments
func (h *handlers) GetAllByNewsID() echo.HandlerFunc {
	panic("implement me")
}
