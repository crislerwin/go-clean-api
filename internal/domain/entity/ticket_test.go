package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTicket(t *testing.T) {
	eventID := uuid.New()
	orderID := uuid.New()

	t.Run("should create a valid ticket", func(t *testing.T) {
		ticket, err := NewTicket(eventID, orderID, 10.0)
		assert.NoError(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, TicketStatusValid, ticket.Status)
		assert.Equal(t, 10.0, ticket.Price)
	})

	t.Run("should create a valid ticket with zero price", func(t *testing.T) {
		ticket, err := NewTicket(eventID, orderID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, TicketStatusValid, ticket.Status)
		assert.Equal(t, 0.0, ticket.Price)
	})

	t.Run("should return error when price is invalid", func(t *testing.T) {
		ticket, err := NewTicket(eventID, orderID, -10.0)
		assert.ErrorIs(t, err, ErrPriceInvalid)
		assert.Nil(t, ticket)
	})
}

func TestTicket_Validate(t *testing.T) {
	t.Run("should validate correctly", func(t *testing.T) {
		ticket := &Ticket{
			Price:  10.0,
			Status: TicketStatusValid,
		}
		assert.NoError(t, ticket.Validate())
	})

	t.Run("should return error for invalid price", func(t *testing.T) {
		ticket := &Ticket{
			Price:  -1,
			Status: TicketStatusValid,
		}
		assert.ErrorIs(t, ticket.Validate(), ErrPriceInvalid)
	})

	t.Run("should return error for invalid status", func(t *testing.T) {
		ticket := &Ticket{
			Price:  10.0,
			Status: "",
		}
		assert.ErrorIs(t, ticket.Validate(), ErrTicketInvalid)
	})
}

func TestTicket_Use(t *testing.T) {
	t.Run("should mark ticket as used", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusValid}
		err := ticket.Use()
		assert.NoError(t, err)
		assert.Equal(t, TicketStatusUsed, ticket.Status)
	})

	t.Run("should fail if already used", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusUsed}
		err := ticket.Use()
		assert.ErrorIs(t, err, ErrTicketAlreadyUsed)
	})

	t.Run("should fail if cancelled", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusCancelled}
		err := ticket.Use()
		assert.ErrorIs(t, err, ErrTicketAlreadyCanceled)
	})
}

func TestTicket_Cancel(t *testing.T) {
	t.Run("should mark ticket as cancelled", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusValid}
		err := ticket.Cancel()
		assert.NoError(t, err)
		assert.Equal(t, TicketStatusCancelled, ticket.Status)
	})

	t.Run("should fail if already used", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusUsed}
		err := ticket.Cancel()
		assert.ErrorIs(t, err, ErrTicketAlreadyUsed)
	})

	t.Run("should fail if already cancelled", func(t *testing.T) {
		ticket := &Ticket{Status: TicketStatusCancelled}
		err := ticket.Cancel()
		assert.ErrorIs(t, err, ErrTicketAlreadyCanceled)
	})
}
