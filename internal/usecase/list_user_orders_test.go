package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListUserOrdersUseCase_Execute(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockEventRepo := new(MockEventRepository)
	useCase := NewListUserOrdersUseCase(mockOrderRepo, mockEventRepo)

	ctx := context.Background()
	userID := "user-123"
	eventID := uuid.New()
	orders := []*entity.Order{
		{
			ID:          uuid.New(),
			EventID:     eventID,
			UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"), // Simplify
			TotalAmount: 100.0,
			Quantity:    2,
			Status:      "PENDING",
			CreatedAt:   time.Now(),
		},
	}
	event := &entity.Event{
		ID:   eventID,
		Name: "Concert",
		Date: time.Now().Add(24 * time.Hour),
	}

	mockOrderRepo.On("GetByUserID", ctx, userID).Return(orders, nil)
	mockEventRepo.On("GetByID", ctx, eventID.String(), false).Return(event, nil)

	output, err := useCase.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, output, 1)
	assert.Equal(t, "Concert", output[0].EventName)
	assert.Equal(t, 100.0, output[0].TotalAmount)

	mockOrderRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}
