package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type VeterinaryUseCase struct {
	veterinaryRepo  interfaces.VeterinaryRepository
	vaccinationRepo interfaces.VaccinationRepository
	animalRepo      interfaces.AnimalRepository
	auditRepo       interfaces.AuditLogRepository
}

func NewVeterinaryUseCase(
	veterinaryRepo interfaces.VeterinaryRepository,
	vaccinationRepo interfaces.VaccinationRepository,
	animalRepo interfaces.AnimalRepository,
	auditRepo interfaces.AuditLogRepository,
) *VeterinaryUseCase {
	return &VeterinaryUseCase{
		veterinaryRepo:  veterinaryRepo,
		vaccinationRepo: vaccinationRepo,
		animalRepo:      animalRepo,
		auditRepo:       auditRepo,
	}
}

func (uc *VeterinaryUseCase) RecordVisit(ctx context.Context, req *entities.VeterinaryVisitCreateRequest, recordedBy string) (*entities.VeterinaryVisit, error) {
	// Verify animal exists
	animal, err := uc.animalRepo.GetByID(ctx, req.AnimalID)
	if err != nil {
		return nil, err
	}

	visit := entities.NewVeterinaryVisit(
		req.AnimalID,
		animal.Name,
		req.Date,
		req.VisitType,
		req.VetName,
		req.Clinic,
	)

	visit.Reason = req.Reason
	visit.Diagnosis = req.Diagnosis
	visit.Treatment = req.Treatment
	visit.Medications = req.Medications
	visit.Cost = req.Cost
	visit.FollowUpDate = req.FollowUpDate
	visit.Notes = req.Notes
	visit.RecordedBy = recordedBy

	if err := uc.veterinaryRepo.Create(ctx, visit); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(recordedBy, "", "", entities.ActionCreate, "veterinary_visit", visit.ID.Hex(), "Veterinary visit recorded")
	uc.auditRepo.Create(ctx, auditLog)

	return visit, nil
}

func (uc *VeterinaryUseCase) GetVisitByID(ctx context.Context, id string) (*entities.VeterinaryVisit, error) {
	return uc.veterinaryRepo.GetByID(ctx, id)
}

func (uc *VeterinaryUseCase) ListVisits(ctx context.Context, filter *entities.VeterinaryFilter) ([]*entities.VeterinaryVisit, int64, error) {
	return uc.veterinaryRepo.List(ctx, filter)
}

func (uc *VeterinaryUseCase) UpdateVisit(ctx context.Context, id string, req *entities.VeterinaryVisitUpdateRequest, updatedBy string) (*entities.VeterinaryVisit, error) {
	visit, err := uc.veterinaryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if !req.Date.IsZero() {
		visit.Date = req.Date
	}
	if req.VisitType != "" {
		visit.VisitType = req.VisitType
	}
	if req.VetName != "" {
		visit.VetName = req.VetName
	}
	if req.Clinic != "" {
		visit.Clinic = req.Clinic
	}
	if req.Reason != "" {
		visit.Reason = req.Reason
	}
	if req.Diagnosis != "" {
		visit.Diagnosis = req.Diagnosis
	}
	if req.Treatment != "" {
		visit.Treatment = req.Treatment
	}
	if req.Medications != nil {
		visit.Medications = req.Medications
	}
	if req.Cost > 0 {
		visit.Cost = req.Cost
	}
	if req.FollowUpDate != nil {
		visit.FollowUpDate = req.FollowUpDate
	}
	if req.Notes != "" {
		visit.Notes = req.Notes
	}
	visit.UpdatedBy = updatedBy

	if err := uc.veterinaryRepo.Update(ctx, id, visit); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "veterinary_visit", id, "Veterinary visit updated")
	uc.auditRepo.Create(ctx, auditLog)

	return visit, nil
}

func (uc *VeterinaryUseCase) DeleteVisit(ctx context.Context, id string, deletedBy string) error {
	if err := uc.veterinaryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "veterinary_visit", id, "Veterinary visit deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *VeterinaryUseCase) GetAnimalVisits(ctx context.Context, animalID string, limit, offset int) ([]*entities.VeterinaryVisit, int64, error) {
	return uc.veterinaryRepo.GetByAnimalID(ctx, animalID, limit, offset)
}

