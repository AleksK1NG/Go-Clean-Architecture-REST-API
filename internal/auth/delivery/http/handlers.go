package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/labstack/echo"
)

// Auth handlers
type Handlers struct {
	cfg *config.Config
	uc  auth.UseCase
	l   *logger.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, uc auth.UseCase, l *logger.Logger) *Handlers {
	return &Handlers{cfg: cfg, uc: uc, l: l}
}

// Crate new user
func (h *Handlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(201, nil)
	}
}

// Fet user by id
func (h *Handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, nil)
	}
}
