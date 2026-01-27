package e2e

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderFlow(t *testing.T) {
	cleanupDatabase(t)

	// Setup: Create admin, event, and user
	adminToken := createAdminUser(t)
	userToken := createRegularUser(t, "buyer@example.com")
	eventID := createTestEvent(t, adminToken, 10) // Event with 10 tickets

	t.Run("Create Order - Success", func(t *testing.T) {
		orderPayload := map[string]interface{}{
			"event_id": eventID,
			"quantity": 2,
		}

		w, err := makeRequest("POST", "/api/v1/orders", orderPayload, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err = parseResponse(w, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["id"])
		assert.Equal(t, "PENDING", resp["status"])
		assert.NotEmpty(t, resp["total_amount"])
		assert.NotEmpty(t, resp["created_at"])
	})

	t.Run("List My Orders", func(t *testing.T) {
		w, err := makeRequest("GET", "/api/v1/orders", nil, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var orders []map[string]interface{}
		err = parseResponse(w, &orders)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(orders), 1) // At least the order created above
	})

	t.Run("Create Order - No Token", func(t *testing.T) {
		orderPayload := map[string]interface{}{
			"event_id": eventID,
			"quantity": 1,
		}

		w, err := makeRequest("POST", "/api/v1/orders", orderPayload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Create Order - Invalid Event ID", func(t *testing.T) {
		orderPayload := map[string]interface{}{
			"event_id": "00000000-0000-0000-0000-000000000000",
			"quantity": 1,
		}

		w, err := makeRequest("POST", "/api/v1/orders", orderPayload, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Create Order - Zero Quantity", func(t *testing.T) {
		orderPayload := map[string]interface{}{
			"event_id": eventID,
			"quantity": 0,
		}

		w, err := makeRequest("POST", "/api/v1/orders", orderPayload, userToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRaceConditions(t *testing.T) {
	cleanupDatabase(t)

	// Setup: Create admin, event with limited capacity, and multiple users
	adminToken := createAdminUser(t)
	eventID := createTestEvent(t, adminToken, 5) // Only 5 tickets available

	// Create 10 users who will try to buy tickets concurrently
	numBuyers := 10
	tokens := make([]string, numBuyers)
	for i := 0; i < numBuyers; i++ {
		email := fmt.Sprintf("buyer%d@example.com", i)
		tokens[i] = createRegularUser(t, email)
	}

	t.Run("Concurrent Ticket Purchase - Race Condition Handling", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make([]int, numBuyers)
		mu := sync.Mutex{}

		// All users try to buy 2 tickets simultaneously
		for i := 0; i < numBuyers; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				orderPayload := map[string]interface{}{
					"event_id": eventID,
					"quantity": 2,
				}

				w, err := makeRequest("POST", "/api/v1/orders", orderPayload, tokens[idx])
				if err != nil {
					return
				}

				mu.Lock()
				results[idx] = w.Code
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		// Count successful and failed orders
		successCount := 0
		conflictCount := 0

		for _, code := range results {
			switch code {
			case http.StatusCreated:
				successCount++
			case http.StatusConflict:
				conflictCount++
			}
		}

		// With 5 tickets available and each order requesting 2 tickets,
		// we expect some orders to succeed and some to fail
		// The exact number depends on timing, but we should have:
		// - At least 2 successful orders (using 4 tickets minimum)
		// - At least some conflicts (sold out responses)
		assert.GreaterOrEqual(t, successCount, 2, "Expected at least 2 successful orders")
		assert.GreaterOrEqual(t, conflictCount, 1, "Expected at least 1 conflict response")
		assert.Equal(t, numBuyers, successCount+conflictCount, "All requests should complete")
	})
}

// Helper function to create a test event and return its ID.
func createTestEvent(t *testing.T, adminToken string, capacity int) string {
	t.Helper()

	futureDate := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)

	eventPayload := map[string]interface{}{
		"name":         fmt.Sprintf("Test Event %d", time.Now().Unix()),
		"location":     "Test Location",
		"organization": "Test Org",
		"rating":       "Livre",
		"date":         futureDate,
		"capacity":     capacity,
		"price":        50.0,
		"image_url":    "https://example.com/test.jpg",
	}

	w, err := makeRequest("POST", "/api/v1/events", eventPayload, adminToken)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)

	return resp["id"].(string)
}
