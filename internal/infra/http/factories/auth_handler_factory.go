package factories

import (
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

// NewAuthHandlerFactory creates a new AuthHandler with its dependencies
func NewAuthHandlerFactory(db *sqlx.DB) *handler.AuthHandler {
	userRepo := repository.NewUserRepositorySQLx(db)
	loginUseCase := usecase.NewLoginUseCase(userRepo)

	return handler.NewAuthHandler(loginUseCase)
}
