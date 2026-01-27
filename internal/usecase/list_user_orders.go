package usecase

import (
	"context"
	"time"
)

type ListOrdersOutputDTO struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	EventName   string    `json:"event_name"`
	EventDate   time.Time `json:"event_date"`
	Quantity    int       `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListUserOrdersUseCase struct {
	orderRepo OrderRepository
	eventRepo EventRepository
}

func NewListUserOrdersUseCase(orderRepo OrderRepository, eventRepo EventRepository) *ListUserOrdersUseCase {
	return &ListUserOrdersUseCase{
		orderRepo: orderRepo,
		eventRepo: eventRepo,
	}
}

func (uc *ListUserOrdersUseCase) Execute(ctx context.Context, userID string) ([]ListOrdersOutputDTO, error) {
	orders, err := uc.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var output []ListOrdersOutputDTO
	for _, order := range orders {
		var eventName string
		var eventDate time.Time

		// Enricht with Event details
		// NOTE: This N+1 query pattern is acceptable for now given assumed low volume of user orders.
		// For high volume, we should join in repo or use DataLoader/Batching.
		event, err := uc.eventRepo.GetByID(ctx, order.EventID.String())
		if err == nil && event != nil { // Best effort enrichment
			eventName = event.Name
			eventDate = event.Date
		} else {
			// If event not found or deleted, we still show the order but with placeholder/empty details
			// Alternatively, could log warning.
			eventName = "Unknown Event"
		}

		output = append(output, ListOrdersOutputDTO{
			ID:          order.ID.String(),
			EventID:     order.EventID.String(),
			EventName:   eventName,
			EventDate:   eventDate,
			Quantity:    order.Quantity,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		})
	}

	return output, nil
}
