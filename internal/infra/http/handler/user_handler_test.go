package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		userRepo *UserRepoMock
		handler  *UserHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		userRepo := &UserRepoMock{}
		uc := usecase.NewCreateUserUseCase(userRepo)
		handler := NewUserHandler(uc)

		r := gin.New()
		r.POST("/users", handler.CreateUser)

		return &testDeps{
			userRepo: userRepo,
			handler:  handler,
		}, r
	}

	t.Run("should create user successfully", func(t *testing.T) {
		deps, r := setup()
		deps.userRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		payload := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "John Doe", response["name"])
		assert.Equal(t, "john@example.com", response["email"])
		assert.NotEmpty(t, response["id"])
		assert.Empty(t, response["password"]) // Password should not be returned

		deps.userRepo.AssertExpectations(t)
	})

	t.Run("should return bad request on invalid body", func(t *testing.T) {
		_, r := setup()

		payload := map[string]string{
			"email": "john@example.com",
		} // Missing name and password
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal server error on repo failure", func(t *testing.T) {
		deps, r := setup()
		deps.userRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("db error"))

		payload := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.userRepo.AssertExpectations(t)
	})
}
