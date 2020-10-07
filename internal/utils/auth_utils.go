package utils

import (
	"context"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/labstack/echo"
	"log"
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
func ValidateIsOwner(ctx context.Context, creatorId string) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}
	log.Printf("ValidateIsOwner: %s : %s", user.ID.String(), creatorId)
	if user.ID.String() != creatorId {
		return errors.Forbidden
	}

	return nil
}
