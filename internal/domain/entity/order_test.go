package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		eventID        uuid.UUID
		userID         uuid.UUID
		quantity       int
		pricePerTicket float64
		expectedErr    error
		expectedTotal  float64
		expectedStatus OrderStatus
		expectedCount  int
	}{
		{
			name:           "should create a new valid order",
			eventID:        eventID,
			userID:         userID,
			quantity:       2,
			pricePerTicket: 100.0,
			expectedErr:    nil,
			expectedTotal:  200.0,
			expectedStatus: OrderStatusPending,
			expectedCount:  2,
		},
		{
			name:           "should create a free order with status CONFIRMED",
			eventID:        eventID,
			userID:         userID,
			quantity:       2,
			pricePerTicket: 0.0,
			expectedErr:    nil,
			expectedTotal:  0.0,
			expectedStatus: OrderStatusConfirmed,
			expectedCount:  2,
		},
		{
			name:           "should throw error when quantity is zero",
			eventID:        eventID,
			userID:         userID,
			quantity:       0,
			pricePerTicket: 100.0,
			expectedErr:    ErrInvalidQuantity,
			expectedTotal:  0,
			expectedStatus: "",
			expectedCount:  0,
		},
		{
			name:           "should throw error when quantity is negative",
			eventID:        eventID,
			userID:         userID,
			quantity:       -1,
			pricePerTicket: 100.0,
			expectedErr:    ErrInvalidQuantity,
			expectedTotal:  0,
			expectedStatus: "",
			expectedCount:  0,
		},
		{
			name:           "should throw error when price is invalid (negative)",
			eventID:        eventID,
			userID:         userID,
			quantity:       1,
			pricePerTicket: -1.0,
			expectedErr:    ErrPriceInvalid,
			expectedTotal:  0,
			expectedStatus: "",
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.eventID, tt.userID, tt.quantity, tt.pricePerTicket)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.expectedTotal, order.TotalAmount)
				assert.Equal(t, tt.expectedStatus, order.Status)
				assert.Equal(t, tt.expectedCount, len(order.Tickets))
				assert.Equal(t, tt.eventID, order.EventID)
				assert.Equal(t, tt.userID, order.UserID)
				assert.NotEmpty(t, order.ID)
				assert.NotEmpty(t, order.CreatedAt)
			}
		})
	}
}
