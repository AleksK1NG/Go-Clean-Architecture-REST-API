package middleware

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httputil"
)

// Debug dump request middleware
func (mw *MiddlewareManager) DebugMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if mw.cfg.Server.Debug {
			dump, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			logger.Info(fmt.Sprintf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump))
		}
		return next(c)
	}
}
