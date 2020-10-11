package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
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
	panic("implement me")
}

// Delete comment
func (h *handlers) Delete() echo.HandlerFunc {
	panic("implement me")
}

// GetByID comment
func (h *handlers) GetByID() echo.HandlerFunc {
	panic("implement me")
}

// GetAllByNewsID comments
func (h *handlers) GetAllByNewsID() echo.HandlerFunc {
	panic("implement me")
}
