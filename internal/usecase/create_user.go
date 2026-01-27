package usecase

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

type SignUpInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpOutputDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignUpUseCase struct {
	UserRepository UserRepository
}

func NewSignUpUseCase(userRepository UserRepository) *SignUpUseCase {
	return &SignUpUseCase{
		UserRepository: userRepository,
	}
}

func (c *SignUpUseCase) Execute(ctx context.Context, input SignUpInputDTO) (*SignUpOutputDTO, error) {
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	if err := c.UserRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	return &SignUpOutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
