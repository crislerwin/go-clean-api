package factories

import (
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

// NewEventHandlerFactory creates a new EventHandler with its dependencies
func NewEventHandlerFactory(db *sqlx.DB) *handler.EventHandler {
	txManager := database.NewTransactionManager(db)
	eventRepo := repository.NewEventRepositorySqlx(db)

	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)
	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	listUserEventsUseCase := usecase.NewListUserEventsUseCase(eventRepo)

	return handler.NewEventHandler(createEventUseCase, listEventsUseCase, listUserEventsUseCase)
}
