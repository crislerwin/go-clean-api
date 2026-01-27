package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	cleanupDatabase(t)

	t.Run("Complete Auth Flow - Signup, Login, Get User", func(t *testing.T) {
		// Step 1: Sign up a new user
		signupPayload := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
		}

		w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var signupResp map[string]interface{}
		err = parseResponse(w, &signupResp)
		require.NoError(t, err)
		assert.NotEmpty(t, signupResp["id"])
		assert.Equal(t, "John Doe", signupResp["name"])
		assert.Equal(t, "john@example.com", signupResp["email"])

		// Step 2: Login with the new user
		loginPayload := map[string]string{
			"email":    "john@example.com",
			"password": "password123",
		}

		w, err = makeRequest("POST", "/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var loginResp map[string]interface{}
		err = parseResponse(w, &loginResp)
		require.NoError(t, err)
		require.NotNil(t, loginResp["access_token"], "access_token should not be nil")
		assert.NotEmpty(t, loginResp["access_token"])

		token := loginResp["access_token"].(string)

		// Step 3: Get current user info with token
		w, err = makeRequest("GET", "/api/v1/me", nil, token)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var meResp map[string]interface{}
		err = parseResponse(w, &meResp)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", meResp["name"])
		assert.Equal(t, "john@example.com", meResp["email"])
		assert.Equal(t, "user", meResp["role"])
	})

	t.Run("Signup - Duplicate Email", func(t *testing.T) {
		signupPayload := map[string]string{
			"name":     "Jane Doe",
			"email":    "john@example.com", // Same email as before
			"password": "password123",
		}

		w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Signup - Invalid Email", func(t *testing.T) {
		signupPayload := map[string]string{
			"name":     "Invalid User",
			"email":    "not-an-email",
			"password": "password123",
		}

		w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Signup - Short Password", func(t *testing.T) {
		signupPayload := map[string]string{
			"name":     "Short Pass User",
			"email":    "short@example.com",
			"password": "123",
		}

		w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Login - Invalid Credentials", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "john@example.com",
			"password": "wrongpassword",
		}

		w, err := makeRequest("POST", "/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Login - Non-existent User", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}

		w, err := makeRequest("POST", "/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get User - No Token", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/me", nil, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get User - Invalid Token", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/me", nil, "invalid-token")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
