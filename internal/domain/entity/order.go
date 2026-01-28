package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)

type Order struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	UserID      uuid.UUID
	TotalAmount float64
	Quantity    int
	Status      string
	Tickets     []Ticket // One-to-many relationship
	CreatedAt   time.Time
}

func NewOrder(eventID, userID uuid.UUID, quantity int, pricePerTicket float64) (*Order, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	orderID := uuid.New()
	totalAmount := 0.0
	tickets := make([]Ticket, quantity)
	for i := range tickets {
		tickets[i] = *NewTicket(eventID, orderID, pricePerTicket)
		totalAmount += pricePerTicket
	}

	return &Order{
		ID:          orderID,
		EventID:     eventID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Quantity:    quantity,
		Status:      "PENDING",
		Tickets:     tickets,
		CreatedAt:   time.Now(),
	}, nil
}
