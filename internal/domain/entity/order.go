package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
)

type Order struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	UserID      uuid.UUID
	TotalAmount float64
	Quantity    int
	Status      OrderStatus
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
		ticket, err := NewTicket(eventID, orderID, pricePerTicket)
		if err != nil {
			return nil, err
		}
		tickets[i] = *ticket
		totalAmount += pricePerTicket
	}

	status := OrderStatusPending
	if totalAmount == 0 {
		status = OrderStatusConfirmed
	}

	return &Order{
		ID:          orderID,
		EventID:     eventID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Quantity:    quantity,
		Status:      status,
		Tickets:     tickets,
		CreatedAt:   time.Now(),
	}, nil
}
