package repository

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/internal/news"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func SetupRedis() news.RedisRepository {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	newsRedisRepo := NewNewsRedisRepo(client)
	return newsRedisRepo
}

func TestNewsRedisRepo_SetNewsCtx(t *testing.T) {
	newsRedisRepo := SetupRedis()

	t.Run("SetNewsCtx", func(t *testing.T) {
		newsUID := uuid.New()
		key := "key"
		n := &models.NewsBase{
			NewsID:  newsUID,
			Title:   "Title",
			Content: "Content",
		}

		err := newsRedisRepo.SetNewsCtx(context.Background(), key, 10, n)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}

func TestNewsRedisRepo_GetNewsByIDCtx(t *testing.T) {
	newsRedisRepo := SetupRedis()

	t.Run("SetNewsCtx", func(t *testing.T) {
		newsUID := uuid.New()
		key := "key"
		n := &models.NewsBase{
			NewsID:  newsUID,
			Title:   "Title",
			Content: "Content",
		}

		err := newsRedisRepo.SetNewsCtx(context.Background(), key, 10, n)
		require.NoError(t, err)
		require.Nil(t, err)

		newsBase, err := newsRedisRepo.GetNewsByIDCtx(context.Background(), key)
		require.NoError(t, err)
		require.Nil(t, err)
		require.NotNil(t, newsBase)
	})
}

func TestNewsRedisRepo_DeleteNewsCtx(t *testing.T) {
	newsRedisRepo := SetupRedis()

	t.Run("SetNewsCtx", func(t *testing.T) {
		key := "key"

		err := newsRedisRepo.DeleteNewsCtx(context.Background(), key)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}
