package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MedicalConditionRepository defines the interface for medical condition data access
type MedicalConditionRepository interface {
	// Create creates a new medical condition
	Create(ctx context.Context, condition *entities.MedicalCondition) error

	// FindByID finds a medical condition by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.MedicalCondition, error)

	// Update updates a medical condition
	Update(ctx context.Context, condition *entities.MedicalCondition) error

	// Delete deletes a medical condition
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of medical conditions with filtering
	List(ctx context.Context, filter MedicalConditionFilter) ([]*entities.MedicalCondition, int64, error)

	// FindByAnimal finds all medical conditions for an animal
	FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error)

	// FindActiveByAnimal finds active medical conditions for an animal
	FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error)

	// FindChronicConditions finds all chronic conditions
	FindChronicConditions(ctx context.Context) ([]*entities.MedicalCondition, error)

	// EnsureIndexes creates necessary indexes for the medical_conditions collection
	EnsureIndexes(ctx context.Context) error
}

// MedicalConditionFilter defines filter criteria for listing medical conditions
type MedicalConditionFilter struct {
	AnimalID  *primitive.ObjectID
	Status    entities.ConditionStatus
	Severity  entities.ConditionSeverity
	IsChronic *bool
	FromDate  *time.Time
	ToDate    *time.Time
	Limit     int64
	Offset    int64
}

// MedicationRepository defines the interface for medication data access
type MedicationRepository interface {
	// Create creates a new medication record
	Create(ctx context.Context, medication *entities.Medication) error

	// FindByID finds a medication by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Medication, error)

	// Update updates a medication record
	Update(ctx context.Context, medication *entities.Medication) error

	// Delete deletes a medication record
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of medications with filtering
	List(ctx context.Context, filter MedicationFilter) ([]*entities.Medication, int64, error)

	// FindByAnimal finds all medications for an animal
	FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error)

	// FindActiveByAnimal finds active medications for an animal
	FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error)

	// FindByCondition finds medications for a specific condition
	FindByCondition(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.Medication, error)

	// FindDueForRefill finds medications that are due for refill
	FindDueForRefill(ctx context.Context) ([]*entities.Medication, error)

	// FindExpiringSoon finds medications expiring within the specified days
	FindExpiringSoon(ctx context.Context, days int) ([]*entities.Medication, error)

	// AddAdministrationLog adds an administration log entry to a medication
	AddAdministrationLog(ctx context.Context, medicationID primitive.ObjectID, log entities.AdministrationLog) error

	// EnsureIndexes creates necessary indexes for the medications collection
	EnsureIndexes(ctx context.Context) error
}

// MedicationFilter defines filter criteria for listing medications
type MedicationFilter struct {
	AnimalID    *primitive.ObjectID
	ConditionID *primitive.ObjectID
	Status      entities.MedicationStatus
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int64
	Offset      int64
}

// TreatmentPlanRepository defines the interface for treatment plan data access
type TreatmentPlanRepository interface {
	// Create creates a new treatment plan
	Create(ctx context.Context, plan *entities.TreatmentPlan) error

	// FindByID finds a treatment plan by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.TreatmentPlan, error)

	// Update updates a treatment plan
	Update(ctx context.Context, plan *entities.TreatmentPlan) error

	// Delete deletes a treatment plan
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of treatment plans with filtering
	List(ctx context.Context, filter TreatmentPlanFilter) ([]*entities.TreatmentPlan, int64, error)

	// FindByAnimal finds all treatment plans for an animal
	FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error)

	// FindActiveByAnimal finds active treatment plans for an animal
	FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error)

	// FindByCondition finds treatment plans for a specific condition
	FindByCondition(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.TreatmentPlan, error)

	// AddProgressNote adds a progress note to a treatment plan
	AddProgressNote(ctx context.Context, planID primitive.ObjectID, note entities.ProgressNote) error

	// UpdateProcedureStatus updates the status of a procedure in a treatment plan
	UpdateProcedureStatus(ctx context.Context, planID primitive.ObjectID, procedureIndex int, status string, completedDate *time.Time, performedBy *primitive.ObjectID) error

	// UpdateFollowUpStatus updates the status of a follow-up in a treatment plan
	UpdateFollowUpStatus(ctx context.Context, planID primitive.ObjectID, followUpIndex int, completedDate *time.Time, vetVisitID *primitive.ObjectID) error

	// EnsureIndexes creates necessary indexes for the treatment_plans collection
	EnsureIndexes(ctx context.Context) error
}

// TreatmentPlanFilter defines filter criteria for listing treatment plans
type TreatmentPlanFilter struct {
	AnimalID    *primitive.ObjectID
	ConditionID *primitive.ObjectID
	Status      entities.TreatmentPlanStatus
	CreatedBy   *primitive.ObjectID
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int64
	Offset      int64
}
