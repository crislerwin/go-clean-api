package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositorySQLx_Save(t *testing.T) {
	t.Run("should save a user successfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		client := postgres.NewClient(sqlxDB)
		repo := NewUserRepositorySQLx(client)

		user, _ := entity.NewUser("John Doe", "john@example.com", "password123")

		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Name, user.Email, user.Password, user.Role).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Save(context.Background(), user)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		client := postgres.NewClient(sqlxDB)
		repo := NewUserRepositorySQLx(client)

		user, _ := entity.NewUser("John Doe", "john@example.com", "password123")

		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Name, user.Email, user.Password, user.Role).
			WillReturnError(errors.New("db error"))

		err = repo.Save(context.Background(), user)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepositorySQLx_GetByEmail(t *testing.T) {
	t.Run("should find user by email", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		client := postgres.NewClient(sqlxDB)
		repo := NewUserRepositorySQLx(client)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
			AddRow("userid-123", "John Doe", "john@example.com", "hashedpassword", "user")

		mock.ExpectQuery("SELECT id, name, email, password, role FROM users WHERE email = ?").
			WithArgs("john@example.com").
			WillReturnRows(rows)

		user, err := repo.GetByEmail(context.Background(), "john@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "user", user.Role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when db fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		client := postgres.NewClient(sqlxDB)
		repo := NewUserRepositorySQLx(client)

		mock.ExpectQuery("SELECT id, name, email, password, role FROM users WHERE email = ?").
			WithArgs("john@example.com").
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByEmail(context.Background(), "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
