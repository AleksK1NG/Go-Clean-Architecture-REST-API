package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/models"
)

func TestNewsRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsRepo := NewNewsRepository(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		authorUID := uuid.New()
		title := "title"
		content := "content"

		rows := sqlmock.NewRows([]string{"author_id", "title", "content"}).AddRow(authorUID, title, content)

		news := &models.News{
			AuthorID: authorUID,
			Title:    title,
			Content:  content,
		}

		mock.ExpectQuery(createNews).WithArgs(news.AuthorID, news.Title, news.Content, news.Category).WillReturnRows(rows)

		createdNews, err := newsRepo.Create(context.Background(), news)

		require.NoError(t, err)
		require.NotNil(t, createdNews)
		require.Equal(t, news.Title, createdNews.Title)
	})
}

func TestNewsRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsRepo := NewNewsRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		newsUID := uuid.New()
		title := "title"
		content := "content"

		rows := sqlmock.NewRows([]string{"news_id", "title", "content"}).AddRow(newsUID, title, content)

		news := &models.News{
			NewsID:  newsUID,
			Title:   title,
			Content: content,
		}

		mock.ExpectQuery(updateNews).WithArgs(news.Title,
			news.Content,
			news.ImageURL,
			news.Category,
			news.NewsID,
		).WillReturnRows(rows)

		updatedNews, err := newsRepo.Update(context.Background(), news)

		require.NoError(t, err)
		require.NotNil(t, updateNews)
		require.Equal(t, updatedNews, news)
	})
}

func TestNewsRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	newsRepo := NewNewsRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		newsUID := uuid.New()
		mock.ExpectExec(deleteNews).WithArgs(newsUID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := newsRepo.Delete(context.Background(), newsUID)

		require.NoError(t, err)
	})
}
