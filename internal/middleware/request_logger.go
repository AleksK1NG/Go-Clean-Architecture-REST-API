package middleware

import (
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

// Request logger middleware
func RequestLoggerMiddleware(logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			res := next(ctx)
			//logger.Info("LoggerMiddleware",
			//	zap.String("Method", ctx.Request().Method),
			//	zap.String("URI", ctx.Request().RequestURI),
			//	zap.Int("Status", ctx.Response().Status),
			//	zap.String("RequestID", utils.GetRequestID(ctx)),
			//	zap.String("Time", time.Since(start).String()),
			//)
			//logger.Info("12345")
			//logger.Sync()
			s := time.Since(start).String()
			logger.Info(s)
			//log.Printf("Time: %s", s)
			//logger.Info(fmt.Sprintf("Method: %s, URI: %s, Status: %s, RequestID: %s, Time: %s",
			//	ctx.Request().Method,
			//	ctx.Request().RequestURI,
			//	ctx.Response().Status,
			//	utils.GetRequestID(ctx),
			//	s,
			//))
			//logger.Info(ctx.Request().Method)
			//logger.Info(ctx.Request().RequestURI)
			return res
		}
	}
}
