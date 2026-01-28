package model

import (
	"github.com/crislerwin/go-clean-api/internal/domain/entity"
)

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

func NewUserFromEntity(u *entity.User) *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
}
