package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase_Execute(t *testing.T) {
	t.Run("should login successfully", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewLoginUseCase(repo)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			ID:       "user-123",
			Email:    "john@example.com",
			Password: string(hashedPassword),
			Role:     "user",
		}

		repo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)

		input := LoginInputDTO{
			Email:    "john@example.com",
			Password: "password123",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.AccessToken)
		repo.AssertExpectations(t)
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewLoginUseCase(repo)

		repo.On("GetByEmail", mock.Anything, "invalid@example.com").Return((*entity.User)(nil), errors.New("user not found"))

		input := LoginInputDTO{
			Email:    "invalid@example.com",
			Password: "password123",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewLoginUseCase(repo)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			ID:       "user-123",
			Email:    "john@example.com",
			Password: string(hashedPassword),
		}

		repo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)

		input := LoginInputDTO{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}
