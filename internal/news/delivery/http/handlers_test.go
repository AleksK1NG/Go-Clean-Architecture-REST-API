package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news/mock"
	"github.com/AleksK1NG/api-mc/pkg/converter"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

func TestNewsHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsUC := mock.NewMockUseCase(ctrl)
	newsHandlers := NewNewsHandlers(nil, mockNewsUC, apiLogger)

	handlerFunc := newsHandlers.Create()

	userID := uuid.New()

	news := &models.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content some text content",
	}

	buf, err := converter.AnyToBytesBuffer(news)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/news/create", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &models.User{
		UserID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctxWithReqID := utils.GetRequestCtx(ctx)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctxWithReqID, "newsHandlers.Create")
	defer span.Finish()

	mockNews := &models.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsUC.EXPECT().Create(ctxWithTrace, gomock.Any()).Return(mockNews, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestNewsHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsUC := mock.NewMockUseCase(ctrl)
	newsHandlers := NewNewsHandlers(nil, mockNewsUC, apiLogger)

	handlerFunc := newsHandlers.Update()

	userID := uuid.New()

	news := &models.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	buf, err := converter.AnyToBytesBuffer(news)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/news/f8a3cc26-fbe1-4713-98be-a2927201356e", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &models.User{
		UserID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctx.SetParamNames("news_id")
	ctx.SetParamValues("f8a3cc26-fbe1-4713-98be-a2927201356e")
	ctxWithReqID := utils.GetRequestCtx(ctx)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctxWithReqID, "newsHandlers.Update")
	defer span.Finish()

	mockNews := &models.News{
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsUC.EXPECT().Update(ctxWithTrace, gomock.Any()).Return(mockNews, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestNewsHandlers_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsUC := mock.NewMockUseCase(ctrl)
	newsHandlers := NewNewsHandlers(nil, mockNewsUC, apiLogger)

	handlerFunc := newsHandlers.GetByID()

	userID := uuid.New()
	newsID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/news/f8a3cc26-fbe1-4713-98be-a2927201356e", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	u := &models.User{
		UserID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	req = req.WithContext(ctxWithValue)
	e := echo.New()
	ctx := e.NewContext(req, res)
	ctx.SetParamNames("news_id")
	ctx.SetParamValues(newsID.String())
	ctxWithReqID := utils.GetRequestCtx(ctx)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctxWithReqID, "newsHandlers.GetByID")
	defer span.Finish()

	mockNews := &models.NewsBase{
		NewsID:   newsID,
		AuthorID: userID,
		Title:    "TestNewsHandlers_Create title",
		Content:  "TestNewsHandlers_Create title content asdasdsadsadadsad",
	}

	mockNewsUC.EXPECT().GetNewsByID(ctxWithTrace, newsID).Return(mockNews, nil)

	err := handlerFunc(ctx)
	require.NoError(t, err)
}
