package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
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
			"Alfan", "Almunawar", "123456", "alfan@gmail.com", "admin", &gender)

		user := &models.User{
			FirstName: "Alfan",
			LastName:  "Almunawar",
			Email:     "alfan@gmail.com",
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
			uid, "Alfan", "Almunawar", "alfan@gmail.com")

		testUser := &models.User{
			UserID:    uid,
			FirstName: "Alfan",
			LastName:  "Almunawar",
			Email:     "alfan@gmail.com",
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
			uid, "Alfan", "Almuanwar", "alfan@gmail.com")

		testUser := &models.User{
			UserID:    uid,
			FirstName: "Alfan",
			LastName:  "Almuanwar",
			Email:     "alfan@gmail.com",
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
			uid, "Alfan", "Almunawar", "alfan@gmai.com")

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
		userName := "Alfan"

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uid, "Alfan", "Almunawar", "alfan@gmail.com")

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
