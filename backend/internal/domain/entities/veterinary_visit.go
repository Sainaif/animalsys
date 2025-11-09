package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VisitType represents the type of veterinary visit
type VisitType string

const (
	VisitTypeCheckup        VisitType = "checkup"
	VisitTypeVaccination    VisitType = "vaccination"
	VisitTypeEmergency      VisitType = "emergency"
	VisitTypeSurgery        VisitType = "surgery"
	VisitTypeDental         VisitType = "dental"
	VisitTypeSpayNeuter     VisitType = "spay_neuter"
	VisitTypeFollowUp       VisitType = "follow_up"
	VisitTypeTreatment      VisitType = "treatment"
	VisitTypeDiagnostic     VisitType = "diagnostic"
)

// VisitStatus represents the status of a veterinary visit
type VisitStatus string

const (
	VisitStatusScheduled  VisitStatus = "scheduled"
	VisitStatusInProgress VisitStatus = "in_progress"
	VisitStatusCompleted  VisitStatus = "completed"
	VisitStatusCancelled  VisitStatus = "cancelled"
	VisitStatusNoShow     VisitStatus = "no_show"
)

// VitalSigns represents an animal's vital signs during a visit
type VitalSigns struct {
	Temperature   float64 `json:"temperature,omitempty" bson:"temperature,omitempty"`       // in Celsius
	Weight        float64 `json:"weight,omitempty" bson:"weight,omitempty"`                 // in kg
	HeartRate     int     `json:"heart_rate,omitempty" bson:"heart_rate,omitempty"`         // beats per minute
	RespiratoryRate int   `json:"respiratory_rate,omitempty" bson:"respiratory_rate,omitempty"` // breaths per minute
	BloodPressure string  `json:"blood_pressure,omitempty" bson:"blood_pressure,omitempty"` // e.g., "120/80"
}

// Prescription represents a medication prescription
type Prescription struct {
	MedicationName string    `json:"medication_name" bson:"medication_name"`
	Dosage         string    `json:"dosage" bson:"dosage"`
	Frequency      string    `json:"frequency" bson:"frequency"` // e.g., "twice daily", "every 12 hours"
	Duration       string    `json:"duration" bson:"duration"`   // e.g., "7 days", "2 weeks"
	Instructions   string    `json:"instructions,omitempty" bson:"instructions,omitempty"`
	StartDate      time.Time `json:"start_date" bson:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

// TestResult represents a medical test result
type TestResult struct {
	TestName    string    `json:"test_name" bson:"test_name"`
	Result      string    `json:"result" bson:"result"`
	ReferenceRange string `json:"reference_range,omitempty" bson:"reference_range,omitempty"`
	Notes       string    `json:"notes,omitempty" bson:"notes,omitempty"`
	TestDate    time.Time `json:"test_date" bson:"test_date"`
}

// VeterinaryVisit represents a veterinary visit or examination
type VeterinaryVisit struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	AnimalID     primitive.ObjectID `json:"animal_id" bson:"animal_id"`
	VisitType    VisitType          `json:"visit_type" bson:"visit_type"`
	Status       VisitStatus        `json:"status" bson:"status"`
	VisitDate    time.Time          `json:"visit_date" bson:"visit_date"`
	ScheduledDate *time.Time        `json:"scheduled_date,omitempty" bson:"scheduled_date,omitempty"`

	// Veterinary Information
	VeterinarianName  string `json:"veterinarian_name" bson:"veterinarian_name"`
	ClinicName        string `json:"clinic_name,omitempty" bson:"clinic_name,omitempty"`
	ClinicAddress     string `json:"clinic_address,omitempty" bson:"clinic_address,omitempty"`
	ClinicPhone       string `json:"clinic_phone,omitempty" bson:"clinic_phone,omitempty"`

	// Examination Details
	ChiefComplaint    string      `json:"chief_complaint,omitempty" bson:"chief_complaint,omitempty"` // reason for visit
	VitalSigns        VitalSigns  `json:"vital_signs" bson:"vital_signs"`
	PhysicalExamNotes string      `json:"physical_exam_notes,omitempty" bson:"physical_exam_notes,omitempty"`
	Diagnosis         string      `json:"diagnosis,omitempty" bson:"diagnosis,omitempty"`
	Treatment         string      `json:"treatment,omitempty" bson:"treatment,omitempty"`

	// Medical Details
	Prescriptions     []Prescription `json:"prescriptions,omitempty" bson:"prescriptions,omitempty"`
	TestsOrdered      []string       `json:"tests_ordered,omitempty" bson:"tests_ordered,omitempty"`
	TestResults       []TestResult   `json:"test_results,omitempty" bson:"test_results,omitempty"`
	VaccinationsGiven []string       `json:"vaccinations_given,omitempty" bson:"vaccinations_given,omitempty"`

	// Follow-up
	FollowUpRequired  bool       `json:"follow_up_required" bson:"follow_up_required"`
	FollowUpDate      *time.Time `json:"follow_up_date,omitempty" bson:"follow_up_date,omitempty"`
	FollowUpNotes     string     `json:"follow_up_notes,omitempty" bson:"follow_up_notes,omitempty"`

	// Financial
	Cost              float64    `json:"cost" bson:"cost"`
	PaymentStatus     string     `json:"payment_status,omitempty" bson:"payment_status,omitempty"` // paid, pending, partial

	// Attachments
	Documents         []string   `json:"documents,omitempty" bson:"documents,omitempty"` // URLs to uploaded documents/images

	// Notes
	InternalNotes     string     `json:"internal_notes,omitempty" bson:"internal_notes,omitempty"`

	// Metadata
	CreatedBy         primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy         primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsCompleted checks if the visit is completed
func (v *VeterinaryVisit) IsCompleted() bool {
	return v.Status == VisitStatusCompleted
}

// IsUpcoming checks if the visit is scheduled for the future
func (v *VeterinaryVisit) IsUpcoming() bool {
	return v.Status == VisitStatusScheduled && v.ScheduledDate != nil && v.ScheduledDate.After(time.Now())
}

// RequiresFollowUp checks if a follow-up visit is required
func (v *VeterinaryVisit) RequiresFollowUp() bool {
	return v.FollowUpRequired && v.FollowUpDate != nil
}

// AddPrescription adds a prescription to the visit
func (v *VeterinaryVisit) AddPrescription(prescription Prescription) {
	v.Prescriptions = append(v.Prescriptions, prescription)
}

// AddTestResult adds a test result to the visit
func (v *VeterinaryVisit) AddTestResult(result TestResult) {
	v.TestResults = append(v.TestResults, result)
}

// MarkAsCompleted marks the visit as completed
func (v *VeterinaryVisit) MarkAsCompleted() {
	v.Status = VisitStatusCompleted
	if v.VisitDate.IsZero() {
		v.VisitDate = time.Now()
	}
}
