package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventHandler_CreateEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		eventRepo *EventRepoMock
		handler   *EventHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		eventRepo := &EventRepoMock{}
		uc := usecase.NewCreateEventUseCase(eventRepo, nil)
		listUc := usecase.NewListEventsUseCase(eventRepo)
		listUserUc := usecase.NewListUserEventsUseCase(eventRepo)
		handler := NewEventHandler(uc, listUc, listUserUc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", "5334d5d7-e5d8-4d56-9257-2b7b5e5d3c8a")
			c.Next()
		})
		r.POST("/events", handler.CreateEvent)

		return &testDeps{
			eventRepo: eventRepo,
			handler:   handler,
		}, r
	}

	futureDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	tests := []struct {
		name           string
		payload        map[string]any
		setupMocks     func(deps *testDeps)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "🟢 Success: Should return 201 when payload is valid",
			payload: map[string]any{
				"name":         "Rock in Rio",
				"location":     "Rio de Janeiro",
				"organization": "Live Nation",
				"rating":       "Livre",
				"date":         futureDate,
				"capacity":     100000,
				"price":        100.0,
				"image_url":    "http://example.com/image.jpg",
			},
			setupMocks: func(deps *testDeps) {
				deps.eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Event")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "", // Body will contain ID and Name, handled in assertions if needed, but here we check status
		},
		{
			name: "🔴 Fail: Should return 400 when body is invalid",
			payload: map[string]any{
				"name": "Rock in Rio",
				// Missing fields
			},
			setupMocks:     func(deps *testDeps) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "🔴 Fail: Should return 400 when date is in past",
			payload: map[string]any{
				"name":         "Old Event",
				"location":     "Rio",
				"organization": "Org",
				"rating":       "Livre",
				"date":         time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"capacity":     100,
				"price":        10.0,
				"image_url":    "url",
			},
			setupMocks:     func(deps *testDeps) {}, // No mock needed as usecase fails before repo call
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "event date must be in the future",
		},
		{
			name: "🔴 Fail: Should return 500 when repository fails",
			payload: map[string]any{
				"name":         "Rock in Rio",
				"location":     "Rio de Janeiro",
				"organization": "Live Nation",
				"rating":       "Livre",
				"date":         futureDate,
				"capacity":     100000,
				"price":        100.0,
				"image_url":    "http://example.com/image.jpg",
			},
			setupMocks: func(deps *testDeps) {
				deps.eventRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Event")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, r := setup()
			tt.setupMocks(deps)

			jsonBody, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(jsonBody))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
			deps.eventRepo.AssertExpectations(t)
		})
	}
}

func TestEventHandler_ListEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		eventRepo *EventRepoMock
		handler   *EventHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		eventRepo := &EventRepoMock{}
		uc := usecase.NewCreateEventUseCase(eventRepo, nil)
		listUc := usecase.NewListEventsUseCase(eventRepo)
		listUserUc := usecase.NewListUserEventsUseCase(eventRepo)
		handler := NewEventHandler(uc, listUc, listUserUc)

		r := gin.New()
		r.GET("/events", handler.ListEvents)

		return &testDeps{
			eventRepo: eventRepo,
			handler:   handler,
		}, r
	}

	t.Run("should list events successfully", func(t *testing.T) {
		deps, r := setup()

		mockEvents := []*entity.Event{
			{Name: "Event 1"},
			{Name: "Event 2"},
		}

		deps.eventRepo.On("ListAll", mock.Anything).Return(mockEvents, nil)

		req, _ := http.NewRequest("GET", "/events", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Event 1")
		assert.Contains(t, w.Body.String(), "Event 2")
		deps.eventRepo.AssertExpectations(t)
	})

	t.Run("should return internal server error when usecase fails", func(t *testing.T) {
		deps, r := setup()

		deps.eventRepo.On("ListAll", mock.Anything).Return(([]*entity.Event)(nil), assert.AnError)

		req, _ := http.NewRequest("GET", "/events", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.eventRepo.AssertExpectations(t)
	})
}

func TestEventHandler_ListMyEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		eventRepo *EventRepoMock
		handler   *EventHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		eventRepo := &EventRepoMock{}
		uc := usecase.NewCreateEventUseCase(eventRepo, nil)
		listUc := usecase.NewListEventsUseCase(eventRepo)
		listUserUc := usecase.NewListUserEventsUseCase(eventRepo) // Added usecase
		handler := NewEventHandler(uc, listUc, listUserUc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", "5334d5d7-e5d8-4d56-9257-2b7b5e5d3c8a")
			c.Next()
		})
		r.GET("/me/events", handler.ListMyEvents)

		return &testDeps{
			eventRepo: eventRepo,
			handler:   handler,
		}, r
	}

	t.Run("should list user events successfully", func(t *testing.T) {
		deps, r := setup()

		mockEvents := []*entity.Event{
			{Name: "My Event 1"},
			{Name: "My Event 2"},
		}

		deps.eventRepo.On("ListByUserID", mock.Anything, "5334d5d7-e5d8-4d56-9257-2b7b5e5d3c8a").Return(mockEvents, nil)

		req, _ := http.NewRequest("GET", "/me/events", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "My Event 1")
		assert.Contains(t, w.Body.String(), "My Event 2")
		deps.eventRepo.AssertExpectations(t)
	})

	t.Run("should return internal server error when usecase fails", func(t *testing.T) {
		deps, r := setup()

		deps.eventRepo.On("ListByUserID", mock.Anything, "5334d5d7-e5d8-4d56-9257-2b7b5e5d3c8a").Return(([]*entity.Event)(nil), assert.AnError)

		req, _ := http.NewRequest("GET", "/me/events", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.eventRepo.AssertExpectations(t)
	})
}
