package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name         string
		eventName    string
		location     string
		organization string
		rating       string
		date         time.Time
		capacity     int
		price        float64
		imageURL     string
		partnerID    int
		expectedErr  error
	}{
		{
			name:         "should create a valid event",
			eventName:    "Rock in Rio",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         futureDate,
			capacity:     100000,
			price:        100.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  nil,
		},
		{
			name:         "should return error when name is empty",
			eventName:    "",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         futureDate,
			capacity:     100000,
			price:        100.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  ErrInvalidEventData,
		},
		{
			name:         "should return error when capacity is zero",
			eventName:    "Rock in Rio",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         futureDate,
			capacity:     0,
			price:        100.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  ErrInvalidEventData,
		},
		{
			name:         "should return error when capacity is negative",
			eventName:    "Rock in Rio",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         futureDate,
			capacity:     -1,
			price:        100.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  ErrInvalidEventData,
		},
		{
			name:         "should return error when price is negative",
			eventName:    "Rock in Rio",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         futureDate,
			capacity:     100000,
			price:        -10.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  ErrInvalidEventData,
		},
		{
			name:         "should return error when date is in the past",
			eventName:    "Rock in Rio",
			location:     "Rio de Janeiro",
			organization: "Live Nation",
			rating:       "Livre",
			date:         pastDate,
			capacity:     100000,
			price:        100.0,
			imageURL:     "http://example.com/image.jpg",
			partnerID:    1,
			expectedErr:  ErrDateInPast,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewEvent(tt.eventName, tt.location, tt.organization, tt.rating, tt.date, tt.capacity, tt.price, tt.imageURL, tt.partnerID)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.eventName, event.Name)
				assert.Equal(t, tt.location, event.Location)
				assert.Equal(t, tt.organization, event.Organization)
				assert.Equal(t, tt.rating, event.Rating)
				assert.Equal(t, tt.date, event.Date)
				assert.Equal(t, tt.capacity, event.Capacity)
				assert.Equal(t, tt.price, event.Price)
				assert.Equal(t, tt.imageURL, event.ImageURL)
				assert.Equal(t, tt.partnerID, event.PartnerID)
				assert.NotEmpty(t, event.ID)
			}
		})
	}
}
