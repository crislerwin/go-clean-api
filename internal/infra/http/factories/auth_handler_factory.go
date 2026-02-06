package factories

import (
	"github.com/crislerwin/go-clean-api/internal/infra/database/postgres"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

// NewAuthHandlerFactory creates a new AuthHandler with its dependencies.
func NewAuthHandlerFactory(db *sqlx.DB) *handler.AuthHandler {
	client := postgres.NewClient(db)
	userRepo := repository.NewUserRepositorySQLx(client)
	loginUseCase := usecase.NewLoginUseCase(userRepo)

	return handler.NewAuthHandler(loginUseCase)
}
