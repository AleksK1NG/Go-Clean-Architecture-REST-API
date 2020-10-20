package middleware

import (
	"github.com/AleksK1NG/api-mc/pkg/sanitize"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

// Sanitize and read request body to ctx for next use in easy json
func Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()
		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("body", sanBody)
		return next(ctx)
	}
}
