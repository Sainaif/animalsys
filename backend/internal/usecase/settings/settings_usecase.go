package settings

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SettingsUseCase struct {
	settingsRepo repositories.SettingsRepository
	auditLogRepo repositories.AuditLogRepository
}

func NewSettingsUseCase(
	settingsRepo repositories.SettingsRepository,
	auditLogRepo repositories.AuditLogRepository,
) *SettingsUseCase {
	return &SettingsUseCase{
		settingsRepo: settingsRepo,
		auditLogRepo: auditLogRepo,
	}
}

// GetSettings retrieves the foundation settings
func (uc *SettingsUseCase) GetSettings(ctx context.Context) (*entities.FoundationSettings, error) {
	return uc.settingsRepo.Get(ctx)
}

// GetSettingsByID retrieves settings by ID
func (uc *SettingsUseCase) GetSettingsByID(ctx context.Context, id primitive.ObjectID) (*entities.FoundationSettings, error) {
	return uc.settingsRepo.GetByID(ctx, id)
}

// UpdateSettings updates the foundation settings
func (uc *SettingsUseCase) UpdateSettings(ctx context.Context, settings *entities.FoundationSettings, userID primitive.ObjectID) error {
	// Validate required fields
	if settings.Name == "" {
		return errors.NewBadRequest("Foundation name is required")
	}

	if settings.ContactInfo.Email == "" {
		return errors.NewBadRequest("Contact email is required")
	}

	settings.UpdatedBy = userID

	if err := uc.settingsRepo.Update(ctx, settings); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "settings", "", "").
			WithEntityID(settings.ID))

	return nil
}

// InitializeSettings creates initial foundation settings
func (uc *SettingsUseCase) InitializeSettings(ctx context.Context, name string, userID primitive.ObjectID) (*entities.FoundationSettings, error) {
	if name == "" {
		return nil, errors.NewBadRequest("Foundation name is required")
	}

	// Check if settings already exist
	existing, err := uc.settingsRepo.Get(ctx)
	if err == nil && existing != nil {
		return nil, errors.NewBadRequest("Settings already initialized")
	}

	settings := entities.NewFoundationSettings(name, userID)
	if err := uc.settingsRepo.Create(ctx, settings); err != nil {
		return nil, err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "settings", "", "").
			WithEntityID(settings.ID))

	return settings, nil
}

// UpdateEmailSettings updates only email settings
func (uc *SettingsUseCase) UpdateEmailSettings(ctx context.Context, emailSettings entities.EmailSettings, userID primitive.ObjectID) error {
	if emailSettings.FromEmail == "" {
		return errors.NewBadRequest("From email is required")
	}

	if err := uc.settingsRepo.UpdateEmailSettings(ctx, emailSettings, userID); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "settings", "email_settings", ""))

	return nil
}

// UpdateNotificationSettings updates only notification settings
func (uc *SettingsUseCase) UpdateNotificationSettings(ctx context.Context, notificationSettings entities.NotificationSettings, userID primitive.ObjectID) error {
	if err := uc.settingsRepo.UpdateNotificationSettings(ctx, notificationSettings, userID); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "settings", "notification_settings", ""))

	return nil
}

// UpdateFeatureFlags updates only feature flags
func (uc *SettingsUseCase) UpdateFeatureFlags(ctx context.Context, features entities.FeatureFlags, userID primitive.ObjectID) error {
	if err := uc.settingsRepo.UpdateFeatureFlags(ctx, features, userID); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "settings", "features", ""))

	return nil
}

// UpdateBranding updates only branding settings
func (uc *SettingsUseCase) UpdateBranding(ctx context.Context, branding entities.Branding, userID primitive.ObjectID) error {
	if err := uc.settingsRepo.UpdateBranding(ctx, branding, userID); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "settings", "branding", ""))

	return nil
}

// GetContactInfo returns only contact information
func (uc *SettingsUseCase) GetContactInfo(ctx context.Context) (*entities.ContactDetails, error) {
	return uc.settingsRepo.GetContactInfo(ctx)
}

// GetOperatingHours returns operating hours
func (uc *SettingsUseCase) GetOperatingHours(ctx context.Context) (map[string]entities.OperatingHour, error) {
	return uc.settingsRepo.GetOperatingHours(ctx)
}
