package entity

import (
	"errors"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusValid     TicketStatus = "VALID"
	TicketStatusUsed      TicketStatus = "USED"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

var (
	ErrPriceInvalid          = errors.New("price must be non-negative")
	ErrTicketAlreadyUsed     = errors.New("ticket is already used")
	ErrTicketAlreadyCanceled = errors.New("ticket is already canceled")
	ErrTicketInvalid         = errors.New("ticket is invalid")
)

type Ticket struct {
	ID      uuid.UUID
	EventID uuid.UUID
	OrderID uuid.UUID
	Price   float64
	Status  TicketStatus
}

func NewTicket(eventID, orderID uuid.UUID, price float64) (*Ticket, error) {
	if price < 0 {
		return nil, ErrPriceInvalid
	}
	return &Ticket{
		ID:      uuid.New(),
		EventID: eventID,
		OrderID: orderID,
		Price:   price,
		Status:  TicketStatusValid,
	}, nil
}

func (t *Ticket) Validate() error {
	if t.Price < 0 {
		return ErrPriceInvalid
	}
	if t.Status == "" {
		return ErrTicketInvalid
	}
	return nil
}

func (t *Ticket) Use() error {
	if t.Status == TicketStatusUsed {
		return ErrTicketAlreadyUsed
	}
	if t.Status == TicketStatusCancelled {
		return ErrTicketAlreadyCanceled
	}
	t.Status = TicketStatusUsed
	return nil
}

func (t *Ticket) Cancel() error {
	if t.Status == TicketStatusUsed {
		return ErrTicketAlreadyUsed
	}
	if t.Status == TicketStatusCancelled {
		return ErrTicketAlreadyCanceled
	}
	t.Status = TicketStatusCancelled
	return nil
}
