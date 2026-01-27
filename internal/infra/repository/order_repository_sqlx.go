package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/database/postgres"
	"github.com/jmoiron/sqlx"
)

type OrderRepositorySQLx struct {
	db *sqlx.DB
}

func NewOrderRepositorySQLx(db *sqlx.DB) *OrderRepositorySQLx {
	return &OrderRepositorySQLx{db: db}
}

func (r *OrderRepositorySQLx) Save(ctx context.Context, order *entity.Order) error {
	executor := database.GetExecutor(ctx, r.db)
	orderQuery := `
		INSERT INTO orders (id, event_id, user_id, quantity, total_amount, status, created_at)
		VALUES (:id, :event_id, :user_id, :quantity, :total_amount, :status, :created_at)
		`
	ticketQuery := `
		INSERT INTO tickets (id, event_id, order_id, price, status)
		VALUES (:id, :event_id, :order_id, :price, :status)
		`

	_, err := sqlx.NamedExecContext(ctx, executor, orderQuery, order)
	if err != nil {
		return postgres.TranslateError(err)
	}

	for i := range order.Tickets {
		order.Tickets[i].OrderID = order.ID
	}

	_, err = sqlx.NamedExecContext(ctx, executor, ticketQuery, order.Tickets)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositorySQLx) GetByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	query := `
		SELECT id, event_id, user_id, quantity, total_amount, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var orders []*entity.Order
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, postgres.TranslateError(err)
	}

	return orders, nil
}
