package model

import (
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	Name         string    `db:"name"`
	Location     string    `db:"location"`
	Organization string    `db:"organization"`
	Rating       string    `db:"rating"`
	Date         time.Time `db:"date"`
	ImageURL     string    `db:"image_url"`
	Capacity     int       `db:"capacity"`
	Price        float64   `db:"price"`
}

func NewEventFromEntity(e *entity.Event) *Event {
	return &Event{
		ID:           e.ID,
		UserID:       e.UserID,
		Name:         e.Name,
		Location:     e.Location,
		Organization: e.Organization,
		Rating:       e.Rating,
		Date:         e.Date,
		ImageURL:     e.ImageURL,
		Capacity:     e.Capacity,
		Price:        e.Price,
	}
}

func (e *Event) ToEntity() *entity.Event {
	return &entity.Event{
		ID:           e.ID,
		UserID:       e.UserID,
		Name:         e.Name,
		Location:     e.Location,
		Organization: e.Organization,
		Rating:       e.Rating,
		Date:         e.Date,
		ImageURL:     e.ImageURL,
		Capacity:     e.Capacity,
		Price:        e.Price,
	}
}
