package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create a new valid user", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "password123")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Password)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		user, err := NewUser("", "john@example.com", "password123")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "name is required", err.Error())
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		user, err := NewUser("John Doe", "", "password123")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "password is required", err.Error())
	})

	t.Run("should return error when password is too short", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "123")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "password must be at least 6 characters", err.Error())
	})
}
