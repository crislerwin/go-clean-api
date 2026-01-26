package repository

import (
	"context"

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
	SELECT total_capacity
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
