package factories

import (
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

// NewOrderHandlerFactory creates a new OrderHandler with its dependencies.
func NewOrderHandlerFactory(db *sqlx.DB) *handler.OrderHandler {
	txManager := database.NewTransactionManager(db)
	orderRepo := repository.NewOrderRepositorySQLx(db)
	eventRepo := repository.NewEventRepositorySqlx(db)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)
	updateOrderStatusUseCase := usecase.NewUpdateOrderStatusUseCase(orderRepo)

	return handler.NewOrderHandler(createOrderUseCase, listUserOrdersUseCase, updateOrderStatusUseCase)
}
