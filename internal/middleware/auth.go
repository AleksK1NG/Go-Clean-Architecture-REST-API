package middleware

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/errors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Auth by sessions stored in redis
func AuthSessionMiddleware(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(cfg.Session.Name)
			if err != nil {
				log.Error(
					"AuthSessionMiddleware",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("error", err.Error()),
				)
				return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			}

			sess, err := sessUC.GetSessionByID(c.Request().Context(), cookie.Value)
			if err != nil {
				log.Error(
					"GetSessionByID",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("cookieValue", cookie.Value),
					zap.String("error", err.Error()),
				)
				return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			}

			user, err := authUC.GetByID(c.Request().Context(), sess.UserID)
			if err != nil {
				log.Error(
					"GetByID",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("error", err.Error()),
				)
				return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			}

			c.Set("user", user)
			ctx := context.WithValue(c.Request().Context(), dto.UserCtxKey{}, user)
			c.Request().WithContext(ctx)
			c.SetRequest(c.Request().WithContext(ctx))

			log.Info(
				"AuthSessionMiddleware",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("IP", utils.GetIPAddress(c)),
				zap.String("userID", user.ID.String()),
			)

			return next(c)
		}
	}
}

// JWT way of auth using cookie or Authorization header
func AuthJWTMiddleware(authUC auth.UseCase, config *config.Config, logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			logger.Info("auth middleware", zap.String("bearerHeader", bearerHeader))

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := validateJWTToken(tokenString, authUC, c, config); err != nil {
					logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
				}

				return next(c)
			} else {
				cookie, err := c.Cookie("jwt-token")
				if err != nil {
					logger.Error("middleware cookie", zap.String("cookieJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
				}

				if err = validateJWTToken(cookie.Value, authUC, c, config); err != nil {
					logger.Error("cookie JWT validate", zap.String("cookieJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
				}
				return next(c)
			}
		}
	}
}

// Admin role
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok || *user.Role != "admin" {
			return c.JSON(http.StatusForbidden, errors.NewUnauthorizedError(errors.PermissionDenied))
		}
		return next(c)
	}
}

// Role based auth middleware, using ctx user
func OwnerOrAdminMiddleware(logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user, ok := c.Get("user").(*models.User)
			if !ok {
				logger.Error(
					"OwnerOrAdminMiddleware",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("ERROR", "invalid user ctx"),
				)
				return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			}

			logger.Info(
				"OwnerOrAdminMiddleware",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("userID", user.ID.String()),
			)

			if *user.Role == "admin" {
				return next(c)
			}

			if user.ID.String() != c.Param("user_id") {
				logger.Error(
					"OwnerOrAdminMiddleware",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("userID", user.ID.String()),
					zap.String("ERROR", "ctx userID != param /:user_id"),
				)
				return c.JSON(http.StatusForbidden, errors.NewForbiddenError(errors.Forbidden))
			}

			return next(c)
		}
	}
}

// Role based auth middleware, using ctx user
func RoleBasedAuthMiddleware(roles []string, logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			logger.Info("RoleBasedAuthMiddleware", zap.String("reqID", utils.GetRequestID(c)))

			user, ok := c.Get("user").(*models.User)
			if !ok {
				logger.Error(
					"RoleBasedAuthMiddleware",
					zap.String("reqID", utils.GetRequestID(c)),
					zap.String("ERROR", "invalid user ctx"),
				)
				return c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(errors.Unauthorized))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			logger.Error(
				"not allowed user role",
				zap.String("reqID", utils.GetRequestID(c)),
				zap.String("userID", user.ID.String()),
				zap.String("ERROR", "not allowed role"),
			)

			return c.JSON(http.StatusForbidden, errors.NewForbiddenError(errors.PermissionDenied))
		}
	}
}

func validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, config *config.Config) error {
	if tokenString == "" {
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
		return err
	}

	if !token.Valid {
		return errors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return errors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), dto.UserCtxKey{}, u)
		c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
