package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (m *OrderRepoMock) Save(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepoMock) GetByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Order), args.Error(1)
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

func (m *EventRepoMock) Create(ctx context.Context, event *entity.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventRepoMock) GetByID(ctx context.Context, eventID string) (*entity.Event, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *EventRepoMock) ListAll(ctx context.Context) ([]*entity.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *EventRepoMock) ListByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
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
		createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
		listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)
		handler := NewOrderHandler(createOrderUseCase, listUserOrdersUseCase)

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
		listUc := usecase.NewListUserOrdersUseCase(nil, nil)
		handler := NewOrderHandler(uc, listUc)
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
		// Mock GetByID
		mockEvent := &entity.Event{
			ID:       uuid.MustParse(eventID),
			Price:    100.0,
			Capacity: 10,
		}
		deps.eventRepo.On("GetByID", mock.Anything, eventID).Return(mockEvent, nil)

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

		mockEvent := &entity.Event{
			ID:       uuid.MustParse(eventID),
			Price:    100.0,
			Capacity: 10,
		}
		deps.eventRepo.On("GetByID", mock.Anything, eventID).Return(mockEvent, nil)
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
		deps.eventRepo.On("GetByID", mock.Anything, eventID).Return((*entity.Event)(nil), usecase.ErrEventNotFound)

		payload := map[string]any{"event_id": eventID, "quantity": 1}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "event not found")
	})
}

func TestOrderHandler_ListMyOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		orderRepo *OrderRepoMock
		eventRepo *EventRepoMock
		handler   *OrderHandler
	}

	setup := func() (*testDeps, *gin.Engine, string) {
		orderRepo := &OrderRepoMock{}
		eventRepo := &EventRepoMock{}
		// TxManager not needed for list
		createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventRepo, nil)
		listUserOrdersUseCase := usecase.NewListUserOrdersUseCase(orderRepo, eventRepo)
		handler := NewOrderHandler(createOrderUseCase, listUserOrdersUseCase)

		userID := "user-123"

		r := gin.New()
		r.Use(func(c *gin.Context) {
			if userID != "" {
				c.Set("userID", userID)
			}
			c.Next()
		})
		r.GET("/orders", handler.ListMyOrders)

		return &testDeps{
			orderRepo: orderRepo,
			eventRepo: eventRepo,
			handler:   handler,
		}, r, userID
	}

	t.Run("🟢 Success: List orders", func(t *testing.T) {
		deps, r, userID := setup()
		eventID := uuid.New()
		now := time.Now()

		mockOrders := []*entity.Order{
			{
				ID:          uuid.New(),
				EventID:     eventID,
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"), // not used here
				TotalAmount: 50.0,
				Quantity:    1,
				Status:      "PENDING",
				CreatedAt:   now,
			},
		}
		mockEvent := &entity.Event{Name: "Show", Date: now}

		deps.orderRepo.On("GetByUserID", mock.Anything, userID).Return(mockOrders, nil)
		deps.eventRepo.On("GetByID", mock.Anything, eventID.String()).Return(mockEvent, nil)

		req, _ := http.NewRequest("GET", "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Show")
		assert.Contains(t, w.Body.String(), "50")
	})
}
