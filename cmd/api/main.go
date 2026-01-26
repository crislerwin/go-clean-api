package main

import (
	"log"
	"os"

	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/http/middleware"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://ticket_user:ticket_pass@localhost:5432/ticket_db?sslmode=disable"
	}

	db, err := sqlx.Connect("pgx", connStr)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	txManager := database.NewTransactionManager(db)
	eventRepo := repository.NewEventRepositorySqlx(db)

	orderRepo := repository.NewOrderRepositorySQLx(db)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)

	orderHandler := handler.NewOrderHandler(createOrderUseCase)
	eventHandler := handler.NewEventHandler(createEventUseCase)

	r := gin.Default()

	api := r.Group("/api/v1")

	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/orders", orderHandler.CreateOrder)
		api.POST("/events", eventHandler.CreateEvent)
	}

	log.Println("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
