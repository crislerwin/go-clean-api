package entity

import "github.com/google/uuid"

type Ticket struct {
	ID      uuid.UUID
	EventID uuid.UUID
	OrderID uuid.UUID
	Price   float64
	Status  string
}

func NewTicket(eventID, orderID uuid.UUID, price float64) *Ticket {
	return &Ticket{
		ID:      uuid.New(),
		EventID: eventID,
		OrderID: orderID,
		Price:   price,
		Status:  "VALID",
	}
}
