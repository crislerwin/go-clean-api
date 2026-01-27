package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/crislerwin/go-clean-api/docs" // Swagger docs
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/http/middleware"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go Clean API
// @version         1.0
// @description     A robust Ticketing System API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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

	// UseCase
	signUpUseCase := usecase.NewSignUpUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)

	// Handler
	userHandler := handler.NewUserHandler(signUpUseCase, usecase.NewGetUserUseCase(userRepo))
	authHandler := handler.NewAuthHandler(loginUseCase)
	eventHandler := handler.NewEventHandler(createEventUseCase)
	orderHandler := handler.NewOrderHandler(createOrderUseCase, listUserOrdersUseCase)

	// Router
	r := gin.Default()

	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			slog.Error("failed to set trusted proxies", "error", err)
		}
	}

	api := r.Group("/api/v1")
	// Swagger Redirect
	api.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	api.POST("/signup", userHandler.SignUp)
	api.POST("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/me", userHandler.Me)
	protected.POST("/orders", orderHandler.CreateOrder)
	protected.GET("/orders", orderHandler.ListMyOrders)
	protected.POST("/events", middleware.RoleMiddleware("admin"), eventHandler.CreateEvent)

	slog.Info("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

}
