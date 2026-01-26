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
	"golang.org/x/crypto/bcrypt"
)

type AuthUserRepoMock struct {
	mock.Mock
}

func (m *AuthUserRepoMock) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *AuthUserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testDeps struct {
		userRepo *AuthUserRepoMock
		handler  *AuthHandler
	}

	setup := func() (*testDeps, *gin.Engine) {
		userRepo := &AuthUserRepoMock{}
		loginUC := usecase.NewLoginUseCase(userRepo)
		handler := NewAuthHandler(loginUC)

		r := gin.New()
		r.POST("/login", handler.Login)

		return &testDeps{
			userRepo: userRepo,
			handler:  handler,
		}, r
	}

	t.Run("should login successfully", func(t *testing.T) {
		deps, r := setup()

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &entity.User{
			ID:       "user-123",
			Email:    "john@example.com",
			Password: string(hashedPassword),
		}

		deps.userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)

		payload := map[string]string{
			"email":    "john@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")
		deps.userRepo.AssertExpectations(t)
	})

	t.Run("should return unauthorized on invalid credentials", func(t *testing.T) {
		deps, r := setup()
		deps.userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return((*entity.User)(nil), errors.New("user not found"))

		payload := map[string]string{
			"email":    "john@example.com",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})

	t.Run("should return bad request on invalid body", func(t *testing.T) {
		_, r := setup()

		payload := map[string]string{
			"email": "john@example.com",
		}
		jsonBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
