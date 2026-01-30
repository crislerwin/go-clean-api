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

// MockOrderRepositoryUpdate avoids conflict with other test files
type MockOrderRepositoryUpdate struct {
	mock.Mock
}

func (m *MockOrderRepositoryUpdate) Save(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepositoryUpdate) GetByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *MockOrderRepositoryUpdate) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepositoryUpdate) Update(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestUpdateOrderStatusUseCase_Execute(t *testing.T) {
	orderID := uuid.New()
	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		input         UpdateOrderStatusInputDTO
		setupMocks    func(*MockOrderRepositoryUpdate)
		expectedError error
	}{
		{
			name: "Success - Paid",
			input: UpdateOrderStatusInputDTO{
				OrderID: orderID.String(),
				Status:  "PAID",
			},
			setupMocks: func(repo *MockOrderRepositoryUpdate) {
				order := &entity.Order{
					ID:        orderID,
					EventID:   eventID,
					UserID:    userID,
					Status:    entity.OrderStatusPending,
					Tickets:   []entity.Ticket{},
					CreatedAt: time.Now(),
				}
				repo.On("GetByID", mock.Anything, orderID.String()).Return(order, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(o *entity.Order) bool {
					return o.ID == orderID && o.Status == entity.OrderStatusPaid
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success - Rejected (Cancel Tickets)",
			input: UpdateOrderStatusInputDTO{
				OrderID: orderID.String(),
				Status:  "REJECTED",
			},
			setupMocks: func(repo *MockOrderRepositoryUpdate) {
				ticket1, _ := entity.NewTicket(eventID, orderID, 10.0)
				ticket2, _ := entity.NewTicket(eventID, orderID, 10.0)

				order := &entity.Order{
					ID:        orderID,
					EventID:   eventID,
					UserID:    userID,
					Status:    entity.OrderStatusPending,
					Tickets:   []entity.Ticket{*ticket1, *ticket2},
					CreatedAt: time.Now(),
				}
				repo.On("GetByID", mock.Anything, orderID.String()).Return(order, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(o *entity.Order) bool {
					return o.ID == orderID &&
						o.Status == entity.OrderStatusRejected &&
						o.Tickets[0].Status == entity.TicketStatusCancelled &&
						o.Tickets[1].Status == entity.TicketStatusCancelled
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Order Not Found",
			input: UpdateOrderStatusInputDTO{
				OrderID: orderID.String(),
				Status:  "PAID",
			},
			setupMocks: func(repo *MockOrderRepositoryUpdate) {
				repo.On("GetByID", mock.Anything, orderID.String()).Return(nil, errors.New("not found"))
			},
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepositoryUpdate)
			tt.setupMocks(repo)

			useCase := NewUpdateOrderStatusUseCase(repo)
			err := useCase.Execute(context.TODO(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
