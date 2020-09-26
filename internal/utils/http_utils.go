package utils

import (
	"context"
	"github.com/labstack/echo"
	"time"
)

func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	ctx = context.WithValue(ctx, "ReqID", GetRequestID(c))
	return ctx, cancel
}
