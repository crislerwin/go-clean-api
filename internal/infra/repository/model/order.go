package model

import (
	"database/sql"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID      `db:"id"`
	EventID     uuid.UUID      `db:"event_id"`
	UserID      uuid.UUID      `db:"user_id"`
	TotalAmount float64        `db:"total_amount"`
	Quantity    int            `db:"quantity"`
	Status      string         `db:"status"`
	Reason      sql.NullString `db:"reason"`
	CreatedAt   time.Time      `db:"created_at"`
}

type Ticket struct {
	ID      uuid.UUID `db:"id"`
	EventID uuid.UUID `db:"event_id"`
	OrderID uuid.UUID `db:"order_id"`
	Price   float64   `db:"price"`
	Status  string    `db:"status"`
}

func NewOrderFromEntity(o *entity.Order) *Order {
	return &Order{
		ID:          o.ID,
		EventID:     o.EventID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Quantity:    o.Quantity,
		Status:      string(o.Status),
		Reason:      sql.NullString{String: o.Reason, Valid: o.Reason != ""},
		CreatedAt:   o.CreatedAt,
	}
}

func NewTicketsFromEntity(tickets []entity.Ticket) []Ticket {
	modelTickets := make([]Ticket, len(tickets))
	for i, t := range tickets {
		modelTickets[i] = Ticket{
			ID:      t.ID,
			EventID: t.EventID,
			OrderID: t.OrderID,
			Price:   t.Price,
			Status:  string(t.Status),
		}
	}
	return modelTickets
}

func (o *Order) ToEntity() *entity.Order {
	return &entity.Order{
		ID:          o.ID,
		EventID:     o.EventID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Quantity:    o.Quantity,
		Status:      entity.OrderStatus(o.Status),
		Reason:      o.Reason.String,
		CreatedAt:   o.CreatedAt,
		// Note: Tickets are usually loaded separately or joined.
		// For basic mapping, we might leave Tickets empty here if not loaded.
	}
}
