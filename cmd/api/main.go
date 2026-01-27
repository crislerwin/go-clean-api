package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/http/middleware"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://ticket_user:ticket_pass@localhost:5432/ticket_db?sslmode=disable"
	}

	db, err := sqlx.Connect("pgx", connStr)

	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	txManager := database.NewTransactionManager(db)
	eventRepo := repository.NewEventRepositorySqlx(db)
	userRepo := repository.NewUserRepositorySQLx(db)

	orderRepo := repository.NewOrderRepositorySQLx(db)

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)

	orderHandler := handler.NewOrderHandler(createOrderUseCase)
	eventHandler := handler.NewEventHandler(createEventUseCase)
	userHandler := handler.NewUserHandler(createUserUseCase)
	authHandler := handler.NewAuthHandler(loginUseCase)

	r := gin.Default()

	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			slog.Error("failed to set trusted proxies", "error", err)
		}
	}

	api := r.Group("/api/v1")

	// Public routes
	api.POST("/users", userHandler.CreateUser)
	api.POST("/login", authHandler.Login)

	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/orders", orderHandler.CreateOrder)
		api.POST("/events", middleware.RoleMiddleware("admin"), eventHandler.CreateEvent)
	}

	slog.Info("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

}
