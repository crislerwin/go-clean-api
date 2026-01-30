package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/crislerwin/go-clean-api/internal/infra/database/postgres"
	"github.com/crislerwin/go-clean-api/internal/infra/repository/model"
	"github.com/crislerwin/go-clean-api/internal/usecase"
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

	orderModel := model.NewOrderFromEntity(order)
	_, err := sqlx.NamedExecContext(ctx, executor, orderQuery, orderModel)
	if err != nil {
		return postgres.TranslateError(err)
	}

	for i := range order.Tickets {
		order.Tickets[i].OrderID = order.ID
	}
	ticketsModel := model.NewTicketsFromEntity(order.Tickets)

	_, err = sqlx.NamedExecContext(ctx, executor, ticketQuery, ticketsModel)
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

	var orderModels []*model.Order
	err := r.db.SelectContext(ctx, &orderModels, query, userID)
	if err != nil {
		return nil, postgres.TranslateError(err)
	}

	var orders []*entity.Order
	for _, om := range orderModels {
		orders = append(orders, om.ToEntity())
	}

	return orders, nil
}

func (r *OrderRepositorySQLx) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	queryOrder := `
		SELECT id, event_id, user_id, quantity, total_amount, status, created_at
		FROM orders
		WHERE id = $1
	`
	queryTickets := `
		SELECT id, event_id, order_id, price, status
		FROM tickets
		WHERE order_id = $1
	`

	executor := database.GetExecutor(ctx, r.db)

	var orderModel model.Order
	err := sqlx.GetContext(ctx, executor, &orderModel, queryOrder, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrOrderNotFound
		}
		return nil, postgres.TranslateError(err)
	}

	var ticketModels []model.Ticket
	err = sqlx.SelectContext(ctx, executor, &ticketModels, queryTickets, id)
	if err != nil {
		return nil, postgres.TranslateError(err)
	}

	order := orderModel.ToEntity()
	tickets := make([]entity.Ticket, len(ticketModels))
	for i, tm := range ticketModels {
		tickets[i] = entity.Ticket{
			ID:      tm.ID,
			EventID: tm.EventID,
			OrderID: tm.OrderID,
			Price:   tm.Price,
			Status:  entity.TicketStatus(tm.Status),
		}
	}
	order.Tickets = tickets

	return order, nil
}

func (r *OrderRepositorySQLx) Update(ctx context.Context, order *entity.Order) error {
	executor := database.GetExecutor(ctx, r.db)

	queryOrder := `
		UPDATE orders
		SET status = :status
		WHERE id = :id
	`

	orderModel := model.NewOrderFromEntity(order)
	_, err := sqlx.NamedExecContext(ctx, executor, queryOrder, orderModel)
	if err != nil {
		return postgres.TranslateError(err)
	}

	queryTickets := `
		UPDATE tickets
		SET status = :status
		WHERE id = :id
	`

	ticketsModel := model.NewTicketsFromEntity(order.Tickets)

	for _, t := range ticketsModel {
		_, err = sqlx.NamedExecContext(ctx, executor, queryTickets, t)
		if err != nil {
			return postgres.TranslateError(err)
		}
	}

	return nil
}
