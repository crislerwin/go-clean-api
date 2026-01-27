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

func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepoMock) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestUserHandler_SignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		userRepo *UserRepoMock
		handler  *UserHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		userRepo := &UserRepoMock{}
		useCase := usecase.NewSignUpUseCase(userRepo)
		getUserUseCase := usecase.NewGetUserUseCase(userRepo)
		handler := NewUserHandler(useCase, getUserUseCase)

		r := gin.New()
		r.POST("/signup", handler.SignUp)

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
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
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
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
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
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.userRepo.AssertExpectations(t)
		deps.userRepo.AssertExpectations(t)
	})
}

func TestUserHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		userRepo *UserRepoMock
		handler  *UserHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		userRepo := &UserRepoMock{}
		signUpUC := usecase.NewSignUpUseCase(userRepo)
		getUserUC := usecase.NewGetUserUseCase(userRepo)
		handler := NewUserHandler(signUpUC, getUserUC)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", "user-123")
		})
		r.GET("/me", handler.Me)

		return &testDeps{
			userRepo: userRepo,
			handler:  handler,
		}, r
	}

	t.Run("should get user info successfully", func(t *testing.T) {
		deps, r := setup()

		user := &entity.User{
			ID:    "user-123",
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "user",
		}
		deps.userRepo.On("GetByID", mock.Anything, "user-123").Return(user, nil)

		req, _ := http.NewRequest("GET", "/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", response["id"])
		assert.Equal(t, "John Doe", response["name"])
		assert.Equal(t, "john@example.com", response["email"])

		deps.userRepo.AssertExpectations(t)
	})

	t.Run("should return internal server error if user not found (should likely not happen with valid token)", func(t *testing.T) {
		deps, r := setup()

		deps.userRepo.On("GetByID", mock.Anything, "user-123").Return((*entity.User)(nil), errors.New("user not found"))

		req, _ := http.NewRequest("GET", "/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		deps.userRepo.AssertExpectations(t)
	})
}
