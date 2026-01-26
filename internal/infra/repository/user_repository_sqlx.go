package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepositorySQLx struct {
	db *sqlx.DB
}

func NewUserRepositorySQLx(db *sqlx.DB) *UserRepositorySQLx {
	return &UserRepositorySQLx{
		db: db,
	}
}

func (r *UserRepositorySQLx) Save(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password)
	return err
}
