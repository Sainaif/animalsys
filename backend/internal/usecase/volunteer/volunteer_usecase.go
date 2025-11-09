package volunteer

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VolunteerUseCase handles volunteer-related business logic
type VolunteerUseCase struct {
	volunteerRepo  repositories.VolunteerRepository
	assignmentRepo repositories.VolunteerAssignmentRepository
	auditLogRepo   repositories.AuditLogRepository
}

// NewVolunteerUseCase creates a new volunteer use case
func NewVolunteerUseCase(
	volunteerRepo repositories.VolunteerRepository,
	assignmentRepo repositories.VolunteerAssignmentRepository,
	auditLogRepo repositories.AuditLogRepository,
) *VolunteerUseCase {
	return &VolunteerUseCase{
		volunteerRepo:  volunteerRepo,
		assignmentRepo: assignmentRepo,
		auditLogRepo:   auditLogRepo,
	}
}

// CreateVolunteer creates a new volunteer
func (uc *VolunteerUseCase) CreateVolunteer(ctx context.Context, volunteer *entities.Volunteer, userID primitive.ObjectID) error {
	// Validate volunteer
	if err := uc.validateVolunteer(volunteer); err != nil {
		return err
	}

	// Check for duplicate email
	existing, _ := uc.volunteerRepo.FindByEmail(ctx, volunteer.Email)
	if existing != nil {
		return errors.NewConflict("Volunteer with this email already exists")
	}

	// Create volunteer
	if err := uc.volunteerRepo.Create(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "volunteer", "", "").
			WithEntityID(volunteer.ID))

	return nil
}

// UpdateVolunteer updates a volunteer
func (uc *VolunteerUseCase) UpdateVolunteer(ctx context.Context, volunteer *entities.Volunteer, userID primitive.ObjectID) error {
	// Check if volunteer exists
	existingVolunteer, err := uc.volunteerRepo.FindByID(ctx, volunteer.ID)
	if err != nil {
		return err
	}

	// Validate volunteer
	if err := uc.validateVolunteer(volunteer); err != nil {
		return err
	}

	// Check for duplicate email (excluding self)
	existing, _ := uc.volunteerRepo.FindByEmail(ctx, volunteer.Email)
	if existing != nil && existing.ID != volunteer.ID {
		return errors.NewConflict("Volunteer with this email already exists")
	}

	// Update expired certifications
	volunteer.UpdateExpiredCertifications()

	// Update volunteer
	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "", "").
			WithEntityID(volunteer.ID).
			WithChanges(map[string]interface{}{
				"old_status": existingVolunteer.Status,
				"new_status": volunteer.Status,
			}))

	return nil
}

// DeleteVolunteer deletes a volunteer
func (uc *VolunteerUseCase) DeleteVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if volunteer exists
	_, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	// Check for active assignments
	activeAssignments, err := uc.assignmentRepo.GetActiveAssignments(ctx, volunteerID)
	if err == nil && len(activeAssignments) > 0 {
		return errors.NewBadRequest("Cannot delete volunteer with active assignments")
	}

	// Delete volunteer
	if err := uc.volunteerRepo.Delete(ctx, volunteerID); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "volunteer", "", "").
			WithEntityID(volunteerID))

	return nil
}

// GetVolunteer gets a volunteer by ID
func (uc *VolunteerUseCase) GetVolunteer(ctx context.Context, volunteerID primitive.ObjectID) (*entities.Volunteer, error) {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return nil, err
	}

	// Update expired certifications
	volunteer.UpdateExpiredCertifications()

	return volunteer, nil
}

// GetVolunteerByEmail gets a volunteer by email
func (uc *VolunteerUseCase) GetVolunteerByEmail(ctx context.Context, email string) (*entities.Volunteer, error) {
	return uc.volunteerRepo.FindByEmail(ctx, email)
}

// ListVolunteers lists volunteers with filters
func (uc *VolunteerUseCase) ListVolunteers(ctx context.Context, filter *repositories.VolunteerFilter) ([]*entities.Volunteer, int64, error) {
	return uc.volunteerRepo.List(ctx, filter)
}

// GetActiveVolunteers gets all active volunteers
func (uc *VolunteerUseCase) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	return uc.volunteerRepo.GetActiveVolunteers(ctx)
}

// GetVolunteersBySkill gets volunteers with a specific skill
func (uc *VolunteerUseCase) GetVolunteersBySkill(ctx context.Context, skill string) ([]*entities.Volunteer, error) {
	return uc.volunteerRepo.GetVolunteersBySkill(ctx, skill)
}

// GetVolunteersNeedingBackgroundCheck gets volunteers needing background check
func (uc *VolunteerUseCase) GetVolunteersNeedingBackgroundCheck(ctx context.Context) ([]*entities.Volunteer, error) {
	return uc.volunteerRepo.GetVolunteersNeedingBackgroundCheck(ctx)
}

