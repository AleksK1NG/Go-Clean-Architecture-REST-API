package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

// Auth handlers
type handlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	log    *logger.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log *logger.Logger) *handlers {
	return &handlers{cfg, authUC, log}
}

// Crate new user
func (h *handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Create user", zap.String("ReqID", utils.GetRequestID(c)))

		var user models.User
		if err := c.Bind(&user); err != nil {
			h.log.Error("Create c.Bind", zap.String("ReqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(http.StatusBadRequest, errors.BadRequest)
		}

		createdUser, err := h.authUC.Create(ctx, &user)
		if err != nil {
			h.log.Error("auth repo create", zap.String("reqID", utils.GetRequestID(c)), zap.String("Error:", err.Error()))
			return c.JSON(errors.ErrorResponse(err))
		}

		h.log.Info("Created user", zap.String("reqID", utils.GetRequestID(c)), zap.String("ID", createdUser.ID.String()))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Fet user by id
func (h *handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errors.BadQueryParams)
		}
		return c.JSON(http.StatusOK, paginationQuery)
	}
}
