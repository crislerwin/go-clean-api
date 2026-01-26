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
	ID          uuid.UUID `db:"id"`
	EventID     uuid.UUID `db:"event_id"`
	UserID      uuid.UUID `db:"user_id"`
	TotalAmount float64   `db:"total_amount"`
	Quantity    int       `db:"quantity"`
	Status      string    `db:"status"`
	Tickets     []Ticket  `db:"-"` // One-to-many relationship, not a column
	CreatedAt   time.Time `db:"created_at"`
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
