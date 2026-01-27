package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEventData = errors.New("invalid event data")
	ErrDateInPast       = errors.New("event date must be in the future")
)

type Event struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	Name         string    `db:"name"`
	Location     string    `db:"location"`
	Organization string    `db:"organization"`
	Rating       string    `db:"rating"` // "Livre", "18+", etc
	Date         time.Time `db:"date"`
	ImageURL     string    `db:"image_url"`
	Capacity     int       `db:"capacity"`
	Price        float64   `db:"price"`
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
