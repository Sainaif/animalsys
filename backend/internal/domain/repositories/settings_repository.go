package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingsRepository defines the interface for settings data access
type SettingsRepository interface {
	// Get retrieves the foundation settings
	// There should only be one settings document in the system
	Get(ctx context.Context) (*entities.FoundationSettings, error)

	// GetByID retrieves settings by ID
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.FoundationSettings, error)

	// Update updates the foundation settings
	// Uses optimistic locking via version field
	Update(ctx context.Context, settings *entities.FoundationSettings) error

	// Create creates initial foundation settings
	// Should only be called once during system initialization
	Create(ctx context.Context, settings *entities.FoundationSettings) error

	// UpdateEmailSettings updates only email settings
	UpdateEmailSettings(ctx context.Context, emailSettings entities.EmailSettings, updatedBy primitive.ObjectID) error

	// UpdateNotificationSettings updates only notification settings
	UpdateNotificationSettings(ctx context.Context, notificationSettings entities.NotificationSettings, updatedBy primitive.ObjectID) error

	// UpdateFeatureFlags updates only feature flags
	UpdateFeatureFlags(ctx context.Context, features entities.FeatureFlags, updatedBy primitive.ObjectID) error

	// UpdateBranding updates only branding settings
	UpdateBranding(ctx context.Context, branding entities.Branding, updatedBy primitive.ObjectID) error

	// GetContactInfo returns only contact information
	GetContactInfo(ctx context.Context) (*entities.ContactDetails, error)

	// GetOperatingHours returns operating hours
	GetOperatingHours(ctx context.Context) (map[string]entities.OperatingHour, error)

	// EnsureIndexes creates necessary indexes for the settings collection
	EnsureIndexes(ctx context.Context) error
}
