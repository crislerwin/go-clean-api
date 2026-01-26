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

	event, err := uc.eventRepo.GetByID(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	order, err := entity.NewOrder(eventID, userID, input.Quantity, event.Price)
	if err != nil {
		return nil, err
	}

	err = uc.txManager.Do(ctx, func(ctxTx context.Context) error {
		// We already have the event matching the ID, and we can check capacity directly if we trust the fetched event.
		// However, for consistency with current logic which might use checking capacity:
		// But GetByID returns the event struct which has Capacity.
		// Optimization: Use event.Capacity directly instead of calling GetTotalCapacity again.

		sold, err := uc.eventRepo.GetSoldTicketsCount(ctxTx, input.EventID)
		if err != nil {
			return err
		}

		if (sold + input.Quantity) > event.Capacity {
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
