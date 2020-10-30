package middleware

import (
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

// Request logger middleware
func (mw *MiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()

		logger.Infof("Method: %s, URI: %s, Status: %v, Size: %v, Time: %s", req.Method, req.URL, status, size, s)
		return err
	}
}
