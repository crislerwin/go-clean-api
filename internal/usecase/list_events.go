package usecase

import (
	"context"
	"time"
)

type ListEventsOutputDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	Organization string    `json:"organization"`
	Rating       string    `json:"rating"`
	Date         time.Time `json:"date"`
	ImageURL     string    `json:"image_url"`
	Capacity     int       `json:"capacity"`
	Price        float64   `json:"price"`
}

type ListEventsUseCase struct {
	eventRepo EventRepository
}

func NewListEventsUseCase(eventRepo EventRepository) *ListEventsUseCase {
	return &ListEventsUseCase{
		eventRepo: eventRepo,
	}
}

func (uc *ListEventsUseCase) Execute(ctx context.Context) ([]ListEventsOutputDTO, error) {
	events, err := uc.eventRepo.ListAll(ctx)
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
	// Ensure we return empty slice instead of null
	if output == nil {
		output = []ListEventsOutputDTO{}
	}

	return output, nil
}
