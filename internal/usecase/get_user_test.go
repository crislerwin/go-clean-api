package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepoMock) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestGetUserUseCase_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(UserRepoMock)
		useCase := NewGetUserUseCase(repo)

		user, _ := entity.NewUser("John Doe", "john@example.com", "password123")
		repo.On("GetByID", mock.Anything, user.ID).Return(user, nil)

		output, err := useCase.Execute(context.TODO(), user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, output.ID)
		assert.Equal(t, user.Name, output.Name)
		assert.Equal(t, user.Email, output.Email)
		assert.Equal(t, user.Role, output.Role)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(UserRepoMock)
		useCase := NewGetUserUseCase(repo)

		repo.On("GetByID", mock.Anything, "invalid-id").Return(nil, errors.New("user not found"))

		output, err := useCase.Execute(context.TODO(), "invalid-id")

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "user not found", err.Error())
		repo.AssertExpectations(t)
	})
}
