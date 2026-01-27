package usecase

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	GetByUserID(ctx context.Context, userID string) ([]*entity.Order, error)
}

type EventRepository interface {
	GetTotalCapacity(ctx context.Context, eventID string) (int, error)
	GetSoldTicketsCount(ctx context.Context, eventID string) (int, error)
	Create(ctx context.Context, event *entity.Event) error
	GetByID(ctx context.Context, eventID string) (*entity.Event, error)
	ListAll(ctx context.Context) ([]*entity.Event, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Event, error)
}

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
}
