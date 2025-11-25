package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserPreferences represents user preferences entity
type UserPreferences struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	Gender               string
	Min_Age              int
	Max_Age              int
	Bio                  string
	ProfilePictureURL    string
	NotificationSettings map[string]any
	LookingFor           string
	PreferredGender      string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// UserPreferencesRepository defines the interface for user preferences persistence
type UserPreferencesRepository interface {
	CreatePreferences(prefs *UserPreferences) error
	FindByUserID(userID uuid.UUID) (*UserPreferences, error)
	UpdatePreferences(prefs *UserPreferences) error
	DeletePreferences(id uuid.UUID) error
}
