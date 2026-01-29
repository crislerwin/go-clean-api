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
	Description  *string   `db:"description"`
}

func NewEventFromEntity(e *entity.Event) *Event {
	var desc *string
	if e.Description != "" {
		desc = &e.Description
	}
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
		Description:  desc,
	}
}

func (e *Event) ToEntity() *entity.Event {
	desc := ""
	if e.Description != nil {
		desc = *e.Description
	}
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
		Description:  desc,
	}
}
