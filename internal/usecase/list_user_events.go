package usecase

import (
	"context"
)

type ListUserEventsUseCase struct {
	eventRepo EventRepository
}

func NewListUserEventsUseCase(eventRepo EventRepository) *ListUserEventsUseCase {
	return &ListUserEventsUseCase{
		eventRepo: eventRepo,
	}
}

func (uc *ListUserEventsUseCase) Execute(ctx context.Context, userID string) ([]ListEventsOutputDTO, error) {
	events, err := uc.eventRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var output []ListEventsOutputDTO
	for _, event := range events {
		output = append(output, ListEventsOutputDTO{
			ID:           event.ID.String(),
			Name:         event.Name,
			Location:     event.Location,
			Organization: event.Organization,
			Rating:       event.Rating,
			Date:         event.Date,
			ImageURL:     event.ImageURL,
			Capacity:     event.Capacity,
			Price:        event.Price,
		})
	}

	if output == nil {
		output = []ListEventsOutputDTO{}
	}

	return output, nil
}
