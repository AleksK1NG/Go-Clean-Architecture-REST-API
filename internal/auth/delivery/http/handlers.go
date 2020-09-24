package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/labstack/echo"
	"net/http"
)

// Auth handlers
type Handlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	log    *logger.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log *logger.Logger) *Handlers {
	return &Handlers{cfg, authUC, log}
}

// Crate new user
func (h *Handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {

		if err := h.authUC.Create(); err != nil {
			return c.JSON(errors.ParseErrors(err).Status(), errors.ParseErrors(err))
		}

		return c.JSON(http.StatusCreated, "Ok")
	}
}

// Fet user by id
func (h *Handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Ok")
	}
}
