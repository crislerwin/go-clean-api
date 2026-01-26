package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ValidationMockUserRepository struct {
	mock.Mock
}

func (m *ValidationMockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *ValidationMockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUserUseCase_Execute(t *testing.T) {
	t.Run("should create a valid user", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		useCase := NewCreateUserUseCase(repo)
		input := CreateUserInputDTO{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, input.Name, output.Name)
		assert.Equal(t, input.Email, output.Email)
		assert.NotEmpty(t, output.ID)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when repo fails", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("db error"))

		useCase := NewCreateUserUseCase(repo)
		input := CreateUserInputDTO{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		repo.AssertExpectations(t)
	})
}
