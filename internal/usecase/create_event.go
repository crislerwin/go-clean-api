package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventSoldOut  = errors.New("event sold out")
)

type CreateEventInputDTO struct {
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Rating       string  `json:"rating"`
	Date         string  `json:"date"` // Recebemos string (ISO8601) e convertemos
	Capacity     int     `json:"capacity"`
	Price        float64 `json:"price"`
	ImageURL     string  `json:"image_url"`
}

type CreateEventOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateEventUseCase struct {
	eventRepo EventRepository
}

func NewCreateEventUseCase(eventRepo EventRepository, txManager TransactionManager) *CreateEventUseCase {
	return &CreateEventUseCase{
		eventRepo: eventRepo,
	}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, input CreateEventInputDTO) (*CreateEventOutputDTO, error) {
	parsedDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return nil, err
	}

	event, err := entity.NewEvent(
		input.Name,
		input.Location,
		input.Organization,
		input.Rating,
		parsedDate,
		input.Capacity,
		input.Price,
		input.ImageURL,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return &CreateEventOutputDTO{
		ID:   event.ID.String(),
		Name: event.Name,
	}, nil
}
