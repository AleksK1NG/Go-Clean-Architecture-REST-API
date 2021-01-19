package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/comments/mock"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

func TestCommentsUC_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo, apiLogger)

	comm := &models.Comment{}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "commentsUC.Create")
	defer span.Finish()

	mockCommRepo.EXPECT().Create(ctx, gomock.Eq(comm)).Return(comm, nil)

	createdComment, err := commUC.Create(context.Background(), comm)
	require.NoError(t, err)
	require.NotNil(t, createdComment)
}

func TestCommentsUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo, apiLogger)

	authorUID := uuid.New()

	comm := &models.Comment{
		CommentID: uuid.New(),
		AuthorID:  authorUID,
	}

	baseComm := &models.CommentBase{
		AuthorID: authorUID,
	}

	user := &models.User{
		UserID: authorUID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "commentsUC.Update")
	defer span.Finish()

	mockCommRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(comm.CommentID)).Return(baseComm, nil)
	mockCommRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(comm)).Return(comm, nil)

	updatedComment, err := commUC.Update(ctx, comm)
	require.NoError(t, err)
	require.NotNil(t, updatedComment)
}

func TestCommentsUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo, apiLogger)

	authorUID := uuid.New()

	comm := &models.Comment{
		CommentID: uuid.New(),
		AuthorID:  authorUID,
	}

	baseComm := &models.CommentBase{
		AuthorID: authorUID,
	}

	user := &models.User{
		UserID: authorUID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "commentsUC.Delete")
	defer span.Finish()

	mockCommRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(comm.CommentID)).Return(baseComm, nil)
	mockCommRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(comm.CommentID)).Return(nil)

	err := commUC.Delete(ctx, comm.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestCommentsUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo, apiLogger)

	comm := &models.Comment{
		CommentID: uuid.New(),
	}

	baseComm := &models.CommentBase{}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "commentsUC.GetByID")
	defer span.Finish()

	mockCommRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(comm.CommentID)).Return(baseComm, nil)

	commentBase, err := commUC.GetByID(ctx, comm.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, commentBase)
}

func TestCommentsUC_GetAllByNewsID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo, apiLogger)

	newsUID := uuid.New()

	comm := &models.Comment{
		CommentID: uuid.New(),
		NewsID:    newsUID,
	}

	commentsList := &models.CommentsList{}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "commentsUC.GetAllByNewsID")
	defer span.Finish()

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	mockCommRepo.EXPECT().GetAllByNewsID(ctxWithTrace, gomock.Eq(comm.NewsID), query).Return(commentsList, nil)

	commList, err := commUC.GetAllByNewsID(ctx, comm.NewsID, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, commList)
}
