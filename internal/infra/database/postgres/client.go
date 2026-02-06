package postgres

import (
	"context"
	"database/sql"

	"github.com/crislerwin/go-clean-api/internal/infra/database"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(db *sqlx.DB) *Client {
	return &Client{db: db}
}

func (c *Client) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	executor := database.GetExecutor(ctx, c.db)
	return sqlx.GetContext(ctx, executor, dest, query, args...)
}

func (c *Client) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	executor := database.GetExecutor(ctx, c.db)
	return sqlx.SelectContext(ctx, executor, dest, query, args...)
}

func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	executor := database.GetExecutor(ctx, c.db)
	return executor.ExecContext(ctx, query, args...)
}

func (c *Client) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	executor := database.GetExecutor(ctx, c.db)
	return sqlx.NamedExecContext(ctx, executor, query, arg)
}
