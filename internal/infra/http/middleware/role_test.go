package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(requiredRole string) *gin.Engine {
		r := gin.New()
		r.Use(RoleMiddleware(requiredRole))
		r.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}

	t.Run("should allow request with correct role", func(t *testing.T) {
		r := setup("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		// Simulate AuthMiddleware setting the role
		r.Use(func(c *gin.Context) {
			c.Set("userRole", "admin")
			c.Next()
		})

		// Need to bypass the mock middleware for setup?
		// Actually, I can just set context manually in a custom handler wrapper or mock previous middleware.
		// Easier: just set the context in a middleware *before* RoleMiddleware.

		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("userRole", "admin")
			c.Next()
		})
		r2.Use(RoleMiddleware("admin"))
		r2.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		r2.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should deny request with incorrect role", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("userRole", "user")
			c.Next()
		})
		r2.Use(RoleMiddleware("admin"))
		r2.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		r2.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("should deny request when role is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		r2 := gin.New()
		r2.Use(RoleMiddleware("admin"))
		r2.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		r2.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
