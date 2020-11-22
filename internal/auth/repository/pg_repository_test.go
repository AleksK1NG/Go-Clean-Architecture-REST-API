package repository

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRepo_Register(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register", func(t *testing.T) {
		//uid := uuid.New()
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
		if err != nil {
			fmt.Printf("ERROR: %s \n", err.Error())
		}
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})

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

		fmt.Printf("test user: %s \n", updatedUser.FirstName)
		fmt.Printf("user: %s \n", user.FirstName)
	})

	t.Run("Delete", func(t *testing.T) {
		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authRepo.Delete(context.Background(), uid)
		if err != nil {
			fmt.Printf("test user: %s \n", err.Error())
		}
		require.Nil(t, err)
	})

}
