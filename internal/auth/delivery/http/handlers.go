package http

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/dto"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session"
	"github.com/AleksK1NG/api-mc/internal/utils"
	"github.com/AleksK1NG/api-mc/pkg/httpErrors"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// Auth handlers
type handlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	sessUC session.UCSession
	log    *logger.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, sessUC session.UCSession, log *logger.Logger) auth.Handlers {
	return &handlers{cfg, authUC, sessUC, log}
}

// Register new user
func (h *handlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Register user", zap.String("ReqID", utils.GetRequestID(c)))

		user := &models.User{}
		if err := c.Bind(user); err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: createdUser.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		h.log.Info(
			"CreatedUser",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("Session", sess),
			zap.String("ID", createdUser.User.UserID.String()),
		)

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Login user
func (h *handlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Login", zap.String("ReqID", utils.GetRequestID(c)))

		loginDTO := &dto.LoginDTO{}
		if err := c.Bind(&loginDTO); err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		userWithToken, err := h.authUC.Login(ctx, loginDTO)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: userWithToken.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		h.log.Info(
			"Login",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.String("Session", sess),
			zap.String("User ID", userWithToken.User.UserID.String()),
		)

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout user
func (h *handlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		cookie, err := c.Cookie("session-id")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		h.log.Info(
			"Logout user",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.String("cookie", cookie.Value),
		)

		if err := h.sessUC.DeleteByID(ctx, cookie.Value); err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// Update existing user
func (h *handlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Update", zap.String("ReqID", utils.GetRequestID(c)))

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		user := &models.UserUpdate{}
		user.ID = uID

		if err = c.Bind(user); err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		updatedUser, err := h.authUC.Update(ctx, user)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		h.log.Info(
			"Update",
			zap.String("reqID", utils.GetRequestID(c)),
			zap.String("ID", updatedUser.UserID.String()),
		)

		return c.JSON(http.StatusCreated, updatedUser)
	}
}

// Get user by id
func (h *handlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetUserByID", zap.String("ReqID", utils.GetRequestID(c)))

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		user, err := h.authUC.GetByID(ctx, uID)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Delete user handler
func (h *handlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("Delete", zap.String("ReqID", utils.GetRequestID(c)))

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		if err := h.authUC.Delete(ctx, uID); err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// Find users by name
func (h *handlers) FindByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info(
			"FindByName",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.String("name", c.QueryParam("name")),
		)

		if c.QueryParam("name") == "" {
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("name is required"))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		response, err := h.authUC.FindByName(ctx, &dto.FindUserQuery{
			Name: c.QueryParam("name"),
			PQ:   *paginationQuery,
		})
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		h.log.Info(
			"FindByName",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.Int("Found", len(response.Users)),
		)

		return c.JSON(http.StatusOK, response)
	}
}

// Gat all users with pagination page and size query params
func (h *handlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := utils.GetCtxWithReqID(c)
		defer cancel()

		h.log.Info("GetUsers", zap.String("ReqID", utils.GetRequestID(c)))

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		usersList, err := h.authUC.GetUsers(ctx, paginationQuery)
		if err != nil {
			return utils.ErrResponseWithLog(c, h.log, err)
		}

		h.log.Info(
			"GetUsers",
			zap.String("ReqID", utils.GetRequestID(c)),
			zap.Int("Found", len(usersList.Users)),
			zap.String("Query", fmt.Sprintf("%#v", paginationQuery)),
		)

		return c.JSON(http.StatusOK, usersList)
	}
}

// Load current user from ctx with auth middleware
func (h *handlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			return utils.ErrResponseWithLog(c, h.log, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		h.log.Info("GetMe", zap.String(
			"ReqID", utils.GetRequestID(c)),
			zap.String("userId", user.UserID.String()),
		)

		return c.JSON(http.StatusOK, c.Get("user"))
	}
}
