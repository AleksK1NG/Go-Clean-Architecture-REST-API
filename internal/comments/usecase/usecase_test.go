package usecase

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/comments/mock"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommentsUC_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo)

	comm := &models.Comment{}

	mockCommRepo.EXPECT().Create(context.Background(), gomock.Eq(comm)).Return(comm, nil)

	createdComment, err := commUC.Create(context.Background(), comm)
	require.NoError(t, err)
	require.NotNil(t, createdComment)
}

func TestCommentsUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo)

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

	mockCommRepo.EXPECT().GetByID(ctx, gomock.Eq(comm.CommentID)).Return(baseComm, nil)
	mockCommRepo.EXPECT().Update(ctx, gomock.Eq(comm)).Return(comm, nil)

	updatedComment, err := commUC.Update(ctx, comm)
	require.NoError(t, err)
	require.NotNil(t, updatedComment)
}

func TestCommentsUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo)

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

	mockCommRepo.EXPECT().GetByID(ctx, gomock.Eq(comm.CommentID)).Return(baseComm, nil)
	mockCommRepo.EXPECT().Delete(ctx, gomock.Eq(comm.CommentID)).Return(nil)

	err := commUC.Delete(ctx, comm.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestCommentsUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo)

	comm := &models.Comment{
		CommentID: uuid.New(),
	}

	baseComm := &models.CommentBase{}

	ctx := context.Background()

	mockCommRepo.EXPECT().GetByID(ctx, gomock.Eq(comm.CommentID)).Return(baseComm, nil)

	commentBase, err := commUC.GetByID(ctx, comm.CommentID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, commentBase)
}

func TestCommentsUC_GetAllByNewsID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommRepo := mock.NewMockRepository(ctrl)
	commUC := NewCommentsUseCase(nil, mockCommRepo)

	newsUID := uuid.New()

	comm := &models.Comment{
		CommentID: uuid.New(),
		NewsID:    newsUID,
	}

	commentsList := &models.CommentsList{}

	ctx := context.Background()
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	mockCommRepo.EXPECT().GetAllByNewsID(ctx, gomock.Eq(comm.NewsID), query).Return(commentsList, nil)

	commList, err := commUC.GetAllByNewsID(ctx, comm.NewsID, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, commList)
}
