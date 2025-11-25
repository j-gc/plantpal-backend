package persistence

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/j-gc/plantpal-backend/internal/modules/auth/domain"
)

// UserRepository implements the domain.UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password_hash = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
