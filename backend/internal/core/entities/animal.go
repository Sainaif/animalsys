package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnimalStatus represents the status of an animal
type AnimalStatus string

const (
	AnimalStatusAvailable    AnimalStatus = "available"
	AnimalStatusAdopted      AnimalStatus = "adopted"
	AnimalStatusMedicalCare  AnimalStatus = "medical_care"
	AnimalStatusReserved     AnimalStatus = "reserved"
	AnimalStatusDeceased     AnimalStatus = "deceased"
)

// AnimalGender represents the gender of an animal
type AnimalGender string

const (
	AnimalGenderMale   AnimalGender = "male"
	AnimalGenderFemale AnimalGender = "female"
)

// AnimalSize represents the size of an animal
type AnimalSize string

const (
	AnimalSizeSmall  AnimalSize = "small"
	AnimalSizeMedium AnimalSize = "medium"
	AnimalSizeLarge  AnimalSize = "large"
)

// Animal represents an animal in the shelter
type Animal struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Species          string             `bson:"species" json:"species"`
	Breed            string             `bson:"breed,omitempty" json:"breed,omitempty"`
	Gender           AnimalGender       `bson:"gender" json:"gender"`
	DateOfBirth      *time.Time         `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
	AgeYears         int                `bson:"age_years,omitempty" json:"age_years,omitempty"`
	AgeMonths        int                `bson:"age_months,omitempty" json:"age_months,omitempty"`
	Color            string             `bson:"color,omitempty" json:"color,omitempty"`
	Size             AnimalSize         `bson:"size,omitempty" json:"size,omitempty"`
	Weight           float64            `bson:"weight,omitempty" json:"weight,omitempty"` // in kg
	Status           AnimalStatus       `bson:"status" json:"status"`
	IntakeDate       time.Time          `bson:"intake_date" json:"intake_date"`
	IntakeReason     string             `bson:"intake_reason,omitempty" json:"intake_reason,omitempty"`
	Description      string             `bson:"description,omitempty" json:"description,omitempty"`
	Characteristics  []string           `bson:"characteristics,omitempty" json:"characteristics,omitempty"`
	SpecialNeeds     string             `bson:"special_needs,omitempty" json:"special_needs,omitempty"`
	Photos           []string           `bson:"photos,omitempty" json:"photos,omitempty"`
	MedicalHistory   []MedicalRecord    `bson:"medical_history,omitempty" json:"medical_history,omitempty"`
	Microchipped     bool               `bson:"microchipped" json:"microchipped"`
	MicrochipNumber  string             `bson:"microchip_number,omitempty" json:"microchip_number,omitempty"`
	Neutered         bool               `bson:"neutered" json:"neutered"`
	Vaccinated       bool               `bson:"vaccinated" json:"vaccinated"`
	AdoptionFee      float64            `bson:"adoption_fee,omitempty" json:"adoption_fee,omitempty"`
	AdoptedBy        string             `bson:"adopted_by,omitempty" json:"adopted_by,omitempty"`
	AdoptionDate     *time.Time         `bson:"adoption_date,omitempty" json:"adoption_date,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy        string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy        string             `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// MedicalRecord represents a medical record entry
type MedicalRecord struct {
	Date        time.Time `bson:"date" json:"date"`
	Type        string    `bson:"type" json:"type"` // vaccination, checkup, treatment, surgery, etc.
	Description string    `bson:"description" json:"description"`
	Veterinarian string   `bson:"veterinarian,omitempty" json:"veterinarian,omitempty"`
	Cost        float64   `bson:"cost,omitempty" json:"cost,omitempty"`
	Notes       string    `bson:"notes,omitempty" json:"notes,omitempty"`
}

