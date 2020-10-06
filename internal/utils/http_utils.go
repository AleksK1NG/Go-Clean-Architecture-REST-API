package utils

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	ctx = context.WithValue(ctx, "ReqID", GetRequestID(c))
	return ctx, cancel
}

// Configure jwt cookie
func ConfigureJWTCookie(cfg *config.Config, jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Cookie.Name,
		Value: jwtToken,
		Path:  "/",
		// Domain: "/",
		// Expires:    time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.Cookie.MaxAge,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HttpOnly,
		SameSite:   0,
	}
}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		return nil, errors.Unauthorized
	}

	return user, nil
}
