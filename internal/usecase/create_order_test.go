package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Reuse MockEventRepository from create_event_test.go (if in same package)
// Since they are in the same package 'usecase' and folder, it is available.

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

type MockTransactionManager struct{}

func (m *MockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateOrderUseCase_Execute(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		input         OrderInputDTO
		setupMocks    func(*MockEventRepository, *MockOrderRepository)
		expectedError error
	}{
		{
			name: "Success",
			input: OrderInputDTO{
				EventID:  eventID.String(),
				UserID:   userID.String(),
				Quantity: 2,
			},
			setupMocks: func(eventRepo *MockEventRepository, orderRepo *MockOrderRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID.String()).Return(&entity.Event{
					ID:       eventID,
					Capacity: 10,
					Price:    10.0,
				}, nil)
				eventRepo.On("GetSoldTicketsCount", mock.Anything, eventID.String()).Return(5, nil)
				orderRepo.On("Save", mock.Anything, mock.MatchedBy(func(o *entity.Order) bool {
					return o.EventID == eventID && o.UserID == userID && len(o.Tickets) == 2
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Event Not Found",
			input: OrderInputDTO{
				EventID:  eventID.String(),
				UserID:   userID.String(),
				Quantity: 1,
			},
			setupMocks: func(eventRepo *MockEventRepository, orderRepo *MockOrderRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID.String()).Return(nil, errors.New("not found"))
			},
			expectedError: errors.New("not found"),
		},
		{
			name: "Event Sold Out",
			input: OrderInputDTO{
				EventID:  eventID.String(),
				UserID:   userID.String(),
				Quantity: 6,
			},
			setupMocks: func(eventRepo *MockEventRepository, orderRepo *MockOrderRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID.String()).Return(&entity.Event{
					ID:       eventID,
					Capacity: 10,
				}, nil)
				eventRepo.On("GetSoldTicketsCount", mock.Anything, eventID.String()).Return(5, nil)
			},
			expectedError: ErrEventSoldOut,
		},
		{
			name: "Invalid Quantity",
			input: OrderInputDTO{
				EventID:  eventID.String(),
				UserID:   userID.String(),
				Quantity: 0,
			},
			setupMocks: func(eventRepo *MockEventRepository, orderRepo *MockOrderRepository) {
				eventRepo.On("GetByID", mock.Anything, eventID.String()).Return(&entity.Event{
					ID:       eventID,
					Capacity: 10,
				}, nil)
			},
			expectedError: entity.ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventRepo := new(MockEventRepository)
			orderRepo := new(MockOrderRepository)
			txManager := &MockTransactionManager{}

			tt.setupMocks(eventRepo, orderRepo)

			useCase := NewCreateOrderUseCase(orderRepo, eventRepo, txManager)
			output, err := useCase.Execute(context.TODO(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, "PENDING", output.Status)
				assert.WithinDuration(t, time.Now(), output.CreatedAt, time.Second)
			}

			eventRepo.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
		})
	}
}
