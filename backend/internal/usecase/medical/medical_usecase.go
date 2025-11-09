package medical

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedicalUseCase struct {
	conditionRepo    repositories.MedicalConditionRepository
	medicationRepo   repositories.MedicationRepository
	treatmentPlanRepo repositories.TreatmentPlanRepository
	animalRepo       repositories.AnimalRepository
	auditLogRepo     repositories.AuditLogRepository
}

func NewMedicalUseCase(
	conditionRepo repositories.MedicalConditionRepository,
	medicationRepo repositories.MedicationRepository,
	treatmentPlanRepo repositories.TreatmentPlanRepository,
	animalRepo repositories.AnimalRepository,
	auditLogRepo repositories.AuditLogRepository,
) *MedicalUseCase {
	return &MedicalUseCase{
		conditionRepo:    conditionRepo,
		medicationRepo:   medicationRepo,
		treatmentPlanRepo: treatmentPlanRepo,
		animalRepo:       animalRepo,
		auditLogRepo:     auditLogRepo,
	}
}

// ===== Medical Condition Methods =====

// CreateCondition creates a new medical condition
func (uc *MedicalUseCase) CreateCondition(ctx context.Context, condition *entities.MedicalCondition, userID primitive.ObjectID) error {
	// Verify animal exists
	_, err := uc.animalRepo.FindByID(ctx, condition.AnimalID)
	if err != nil {
		return err
	}

	if err := uc.conditionRepo.Create(ctx, condition); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionCreate,
		EntityType: "medical_condition",
		EntityID:   &condition.ID,
	})

	return nil
}

// GetCondition retrieves a medical condition by ID
func (uc *MedicalUseCase) GetCondition(ctx context.Context, id primitive.ObjectID) (*entities.MedicalCondition, error) {
	return uc.conditionRepo.FindByID(ctx, id)
}

// UpdateCondition updates a medical condition
func (uc *MedicalUseCase) UpdateCondition(ctx context.Context, condition *entities.MedicalCondition, userID primitive.ObjectID) error {
	if err := uc.conditionRepo.Update(ctx, condition); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "medical_condition",
		EntityID:   &condition.ID,
	})

	return nil
}

// DeleteCondition deletes a medical condition
func (uc *MedicalUseCase) DeleteCondition(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := uc.conditionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.conditionRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionDelete,
		EntityType: "medical_condition",
		EntityID:   &id,
	})

	return nil
}

// ListConditions lists medical conditions with filtering
func (uc *MedicalUseCase) ListConditions(ctx context.Context, filter repositories.MedicalConditionFilter) ([]*entities.MedicalCondition, int64, error) {
	return uc.conditionRepo.List(ctx, filter)
}

// GetAnimalConditions gets all conditions for an animal
func (uc *MedicalUseCase) GetAnimalConditions(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error) {
	return uc.conditionRepo.FindByAnimal(ctx, animalID)
}

// GetActiveConditions gets active conditions for an animal
func (uc *MedicalUseCase) GetActiveConditions(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error) {
	return uc.conditionRepo.FindActiveByAnimal(ctx, animalID)
}

// GetChronicConditions gets all chronic conditions
func (uc *MedicalUseCase) GetChronicConditions(ctx context.Context) ([]*entities.MedicalCondition, error) {
	return uc.conditionRepo.FindChronicConditions(ctx)
}

// ResolveCondition marks a condition as resolved
func (uc *MedicalUseCase) ResolveCondition(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	condition, err := uc.conditionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	condition.Status = entities.ConditionStatusResolved
	condition.ResolvedDate = &now

	if err := uc.conditionRepo.Update(ctx, condition); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "medical_condition",
		EntityID:   &id,
	})

	return nil
}

// ===== Medication Methods =====

