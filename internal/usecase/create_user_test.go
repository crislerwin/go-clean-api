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

func (m *ValidationMockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestSignUpUseCase_Execute(t *testing.T) {
	t.Run("should create a new valid user", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewSignUpUseCase(repo)
		input := SignUpInputDTO{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		repo.On("Save", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Name == input.Name && u.Email == input.Email
		})).Return(nil)

		output, err := useCase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
		assert.Equal(t, input.Name, output.Name)
		assert.Equal(t, input.Email, output.Email)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when entity creation fails", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewSignUpUseCase(repo)
		input := SignUpInputDTO{
			Name:     "", // Invalid name
			Email:    "john@example.com",
			Password: "password123",
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		// We can assert generic error message or specific if we want
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repo := new(ValidationMockUserRepository)
		useCase := NewSignUpUseCase(repo)
		input := SignUpInputDTO{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "db error", err.Error())
	})
}
