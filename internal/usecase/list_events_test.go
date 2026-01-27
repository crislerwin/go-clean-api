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

type EventListRepoMock struct {
	mock.Mock
}

func (m *EventListRepoMock) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *EventListRepoMock) GetSoldTicketsCount(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *EventListRepoMock) Create(ctx context.Context, event *entity.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventListRepoMock) GetByID(ctx context.Context, eventID string, forUpdate bool) (*entity.Event, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *EventListRepoMock) ListAll(ctx context.Context) ([]*entity.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func (m *EventListRepoMock) ListByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Event), args.Error(1)
}

func TestListEventsUseCase_Execute(t *testing.T) {
	t.Run("should list events successfully", func(t *testing.T) {
		repo := new(EventListRepoMock)
		useCase := NewListEventsUseCase(repo)

		event, _ := entity.NewEvent(uuid.New(), "Event 1", "Location", "Org", "Livre", time.Now().Add(time.Hour), 100, 10.0, "img")
		events := []*entity.Event{event}

		repo.On("ListAll", mock.Anything).Return(events, nil)

		output, err := useCase.Execute(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, output, 1)
		assert.Equal(t, event.ID.String(), output[0].ID)
		repo.AssertExpectations(t)
	})

	t.Run("should return empty list when no events found", func(t *testing.T) {
		repo := new(EventListRepoMock)
		useCase := NewListEventsUseCase(repo)

		repo.On("ListAll", mock.Anything).Return([]*entity.Event{}, nil)

		output, err := useCase.Execute(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, output, 0)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repo := new(EventListRepoMock)
		useCase := NewListEventsUseCase(repo)

		repo.On("ListAll", mock.Anything).Return(([]*entity.Event)(nil), errors.New("db error"))

		output, err := useCase.Execute(context.TODO())

		assert.Error(t, err)
		assert.Nil(t, output)
		repo.AssertExpectations(t)
	})
}
