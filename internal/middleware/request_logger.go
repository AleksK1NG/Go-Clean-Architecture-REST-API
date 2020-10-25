package middleware

import (
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

// Request logger middleware
func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			res := next(ctx)

			s := time.Since(start).String()
			logger.Infof("TimeSince: %s", s)
			return res
		}
	}
}
