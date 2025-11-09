package donor

import (
	"context"
	"strings"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorUseCase handles donor business logic
type DonorUseCase struct {
	donorRepo    repositories.DonorRepository
	auditLogRepo repositories.AuditLogRepository
}

// NewDonorUseCase creates a new donor use case
func NewDonorUseCase(
	donorRepo repositories.DonorRepository,
	auditLogRepo repositories.AuditLogRepository,
) *DonorUseCase {
	return &DonorUseCase{
		donorRepo:    donorRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateDonor creates a new donor
func (uc *DonorUseCase) CreateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error {
	// Validate donor
	if err := uc.validateDonor(donor); err != nil {
		return err
	}

	// Check for duplicate email
	if donor.Contact.Email != "" {
		existing, _ := uc.donorRepo.FindByEmail(ctx, donor.Contact.Email)
		if existing != nil {
			return errors.NewBadRequest("A donor with this email already exists")
		}
	}

	donor.CreatedBy = userID
	donor.UpdatedBy = userID

	if err := uc.donorRepo.Create(ctx, donor); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionCreate, "donor", "", "").
		WithEntityID(donor.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetDonor gets a donor by ID
func (uc *DonorUseCase) GetDonor(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error) {
	return uc.donorRepo.FindByID(ctx, id)
}

// UpdateDonor updates a donor
func (uc *DonorUseCase) UpdateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error {
	// Validate donor
	if err := uc.validateDonor(donor); err != nil {
		return err
	}

	// Check if donor exists
	existing, err := uc.donorRepo.FindByID(ctx, donor.ID)
	if err != nil {
		return err
	}

	// Check for duplicate email (excluding current donor)
	if donor.Contact.Email != "" && donor.Contact.Email != existing.Contact.Email {
		emailCheck, _ := uc.donorRepo.FindByEmail(ctx, donor.Contact.Email)
		if emailCheck != nil && emailCheck.ID != donor.ID {
			return errors.NewBadRequest("A donor with this email already exists")
		}
	}

	donor.UpdatedBy = userID

	if err := uc.donorRepo.Update(ctx, donor); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donor", "", "").
		WithEntityID(donor.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// DeleteDonor deletes a donor
func (uc *DonorUseCase) DeleteDonor(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if donor exists
	donor, err := uc.donorRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if donor has donations (business rule: can't delete donors with donation history)
	if donor.DonationCount > 0 {
		return errors.NewBadRequest("Cannot delete donor with donation history. Consider marking as inactive instead.")
	}

	if err := uc.donorRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionDelete, "donor", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListDonors lists donors with filters
func (uc *DonorUseCase) ListDonors(ctx context.Context, filter *repositories.DonorFilter) ([]*entities.Donor, int64, error) {
	return uc.donorRepo.List(ctx, filter)
}

// GetMajorDonors gets all major donors
func (uc *DonorUseCase) GetMajorDonors(ctx context.Context) ([]*entities.Donor, error) {
	return uc.donorRepo.GetMajorDonors(ctx)
}

// GetLapsedDonors gets donors who haven't donated in a while
func (uc *DonorUseCase) GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error) {
	if days <= 0 {
		days = 365 // Default to 1 year
	}
	return uc.donorRepo.GetLapsedDonors(ctx, days)
}

// GetDonorStatistics gets donor statistics
func (uc *DonorUseCase) GetDonorStatistics(ctx context.Context) (*repositories.DonorStatistics, error) {
	return uc.donorRepo.GetDonorStatistics(ctx)
}

// UpdateDonorEngagement updates donor engagement metrics
func (uc *DonorUseCase) UpdateDonorEngagement(ctx context.Context, id primitive.ObjectID, volunteerHours int, eventsAttended int, userID primitive.ObjectID) error {
	donor, err := uc.donorRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if volunteerHours > 0 {
		donor.VolunteerHours += volunteerHours
	}
	if eventsAttended > 0 {
		donor.EventsAttended += eventsAttended
	}

	donor.UpdatedBy = userID

	if err := uc.donorRepo.Update(ctx, donor); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donor", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// validateDonor validates donor data
func (uc *DonorUseCase) validateDonor(donor *entities.Donor) error {
	if donor.Type == "" {
		return errors.NewBadRequest("Donor type is required")
	}

	// Individual donor must have first and last name
	if donor.Type == entities.DonorTypeIndividual {
		if strings.TrimSpace(donor.FirstName) == "" || strings.TrimSpace(donor.LastName) == "" {
			return errors.NewBadRequest("First name and last name are required for individual donors")
		}
	}

	// Organization donor must have organization name
	if donor.Type == entities.DonorTypeOrganization ||
		donor.Type == entities.DonorTypeCorporate ||
		donor.Type == entities.DonorTypeFoundation {
		if strings.TrimSpace(donor.OrganizationName) == "" {
			return errors.NewBadRequest("Organization name is required for organization donors")
		}
	}

	// At least one contact method should be provided
	if donor.Contact.Email == "" && donor.Contact.Phone == "" {
		return errors.NewBadRequest("At least one contact method (email or phone) is required")
	}

	return nil
}
