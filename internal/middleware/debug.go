package middleware

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
)

// Debug dump request middleware
func DebugMiddleware(isDebug bool, logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isDebug {
				dump, err := httputil.DumpRequest(c.Request(), true)
				if err != nil {
					return c.NoContent(http.StatusInternalServerError)
				}
				logger.Info(
					"DebugMiddleware",
					zap.String("DEBUG", fmt.Sprintf("\nRequest dump begin :--------------\n\n%s\n\nRequest dump end :--------------", dump)),
				)
			}
			return next(c)
		}
	}
}
