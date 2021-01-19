package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
)

func TestAuthRepo_Register(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register", func(t *testing.T) {
		gender := "male"
		role := "admin"

		rows := sqlmock.NewRows([]string{"first_name", "last_name", "password", "email", "role", "gender"}).AddRow(
			"Alex", "Bryksin", "123456", "alex@gmail.com", "admin", &gender)

		user := &models.User{
			FirstName: "Alex",
			LastName:  "Bryksin",
			Email:     "alex@gmail.com",
			Password:  "123456",
			Role:      &role,
			Gender:    &gender,
		}

		mock.ExpectQuery(createUserQuery).WithArgs(&user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City,
			&user.Gender, &user.Postcode, &user.Birthday).WillReturnRows(rows)

		createdUser, err := authRepo.Register(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func TestAuthRepo_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		uid := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uid, "Alex", "Bryksin", "alex@mail.ru")

		testUser := &models.User{
			UserID:    uid,
			FirstName: "Alex",
			LastName:  "Bryksin",
			Email:     "alex@mail.ru",
		}

		mock.ExpectQuery(getUserQuery).
			WithArgs(uid).
			WillReturnRows(rows)

		user, err := authRepo.GetByID(context.Background(), uid)
		require.NoError(t, err)
		require.Equal(t, user.FirstName, testUser.FirstName)
		fmt.Printf("test user: %s \n", testUser.FirstName)
		fmt.Printf("user: %s \n", user.FirstName)
	})
}

func TestAuthRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {

		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authRepo.Delete(context.Background(), uid)
		require.Nil(t, err)
	})

	t.Run("Delete No rows", func(t *testing.T) {

		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 0))

		err := authRepo.Delete(context.Background(), uid)

		require.NotNil(t, err)
	})
}

func TestAuthRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		gender := "male"
		role := "admin"

		rows := sqlmock.NewRows([]string{"first_name", "last_name", "password", "email", "role", "gender"}).AddRow(
			"Alex", "Bryksin", "123456", "alex@gmail.com", "admin", &gender)

		user := &models.User{
			FirstName: "Alex",
			LastName:  "Bryksin",
			Email:     "alex@gmail.com",
			Password:  "123456",
			Role:      &role,
			Gender:    &gender,
		}

		mock.ExpectQuery(updateUserQuery).WithArgs(&user.FirstName, &user.LastName, &user.Email,
			&user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City, &user.Gender,
			&user.Postcode, &user.Birthday, &user.UserID).WillReturnRows(rows)

		updatedUser, err := authRepo.Update(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, user, updatedUser)
	})
}

func TestAuthRepo_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uid := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uid, "Alex", "Bryksin", "alex@mail.ru")

		testUser := &models.User{
			UserID:    uid,
			FirstName: "Alex",
			LastName:  "Bryksin",
			Email:     "alex@mail.ru",
		}

		mock.ExpectQuery(findUserByEmail).WithArgs(testUser.Email).WillReturnRows(rows)

		foundUser, err := authRepo.FindByEmail(context.Background(), testUser)

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		require.Equal(t, foundUser.FirstName, testUser.FirstName)
	})
}

func TestAuthRepo_GetUsers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uid := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uid, "Alex", "Bryksin", "alex@mail.ru")

		mock.ExpectQuery(getTotal).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getUsers).WithArgs("", 0, 10).WillReturnRows(rows)

		users, err := authRepo.GetUsers(context.Background(), &utils.PaginationQuery{
			Size:    10,
			Page:    1,
			OrderBy: "",
		})
		require.NoError(t, err)
		require.NotNil(t, users)
	})

}

func TestAuthRepo_FindByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByName", func(t *testing.T) {
		uid := uuid.New()
		userName := "Alex"

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uid, "Alex", "Bryksin", "alex@mail.ru")

		mock.ExpectQuery(getTotalCount).WillReturnRows(totalCountRows)
		mock.ExpectQuery(findUsers).WithArgs("", 0, 10).WillReturnRows(rows)

		usersList, err := authRepo.FindByName(context.Background(), userName, &utils.PaginationQuery{
			Size:    10,
			Page:    1,
			OrderBy: "",
		})

		require.NoError(t, err)
		require.NotNil(t, usersList)
	})
}
