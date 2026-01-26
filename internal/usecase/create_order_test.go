package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type OrderRepositoryMock struct {
	SaveCalled bool
}

type TransactionManagerMock struct{}

func (m *TransactionManagerMock) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
func (m *OrderRepositoryMock) Save(ctx context.Context, order *entity.Order) error {
	m.SaveCalled = true
	return nil
}

type EventRepositoryMock struct {
	Capacity int
	Sold     int
}

func (m *EventRepositoryMock) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {
	return m.Capacity, nil
}
func (m *EventRepositoryMock) GetSoldTicketsCount(ctx context.Context, eventID string) (int, error) {
	return m.Sold, nil
}
func TestCreateOrderUseCase_Execute(t *testing.T) {
	t.Run("should create order successfully when stock is available", func(t *testing.T) {
		// Arrange
		orderRepo := &OrderRepositoryMock{}
		eventRepo := &EventRepositoryMock{Capacity: 10, Sold: 5}
		txManager := &TransactionManagerMock{}
		useCase := NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
		input := OrderInputDTO{
			EventID:        uuid.New().String(),
			UserID:         uuid.New().String(),
			Quantity:       2,
			PricePerTicket: 100.0,
		}

		// Act
		output, err := useCase.Execute(context.Background(), input)

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "PENDING", output.Status)
		assert.True(t, orderRepo.SaveCalled)
	})
	t.Run("should return error when event is sold out", func(t *testing.T) {
		// Arrange
		orderRepo := &OrderRepositoryMock{}

		eventRepo := &EventRepositoryMock{Capacity: 10, Sold: 10}

		txManager := &TransactionManagerMock{}
		useCase := NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
		input := OrderInputDTO{
			EventID:        uuid.New().String(),
			UserID:         uuid.New().String(),
			Quantity:       1,
			PricePerTicket: 2,
		}
		output, err := useCase.Execute(context.Background(), input)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrEventSoldOut))
		assert.Nil(t, output)
		assert.False(t, orderRepo.SaveCalled)
	})
}
