package partner

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartnerUseCase struct {
	partnerRepo  repositories.PartnerRepository
	auditLogRepo repositories.AuditLogRepository
}

func NewPartnerUseCase(
	partnerRepo repositories.PartnerRepository,
	auditLogRepo repositories.AuditLogRepository,
) *PartnerUseCase {
	return &PartnerUseCase{
		partnerRepo:  partnerRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreatePartner creates a new partner
func (uc *PartnerUseCase) CreatePartner(ctx context.Context, partner *entities.Partner, userID primitive.ObjectID) error {
	// Validate required fields
	if partner.Name == "" {
		return errors.NewBadRequest("Partner name is required")
	}

	if partner.Type == "" {
		return errors.NewBadRequest("Partner type is required")
	}

	// Validate partner type
	validTypes := map[entities.PartnerType]bool{
		entities.PartnerTypeRescue:     true,
		entities.PartnerTypeShelter:    true,
		entities.PartnerTypeVeterinary: true,
		entities.PartnerTypeFoster:     true,
		entities.PartnerTypeTransport:  true,
		entities.PartnerTypeSanctuary:  true,
		entities.PartnerTypeGovernment: true,
		entities.PartnerTypeCorporate:  true,
		entities.PartnerTypeOther:      true,
	}
	if !validTypes[partner.Type] {
		return errors.NewBadRequest("Invalid partner type")
	}

	// Set creation metadata
	partner.CreatedBy = userID
	now := time.Now()
	partner.CreatedAt = now
	partner.UpdatedAt = now

	// Set default status if not provided
	if partner.Status == "" {
		partner.Status = entities.PartnerStatusPending
	}

	// Set partner since date if not provided
	if partner.PartnerSince.IsZero() {
		partner.PartnerSince = now
	}

	// Initialize arrays if nil
	if partner.ServicesProvided == nil {
		partner.ServicesProvided = []string{}
	}
	if partner.Specializations == nil {
		partner.Specializations = []string{}
	}
	if partner.Documents == nil {
		partner.Documents = []string{}
	}
	if partner.Tags == nil {
		partner.Tags = []string{}
	}
	if partner.SocialMedia == nil {
		partner.SocialMedia = make(map[string]string)
	}

	if err := uc.partnerRepo.Create(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "partner", partner.Name, "").
			WithEntityID(partner.ID))

	return nil
}

// GetPartnerByID retrieves a partner by ID
func (uc *PartnerUseCase) GetPartnerByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	partner, err := uc.partnerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

// UpdatePartner updates a partner
func (uc *PartnerUseCase) UpdatePartner(ctx context.Context, partner *entities.Partner, userID primitive.ObjectID) error {
	// Validate required fields
	if partner.Name == "" {
		return errors.NewBadRequest("Partner name is required")
	}

	// Check if partner exists
	existing, err := uc.partnerRepo.FindByID(ctx, partner.ID)
	if err != nil {
		return err
	}

	// Preserve creation info
	partner.CreatedBy = existing.CreatedBy
	partner.CreatedAt = existing.CreatedAt
	partner.UpdatedAt = time.Now()

	// Preserve statistics
	partner.Rating = existing.Rating
	partner.TotalRatings = existing.TotalRatings
	partner.TotalTransfersIn = existing.TotalTransfersIn
	partner.TotalTransfersOut = existing.TotalTransfersOut
	partner.SuccessfulPlacements = existing.SuccessfulPlacements

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "").
			WithEntityID(partner.ID))

	return nil
}

// DeletePartner deletes a partner
func (uc *PartnerUseCase) DeletePartner(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if partner exists
	partner, err := uc.partnerRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.partnerRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "partner", partner.Name, "").
			WithEntityID(id))

	return nil
}

// ListPartners lists partners with filtering and pagination
func (uc *PartnerUseCase) ListPartners(ctx context.Context, filter *repositories.PartnerFilter) ([]*entities.Partner, int64, error) {
	return uc.partnerRepo.List(ctx, filter)
}

