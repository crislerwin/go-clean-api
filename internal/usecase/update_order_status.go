package usecase

import (
	"context"
	"errors"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/domain/repository"
)

var ErrOrderNotFound = errors.New("order not found")

type UpdateOrderStatusInputDTO struct {
	OrderID string
	Status  string
	Reason  string
}

type UpdateOrderStatusUseCase struct {
	OrderRepo repository.OrderRepository
}

func NewUpdateOrderStatusUseCase(orderRepo repository.OrderRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		OrderRepo: orderRepo,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, input UpdateOrderStatusInputDTO) error {
	order, err := uc.OrderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return err
	}

	switch input.Status {
	case "PAID":
		order.Status = entity.OrderStatusPaid
	case "REJECTED":
		order.Status = entity.OrderStatusRejected
		order.Reason = input.Reason
		for i := range order.Tickets {
			if err := order.Tickets[i].Cancel(); err != nil {
				return err
			}
		}
	}

	return uc.OrderRepo.Update(ctx, order)
}
