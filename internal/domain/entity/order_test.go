package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	t.Run("should create a new valid order", func(t *testing.T) {
		eventID := uuid.New()
		userID := uuid.New()
		quantity := 2
		pricePerTicket := 100.0

		order, err := NewOrder(eventID, userID, quantity, pricePerTicket)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if order.TotalAmount != 200.00 {
			t.Errorf("Expected total amount 200.0, got %v", order.TotalAmount)
		}
		if order.Status != "PENDING" {
			t.Errorf("Expected status PENDING, got %s", order.Status)
		}
		if len(order.Tickets) != 2 {

			t.Errorf("Expected 2 tickets, got %d", len(order.Tickets))
		}

	})

	t.Run("should throw error when quantity is zero", func(t *testing.T) {
		eventID := uuid.New()
		userID := uuid.New()
		quantity := 0
		price := 10.0

		_, err := NewOrder(eventID, userID, quantity, price)

		if err == nil {
			t.Error("Expected error for zero quantity, got nil")
		}
	})
}
