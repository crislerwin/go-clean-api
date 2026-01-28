package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEventData = errors.New("invalid event data")
	ErrDateInPast       = errors.New("event date must be in the future")
	ErrEventSoldOut     = errors.New("event sold out")
)

type Event struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	Location     string
	Organization string
	Rating       string // "Livre", "18+", etc
	Date         time.Time
	ImageURL     string
	Capacity     int
	Price        float64
}

// Factory para garantir integridade.
func NewEvent(userID uuid.UUID, name, location, organization string, rating string, date time.Time, capacity int, price float64, imageURL string) (*Event, error) {
	if name == "" || capacity <= 0 || price < 0 {
		return nil, ErrInvalidEventData
	}

	if date.Before(time.Now()) {
		return nil, ErrDateInPast
	}

	return &Event{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		Location:     location,
		Organization: organization,
		Rating:       rating,
		Date:         date,
		Capacity:     capacity,
		Price:        price,
		ImageURL:     imageURL,
	}, nil
}

func (e *Event) CanSell(quantity int, soldTickets int) error {
	if (soldTickets + quantity) > e.Capacity {
		return ErrEventSoldOut
	}
	return nil
}
