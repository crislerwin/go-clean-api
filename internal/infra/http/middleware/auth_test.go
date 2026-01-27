package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/infra/http/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func() *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware())
		r.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})
		return r
	}

	t.Run("should allow request with valid token", func(t *testing.T) {
		r := setup()
		token, _ := auth.GenerateToken("user-123", "user")

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user-123")
	})

	t.Run("should return unauthorized when header is missing", func(t *testing.T) {
		r := setup()
		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("should return unauthorized when format is invalid", func(t *testing.T) {
		r := setup()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("should return unauthorized when token is invalid", func(t *testing.T) {
		r := setup()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})
}