func (uc *VeterinaryUseCase) GetUpcomingFollowUps(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	return uc.veterinaryRepo.GetUpcomingFollowUps(ctx, days)
}

func (uc *VeterinaryUseCase) RecordVaccination(ctx context.Context, req *entities.VaccinationCreateRequest, recordedBy string) (*entities.Vaccination, error) {
	// Verify animal exists
	animal, err := uc.animalRepo.GetByID(ctx, req.AnimalID)
	if err != nil {
		return nil, err
	}

	vaccination := entities.NewVaccination(
		req.AnimalID,
		animal.Name,
		req.VaccineName,
		req.Date,
		req.VetName,
		req.NextDueDate,
	)

	vaccination.BatchNumber = req.BatchNumber
	vaccination.Clinic = req.Clinic
	vaccination.Cost = req.Cost
	vaccination.Notes = req.Notes
	vaccination.RecordedBy = recordedBy

	if err := uc.vaccinationRepo.Create(ctx, vaccination); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(recordedBy, "", "", entities.ActionCreate, "vaccination", vaccination.ID.Hex(), "Vaccination recorded")
	uc.auditRepo.Create(ctx, auditLog)

	return vaccination, nil
}

func (uc *VeterinaryUseCase) GetVaccinationByID(ctx context.Context, id string) (*entities.Vaccination, error) {
	return uc.vaccinationRepo.GetByID(ctx, id)
}

func (uc *VeterinaryUseCase) ListVaccinations(ctx context.Context, filter *entities.VaccinationFilter) ([]*entities.Vaccination, int64, error) {
	return uc.vaccinationRepo.List(ctx, filter)
}

func (uc *VeterinaryUseCase) UpdateVaccination(ctx context.Context, id string, req *entities.VaccinationUpdateRequest, updatedBy string) (*entities.Vaccination, error) {
	vaccination, err := uc.vaccinationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.VaccineName != "" {
		vaccination.VaccineName = req.VaccineName
	}
	if !req.Date.IsZero() {
		vaccination.Date = req.Date
	}
	if req.VetName != "" {
		vaccination.VetName = req.VetName
	}
	if req.BatchNumber != "" {
		vaccination.BatchNumber = req.BatchNumber
	}
	if req.Clinic != "" {
		vaccination.Clinic = req.Clinic
	}
	if !req.NextDueDate.IsZero() {
		vaccination.NextDueDate = req.NextDueDate
	}
	if req.Cost > 0 {
		vaccination.Cost = req.Cost
	}
	if req.Notes != "" {
		vaccination.Notes = req.Notes
	}
	vaccination.UpdatedBy = updatedBy

	if err := uc.vaccinationRepo.Update(ctx, id, vaccination); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "vaccination", id, "Vaccination updated")
	uc.auditRepo.Create(ctx, auditLog)

	return vaccination, nil
}

func (uc *VeterinaryUseCase) DeleteVaccination(ctx context.Context, id string, deletedBy string) error {
	if err := uc.vaccinationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "vaccination", id, "Vaccination deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *VeterinaryUseCase) GetAnimalVaccinations(ctx context.Context, animalID string, limit, offset int) ([]*entities.Vaccination, int64, error) {
	return uc.vaccinationRepo.GetByAnimalID(ctx, animalID, limit, offset)
}

func (uc *VeterinaryUseCase) GetUpcomingVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	return uc.vaccinationRepo.GetUpcomingVaccinations(ctx, days)
}

func (uc *VeterinaryUseCase) GetVeterinaryStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	visits, _, err := uc.veterinaryRepo.GetByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return nil, err
	}

	totalCost := 0.0
	byType := make(map[string]int)
	byVet := make(map[string]int)

	for _, visit := range visits {
		totalCost += visit.Cost
		byType[string(visit.VisitType)]++
		byVet[visit.VetName]++
	}

	stats := map[string]interface{}{
		"period": map[string]string{
			"start": startDate,
			"end":   endDate,
		},
		"total_visits": len(visits),
		"total_cost":   totalCost,
		"by_type":      byType,
		"by_vet":       byVet,
	}

	return stats, nil
}
