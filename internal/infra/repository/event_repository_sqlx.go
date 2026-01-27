package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/jmoiron/sqlx"
)

type EventRepositorySQLx struct {
	db *sqlx.DB
}

func NewEventRepositorySqlx(db *sqlx.DB) *EventRepositorySQLx {
	return &EventRepositorySQLx{db: db}
}

func (r *EventRepositorySQLx) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {

	query := `
	SELECT capacity
	FROM events
	WHERE id = $1
	FOR UPDATE
	`

	executor := database.GetExecutor(ctx, r.db)

	var total int

	err := sqlx.GetContext(ctx, executor, &total, query, eventID)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *EventRepositorySQLx) GetSoldTicketsCount(ctx context.Context, eventID string) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM tickets t
	JOIN orders o ON t.order_id = o.id
	WHERE t.event_id = $1
	AND o.status IN ('PAID', 'PENDING')
	`

	executor := database.GetExecutor(ctx, r.db)

	var count int

	err := sqlx.GetContext(ctx, executor, &count, query, eventID)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *EventRepositorySQLx) Create(ctx context.Context, event *entity.Event) error {
	query := `
	INSERT INTO events (id, name, location, organization, rating, date, capacity, price, image_url)
	VALUES (:id, :name, :location, :organization, :rating, :date, :capacity, :price, :image_url)
	`

	executor := database.GetExecutor(ctx, r.db)

	_, err := sqlx.NamedExecContext(ctx, executor, query, event)

	return err
}

func (r *EventRepositorySQLx) GetByID(ctx context.Context, eventID string) (*entity.Event, error) {
	query := `
		SELECT id, name, location, organization, rating, date, capacity, price, image_url
		FROM events
		WHERE id = $1
	`
	executor := database.GetExecutor(ctx, r.db)
	var event entity.Event
	err := sqlx.GetContext(ctx, executor, &event, query, eventID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepositorySQLx) ListAll(ctx context.Context) ([]*entity.Event, error) {
	query := `
		SELECT id, name, location, organization, rating, date, capacity, price, image_url
		FROM events
	`
	executor := database.GetExecutor(ctx, r.db)
	var events []*entity.Event
	err := sqlx.SelectContext(ctx, executor, &events, query)
	if err != nil {
		return nil, err
	}
	return events, nil
}
