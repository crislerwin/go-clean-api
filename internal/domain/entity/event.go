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
	ID           uuid.UUID
	Name         string
	Location     string
	Organization string
	Rating       string // "Livre", "18+", etc
	Date         time.Time
	ImageURL     string
	Capacity     int
	Price        float64
	PartnerID    int // Quem criou o evento (se for multi-tenant)
}

// Factory para garantir integridade
func NewEvent(name, location, organization string, rating string, date time.Time, capacity int, price float64, imageURL string, partnerID int) (*Event, error) {
	if name == "" || capacity <= 0 || price < 0 {
		return nil, ErrInvalidEventData
	}

	if date.Before(time.Now()) {
		return nil, ErrDateInPast
	}

	return &Event{
		ID:           uuid.New(),
		Name:         name,
		Location:     location,
		Organization: organization,
		Rating:       rating,
		Date:         date,
		Capacity:     capacity,
		Price:        price,
		ImageURL:     imageURL,
		PartnerID:    partnerID,
	}, nil
}
