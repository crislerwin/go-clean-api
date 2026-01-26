package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"

	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventSoldOut  = errors.New("event sold out")
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

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
}

type EventRepository interface {
	GetTotalCapacity(ctx context.Context, eventID string) (int, error)
	GetSoldTicketsCount(ctx context.Context, eventID string) (int, error)
}

type CreateOrderUseCase struct {
	orderRepo OrderRepository
	eventRepo EventRepository
	txManager TransactionManager
}

func NewCreateOrderUseCase(
	orderRepo OrderRepository,
	eventRepo EventRepository,
	txManager TransactionManager,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo: orderRepo,
		eventRepo: eventRepo,
		txManager: txManager,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input OrderInputDTO) (*OrderOutputDTO, error) {
	eventID, err := uuid.Parse(input.EventID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, err
	}
	order, err := entity.NewOrder(eventID, userID, input.Quantity, input.PricePerTicket)
	if err != nil {
		return nil, err
	}

	err = uc.txManager.Do(ctx, func(ctxTx context.Context) error {
		total, err := uc.eventRepo.GetTotalCapacity(ctx, input.EventID)
		if err != nil {
			return err
		}
		sold, err := uc.eventRepo.GetSoldTicketsCount(ctxTx, input.EventID)

		if err != nil {
			return err
		}
		if (sold + input.Quantity) > total {
			return ErrEventSoldOut
		}

		if err := uc.orderRepo.Save(ctxTx, order); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &OrderOutputDTO{
		ID:          order.ID.String(),
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}, nil
}
