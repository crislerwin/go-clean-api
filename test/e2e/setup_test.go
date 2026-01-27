package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/http/handler"
	"github.com/crislerwin/go-clean-api/internal/infra/http/middleware"
	"github.com/crislerwin/go-clean-api/internal/infra/repository"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	testDB     *sqlx.DB
	testServer *gin.Engine
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Get test database URL from environment
	// Uses a separate test database to avoid polluting dev data
	// Matches docker-compose db_test service on port 5433
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://test_user:test_pass@localhost:5433/ticket_db_test?sslmode=disable"
	}

	var err error
	testDB, err = sqlx.Connect("pgx", testDBURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	// Setup test server
	testServer = setupTestServer(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	os.Exit(code)
}

// setupTestServer initializes the Gin server with all routes
func setupTestServer(db *sqlx.DB) *gin.Engine {
	r := gin.New()

	// Setup repositories
	txManager := database.NewTransactionManager(db)
	eventRepo := repository.NewEventRepositorySqlx(db)
	userRepo := repository.NewUserRepositorySQLx(db)
	orderRepo := repository.NewOrderRepositorySQLx(db)

	// Setup use cases
	signUpUseCase := usecase.NewSignUpUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)
	createEventUseCase := usecase.NewCreateEventUseCase(eventRepo, txManager)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
	listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)
	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	listUserEventsUseCase := usecase.NewListUserEventsUseCase(eventRepo)

	// Setup handlers
	userHandler := handler.NewUserHandler(signUpUseCase, usecase.NewGetUserUseCase(userRepo))
	authHandler := handler.NewAuthHandler(loginUseCase)
	eventHandler := handler.NewEventHandler(createEventUseCase, listEventsUseCase, listUserEventsUseCase)
	orderHandler := handler.NewOrderHandler(createOrderUseCase, listUserOrdersUseCase)

	// Setup routes
	api := r.Group("/api/v1")

	// Public routes
	api.POST("/signup", userHandler.SignUp)
	api.POST("/login", authHandler.Login)
	api.GET("/events", eventHandler.ListEvents)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/me", userHandler.Me)
	protected.GET("/me/events", eventHandler.ListMyEvents)
	protected.POST("/orders", orderHandler.CreateOrder)
	protected.GET("/orders", orderHandler.ListMyOrders)
	protected.POST("/events", middleware.RoleMiddleware("admin"), eventHandler.CreateEvent)

	return r
}

// cleanupDatabase clears all test data between tests
// This ensures test isolation while using a dedicated test database
func cleanupDatabase(t *testing.T) {
	t.Helper()

	// Delete in order to respect foreign key constraints
	queries := []string{
		"DELETE FROM tickets",
		"DELETE FROM orders",
		"DELETE FROM events",
		"DELETE FROM users",
	}

	for _, query := range queries {
		if _, err := testDB.Exec(query); err != nil {
			t.Fatalf("failed to cleanup database: %v", err)
		}
	}
}

// makeRequest is a helper to make HTTP requests
func makeRequest(method, path string, body interface{}, token string) (*httptest.ResponseRecorder, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testServer.ServeHTTP(w, req)

	return w, nil
}

// parseResponse is a helper to parse JSON responses
func parseResponse(w *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), v)
}
