package middleware

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"strings"
)

func AuthJWTMiddleware(authUC auth.UseCase, config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			bearerHeader := ctx.Request().Header.Get("Authorization")
			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					return errors.Unauthorized
				}

				tokenString := headerParts[1]

				if err := validateJWTToken(tokenString, authUC, ctx, config); err != nil {
					return err
				}
				return next(ctx)
			} else {
				cookie, err := ctx.Cookie("jwt-token")
				if err != nil {
					return err
				}

				if err = validateJWTToken(cookie.Value, authUC, ctx, config); err != nil {
					return err
				}
				return next(ctx)
			}
		}
	}
}

// func AuthRoleMiddleware(roles ...string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user, ok := c.Get("user").(*models.User)
// 			if !ok {
// 				return c.JSON(http.StatusBadRequest, echo.ErrUnauthorized)
// 			}
//
// 			for _, role := range roles {
// 				if user.Role == role {
// 					return next(c)
// 				}
// 			}
// 			return c.JSON(http.StatusBadRequest, echo.ErrUnauthorized)
// 		}
// 	}
// }

// func AuthExtractUserMiddleware(s *usecases.UseCases, config *config.Config) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			bearerHeader := c.Request().Header.Get("Authorization")
// 			if bearerHeader != "" {
// 				headerParts := strings.Split(bearerHeader, " ")
// 				if len(headerParts) != 2 {
// 					return next(c)
// 				}
//
// 				tokenString := headerParts[1]
//
// 				if err := valiadateJWTToken(tokenString, s, c, config); err != nil {
// 					return next(c)
// 				}
// 				// case with user in context
// 				return next(c)
// 			} else {
// 				cookie, err := c.Cookie("jwt-token")
// 				if err != nil {
// 					return next(c)
// 				}
//
// 				if err = valiadateJWTToken(cookie.Value, s, c, config); err != nil {
// 					return next(c)
// 				}
// 				// case with user in context
// 				return next(c)
// 			}
// 		}
// 	}
// }

// Check user from token in context validation, for usage with AuthExtractUserMiddleware
// func AuthChekUserMiddleware() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			user, ok := c.Get("user").(*models.User)
// 			if !ok {
// 				return c.JSON(http.StatusBadRequest, echo.ErrUnauthorized)
// 			}
// 			c.Set("user", user)
// 			return next(c)
// 		}
// 	}
// }

func validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, config *config.Config) error {
	if tokenString == "" {
		return errors.Unauthorized
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(config.Server.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.Unauthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)

		if !ok {
			return errors.Unauthorized
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		ctx := context.WithValue(c.Request().Context(), "user", u)
		c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
