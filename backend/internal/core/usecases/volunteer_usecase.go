package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type VolunteerUseCase struct {
	volunteerRepo     interfaces.VolunteerRepository
	volunteerHourRepo interfaces.VolunteerHourRepository
	auditRepo         interfaces.AuditLogRepository
}

func NewVolunteerUseCase(
	volunteerRepo interfaces.VolunteerRepository,
	volunteerHourRepo interfaces.VolunteerHourRepository,
	auditRepo interfaces.AuditLogRepository,
) *VolunteerUseCase {
	return &VolunteerUseCase{
		volunteerRepo:     volunteerRepo,
		volunteerHourRepo: volunteerHourRepo,
		auditRepo:         auditRepo,
	}
}

func (uc *VolunteerUseCase) Create(ctx context.Context, req *entities.VolunteerCreateRequest, createdBy string) (*entities.Volunteer, error) {
	volunteer := entities.NewVolunteer(
		req.UserID,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Phone,
	)

	volunteer.Address = req.Address
	volunteer.EmergencyContact = req.EmergencyContact
	volunteer.EmergencyPhone = req.EmergencyPhone
	volunteer.Skills = req.Skills
	volunteer.Availability = req.Availability
	volunteer.Notes = req.Notes
	volunteer.CreatedBy = createdBy

	if err := uc.volunteerRepo.Create(ctx, volunteer); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "volunteer", volunteer.ID.Hex(), "Volunteer created")
	uc.auditRepo.Create(ctx, auditLog)

	return volunteer, nil
}

func (uc *VolunteerUseCase) GetByID(ctx context.Context, id string) (*entities.Volunteer, error) {
	return uc.volunteerRepo.GetByID(ctx, id)
}

func (uc *VolunteerUseCase) List(ctx context.Context, filter *entities.VolunteerFilter) ([]*entities.Volunteer, int64, error) {
	return uc.volunteerRepo.List(ctx, filter)
}

func (uc *VolunteerUseCase) Update(ctx context.Context, id string, req *entities.VolunteerUpdateRequest, updatedBy string) (*entities.Volunteer, error) {
	volunteer, err := uc.volunteerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		volunteer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		volunteer.LastName = req.LastName
	}
	if req.Phone != "" {
		volunteer.Phone = req.Phone
	}
	if req.Address != "" {
		volunteer.Address = req.Address
	}
	if req.EmergencyContact != "" {
		volunteer.EmergencyContact = req.EmergencyContact
	}
	if req.EmergencyPhone != "" {
		volunteer.EmergencyPhone = req.EmergencyPhone
	}
	if req.Skills != nil {
		volunteer.Skills = req.Skills
	}
	if req.Availability != nil {
		volunteer.Availability = req.Availability
	}
	if req.Status != "" {
		volunteer.Status = req.Status
	}
	if req.Notes != "" {
		volunteer.Notes = req.Notes
	}
	volunteer.UpdatedBy = updatedBy

	if err := uc.volunteerRepo.Update(ctx, id, volunteer); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "volunteer", id, "Volunteer updated")
	uc.auditRepo.Create(ctx, auditLog)

	return volunteer, nil
}

func (uc *VolunteerUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.volunteerRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "volunteer", id, "Volunteer deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *VolunteerUseCase) AddTraining(ctx context.Context, id string, training entities.Training, addedBy string) error {
	if err := uc.volunteerRepo.AddTraining(ctx, id, training); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(addedBy, "", "", entities.ActionUpdate, "volunteer", id, "Training added")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *VolunteerUseCase) LogHours(ctx context.Context, req *entities.VolunteerHourCreateRequest, loggedBy string) (*entities.VolunteerHour, error) {
	// Verify volunteer exists
	volunteer, err := uc.volunteerRepo.GetByID(ctx, req.VolunteerID)
	if err != nil {
		return nil, err
	}

	hour := entities.NewVolunteerHour(
		req.VolunteerID,
		volunteer.FirstName+" "+volunteer.LastName,
		req.Date,
		req.Hours,
		req.Activity,
	)
	hour.Description = req.Description
	hour.LoggedBy = loggedBy

	if err := uc.volunteerHourRepo.Create(ctx, hour); err != nil {
		return nil, err
	}

	// Update volunteer total hours
	uc.volunteerRepo.UpdateTotalHours(ctx, req.VolunteerID, req.Hours)

	// Audit
	auditLog := entities.NewAuditLog(loggedBy, "", "", entities.ActionCreate, "volunteer_hours", hour.ID.Hex(), "Volunteer hours logged")
	uc.auditRepo.Create(ctx, auditLog)

	return hour, nil
}

func (uc *VolunteerUseCase) GetVolunteerHours(ctx context.Context, volunteerID string, startDate, endDate string, limit, offset int) ([]*entities.VolunteerHour, int64, error) {
	return uc.volunteerHourRepo.GetByVolunteerID(ctx, volunteerID, startDate, endDate, limit, offset)
}

func (uc *VolunteerUseCase) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	return uc.volunteerRepo.GetActiveVolunteers(ctx)
}

func (uc *VolunteerUseCase) GetVolunteerStatistics(ctx context.Context, volunteerID string, startDate, endDate string) (map[string]interface{}, error) {
	volunteer, err := uc.volunteerRepo.GetByID(ctx, volunteerID)
	if err != nil {
		return nil, err
	}

	hours, _, err := uc.volunteerHourRepo.GetByVolunteerID(ctx, volunteerID, startDate, endDate, 0, 0)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalHours := 0.0
	activities := make(map[string]float64)
	for _, hour := range hours {
		totalHours += hour.Hours
		activities[hour.Activity] += hour.Hours
	}

	stats := map[string]interface{}{
		"volunteer_id":   volunteerID,
		"volunteer_name": volunteer.FirstName + " " + volunteer.LastName,
		"total_hours":    totalHours,
		"activities":     activities,
		"start_date":     startDate,
		"end_date":       endDate,
	}

	return stats, nil
}
