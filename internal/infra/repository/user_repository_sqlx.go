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
	query := `INSERT INTO users (id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *UserRepositorySQLx) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE email = $1`
	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}
