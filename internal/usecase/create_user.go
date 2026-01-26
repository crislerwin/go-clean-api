package usecase

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

type CreateUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutputDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserUseCase struct {
	UserRepository UserRepository
}

func NewCreateUserUseCase(userRepository UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (c *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInputDTO) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	if err := c.UserRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserOutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
