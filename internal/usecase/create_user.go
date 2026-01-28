package usecase

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/domain/repository"
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
	userRepo repository.UserRepository
}

func NewSignUpUseCase(userRepo repository.UserRepository) *SignUpUseCase {
	return &SignUpUseCase{
		userRepo: userRepo,
	}
}

func (c *SignUpUseCase) Execute(ctx context.Context, input SignUpInputDTO) (*SignUpOutputDTO, error) {
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	if err := c.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &SignUpOutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
