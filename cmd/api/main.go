package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/crislerwin/go-clean-api/docs" // Swagger docs
	"github.com/crislerwin/go-clean-api/internal/di"
	"github.com/crislerwin/go-clean-api/internal/infra/http/middleware"
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

	container := di.NewContainer(db)

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
	api.POST("/signup", container.UserHandler.SignUp)
	api.POST("/login", container.AuthHandler.Login)
	api.GET("/events", container.EventHandler.ListEvents)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/me", container.UserHandler.Me)
	protected.GET("/me/events", container.EventHandler.ListMyEvents)
	protected.POST("/orders", container.OrderHandler.CreateOrder)
	protected.GET("/orders", container.OrderHandler.ListMyOrders)
	protected.POST("/events", middleware.RoleMiddleware("admin"), container.EventHandler.CreateEvent)

	slog.Info("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

}
