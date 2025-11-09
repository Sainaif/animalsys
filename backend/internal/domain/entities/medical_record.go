package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MedicalCondition represents a diagnosed medical condition for an animal
type MedicalCondition struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AnimalID        primitive.ObjectID `bson:"animal_id" json:"animal_id"`
	ConditionName   string             `bson:"condition_name" json:"condition_name"`
	DiagnosisDate   time.Time          `bson:"diagnosis_date" json:"diagnosis_date"`
	DiagnosedBy     primitive.ObjectID `bson:"diagnosed_by" json:"diagnosed_by"` // User ID of diagnosing vet
	Severity        ConditionSeverity  `bson:"severity" json:"severity"`
	Status          ConditionStatus    `bson:"status" json:"status"`
	Description     string             `bson:"description" json:"description"`
	Symptoms        []string           `bson:"symptoms" json:"symptoms"`
	TreatmentNotes  string             `bson:"treatment_notes" json:"treatment_notes"`
	ResolvedDate    *time.Time         `bson:"resolved_date,omitempty" json:"resolved_date,omitempty"`
	IsChronic       bool               `bson:"is_chronic" json:"is_chronic"`
	RequiresMonitor bool               `bson:"requires_monitoring" json:"requires_monitoring"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// ConditionSeverity represents the severity of a medical condition
type ConditionSeverity string

const (
	SeverityMild     ConditionSeverity = "mild"
	SeverityModerate ConditionSeverity = "moderate"
	SeveritySevere   ConditionSeverity = "severe"
	SeverityCritical ConditionSeverity = "critical"
)

// ConditionStatus represents the status of a medical condition
type ConditionStatus string

const (
	ConditionStatusActive    ConditionStatus = "active"
	ConditionStatusTreating  ConditionStatus = "treating"
	ConditionStatusMonitored ConditionStatus = "monitored"
	ConditionStatusResolved  ConditionStatus = "resolved"
	ConditionStatusChronic   ConditionStatus = "chronic"
)

// Medication represents a medication prescribed to an animal
type Medication struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	AnimalID           primitive.ObjectID  `bson:"animal_id" json:"animal_id"`
	ConditionID        *primitive.ObjectID `bson:"condition_id,omitempty" json:"condition_id,omitempty"`
	MedicationName     string              `bson:"medication_name" json:"medication_name"`
	Dosage             string              `bson:"dosage" json:"dosage"`
	Unit               string              `bson:"unit" json:"unit"` // mg, ml, tablets, etc.
	Frequency          string              `bson:"frequency" json:"frequency"`
	Route              MedicationRoute     `bson:"route" json:"route"`
	StartDate          time.Time           `bson:"start_date" json:"start_date"`
	EndDate            *time.Time          `bson:"end_date,omitempty" json:"end_date,omitempty"`
	PrescribedBy       primitive.ObjectID  `bson:"prescribed_by" json:"prescribed_by"` // User ID
	Instructions       string              `bson:"instructions" json:"instructions"`
	SideEffects        []string            `bson:"side_effects" json:"side_effects"`
	Status             MedicationStatus    `bson:"status" json:"status"`
	RefillsRemaining   int                 `bson:"refills_remaining" json:"refills_remaining"`
	LastRefillDate     *time.Time          `bson:"last_refill_date,omitempty" json:"last_refill_date,omitempty"`
	NextRefillDue      *time.Time          `bson:"next_refill_due,omitempty" json:"next_refill_due,omitempty"`
	Cost               float64             `bson:"cost" json:"cost"`
	Notes              string              `bson:"notes" json:"notes"`
	AdministrationLogs []AdministrationLog `bson:"administration_logs" json:"administration_logs"`
	CreatedAt          time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time           `bson:"updated_at" json:"updated_at"`
}

// MedicationRoute represents how medication is administered
type MedicationRoute string

const (
	RouteOral       MedicationRoute = "oral"
	RouteTopical    MedicationRoute = "topical"
	RouteInjection  MedicationRoute = "injection"
	RouteIntravenous MedicationRoute = "intravenous"
	RouteSubcutaneous MedicationRoute = "subcutaneous"
	RouteIntramuscular MedicationRoute = "intramuscular"
	RouteInhalation MedicationRoute = "inhalation"
	RouteOphthalmic MedicationRoute = "ophthalmic"
	RouteOtic       MedicationRoute = "otic"
)

// MedicationStatus represents the current status of a medication
type MedicationStatus string

const (
	MedicationStatusActive      MedicationStatus = "active"
	MedicationStatusCompleted   MedicationStatus = "completed"
	MedicationStatusDiscontinued MedicationStatus = "discontinued"
	MedicationStatusOnHold      MedicationStatus = "on_hold"
)

// AdministrationLog tracks when medication was given
type AdministrationLog struct {
	AdministeredAt time.Time          `bson:"administered_at" json:"administered_at"`
	AdministeredBy primitive.ObjectID `bson:"administered_by" json:"administered_by"` // User ID
	DosageGiven    string             `bson:"dosage_given" json:"dosage_given"`
	Notes          string             `bson:"notes" json:"notes"`
}

// TreatmentPlan represents a comprehensive treatment plan for a condition
type TreatmentPlan struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	AnimalID        primitive.ObjectID   `bson:"animal_id" json:"animal_id"`
	ConditionID     primitive.ObjectID   `bson:"condition_id" json:"condition_id"`
	PlanName        string               `bson:"plan_name" json:"plan_name"`
	Description     string               `bson:"description" json:"description"`
	StartDate       time.Time            `bson:"start_date" json:"start_date"`
	EndDate         *time.Time           `bson:"end_date,omitempty" json:"end_date,omitempty"`
	CreatedBy       primitive.ObjectID   `bson:"created_by" json:"created_by"` // User ID
	Status          TreatmentPlanStatus  `bson:"status" json:"status"`
	Goals           []string             `bson:"goals" json:"goals"`
	Medications     []primitive.ObjectID `bson:"medications" json:"medications"` // Medication IDs
	Procedures      []PlannedProcedure          `bson:"procedures" json:"procedures"`
	DietaryPlan     string                      `bson:"dietary_plan" json:"dietary_plan"`
	ExercisePlan    string                      `bson:"exercise_plan" json:"exercise_plan"`
	MonitoringPlan  string                      `bson:"monitoring_plan" json:"monitoring_plan"`
	FollowUpSchedule []TreatmentFollowUpSchedule `bson:"follow_up_schedule" json:"follow_up_schedule"`
	Notes           string                      `bson:"notes" json:"notes"`
	Progress        []ProgressNote       `bson:"progress" json:"progress"`
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
}

// TreatmentPlanStatus represents the status of a treatment plan
type TreatmentPlanStatus string

const (
	TreatmentPlanStatusDraft      TreatmentPlanStatus = "draft"
	TreatmentPlanStatusActive     TreatmentPlanStatus = "active"
	TreatmentPlanStatusCompleted  TreatmentPlanStatus = "completed"
	TreatmentPlanStatusCancelled  TreatmentPlanStatus = "cancelled"
	TreatmentPlanStatusOnHold     TreatmentPlanStatus = "on_hold"
)

// PlannedProcedure represents a procedure planned as part of treatment
type PlannedProcedure struct {
	ProcedureName   string     `bson:"procedure_name" json:"procedure_name"`
	Description     string     `bson:"description" json:"description"`
	ScheduledDate   time.Time  `bson:"scheduled_date" json:"scheduled_date"`
	CompletedDate   *time.Time `bson:"completed_date,omitempty" json:"completed_date,omitempty"`
	PerformedBy     *primitive.ObjectID `bson:"performed_by,omitempty" json:"performed_by,omitempty"`
	Status          string     `bson:"status" json:"status"` // scheduled, completed, cancelled
	Cost            float64    `bson:"cost" json:"cost"`
	Notes           string     `bson:"notes" json:"notes"`
}

// TreatmentFollowUpSchedule represents scheduled follow-up appointments for treatment
type TreatmentFollowUpSchedule struct {
	ScheduledDate   time.Time           `bson:"scheduled_date" json:"scheduled_date"`
	Purpose         string              `bson:"purpose" json:"purpose"`
	CompletedDate   *time.Time          `bson:"completed_date,omitempty" json:"completed_date,omitempty"`
	VetVisitID      *primitive.ObjectID `bson:"vet_visit_id,omitempty" json:"vet_visit_id,omitempty"`
	Notes           string              `bson:"notes" json:"notes"`
}

// ProgressNote represents a note about treatment progress
type ProgressNote struct {
	Date      time.Time          `bson:"date" json:"date"`
	RecordedBy primitive.ObjectID `bson:"recorded_by" json:"recorded_by"` // User ID
	Note      string             `bson:"note" json:"note"`
	Improvement string           `bson:"improvement" json:"improvement"` // improving, stable, declining
}
