package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news/mock"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/AleksK1NG/api-mc/pkg/utils"
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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.Create")
	defer span.Finish()

	mockNewsRepo.EXPECT().Create(ctxWithTrace, gomock.Eq(news)).Return(news, nil)

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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.Update")
	defer span.Finish()

	mockNewsRepo.EXPECT().GetNewsByID(ctxWithTrace, gomock.Eq(news.NewsID)).Return(newsBase, nil)
	mockNewsRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(news)).Return(news, nil)
	mockRedisRepo.EXPECT().DeleteNewsCtx(ctxWithTrace, gomock.Eq(cacheKey)).Return(nil)

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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.GetNewsByID")
	defer span.Finish()
	cacheKey := fmt.Sprintf("%s: %s", basePrefix, newsUID)

	mockRedisRepo.EXPECT().GetNewsByIDCtx(ctxWithTrace, gomock.Eq(cacheKey)).Return(nil, nil)
	mockNewsRepo.EXPECT().GetNewsByID(ctxWithTrace, gomock.Eq(newsUID)).Return(newsBase, nil)
	mockRedisRepo.EXPECT().SetNewsCtx(ctxWithTrace, cacheKey, cacheDuration, newsBase).Return(nil)

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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.Delete")
	defer span.Finish()

	mockNewsRepo.EXPECT().GetNewsByID(ctxWithTrace, gomock.Eq(newsBase.NewsID)).Return(newsBase, nil)
	mockNewsRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(newsUID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteNewsCtx(ctxWithTrace, gomock.Eq(cacheKey)).Return(nil)

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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.GetNews")
	defer span.Finish()

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &models.NewsList{}

	mockNewsRepo.EXPECT().GetNews(ctxWithTrace, query).Return(newsList, nil)

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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "newsUC.SearchByTitle")
	defer span.Finish()
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	newsList := &models.NewsList{}
	title := "title"

	mockNewsRepo.EXPECT().SearchByTitle(ctxWithTrace, title, query).Return(newsList, nil)

	news, err := newsUC.SearchByTitle(ctx, title, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, news)
}
