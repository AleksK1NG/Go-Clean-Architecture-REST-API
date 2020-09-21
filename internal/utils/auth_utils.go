package utils

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/labstack/echo"
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