// CreateMedication creates a new medication record
func (uc *MedicalUseCase) CreateMedication(ctx context.Context, medication *entities.Medication, userID primitive.ObjectID) error {
	// Verify animal exists
	_, err := uc.animalRepo.FindByID(ctx, medication.AnimalID)
	if err != nil {
		return err
	}

	// Verify condition exists if provided
	if medication.ConditionID != nil {
		_, err := uc.conditionRepo.FindByID(ctx, *medication.ConditionID)
		if err != nil {
			return err
		}
	}

	if err := uc.medicationRepo.Create(ctx, medication); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionCreate,
		EntityType: "medication",
		EntityID:   &medication.ID,
	})

	return nil
}

// GetMedication retrieves a medication by ID
func (uc *MedicalUseCase) GetMedication(ctx context.Context, id primitive.ObjectID) (*entities.Medication, error) {
	return uc.medicationRepo.FindByID(ctx, id)
}

// UpdateMedication updates a medication record
func (uc *MedicalUseCase) UpdateMedication(ctx context.Context, medication *entities.Medication, userID primitive.ObjectID) error {
	if err := uc.medicationRepo.Update(ctx, medication); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "medication",
		EntityID:   &medication.ID,
	})

	return nil
}

// DeleteMedication deletes a medication record
func (uc *MedicalUseCase) DeleteMedication(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := uc.medicationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.medicationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionDelete,
		EntityType: "medication",
		EntityID:   &id,
	})

	return nil
}

// ListMedications lists medications with filtering
func (uc *MedicalUseCase) ListMedications(ctx context.Context, filter repositories.MedicationFilter) ([]*entities.Medication, int64, error) {
	return uc.medicationRepo.List(ctx, filter)
}

// GetAnimalMedications gets all medications for an animal
func (uc *MedicalUseCase) GetAnimalMedications(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error) {
	return uc.medicationRepo.FindByAnimal(ctx, animalID)
}

// GetActiveMedications gets active medications for an animal
func (uc *MedicalUseCase) GetActiveMedications(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error) {
	return uc.medicationRepo.FindActiveByAnimal(ctx, animalID)
}

// GetConditionMedications gets medications for a specific condition
func (uc *MedicalUseCase) GetConditionMedications(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.Medication, error) {
	return uc.medicationRepo.FindByCondition(ctx, conditionID)
}

// GetMedicationsDueForRefill gets medications due for refill
func (uc *MedicalUseCase) GetMedicationsDueForRefill(ctx context.Context) ([]*entities.Medication, error) {
	return uc.medicationRepo.FindDueForRefill(ctx)
}

// GetExpiringSoonMedications gets medications expiring within specified days
func (uc *MedicalUseCase) GetExpiringSoonMedications(ctx context.Context, days int) ([]*entities.Medication, error) {
	return uc.medicationRepo.FindExpiringSoon(ctx, days)
}

// RecordMedicationAdministration records when medication was given
func (uc *MedicalUseCase) RecordMedicationAdministration(ctx context.Context, medicationID primitive.ObjectID, log entities.AdministrationLog) error {
	return uc.medicationRepo.AddAdministrationLog(ctx, medicationID, log)
}

// RefillMedication processes a medication refill
func (uc *MedicalUseCase) RefillMedication(ctx context.Context, medicationID primitive.ObjectID, userID primitive.ObjectID) error {
	medication, err := uc.medicationRepo.FindByID(ctx, medicationID)
	if err != nil {
		return err
	}

	now := time.Now()
	medication.LastRefillDate = &now
	medication.RefillsRemaining--

	// Calculate next refill due date based on frequency (simplified)
	if medication.EndDate != nil {
		nextRefill := now.AddDate(0, 0, 30) // Default to 30 days
		if nextRefill.Before(*medication.EndDate) {
			medication.NextRefillDue = &nextRefill
		}
	}

	if err := uc.medicationRepo.Update(ctx, medication); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "medication",
		EntityID:   &medicationID,
	})

	return nil
}

// ===== Treatment Plan Methods =====

