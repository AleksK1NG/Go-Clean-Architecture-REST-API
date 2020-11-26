package http

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/comments/mock"
	"github.com/AleksK1NG/api-mc/internal/comments/usecase"
	"github.com/AleksK1NG/api-mc/internal/models"
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

func TestCommentsHandlers_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommUC := mock.NewMockUseCase(ctrl)
	commUC := usecase.NewCommentsUseCase(nil, mockCommUC)

	commHandlers := NewCommentsHandlers(nil, commUC)
	handlerFunc := commHandlers.Create()

	userID := uuid.New()
	newsUID := uuid.New()
	comment := &models.Comment{
		AuthorID: userID,
		Message:  "message Key: 'Comment.Message' Error:Field validation for 'Message' failed on the 'gte' tag",
		NewsID:   newsUID,
	}

	buf, err := converter.AnyToBytesBuffer(comment)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", strings.NewReader(buf.String()))
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

	mockComm := &models.Comment{
		AuthorID: userID,
		NewsID:   comment.NewsID,
		Message:  "message",
	}

	fmt.Printf("COMMENT: %#v\n", comment)
	fmt.Printf("MOCK COMMENT: %#v\n", mockComm)

	mockCommUC.EXPECT().Create(ctxWithReqID, gomock.Any()).Return(mockComm, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestCommentsHandlers_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommUC := mock.NewMockUseCase(ctrl)
	commUC := usecase.NewCommentsUseCase(nil, mockCommUC)

	commHandlers := NewCommentsHandlers(nil, commUC)
	handlerFunc := commHandlers.GetByID()

	r := httptest.NewRequest(http.MethodGet, "/api/v1/comments/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("comment_id")
	c.SetParamValues("5c9a9d67-ad38-499c-9858-086bfdeaf7d2")
	ctx := utils.GetRequestCtx(c)

	comm := &models.CommentBase{}

	mockCommUC.EXPECT().GetByID(ctx, gomock.Any()).Return(comm, nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}

func TestCommentsHandlers_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommUC := mock.NewMockUseCase(ctrl)
	commUC := usecase.NewCommentsUseCase(nil, mockCommUC)

	commHandlers := NewCommentsHandlers(nil, commUC)
	handlerFunc := commHandlers.Delete()

	userID := uuid.New()
	commID := uuid.New()
	comm := &models.CommentBase{
		CommentID: commID,
		AuthorID:  userID,
	}

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	u := &models.User{
		UserID: userID,
	}
	ctxWithValue := context.WithValue(context.Background(), utils.UserCtxKey{}, u)
	r = r.WithContext(ctxWithValue)
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("comment_id")
	c.SetParamValues(commID.String())
	ctx := utils.GetRequestCtx(c)

	mockCommUC.EXPECT().GetByID(ctx, commID).Return(comm, nil)
	mockCommUC.EXPECT().Delete(ctx, commID).Return(nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}
