package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
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
	Description  string  `json:"description"`
}

type CreateEventOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateEventUseCase struct {
	eventRepo repository.EventRepository
}

func NewCreateEventUseCase(eventRepo repository.EventRepository, txManager TransactionManager) *CreateEventUseCase {
	return &CreateEventUseCase{
		eventRepo: eventRepo,
	}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, userID string, input CreateEventInputDTO) (*CreateEventOutputDTO, error) {
	parsedDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	event, err := entity.NewEvent(
		userUUID,
		input.Name,
		input.Location,
		input.Organization,
		input.Rating,
		parsedDate,
		input.Capacity,
		input.Price,
		input.ImageURL,
		input.Description,
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
