package usecase

import (
	"context"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"

	"github.com/google/uuid"
)

type OrderInputDTO struct {
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	Quantity int    `json:"quantity"`
}

type OrderOutputDTO struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
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

	var order *entity.Order

	err = uc.txManager.Do(ctx, func(ctxTx context.Context) error {
		// Lock the event row to prevent race conditions
		event, err := uc.eventRepo.GetByID(ctxTx, input.EventID, true)
		if err != nil {
			return err
		}

		// Check sold tickets count within the same transaction
		sold, err := uc.eventRepo.GetSoldTicketsCount(ctxTx, input.EventID)
		if err != nil {
			return err
		}

		if (sold + input.Quantity) > event.Capacity {
			return ErrEventSoldOut
		}

		// Create the order
		order, err = entity.NewOrder(eventID, userID, input.Quantity, event.Price)
		if err != nil {
			return err
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
