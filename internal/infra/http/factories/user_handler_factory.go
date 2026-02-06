package factories

import (
	"github.com/crislerwin/go-clean-api/internal/infra/database/postgres"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

// NewUserHandlerFactory creates a new UserHandler with its dependencies.
func NewUserHandlerFactory(db *sqlx.DB) *handler.UserHandler {
	client := postgres.NewClient(db)
	userRepo := repository.NewUserRepositorySQLx(client)
	signUpUseCase := usecase.NewSignUpUseCase(userRepo)
	getUserUseCase := usecase.NewGetUserUseCase(userRepo)

	return handler.NewUserHandler(signUpUseCase, getUserUseCase)
}
