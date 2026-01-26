package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txKey string

const transactionKey txKey = "tx"

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, transactionKey, tx)

	err = fn(ctxWithTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

func GetExecutor(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(transactionKey).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
