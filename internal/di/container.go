package di

import (
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	UserHandler  *handler.UserHandler
	AuthHandler  *handler.AuthHandler
	EventHandler *handler.EventHandler
	OrderHandler *handler.OrderHandler
}

func NewContainer(db *sqlx.DB) *Container {
	txManager := database.NewTransactionManager(db)
	eventRepo := repository.NewEventRepositorySqlx(db)
	userRepo := repository.NewUserRepositorySQLx(db)
	orderRepo := repository.NewOrderRepositorySQLx(db)

	// UseCase
	signUpUseCase := usecase.NewSignUpUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)
	updateOrderStatusUseCase := usecase.NewUpdateOrderStatusUseCase(orderRepo)
	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	listUserEventsUseCase := usecase.NewListUserEventsUseCase(eventRepo)
	getUserUseCase := usecase.NewGetUserUseCase(userRepo)

	// Handler
	userHandler := handler.NewUserHandler(signUpUseCase, getUserUseCase)
	authHandler := handler.NewAuthHandler(loginUseCase)
	eventHandler := handler.NewEventHandler(createEventUseCase, listEventsUseCase, listUserEventsUseCase)
	orderHandler := handler.NewOrderHandler(createOrderUseCase, listUserOrdersUseCase, updateOrderStatusUseCase)

	return &Container{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		EventHandler: eventHandler,
		OrderHandler: orderHandler,
	}
}
