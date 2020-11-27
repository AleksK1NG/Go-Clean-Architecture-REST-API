package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/session/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSessionUC_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessRepo := mock.NewMockSessRepository(ctrl)
	sessUC := NewSessionUseCase(mockSessRepo, nil)

	ctx := context.Background()
	sess := &models.Session{}
	sid := "session id"

	mockSessRepo.EXPECT().CreateSession(ctx, gomock.Eq(sess), 10).Return(sid, nil)

	createdSess, err := sessUC.CreateSession(ctx, sess, 10)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotEqual(t, createdSess, "")
}

func TestSessionUC_GetSessionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessRepo := mock.NewMockSessRepository(ctrl)
	sessUC := NewSessionUseCase(mockSessRepo, nil)

	ctx := context.Background()
	sess := &models.Session{}
	sid := "session id"

	mockSessRepo.EXPECT().GetSessionByID(ctx, gomock.Eq(sid)).Return(sess, nil)

	session, err := sessUC.GetSessionByID(ctx, sid)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, session)
}

func TestSessionUC_DeleteByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessRepo := mock.NewMockSessRepository(ctrl)
	sessUC := NewSessionUseCase(mockSessRepo, nil)

	ctx := context.Background()
	sid := "session id"

	mockSessRepo.EXPECT().DeleteByID(ctx, gomock.Eq(sid)).Return(nil)

	err := sessUC.DeleteByID(ctx, sid)
	require.NoError(t, err)
	require.Nil(t, err)
}
