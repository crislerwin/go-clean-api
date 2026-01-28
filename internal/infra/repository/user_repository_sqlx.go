package repository

import (
	"context"

	"github.com/crislerwin/go-clean-api/internal/domain/entity"
	"github.com/crislerwin/go-clean-api/internal/infra/repository/model"
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
	query := `INSERT INTO users (id, name, email, password, role) VALUES (:id, :name, :email, :password, :role)`

	userModel := model.NewUserFromEntity(user)
	_, err := r.db.NamedExecContext(ctx, query, userModel)
	return err
}

func (r *UserRepositorySQLx) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE email = $1`
	var userModel model.User
	if err := r.db.GetContext(ctx, &userModel, query, email); err != nil {
		return nil, err
	}
	return userModel.ToEntity(), nil
}

func (r *UserRepositorySQLx) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE id = $1`
	var userModel model.User
	if err := r.db.GetContext(ctx, &userModel, query, id); err != nil {
		return nil, err
	}
	return userModel.ToEntity(), nil
}
