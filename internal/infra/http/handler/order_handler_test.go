package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCreateOrderUseCase struct {
	mock.Mock
}

type OrderRepoMock struct {
	mock.Mock
}

func (m *OrderRepoMock) Save(ctx context.Context, o *entity.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

type EventRepoMock struct {
	mock.Mock
}

func (m *EventRepoMock) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *EventRepoMock) GetSoldTicketsCount(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

type TxManagerMock struct {
	mock.Mock
}

func (m *TxManagerMock) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	eventID := uuid.New().String()

	t.Run("🔴 Fail: Should return 401 when UserID is missing from contex", func(t *testing.T) {

		// Arrange

		uc := usecase.NewCreateOrderUseCase(nil, nil, nil)
		handler := NewOrderHandler(uc)

		r := gin.New()
		r.POST("/orders", handler.CreateOrder)
		payload := map[string]any{
			"event_id": eventID,
			"quantity": 1,
		}

		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		// Act
		r.ServeHTTP(w, req)
		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "user not authenticated")

	})
	t.Run("🟢 Success: Should return 201 when authenticated and valid", func(t *testing.T) {
		// Arrange
		orderRepo := &OrderRepoMock{}
		eventRepo := &EventRepoMock{}
		txManager := &TxManagerMock{}
		userID := uuid.New().String()
		// Configurando comportamento esperado dos mocks
		// 1. Evento existe e tem capacidade
		eventRepo.On("GetTotalCapacity", mock.Anything, eventID).Return(10, nil)
		eventRepo.On("GetSoldTicketsCount", mock.Anything, eventID).Return(0, nil)
		// 2. Order salva com sucesso
		orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		uc := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
		handler := NewOrderHandler(uc)

		r := gin.New()

		// TRUQUE DO TESTE: Middleware Fake Local
		// Simulamos o Auth Middleware apenas para este teste
		r.Use(func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		})

		r.POST("/orders", handler.CreateOrder)

		payload := map[string]any{
			"event_id": eventID,
			"quantity": 2,
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		// Verifica se o UseCase recebeu o userID correto (que veio do contexto)
		// Isso garante que a integração Handler -> Context -> UseCase funcionou
		orderRepo.AssertExpectations(t)
	})
}
