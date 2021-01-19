package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/models"
)

func TestCommentsRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewCommentsRepository(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		authorUID := uuid.New()
		newsUID := uuid.New()
		message := "message"

		rows := sqlmock.NewRows([]string{"author_id", "news_id", "message"}).AddRow(authorUID, newsUID, message)

		comment := &models.Comment{
			AuthorID: authorUID,
			NewsID:   newsUID,
			Message:  message,
		}

		mock.ExpectQuery(createComment).WithArgs(comment.AuthorID, &comment.NewsID, comment.Message).WillReturnRows(rows)

		createdComment, err := commRepo.Create(context.Background(), comment)

		require.NoError(t, err)
		require.NotNil(t, createdComment)
		require.Equal(t, createdComment, comment)
	})

	t.Run("Create ERR", func(t *testing.T) {
		newsUID := uuid.New()
		message := "message"
		createErr := errors.New("Create comment error")

		comment := &models.Comment{
			NewsID:  newsUID,
			Message: message,
		}

		mock.ExpectQuery(createComment).WithArgs(comment.AuthorID, &comment.NewsID, comment.Message).WillReturnError(createErr)

		createdComment, err := commRepo.Create(context.Background(), comment)

		require.Nil(t, createdComment)
		require.NotNil(t, err)
	})
}

func TestCommentsRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewCommentsRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		commUID := uuid.New()
		newsUID := uuid.New()
		message := "message"

		rows := sqlmock.NewRows([]string{"comment_id", "news_id", "message"}).AddRow(commUID, newsUID, message)

		comment := &models.Comment{
			CommentID: commUID,
			Message:   message,
		}

		mock.ExpectQuery(updateComment).WithArgs(comment.Message, comment.CommentID).WillReturnRows(rows)

		createdComment, err := commRepo.Update(context.Background(), comment)

		require.NoError(t, err)
		require.NotNil(t, createdComment)
		require.Equal(t, createdComment.Message, comment.Message)
	})

	t.Run("Update ERR", func(t *testing.T) {
		commUID := uuid.New()
		message := "message"
		updateErr := errors.New("Create comment error")

		comment := &models.Comment{
			CommentID: commUID,
			Message:   message,
		}

		mock.ExpectQuery(updateComment).WithArgs(comment.Message, comment.CommentID).WillReturnError(updateErr)

		createdComment, err := commRepo.Update(context.Background(), comment)

		require.NotNil(t, err)
		require.Nil(t, createdComment)
	})
}

func TestCommentsRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewCommentsRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		commUID := uuid.New()
		mock.ExpectExec(deleteComment).WithArgs(commUID).WillReturnResult(sqlmock.NewResult(1, 1))
		err := commRepo.Delete(context.Background(), commUID)

		require.NoError(t, err)
	})

	t.Run("Delete Err", func(t *testing.T) {
		commUID := uuid.New()

		mock.ExpectExec(deleteComment).WithArgs(commUID).WillReturnResult(sqlmock.NewResult(1, 0))

		err := commRepo.Delete(context.Background(), commUID)
		require.NotNil(t, err)
	})
}
