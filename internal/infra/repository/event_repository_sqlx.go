package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/repository/model"
)

type EventRepositorySQLx struct {
	client database.Client
}

func NewEventRepositorySqlx(client database.Client) *EventRepositorySQLx {
	return &EventRepositorySQLx{client: client}
}

func (r *EventRepositorySQLx) GetTotalCapacity(ctx context.Context, eventID string) (int, error) {

	query := `
	SELECT capacity
	FROM events
	WHERE id = $1
	FOR UPDATE
	`

	var total int

	err := r.client.Get(ctx, &total, query, eventID)

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

	var count int

	err := r.client.Get(ctx, &count, query, eventID)

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
	eventModel := model.NewEventFromEntity(event)
	_, err := r.client.NamedExec(ctx, query, eventModel)
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

	var eventModel model.Event
	err := r.client.Get(ctx, &eventModel, query, eventID)
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
	var eventModels []*model.Event
	err := r.client.Select(ctx, &eventModels, query)
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
	var eventModels []*model.Event
	err := r.client.Select(ctx, &eventModels, query, userID)
	if err != nil {
		return nil, err
	}

	var events []*entity.Event
	for _, em := range eventModels {
		events = append(events, em.ToEntity())
	}
	return events, nil
}
