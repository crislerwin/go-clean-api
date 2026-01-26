package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		handler := NewEventHandler(uc)

		r := gin.New()
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
