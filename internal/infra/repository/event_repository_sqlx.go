package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/repository/model"
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
	INSERT INTO events (id, user_id, name, location, organization, rating, date, capacity, price, image_url)
	VALUES (:id, :user_id, :name, :location, :organization, :rating, :date, :capacity, :price, :image_url)
	`
	executor := database.GetExecutor(ctx, r.db)
	eventModel := model.NewEventFromEntity(event)
	_, err := sqlx.NamedExecContext(ctx, executor, query, eventModel)
	return err
}

func (r *EventRepositorySQLx) GetByID(ctx context.Context, eventID string, forUpdate bool) (*entity.Event, error) {
	query := `
		SELECT id, user_id, name, location, organization, rating, date, capacity, price, image_url
		FROM events
		WHERE id = $1
	`
	if forUpdate {
		query += " FOR UPDATE"
	}

	executor := database.GetExecutor(ctx, r.db)
	var eventModel model.Event
	err := sqlx.GetContext(ctx, executor, &eventModel, query, eventID)
	if err != nil {
		return nil, err
	}
	return eventModel.ToEntity(), nil
}

func (r *EventRepositorySQLx) ListAll(ctx context.Context) ([]*entity.Event, error) {
	query := `
		SELECT id, user_id, name, location, organization, rating, date, capacity, price, image_url
		FROM events
	`
	executor := database.GetExecutor(ctx, r.db)
	var eventModels []*model.Event
	err := sqlx.SelectContext(ctx, executor, &eventModels, query)
	if err != nil {
		return nil, err
	}

	var events []*entity.Event
	for _, em := range eventModels {
		events = append(events, em.ToEntity())
	}
	return events, nil
}

func (r *EventRepositorySQLx) ListByUserID(ctx context.Context, userID string) ([]*entity.Event, error) {
	query := `
		SELECT id, user_id, name, location, organization, rating, date, capacity, price, image_url
		FROM events
		WHERE user_id = $1
	`
	executor := database.GetExecutor(ctx, r.db)
	var eventModels []*model.Event
	err := sqlx.SelectContext(ctx, executor, &eventModels, query, userID)
	if err != nil {
		return nil, err
	}

	var events []*entity.Event
	for _, em := range eventModels {
		events = append(events, em.ToEntity())
	}
	return events, nil
}