// GetPartnersByType gets partners by type
func (uc *PartnerUseCase) GetPartnersByType(ctx context.Context, partnerType entities.PartnerType) ([]*entities.Partner, error) {
	return uc.partnerRepo.GetByType(ctx, partnerType)
}

// GetPartnersByStatus gets partners by status
func (uc *PartnerUseCase) GetPartnersByStatus(ctx context.Context, status entities.PartnerStatus) ([]*entities.Partner, error) {
	return uc.partnerRepo.GetByStatus(ctx, status)
}

// GetActivePartners gets all active partners
func (uc *PartnerUseCase) GetActivePartners(ctx context.Context) ([]*entities.Partner, error) {
	return uc.partnerRepo.GetActivePartners(ctx)
}

// GetPartnersWithCapacity gets partners with available capacity
func (uc *PartnerUseCase) GetPartnersWithCapacity(ctx context.Context) ([]*entities.Partner, error) {
	return uc.partnerRepo.GetPartnersWithCapacity(ctx)
}

// GetPartnerStatistics gets partner statistics
func (uc *PartnerUseCase) GetPartnerStatistics(ctx context.Context) (*repositories.PartnerStatistics, error) {
	return uc.partnerRepo.GetPartnerStatistics(ctx)
}

// ActivatePartner activates a partner
func (uc *PartnerUseCase) ActivatePartner(ctx context.Context, partnerID, userID primitive.ObjectID) error {
	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	if partner.Status == entities.PartnerStatusActive {
		return errors.NewBadRequest("Partner is already active")
	}

	partner.Activate()

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "activated partner").
			WithEntityID(partnerID))

	return nil
}

// SuspendPartner suspends a partner
func (uc *PartnerUseCase) SuspendPartner(ctx context.Context, partnerID, userID primitive.ObjectID, reason string) error {
	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	if partner.Status == entities.PartnerStatusSuspended {
		return errors.NewBadRequest("Partner is already suspended")
	}

	partner.Suspend()

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "suspended partner: "+reason).
			WithEntityID(partnerID))

	return nil
}

// DeactivatePartner deactivates a partner
func (uc *PartnerUseCase) DeactivatePartner(ctx context.Context, partnerID, userID primitive.ObjectID, reason string) error {
	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	if partner.Status == entities.PartnerStatusInactive {
		return errors.NewBadRequest("Partner is already inactive")
	}

	partner.Deactivate()

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "deactivated partner: "+reason).
			WithEntityID(partnerID))

	return nil
}

// AddRating adds a rating to a partner
func (uc *PartnerUseCase) AddRating(ctx context.Context, partnerID, userID primitive.ObjectID, rating float64) error {
	// Validate rating
	if rating < 0 || rating > 5 {
		return errors.NewBadRequest("Rating must be between 0 and 5")
	}

	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	partner.AddRating(rating)

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "added rating").
			WithEntityID(partnerID))

	return nil
}

// UpdateCapacity updates the current capacity of a partner
func (uc *PartnerUseCase) UpdateCapacity(ctx context.Context, partnerID, userID primitive.ObjectID, currentCapacity int) error {
	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	// Validate capacity
	if currentCapacity < 0 {
		return errors.NewBadRequest("Capacity cannot be negative")
	}

	if partner.MaxCapacity > 0 && currentCapacity > partner.MaxCapacity {
		return errors.NewBadRequest("Current capacity cannot exceed max capacity")
	}

	partner.UpdateCapacity(currentCapacity)

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "updated capacity").
			WithEntityID(partnerID))

	return nil
}

// SetAgreementExpiry sets the agreement expiry date for a partner
func (uc *PartnerUseCase) SetAgreementExpiry(ctx context.Context, partnerID, userID primitive.ObjectID, expiryDate time.Time) error {
	partner, err := uc.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return err
	}

	partner.AgreementExpiry = &expiryDate
	partner.UpdatedAt = time.Now()

	if err := uc.partnerRepo.Update(ctx, partner); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "partner", partner.Name, "set agreement expiry date").
			WithEntityID(partnerID))

	return nil
}
