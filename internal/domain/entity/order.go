package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)

type Ticket struct {
	ID      uuid.UUID `db:"id"`
	EventID uuid.UUID `db:"event_id"`
	OrderID uuid.UUID `db:"order_id"`
	Price   float64   `db:"price"`
	Status  string    `db:"status"`
}

type Order struct {
	ID          uuid.UUID `db:"id"`
	EventID     uuid.UUID `db:"event_id"`
	UserID      uuid.UUID `db:"user_id"`
	TotalAmount float64   `db:"total_amount"`
	Status      string    `db:"status"`
	Tickets     []Ticket  `db:"tickets"`
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
		tickets[i] = Ticket{
			ID:      uuid.New(),
			EventID: eventID,
			OrderID: orderID,
			Price:   pricePerTicket,
			Status:  "VALID",
		}
		totalAmount += pricePerTicket
	}
	return &Order{
		ID:          orderID,
		EventID:     eventID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "PENDING",
		Tickets:     tickets,
		CreatedAt:   time.Now(),
	}, nil
}
