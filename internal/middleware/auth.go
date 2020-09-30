package middleware

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/AleksK1NG/api-mc/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"log"
	"strings"
)

func AuthJWTMiddleware(authUC auth.UseCase, config *config.Config, logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			logger.Info("auth middleware", zap.String("bearerHeader", bearerHeader))

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(errors.ErrorResponse(errors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := validateJWTToken(tokenString, authUC, c, config); err != nil {
					logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(errors.ErrorResponse(err))
				}

				return next(c)
			} else {
				cookie, err := c.Cookie("jwt-token")
				if err != nil {
					logger.Error("middleware cookie", zap.String("cookieJWT", err.Error()))
					return c.JSON(errors.ErrorResponse(err))
				}

				if err = validateJWTToken(cookie.Value, authUC, c, config); err != nil {
					logger.Error("cookie JWT validate", zap.String("cookieJWT", err.Error()))
					return c.JSON(errors.ErrorResponse(err))
				}
				return next(c)
			}
		}
	}
}

func validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, config *config.Config) error {
	if tokenString == "" {
		log.Println(6666666)
		return errors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(config.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		log.Println(1111111)
		return err
	}

	if !token.Valid {
		return errors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			log.Println(5555555)
			return errors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Println(2222222222)
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			log.Println(3333333333)
			return err
		}

		ctx := context.WithValue(c.Request().Context(), "user", u)
		c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
