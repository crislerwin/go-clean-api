package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository mocks the EventRepository interface.
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) GetSoldTicketsCount(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) Create(ctx context.Context, event *entity.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, eventID string) (*entity.Event, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventRepository) ListAll(ctx context.Context) ([]*entity.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *MockEventRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func TestCreateEventUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateEventInputDTO
		setupRepo     func(*MockEventRepository)
		expectedError error
		expectedID    bool // true if ID check is required
	}{
		{
			name: "Success",
			input: CreateEventInputDTO{
				Name:         "Concert",
				Location:     "Stadium",
				Organization: "Org",
				Rating:       "Livre",
				Date:         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Capacity:     100,
				Price:        50.0,
				ImageURL:     "url",
			},
			setupRepo: func(repo *MockEventRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(e *entity.Event) bool {
					return e.Name == "Concert"
				})).Return(nil)
			},
			expectedError: nil,
			expectedID:    true,
		},
		{
			name: "Invalid Date",
			input: CreateEventInputDTO{
				Name: "Concert",
				Date: "invalid-date",
			},
			setupRepo:     func(repo *MockEventRepository) {},
			expectedError: errors.New("parsing time"), // Approximate check
			expectedID:    false,
		},
		{
			name: "Invalid Input (Past Date)",
			input: CreateEventInputDTO{
				Name:         "Concert",
				Location:     "Stadium",
				Organization: "Org",
				Rating:       "Livre",
				Date:         time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				Capacity:     100,
				Price:        50.0,
				ImageURL:     "url",
			},
			setupRepo:     func(repo *MockEventRepository) {},
			expectedError: entity.ErrDateInPast,
			expectedID:    false,
		},
		{
			name: "Repository Error",
			input: CreateEventInputDTO{
				Name:         "Concert",
				Location:     "Stadium",
				Organization: "Org",
				Rating:       "Livre",
				Date:         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Capacity:     100,
				Price:        50.0,
				ImageURL:     "url",
			},
			setupRepo: func(repo *MockEventRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
			expectedID:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			tt.setupRepo(mockRepo)

			useCase := NewCreateEventUseCase(mockRepo, nil) // TransactionManager unused in CreateEvent
			userID := "5334d5d7-e5d8-4d56-9257-2b7b5e5d3c8a"
			output, err := useCase.Execute(context.TODO(), userID, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.name == "Invalid Input (Past Date)" {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				if tt.expectedID {
					assert.NotEmpty(t, output.ID)
					assert.Equal(t, tt.input.Name, output.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
