package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventFlow(t *testing.T) {
	cleanupDatabase(t)

	// Setup: Create admin and regular user
	adminToken := createAdminUser(t)
	userToken := createRegularUser(t, "user@example.com")

	t.Run("Create Event - Admin Success", func(t *testing.T) {
		futureDate := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)

		eventPayload := map[string]interface{}{
			"name":         "Rock in Rio 2025",
			"location":     "Rio de Janeiro",
			"organization": "Live Nation",
			"rating":       "Livre",
			"date":         futureDate,
			"capacity":     100000,
			"price":        150.50,
			"image_url":    "https://example.com/rock-in-rio.jpg",
		}

		w, err := makeRequest("POST", "/api/v1/events", eventPayload, adminToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err = parseResponse(w, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["id"])
		assert.Equal(t, "Rock in Rio 2025", resp["name"])
	})

	t.Run("Create Event - Non-Admin Forbidden", func(t *testing.T) {
		futureDate := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)

		eventPayload := map[string]interface{}{
			"name":         "Unauthorized Event",
			"location":     "São Paulo",
			"organization": "Test Org",
			"rating":       "Livre",
			"date":         futureDate,
			"capacity":     1000,
			"price":        50.0,
			"image_url":    "https://example.com/event.jpg",
		}

		w, err := makeRequest("POST", "/api/v1/events", eventPayload, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Create Event - No Token", func(t *testing.T) {
		futureDate := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)

		eventPayload := map[string]interface{}{
			"name":         "No Auth Event",
			"location":     "São Paulo",
			"organization": "Test Org",
			"rating":       "Livre",
			"date":         futureDate,
			"capacity":     1000,
			"price":        50.0,
			"image_url":    "https://example.com/event.jpg",
		}

		w, err := makeRequest("POST", "/api/v1/events", eventPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Create Event - Past Date", func(t *testing.T) {
		pastDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)

		eventPayload := map[string]interface{}{
			"name":         "Past Event",
			"location":     "Rio",
			"organization": "Org",
			"rating":       "Livre",
			"date":         pastDate,
			"capacity":     100,
			"price":        10.0,
			"image_url":    "https://example.com/event.jpg",
		}

		w, err := makeRequest("POST", "/api/v1/events", eventPayload, adminToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("List All Events - Public", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/events", nil, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var events []map[string]interface{}
		err = parseResponse(w, &events)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 1) // At least the event created above
	})

	t.Run("List My Events - Admin", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/me/events", nil, adminToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var events []map[string]interface{}
		err = parseResponse(w, &events)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 1) // Admin created at least one event
	})

	t.Run("List My Events - Regular User (Empty)", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/me/events", nil, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var events []map[string]interface{}
		err = parseResponse(w, &events)
		require.NoError(t, err)
		assert.Equal(t, 0, len(events)) // Regular user hasn't created any events
	})
}

// Helper function to create an admin user and return token
func createAdminUser(t *testing.T) string {
	t.Helper()

	// First create a regular user
	signupPayload := map[string]string{
		"name":     "Admin User",
		"email":    "admin@example.com",
		"password": "admin123",
	}

	w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, w.Code)

	var signupResp map[string]interface{}
	err = parseResponse(w, &signupResp)
	require.NoError(t, err)

	userID := signupResp["id"].(string)

	// Manually update user role to admin in database
	_, err = testDB.Exec("UPDATE users SET role = 'admin' WHERE id = $1", userID)
	require.NoError(t, err)

	// Login to get token
	loginPayload := map[string]string{
		"email":    "admin@example.com",
		"password": "admin123",
	}

	w, err = makeRequest("POST", "/api/v1/login", loginPayload, "")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = parseResponse(w, &loginResp)
	require.NoError(t, err)

	return loginResp["access_token"].(string)
}

// Helper function to create a regular user and return token
func createRegularUser(t *testing.T, email string) string {
	t.Helper()

	signupPayload := map[string]string{
		"name":     "Regular User",
		"email":    email,
		"password": "user123",
	}

	w, err := makeRequest("POST", "/api/v1/signup", signupPayload, "")
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, w.Code)

	loginPayload := map[string]string{
		"email":    email,
		"password": "user123",
	}

	w, err = makeRequest("POST", "/api/v1/login", loginPayload, "")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = parseResponse(w, &loginResp)
	require.NoError(t, err)

	return loginResp["access_token"].(string)
}
