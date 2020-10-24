package middleware

import (
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

// Request logger middleware
func RequestLoggerMiddleware(logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			res := next(ctx)
			logger.Info("LoggerMiddleware",
				zap.String("Method", ctx.Request().Method),
				zap.String("URI", ctx.Request().RequestURI),
				zap.Int("Status", ctx.Response().Status),
				zap.String("Time", time.Since(start).String()),
			)
			return res
		}
	}
}
