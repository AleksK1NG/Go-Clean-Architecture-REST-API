package utils

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Set Auth Cookie with token
func SetAuthCookieWithToken(c echo.Context, token string, config *config.Config) {
	c.SetCookie(&http.Cookie{
		Name:     config.Cookie.Name,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(1 * 24 * time.Hour),
		MaxAge:   config.Cookie.MaxAge,
		Secure:   config.Cookie.Secure,
		HttpOnly: config.Cookie.HttpOnly,
	})
}

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorId string, logger *logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.UserID.String() != creatorId {
		logger.Error(
			"ValidateIsOwner",
			zap.String("userID", user.UserID.String()),
			zap.String("creatorId", creatorId),
		)
		return errors.Forbidden
	}

	return nil
}
