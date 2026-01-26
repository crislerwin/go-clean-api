package entity

import "github.com/google/uuid"

type Ticket struct {
	ID      uuid.UUID `db:"id"`
	EventID uuid.UUID `db:"event_id"`
	OrderID uuid.UUID `db:"order_id"`
	Price   float64   `db:"price"`
	Status  string    `db:"status"`
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