// AnimalCreateRequest represents animal creation request
type AnimalCreateRequest struct {
	Name            string       `json:"name" validate:"required,min=1,max=100"`
	Species         string       `json:"species" validate:"required"`
	Breed           string       `json:"breed,omitempty"`
	Gender          AnimalGender `json:"gender" validate:"required,oneof=male female"`
	DateOfBirth     *time.Time   `json:"date_of_birth,omitempty"`
	AgeYears        int          `json:"age_years,omitempty" validate:"omitempty,gte=0"`
	AgeMonths       int          `json:"age_months,omitempty" validate:"omitempty,gte=0,lte=11"`
	Color           string       `json:"color,omitempty"`
	Size            AnimalSize   `json:"size,omitempty" validate:"omitempty,oneof=small medium large"`
	Weight          float64      `json:"weight,omitempty" validate:"omitempty,gt=0"`
	IntakeReason    string       `json:"intake_reason,omitempty"`
	Description     string       `json:"description,omitempty"`
	Characteristics []string     `json:"characteristics,omitempty"`
	SpecialNeeds    string       `json:"special_needs,omitempty"`
	Microchipped    bool         `json:"microchipped"`
	MicrochipNumber string       `json:"microchip_number,omitempty"`
	Neutered        bool         `json:"neutered"`
	Vaccinated      bool         `json:"vaccinated"`
	AdoptionFee     float64      `json:"adoption_fee,omitempty" validate:"omitempty,gte=0"`
}

// AnimalUpdateRequest represents animal update request
type AnimalUpdateRequest struct {
	Name            string        `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Breed           string        `json:"breed,omitempty"`
	Color           string        `json:"color,omitempty"`
	Size            AnimalSize    `json:"size,omitempty" validate:"omitempty,oneof=small medium large"`
	Weight          float64       `json:"weight,omitempty" validate:"omitempty,gt=0"`
	Status          AnimalStatus  `json:"status,omitempty" validate:"omitempty,oneof=available adopted medical_care reserved deceased"`
	Description     string        `json:"description,omitempty"`
	Characteristics []string      `json:"characteristics,omitempty"`
	SpecialNeeds    string        `json:"special_needs,omitempty"`
	Microchipped    *bool         `json:"microchipped,omitempty"`
	MicrochipNumber string        `json:"microchip_number,omitempty"`
	Neutered        *bool         `json:"neutered,omitempty"`
	Vaccinated      *bool         `json:"vaccinated,omitempty"`
	AdoptionFee     float64       `json:"adoption_fee,omitempty" validate:"omitempty,gte=0"`
}

// AnimalFilter represents filters for querying animals
type AnimalFilter struct {
	Species     string
	Status      AnimalStatus
	Gender      AnimalGender
	Size        AnimalSize
	MinAge      int
	MaxAge      int
	Neutered    *bool
	Vaccinated  *bool
	Search      string // Search in name, breed, description
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string // asc or desc
}

// NewAnimal creates a new animal
func NewAnimal(name, species string, gender AnimalGender) *Animal {
	now := time.Now()
	return &Animal{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Species:     species,
		Gender:      gender,
		Status:      AnimalStatusAvailable,
		IntakeDate:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddMedicalRecord adds a medical record to the animal
func (a *Animal) AddMedicalRecord(record MedicalRecord) {
	a.MedicalHistory = append(a.MedicalHistory, record)
	a.UpdatedAt = time.Now()
}

// AddPhoto adds a photo URL to the animal
func (a *Animal) AddPhoto(photoURL string) {
	a.Photos = append(a.Photos, photoURL)
	a.UpdatedAt = time.Now()
}

// MarkAsAdopted marks the animal as adopted
func (a *Animal) MarkAsAdopted(adoptedBy string) {
	now := time.Now()
	a.Status = AnimalStatusAdopted
	a.AdoptedBy = adoptedBy
	a.AdoptionDate = &now
	a.UpdatedAt = now
}

// CalculateAge calculates age from date of birth
func (a *Animal) CalculateAge() (years int, months int) {
	if a.DateOfBirth == nil {
		return a.AgeYears, a.AgeMonths
	}

	now := time.Now()
	years = now.Year() - a.DateOfBirth.Year()
	months = int(now.Month() - a.DateOfBirth.Month())

	if months < 0 {
		years--
		months += 12
	}

	return years, months
}
