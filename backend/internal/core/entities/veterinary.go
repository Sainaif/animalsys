package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VeterinaryVisitType represents the type of veterinary visit
type VeterinaryVisitType string

const (
	VisitTypeCheckup     VeterinaryVisitType = "checkup"
	VisitTypeVaccination VeterinaryVisitType = "vaccination"
	VisitTypeTreatment   VeterinaryVisitType = "treatment"
	VisitTypeSurgery     VeterinaryVisitType = "surgery"
	VisitTypeEmergency   VeterinaryVisitType = "emergency"
)

// VeterinaryVisit represents a veterinary visit
type VeterinaryVisit struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	AnimalID      string              `bson:"animal_id" json:"animal_id"`
	AnimalName    string              `bson:"animal_name" json:"animal_name"`
	Type          VeterinaryVisitType `bson:"type" json:"type"`
	VisitDate     time.Time           `bson:"visit_date" json:"visit_date"`
	VeterinarianID string             `bson:"veterinarian_id,omitempty" json:"veterinarian_id,omitempty"`
	Veterinarian  string              `bson:"veterinarian" json:"veterinarian"`
	Clinic        string              `bson:"clinic,omitempty" json:"clinic,omitempty"`
	Diagnosis     string              `bson:"diagnosis,omitempty" json:"diagnosis,omitempty"`
	Treatment     string              `bson:"treatment,omitempty" json:"treatment,omitempty"`
	Medications   []Medication        `bson:"medications,omitempty" json:"medications,omitempty"`
	NextVisitDate *time.Time          `bson:"next_visit_date,omitempty" json:"next_visit_date,omitempty"`
	Cost          float64             `bson:"cost,omitempty" json:"cost,omitempty"`
	Notes         string              `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt     time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updated_at"`
	RecordedBy    string              `bson:"recorded_by" json:"recorded_by"`
}

// Vaccination represents a vaccination record
type Vaccination struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AnimalID        string             `bson:"animal_id" json:"animal_id"`
	AnimalName      string             `bson:"animal_name" json:"animal_name"`
	VaccineName     string             `bson:"vaccine_name" json:"vaccine_name"`
	VaccineType     string             `bson:"vaccine_type" json:"vaccine_type"`
	DateAdministered time.Time         `bson:"date_administered" json:"date_administered"`
	NextDueDate     *time.Time         `bson:"next_due_date,omitempty" json:"next_due_date,omitempty"`
	BatchNumber     string             `bson:"batch_number,omitempty" json:"batch_number,omitempty"`
	Veterinarian    string             `bson:"veterinarian" json:"veterinarian"`
	Clinic          string             `bson:"clinic,omitempty" json:"clinic,omitempty"`
	Reactions       string             `bson:"reactions,omitempty" json:"reactions,omitempty"`
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	RecordedBy      string             `bson:"recorded_by" json:"recorded_by"`
}

// Medication represents a medication
type Medication struct {
	Name        string     `bson:"name" json:"name"`
	Dosage      string     `bson:"dosage" json:"dosage"`
	Frequency   string     `bson:"frequency" json:"frequency"`
	StartDate   time.Time  `bson:"start_date" json:"start_date"`
	EndDate     *time.Time `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Instructions string    `bson:"instructions,omitempty" json:"instructions,omitempty"`
}

// HealthCondition represents an ongoing health condition
type HealthCondition struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AnimalID    string             `bson:"animal_id" json:"animal_id"`
	Condition   string             `bson:"condition" json:"condition"`
	DiagnosedDate time.Time        `bson:"diagnosed_date" json:"diagnosed_date"`
	Status      string             `bson:"status" json:"status"` // active, resolved, chronic
	Treatment   string             `bson:"treatment,omitempty" json:"treatment,omitempty"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// VeterinaryVisitCreateRequest represents veterinary visit creation
type VeterinaryVisitCreateRequest struct {
	AnimalID       string              `json:"animal_id" validate:"required"`
	Type           VeterinaryVisitType `json:"type" validate:"required,oneof=checkup vaccination treatment surgery emergency"`
	VisitDate      time.Time           `json:"visit_date" validate:"required"`
	VeterinarianID string              `json:"veterinarian_id,omitempty"`
	Veterinarian   string              `json:"veterinarian" validate:"required"`
	Clinic         string              `json:"clinic,omitempty"`
	Diagnosis      string              `json:"diagnosis,omitempty"`
	Treatment      string              `json:"treatment,omitempty"`
	Medications    []Medication        `json:"medications,omitempty"`
	NextVisitDate  *time.Time          `json:"next_visit_date,omitempty"`
	Cost           float64             `json:"cost,omitempty" validate:"omitempty,gte=0"`
	Notes          string              `json:"notes,omitempty"`
}

// VaccinationCreateRequest represents vaccination creation
type VaccinationCreateRequest struct {
	AnimalID         string     `json:"animal_id" validate:"required"`
	VaccineName      string     `json:"vaccine_name" validate:"required"`
	VaccineType      string     `json:"vaccine_type" validate:"required"`
	DateAdministered time.Time  `json:"date_administered" validate:"required"`
	NextDueDate      *time.Time `json:"next_due_date,omitempty"`
	BatchNumber      string     `json:"batch_number,omitempty"`
	Veterinarian     string     `json:"veterinarian" validate:"required"`
	Clinic           string     `json:"clinic,omitempty"`
	Reactions        string     `json:"reactions,omitempty"`
	Notes            string     `json:"notes,omitempty"`
}

// NewVeterinaryVisit creates a new veterinary visit
func NewVeterinaryVisit(animalID, animalName string, visitType VeterinaryVisitType, visitDate time.Time, veterinarian, recordedBy string) *VeterinaryVisit {
	now := time.Now()
	return &VeterinaryVisit{
		ID:           primitive.NewObjectID(),
		AnimalID:     animalID,
		AnimalName:   animalName,
		Type:         visitType,
		VisitDate:    visitDate,
		Veterinarian: veterinarian,
		CreatedAt:    now,
		UpdatedAt:    now,
		RecordedBy:   recordedBy,
	}
}

// NewVaccination creates a new vaccination record
func NewVaccination(animalID, animalName, vaccineName, vaccineType string, dateAdministered time.Time, veterinarian, recordedBy string) *Vaccination {
	now := time.Now()
	return &Vaccination{
		ID:               primitive.NewObjectID(),
		AnimalID:         animalID,
		AnimalName:       animalName,
		VaccineName:      vaccineName,
		VaccineType:      vaccineType,
		DateAdministered: dateAdministered,
		Veterinarian:     veterinarian,
		CreatedAt:        now,
		RecordedBy:       recordedBy,
	}
}