// CreateTreatmentPlan creates a new treatment plan
func (uc *MedicalUseCase) CreateTreatmentPlan(ctx context.Context, plan *entities.TreatmentPlan, userID primitive.ObjectID) error {
	// Verify animal exists
	_, err := uc.animalRepo.FindByID(ctx, plan.AnimalID)
	if err != nil {
		return err
	}

	// Verify condition exists
	_, err = uc.conditionRepo.FindByID(ctx, plan.ConditionID)
	if err != nil {
		return err
	}

	if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionCreate,
		EntityType: "treatment_plan",
		EntityID:   &plan.ID,
	})

	return nil
}

// GetTreatmentPlan retrieves a treatment plan by ID
func (uc *MedicalUseCase) GetTreatmentPlan(ctx context.Context, id primitive.ObjectID) (*entities.TreatmentPlan, error) {
	return uc.treatmentPlanRepo.FindByID(ctx, id)
}

// UpdateTreatmentPlan updates a treatment plan
func (uc *MedicalUseCase) UpdateTreatmentPlan(ctx context.Context, plan *entities.TreatmentPlan, userID primitive.ObjectID) error {
	if err := uc.treatmentPlanRepo.Update(ctx, plan); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "treatment_plan",
		EntityID:   &plan.ID,
	})

	return nil
}

// DeleteTreatmentPlan deletes a treatment plan
func (uc *MedicalUseCase) DeleteTreatmentPlan(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := uc.treatmentPlanRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.treatmentPlanRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionDelete,
		EntityType: "treatment_plan",
		EntityID:   &id,
	})

	return nil
}

// ListTreatmentPlans lists treatment plans with filtering
func (uc *MedicalUseCase) ListTreatmentPlans(ctx context.Context, filter repositories.TreatmentPlanFilter) ([]*entities.TreatmentPlan, int64, error) {
	return uc.treatmentPlanRepo.List(ctx, filter)
}

// GetAnimalTreatmentPlans gets all treatment plans for an animal
func (uc *MedicalUseCase) GetAnimalTreatmentPlans(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	return uc.treatmentPlanRepo.FindByAnimal(ctx, animalID)
}

// GetActiveTreatmentPlans gets active treatment plans for an animal
func (uc *MedicalUseCase) GetActiveTreatmentPlans(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	return uc.treatmentPlanRepo.FindActiveByAnimal(ctx, animalID)
}

// GetConditionTreatmentPlans gets treatment plans for a specific condition
func (uc *MedicalUseCase) GetConditionTreatmentPlans(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	return uc.treatmentPlanRepo.FindByCondition(ctx, conditionID)
}

// AddProgressNote adds a progress note to a treatment plan
func (uc *MedicalUseCase) AddProgressNote(ctx context.Context, planID primitive.ObjectID, note entities.ProgressNote) error {
	return uc.treatmentPlanRepo.AddProgressNote(ctx, planID, note)
}

// CompleteTreatmentPlan marks a treatment plan as completed
func (uc *MedicalUseCase) CompleteTreatmentPlan(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	plan, err := uc.treatmentPlanRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	plan.Status = entities.TreatmentPlanStatusCompleted
	plan.EndDate = &now

	if err := uc.treatmentPlanRepo.Update(ctx, plan); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "treatment_plan",
		EntityID:   &id,
	})

	return nil
}

// ActivateTreatmentPlan activates a draft treatment plan
func (uc *MedicalUseCase) ActivateTreatmentPlan(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	plan, err := uc.treatmentPlanRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	plan.Status = entities.TreatmentPlanStatusActive

	if err := uc.treatmentPlanRepo.Update(ctx, plan); err != nil {
		return err
	}

	// Create audit log
	uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		UserID:     userID,
		Action:     entities.ActionUpdate,
		EntityType: "treatment_plan",
		EntityID:   &id,
	})

	return nil
}
