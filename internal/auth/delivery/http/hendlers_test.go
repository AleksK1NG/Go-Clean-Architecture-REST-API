package http

import (
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth/mock"
	"github.com/AleksK1NG/api-mc/internal/models"
	mockSess "github.com/AleksK1NG/api-mc/internal/session/mock"
	"github.com/AleksK1NG/api-mc/pkg/converter"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandlers_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)
	mockSessUC := mockSess.NewMockUCSession(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Expire: 10,
		},
	}
	authHandlers := NewAuthHandlers(cfg, mockAuthUC, mockSessUC)

	gender := "male"
	user := &models.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "email@gmail.com",
		Password:  "123456",
		Gender:    &gender,
	}

	buf, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)

	handlerFunc := authHandlers.Register()

	userUID := uuid.New()
	userWithToken := &models.UserWithToken{
		User: &models.User{
			UserID: userUID,
		},
	}
	sess := &models.Session{
		UserID: userUID,
	}
	session := "session"

	mockAuthUC.EXPECT().Register(ctx, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessUC.EXPECT().CreateSession(ctx, gomock.Eq(sess), 10).Return(session, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthHandlers_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)
	mockSessUC := mockSess.NewMockUCSession(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Expire: 10,
		},
	}

	authHandlers := NewAuthHandlers(cfg, mockAuthUC, mockSessUC)

	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	login := &Login{
		Email:    "email@mail.com",
		Password: "123456",
	}

	user := &models.User{
		Email:    login.Email,
		Password: login.Password,
	}

	buf, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)

	handlerFunc := authHandlers.Login()

	userUID := uuid.New()
	userWithToken := &models.UserWithToken{
		User: &models.User{
			UserID: userUID,
		},
	}
	sess := &models.Session{
		UserID: userUID,
	}
	session := "session"

	mockAuthUC.EXPECT().Login(ctx, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessUC.EXPECT().CreateSession(ctx, gomock.Eq(sess), 10).Return(session, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthHandlers_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockUseCase(ctrl)
	mockSessUC := mockSess.NewMockUCSession(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Expire: 10,
			Name:   "session-id",
		},
	}
	sessionKey := "session-id"
	cookieValue := "cookieValue"

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{Name: sessionKey, Value: cookieValue})

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	ctx := utils.GetRequestCtx(c)

	authHandlers := NewAuthHandlers(cfg, mockAuthUC, mockSessUC)
	logout := authHandlers.Logout()

	cookie, err := req.Cookie(sessionKey)
	require.NoError(t, err)
	require.NotNil(t, cookie)
	require.NotEqual(t, cookie.Value, "")
	require.Equal(t, cookie.Value, cookieValue)

	mockSessUC.EXPECT().DeleteByID(ctx, gomock.Eq(cookie.Value)).Return(nil)

	err = logout(c)
	require.NoError(t, err)
	require.Nil(t, err)
}
