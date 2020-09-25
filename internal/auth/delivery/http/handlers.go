package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/labstack/echo"
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

		if err := h.authUC.Create(); err != nil {
			return c.JSON(errors.ParseErrors(err).Status(), errors.ParseErrors(err))
		}

		return c.JSON(http.StatusCreated, "Ok")
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