// GetTopVolunteers gets top volunteers by hours
func (uc *VolunteerUseCase) GetTopVolunteers(ctx context.Context, limit int) ([]*entities.Volunteer, error) {
	return uc.volunteerRepo.GetTopVolunteers(ctx, limit)
}

// ApproveVolunteer approves a volunteer application
func (uc *VolunteerUseCase) ApproveVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	if volunteer.Status != entities.VolunteerStatusInactive {
		return errors.NewBadRequest("Only inactive volunteers can be approved")
	}

	volunteer.Approve()
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "approved", "").
			WithEntityID(volunteerID))

	return nil
}

// SuspendVolunteer suspends a volunteer
func (uc *VolunteerUseCase) SuspendVolunteer(ctx context.Context, volunteerID primitive.ObjectID, reason string, userID primitive.ObjectID) error {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	volunteer.Status = entities.VolunteerStatusSuspended
	volunteer.Notes += "\nSuspension reason: " + reason
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "suspended", reason).
			WithEntityID(volunteerID))

	return nil
}

// AddHours adds volunteer hours
func (uc *VolunteerUseCase) AddHours(ctx context.Context, volunteerID primitive.ObjectID, hours float64, userID primitive.ObjectID) error {
	if hours <= 0 {
		return errors.NewBadRequest("Hours must be greater than 0")
	}

	if err := uc.volunteerRepo.UpdateHours(ctx, volunteerID, hours); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "added_hours", "").
			WithEntityID(volunteerID).
			WithChanges(map[string]interface{}{
				"hours": hours,
			}))

	return nil
}

// AddCommendation adds a commendation to a volunteer
func (uc *VolunteerUseCase) AddCommendation(ctx context.Context, volunteerID primitive.ObjectID, note string, userID primitive.ObjectID) error {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	volunteer.AddCommendation()
	volunteer.Notes += "\nCommendation: " + note
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "commendation", note).
			WithEntityID(volunteerID))

	return nil
}

// AddWarning adds a warning to a volunteer
func (uc *VolunteerUseCase) AddWarning(ctx context.Context, volunteerID primitive.ObjectID, reason string, userID primitive.ObjectID) error {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	volunteer.AddWarning()
	volunteer.Notes += "\nWarning: " + reason
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "warning", reason).
			WithEntityID(volunteerID))

	return nil
}

// UpdateRating updates volunteer rating
func (uc *VolunteerUseCase) UpdateRating(ctx context.Context, volunteerID primitive.ObjectID, rating float64, userID primitive.ObjectID) error {
	if rating < 1 || rating > 5 {
		return errors.NewBadRequest("Rating must be between 1 and 5")
	}

	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	volunteer.UpdateRating(rating)
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "rating_updated", "").
			WithEntityID(volunteerID).
			WithChanges(map[string]interface{}{
				"new_rating": rating,
			}))

	return nil
}

// AddCertification adds a certification to a volunteer
func (uc *VolunteerUseCase) AddCertification(ctx context.Context, volunteerID primitive.ObjectID, cert entities.Certification, userID primitive.ObjectID) error {
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return err
	}

	cert.ID = primitive.NewObjectID()

	// Check if certification is expired
	if cert.ExpirationDate != nil && cert.ExpirationDate.Before(time.Now()) {
		cert.IsExpired = true
	}

	volunteer.Certifications = append(volunteer.Certifications, cert)
	volunteer.UpdatedBy = userID

	if err := uc.volunteerRepo.Update(ctx, volunteer); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "volunteer", "added_certification", "").
			WithEntityID(volunteerID).
			WithChanges(map[string]interface{}{
				"certification": cert.Name,
			}))

	return nil
}

// GetVolunteerStatistics gets volunteer statistics
func (uc *VolunteerUseCase) GetVolunteerStatistics(ctx context.Context) (*repositories.VolunteerStatistics, error) {
	return uc.volunteerRepo.GetVolunteerStatistics(ctx)
}

// validateVolunteer validates volunteer data
func (uc *VolunteerUseCase) validateVolunteer(volunteer *entities.Volunteer) error {
	if volunteer.FirstName == "" {
		return errors.NewBadRequest("First name is required")
	}

	if volunteer.LastName == "" {
		return errors.NewBadRequest("Last name is required")
	}

	if volunteer.Email == "" {
		return errors.NewBadRequest("Email is required")
	}

	// Validate email format (basic validation)
	if !isValidEmail(volunteer.Email) {
		return errors.NewBadRequest("Invalid email format")
	}

	return nil
}

// isValidEmail is a basic email validation
func isValidEmail(email string) bool {
	// Basic email validation
	return len(email) > 3 && len(email) < 255 &&
		containsChar(email, '@') && containsChar(email, '.')
}

func containsChar(s string, char rune) bool {
	for _, c := range s {
		if c == char {
			return true
		}
	}
	return false
}
