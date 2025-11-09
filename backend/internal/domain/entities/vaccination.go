package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VaccinationType represents common vaccine types
type VaccinationType string

const (
	// Dog Vaccines
	VaccineRabies         VaccinationType = "rabies"
	VaccineDHPP           VaccinationType = "dhpp"           // Distemper, Hepatitis, Parvovirus, Parainfluenza
	VaccineBordetella     VaccinationType = "bordetella"     // Kennel cough
	VaccineLeptospirosis  VaccinationType = "leptospirosis"
	VaccineLyme           VaccinationType = "lyme"
	VaccineCanineInfluenza VaccinationType = "canine_influenza"

	// Cat Vaccines
	VaccineFVRCP          VaccinationType = "fvrcp"          // Feline Viral Rhinotracheitis, Calicivirus, Panleukopenia
	VaccineFeLV           VaccinationType = "felv"           // Feline Leukemia
	VaccineFIP            VaccinationType = "fip"            // Feline Infectious Peritonitis

	// Other
	VaccineOther          VaccinationType = "other"
)

// VaccinationStatus represents the status of a vaccination
type VaccinationStatus string

const (
	VaccinationStatusCurrent  VaccinationStatus = "current"
	VaccinationStatusDue      VaccinationStatus = "due"
	VaccinationStatusOverdue  VaccinationStatus = "overdue"
	VaccinationStatusExpired  VaccinationStatus = "expired"
)

// Vaccination represents a vaccination record for an animal
type Vaccination struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Animal Information
	AnimalID primitive.ObjectID `json:"animal_id" bson:"animal_id"`

	// Vaccine Information
	VaccineType     VaccinationType `json:"vaccine_type" bson:"vaccine_type"`
	VaccineName     string          `json:"vaccine_name" bson:"vaccine_name"`           // Brand or specific name
	Manufacturer    string          `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	LotNumber       string          `json:"lot_number,omitempty" bson:"lot_number,omitempty"`
	DoseNumber      int             `json:"dose_number" bson:"dose_number"`             // 1st, 2nd, 3rd dose, etc.
	TotalDoses      int             `json:"total_doses,omitempty" bson:"total_doses,omitempty"` // Expected total doses

	// Administration
	DateAdministered time.Time `json:"date_administered" bson:"date_administered"`
	NextDueDate      *time.Time `json:"next_due_date,omitempty" bson:"next_due_date,omitempty"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`

	// Veterinary Information
	VeterinarianName string `json:"veterinarian_name" bson:"veterinarian_name"`
	ClinicName       string `json:"clinic_name,omitempty" bson:"clinic_name,omitempty"`
	ClinicAddress    string `json:"clinic_address,omitempty" bson:"clinic_address,omitempty"`
	ClinicPhone      string `json:"clinic_phone,omitempty" bson:"clinic_phone,omitempty"`

	// Additional Details
	AdministrationSite string  `json:"administration_site,omitempty" bson:"administration_site,omitempty"` // e.g., "left shoulder", "right thigh"
	Reactions          string  `json:"reactions,omitempty" bson:"reactions,omitempty"`
	Notes              string  `json:"notes,omitempty" bson:"notes,omitempty"`
	Cost               float64 `json:"cost,omitempty" bson:"cost,omitempty"`

	// Certificate/Documentation
	CertificateNumber string   `json:"certificate_number,omitempty" bson:"certificate_number,omitempty"`
	CertificateURL    string   `json:"certificate_url,omitempty" bson:"certificate_url,omitempty"` // URL to uploaded certificate

	// Reference to Visit (if part of a vet visit)
	VisitID *primitive.ObjectID `json:"visit_id,omitempty" bson:"visit_id,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetStatus returns the current vaccination status based on dates
func (v *Vaccination) GetStatus() VaccinationStatus {
	now := time.Now()

	// Check if expired
	if v.ExpirationDate != nil && v.ExpirationDate.Before(now) {
		return VaccinationStatusExpired
	}

	// Check if next dose is due
	if v.NextDueDate != nil {
		dueDate := *v.NextDueDate

		// Overdue if past due date
		if dueDate.Before(now) {
			return VaccinationStatusOverdue
		}

		// Due if within 30 days
		thirtyDaysFromNow := now.AddDate(0, 0, 30)
		if dueDate.Before(thirtyDaysFromNow) {
			return VaccinationStatusDue
		}
	}

	return VaccinationStatusCurrent
}

// IsDue checks if the vaccination is due or overdue
func (v *Vaccination) IsDue() bool {
	status := v.GetStatus()
	return status == VaccinationStatusDue || status == VaccinationStatusOverdue
}

// IsComplete checks if the vaccination series is complete
func (v *Vaccination) IsComplete() bool {
	if v.TotalDoses == 0 {
		return true // Single dose vaccine
	}
	return v.DoseNumber >= v.TotalDoses
}

// DaysUntilDue returns the number of days until the next dose is due
func (v *Vaccination) DaysUntilDue() int {
	if v.NextDueDate == nil {
		return -1
	}

	duration := time.Until(*v.NextDueDate)
	return int(duration.Hours() / 24)
}

// NewVaccination creates a new vaccination record
func NewVaccination(animalID primitive.ObjectID, vaccineType VaccinationType, vaccineName string, createdBy primitive.ObjectID) *Vaccination {
	now := time.Now()
	return &Vaccination{
		AnimalID:         animalID,
		VaccineType:      vaccineType,
		VaccineName:      vaccineName,
		DoseNumber:       1,
		DateAdministered: now,
		CreatedBy:        createdBy,
		UpdatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
