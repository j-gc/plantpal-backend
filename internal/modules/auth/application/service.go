package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/j-gc/plantpal-backend/internal/modules/auth/domain"
)

var (
	ErrEmailTaken    = errors.New("email already in use")
	ErrInvalidLogin  = errors.New("invalid email or password")
	ErrMissingFields = errors.New("missing required fields")
	ErrInvalidEmail  = errors.New("invalid email format")
)

// PasswordHasher handles password hashing and comparison
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
}

// TokenIssuer issues JWT tokens
type TokenIssuer interface {
	Issue(subject string, ttl time.Duration, claims map[string]any) (string, error)
}

// Service handles user-related business logic
type Service struct {
	repo   domain.UserRepository
	hasher PasswordHasher
	issuer TokenIssuer
}

// NewService creates a new user service
func NewService(repo domain.UserRepository, hasher PasswordHasher, issuer TokenIssuer) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
		issuer: issuer,
	}
}

// RegisterInput represents registration input
type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

// RegisterOutput represents registration output
type RegisterOutput struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
	// Normalize email
	email := strings.TrimSpace(strings.ToLower(in.Email))
	firstName := strings.TrimSpace(in.FirstName)
	lastName := strings.TrimSpace(in.LastName)

	// Validate
	if email == "" || firstName == "" || in.Password == "" {
		return nil, ErrMissingFields
	}

	// Check if email already exists
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	// Hash password
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		FirstName:    firstName,
		LastName:     lastName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &RegisterOutput{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// LoginInput represents login input
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginOutput represents login output
type LoginOutput struct {
	AccessToken string   `json:"access_token"`
	User        UserInfo `json:"user"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	// Normalize email
	email := strings.TrimSpace(strings.ToLower(in.Email))

	// Find user
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidLogin
	}

	// Verify password
	if err := s.hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return nil, ErrInvalidLogin
	}

	// Issue token
	token, err := s.issuer.Issue(user.ID.String(), 24*time.Hour, map[string]any{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken: token,
		User: UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(uid)
}
