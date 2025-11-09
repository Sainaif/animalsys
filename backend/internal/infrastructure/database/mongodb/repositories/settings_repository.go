package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type settingsRepository struct {
	db *mongodb.Database
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *mongodb.Database) repositories.SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Settings)
}

// Get retrieves the foundation settings
func (r *settingsRepository) Get(ctx context.Context) (*entities.FoundationSettings, error) {
	var settings entities.FoundationSettings

	// Get the first (and should be only) settings document
	err := r.collection().FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("Settings not found")
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return &settings, nil
}

// GetByID retrieves settings by ID
func (r *settingsRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.FoundationSettings, error) {
	var settings entities.FoundationSettings

	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("Settings not found")
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return &settings, nil
}

// Create creates initial foundation settings
func (r *settingsRepository) Create(ctx context.Context, settings *entities.FoundationSettings) error {
	if settings.ID.IsZero() {
		settings.ID = primitive.NewObjectID()
	}

	now := time.Now()
	settings.CreatedAt = now
	settings.UpdatedAt = now
	settings.Version = 1

	_, err := r.collection().InsertOne(ctx, settings)
	if err != nil {
		return fmt.Errorf("failed to create settings: %w", err)
	}

	return nil
}

// Update updates the foundation settings with optimistic locking
func (r *settingsRepository) Update(ctx context.Context, settings *entities.FoundationSettings) error {
	oldVersion := settings.Version
	settings.UpdateVersion(settings.UpdatedBy)

	// Use optimistic locking - only update if version matches
	filter := bson.M{
		"_id":     settings.ID,
		"version": oldVersion,
	}

	update := bson.M{
		"$set": settings,
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewConflict("Settings have been modified by another user. Please refresh and try again.")
	}

	return nil
}

// UpdateEmailSettings updates only email settings
func (r *settingsRepository) UpdateEmailSettings(ctx context.Context, emailSettings entities.EmailSettings, updatedBy primitive.ObjectID) error {
	filter := bson.M{}

	update := bson.M{
		"$set": bson.M{
			"email_settings": emailSettings,
			"updated_by":     updatedBy,
			"updated_at":     time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update email settings: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Settings not found")
	}

	return nil
}

// UpdateNotificationSettings updates only notification settings
func (r *settingsRepository) UpdateNotificationSettings(ctx context.Context, notificationSettings entities.NotificationSettings, updatedBy primitive.ObjectID) error {
	filter := bson.M{}

	update := bson.M{
		"$set": bson.M{
			"notification_settings": notificationSettings,
			"updated_by":            updatedBy,
			"updated_at":            time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update notification settings: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Settings not found")
	}

	return nil
}

// UpdateFeatureFlags updates only feature flags
func (r *settingsRepository) UpdateFeatureFlags(ctx context.Context, features entities.FeatureFlags, updatedBy primitive.ObjectID) error {
	filter := bson.M{}

	update := bson.M{
		"$set": bson.M{
			"features":   features,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update feature flags: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Settings not found")
	}

	return nil
}

// UpdateBranding updates only branding settings
func (r *settingsRepository) UpdateBranding(ctx context.Context, branding entities.Branding, updatedBy primitive.ObjectID) error {
	filter := bson.M{}

	update := bson.M{
		"$set": bson.M{
			"branding":   branding,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update branding: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Settings not found")
	}

	return nil
}

// GetContactInfo returns only contact information
func (r *settingsRepository) GetContactInfo(ctx context.Context) (*entities.ContactDetails, error) {
	var result struct {
		ContactInfo entities.ContactDetails `bson:"contact_info"`
	}

	opts := options.FindOne().SetProjection(bson.M{"contact_info": 1})
	err := r.collection().FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("Settings not found")
		}
		return nil, fmt.Errorf("failed to get contact info: %w", err)
	}

	return &result.ContactInfo, nil
}

// GetOperatingHours returns operating hours
func (r *settingsRepository) GetOperatingHours(ctx context.Context) (map[string]entities.OperatingHour, error) {
	var result struct {
		OperatingHours map[string]entities.OperatingHour `bson:"operating_hours"`
	}

	opts := options.FindOne().SetProjection(bson.M{"operating_hours": 1})
	err := r.collection().FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("Settings not found")
		}
		return nil, fmt.Errorf("failed to get operating hours: %w", err)
	}

	return result.OperatingHours, nil
}

// EnsureIndexes creates necessary indexes for the settings collection
func (r *settingsRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "updated_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "version", Value: 1},
			},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create settings indexes: %w", err)
	}

	return nil
}
