package usecase

import (
	"context"
	"errors"

	"github.com/crislerwin/go-clean-api/internal/domain/repository"

	"github.com/crislerwin/go-clean-api/internal/infra/http/auth"
	"golang.org/x/crypto/bcrypt"
)

type LoginInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutputDTO struct {
	AccessToken string `json:"access_token"`
}

type LoginUseCase struct {
	UserRepository repository.UserRepository
}

func NewLoginUseCase(userRepository repository.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		UserRepository: userRepository,
	}
}

func (c *LoginUseCase) Execute(ctx context.Context, input LoginInputDTO) (*LoginOutputDTO, error) {
	user, err := c.UserRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginOutputDTO{
		AccessToken: token,
	}, nil
}
