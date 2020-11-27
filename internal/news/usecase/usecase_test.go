package usecase

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news/mock"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewsUC_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, nil, apiLogger)

	userUID := uuid.New()

	news := &models.News{
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	user := &models.User{
		UserID: userUID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsRepo.EXPECT().Create(ctx, gomock.Eq(news)).Return(news, nil)

	createdNews, err := newsUC.Create(ctx, news)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, createdNews)
}

func TestNewsUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, mockRedisRepo, apiLogger)

	userUID := uuid.New()
	newsUID := uuid.New()
	news := &models.News{
		NewsID:   newsUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 20 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	newsBase := &models.NewsBase{
		NewsID:   newsUID,
		AuthorID: userUID,
		Title:    "Title long text string greater then 55555 characters",
		Content:  "Content long text string greater then 20 characters",
	}

	user := &models.User{
		UserID: userUID,
	}

	cacheKey := fmt.Sprintf("%s: %s", basePrefix, news.NewsID)

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsRepo.EXPECT().GetNewsByID(ctx, gomock.Eq(news.NewsID)).Return(newsBase, nil)
	mockNewsRepo.EXPECT().Update(ctx, gomock.Eq(news)).Return(news, nil)
	mockRedisRepo.EXPECT().DeleteNewsCtx(ctx, gomock.Eq(cacheKey)).Return(nil)

	updatedNews, err := newsUC.Update(ctx, news)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, updatedNews)
}

func TestNewsUC_GetNewsByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, mockRedisRepo, apiLogger)

	newsUID := uuid.New()
	newsBase := &models.NewsBase{
		NewsID: newsUID,
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf("%s: %s", basePrefix, newsUID)

	mockRedisRepo.EXPECT().GetNewsByIDCtx(ctx, gomock.Eq(cacheKey)).Return(nil, nil)
	mockNewsRepo.EXPECT().GetNewsByID(ctx, gomock.Eq(newsUID)).Return(newsBase, nil)
	mockRedisRepo.EXPECT().SetNewsCtx(ctx, cacheKey, cacheDuration, newsBase).Return(nil)

	newsByID, err := newsUC.GetNewsByID(ctx, newsBase.NewsID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, newsByID)
}

func TestNewsUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, mockRedisRepo, apiLogger)

	newsUID := uuid.New()
	userUID := uuid.New()
	newsBase := &models.NewsBase{
		NewsID:   newsUID,
		AuthorID: userUID,
	}
	cacheKey := fmt.Sprintf("%s: %s", basePrefix, newsUID)

	user := &models.User{
		UserID: userUID,
	}

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, user)

	mockNewsRepo.EXPECT().GetNewsByID(ctx, gomock.Eq(newsBase.NewsID)).Return(newsBase, nil)
	mockNewsRepo.EXPECT().Delete(ctx, gomock.Eq(newsUID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteNewsCtx(ctx, gomock.Eq(cacheKey)).Return(nil)

	err := newsUC.Delete(ctx, newsBase.NewsID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestNewsUC_GetNews(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, mockRedisRepo, apiLogger)

	ctx := context.Background()
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &models.NewsList{}

	mockNewsRepo.EXPECT().GetNews(ctx, query).Return(newsList, nil)

	news, err := newsUC.GetNews(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, news)
}

func TestNewsUC_SearchByTitle(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockNewsRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	newsUC := NewNewsUseCase(nil, mockNewsRepo, mockRedisRepo, apiLogger)

	ctx := context.Background()
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &models.NewsList{}
	title := "title"

	mockNewsRepo.EXPECT().SearchByTitle(ctx, title, query).Return(newsList, nil)

	news, err := newsUC.SearchByTitle(ctx, title, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, news)
}
