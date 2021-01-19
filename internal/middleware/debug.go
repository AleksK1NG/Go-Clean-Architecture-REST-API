package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

// Debug dump request middleware
func (mw *MiddlewareManager) DebugMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if mw.cfg.Server.Debug {
			dump, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			mw.logger.Info(fmt.Sprintf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump))
		}
		return next(c)
	}
}
