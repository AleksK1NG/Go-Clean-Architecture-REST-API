package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/comments"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
)

// Comments handlers
type handlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	sessUC session.UCSession
	log    *logger.Logger
}

// Comments handlers constructor
func NewCommentsHandlers(cfg *config.Config, authUC auth.UseCase, sessUC session.UCSession, log *logger.Logger) comments.Handlers {
	return &handlers{cfg: cfg, authUC: authUC, sessUC: sessUC, log: log}
}

// Create comment
func (h *handlers) Create() echo.HandlerFunc {
	panic("implement me")
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
