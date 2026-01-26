package usecase

import (
	"context"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

type OrderInputDTO struct {
	EventID        string  `json:"event_id"`
	UserID         string  `json:"user_id"`
	Quantity       int     `json:"quantity"`
	PricePerTicket float64 `json:"price_per_ticket"`
}

type OrderOutputDTO struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
}

type EventRepository interface {
	GetTotalCapacity(ctx context.Context, eventID string) (int, error)
	GetSoldTicketsCount(ctx context.Context, eventID string) (int, error)
}
