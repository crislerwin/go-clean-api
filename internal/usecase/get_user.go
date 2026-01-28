package usecase

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/repository"
)

type GetUserOutputDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type GetUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewGetUserUseCase(userRepository repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		UserRepository: userRepository,
	}
}

func (c *GetUserUseCase) Execute(ctx context.Context, userID string) (*GetUserOutputDTO, error) {
	user, err := c.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
