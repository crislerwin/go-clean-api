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

	type testDeps struct {
		orderRepo *OrderRepoMock
		eventRepo *EventRepoMock
		txManager *TxManagerMock
		handler   *OrderHandler
	}

	setup := func() (*testDeps, *gin.Engine, string) {
		orderRepo := &OrderRepoMock{}
		eventRepo := &EventRepoMock{}
		txManager := &TxManagerMock{}
		userID := uuid.New().String()

		uc := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
		handler := NewOrderHandler(uc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			if userID != "" {
				c.Set("userID", userID)
			}
			c.Next()
		})
		r.POST("/orders", handler.CreateOrder)

		return &testDeps{
			orderRepo: orderRepo,
			eventRepo: eventRepo,
			txManager: txManager,
			handler:   handler,
		}, r, userID
	}

	t.Run("🔴 Fail: Should return 401 when UserID is missing from context", func(t *testing.T) {
		// Custom setup for unauthenticated case
		uc := usecase.NewCreateOrderUseCase(nil, nil, nil)
		handler := NewOrderHandler(uc)
		r := gin.New()
		r.POST("/orders", handler.CreateOrder) // No auth middleware

		payload := map[string]any{"event_id": eventID, "quantity": 1}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "user not authenticated")
	})

	t.Run("🟢 Success: Should return 201 when authenticated and valid", func(t *testing.T) {
		deps, r, _ := setup()
		deps.eventRepo.On("GetTotalCapacity", mock.Anything, eventID).Return(10, nil)
		deps.eventRepo.On("GetSoldTicketsCount", mock.Anything, eventID).Return(0, nil)
		deps.orderRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		payload := map[string]any{"event_id": eventID, "quantity": 2}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		deps.orderRepo.AssertExpectations(t)
	})

	t.Run("🔴 Fail: Should return 400 when body is invalid", func(t *testing.T) {
		_, r, _ := setup()
		payload := map[string]any{"event_id": eventID} // Missing quantity
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("🔴 Fail: Should return 409 when event is sold out", func(t *testing.T) {
		deps, r, _ := setup()
		deps.eventRepo.On("GetTotalCapacity", mock.Anything, eventID).Return(10, nil)
		deps.eventRepo.On("GetSoldTicketsCount", mock.Anything, eventID).Return(10, nil)

		payload := map[string]any{"event_id": eventID, "quantity": 1}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "event sold out")
	})

	t.Run("🔴 Fail: Should return 404 when event not found", func(t *testing.T) {
		deps, r, _ := setup()
		deps.eventRepo.On("GetTotalCapacity", mock.Anything, eventID).Return(0, usecase.ErrEventNotFound)

		payload := map[string]any{"event_id": eventID, "quantity": 1}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "event not found")
	})
}
