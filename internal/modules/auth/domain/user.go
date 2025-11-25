package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository defines the interface for user persistence
type UserRepository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}
