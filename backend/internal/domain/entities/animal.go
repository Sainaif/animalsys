package entities

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnimalCategory represents the main category of an animal
type AnimalCategory string

const (
	CategoryMammal       AnimalCategory = "mammal"
	CategoryReptile      AnimalCategory = "reptile"
	CategoryBird         AnimalCategory = "bird"
	CategoryAmphibian    AnimalCategory = "amphibian"
	CategoryFish         AnimalCategory = "fish"
	CategoryInvertebrate AnimalCategory = "invertebrate"
	CategoryFarmAnimal   AnimalCategory = "farm_animal"
)

// AnimalStatus represents the current status of an animal
type AnimalStatus string

const (
	AnimalStatusAvailable     AnimalStatus = "available"
	AnimalStatusAdopted       AnimalStatus = "adopted"
	AnimalStatusUnderTreatment AnimalStatus = "under_treatment"
	AnimalStatusQuarantine    AnimalStatus = "quarantine"
	AnimalStatusFostered      AnimalStatus = "fostered"
	AnimalStatusReserved      AnimalStatus = "reserved"
	AnimalStatusDeceased      AnimalStatus = "deceased"
	AnimalStatusTransferred   AnimalStatus = "transferred"
)

// AnimalSex represents the sex of an animal
type AnimalSex string

const (
	SexMale    AnimalSex = "male"
	SexFemale  AnimalSex = "female"
	SexUnknown AnimalSex = "unknown"
)

// AnimalSize represents the size category of an animal
type AnimalSize string

const (
	SizeSmall  AnimalSize = "small"
	SizeMedium AnimalSize = "medium"
	SizeLarge  AnimalSize = "large"
	SizeXLarge AnimalSize = "xlarge"
)

// Temperament represents behavioral characteristics
type Temperament string

const (
	TemperamentFriendly   Temperament = "friendly"
	TemperamentShy        Temperament = "shy"
	TemperamentAggressive Temperament = "aggressive"
	TemperamentPlayful    Temperament = "playful"
	TemperamentCalm       Temperament = "calm"
	TemperamentEnergetic  Temperament = "energetic"
)

// MultilingualName holds names in different languages
type MultilingualName struct {
	English string `json:"en" bson:"en"`
	Polish  string `json:"pl" bson:"pl"`
	Latin   string `json:"latin,omitempty" bson:"latin,omitempty"` // Scientific name
}

// AnimalImages holds image URLs for the animal
type AnimalImages struct {
	Primary     string   `json:"primary" bson:"primary"`                         // Main profile image
	Gallery     []string `json:"gallery,omitempty" bson:"gallery,omitempty"`     // Additional images
	Thumbnails  []string `json:"thumbnails,omitempty" bson:"thumbnails,omitempty"` // Thumbnail versions
}

// MedicalInfo holds medical information about the animal
type MedicalInfo struct {
	Vaccinated       bool      `json:"vaccinated" bson:"vaccinated"`
	Sterilized       bool      `json:"sterilized" bson:"sterilized"`
	Microchipped     bool      `json:"microchipped" bson:"microchipped"`
	MicrochipNumber  string    `json:"microchip_number,omitempty" bson:"microchip_number,omitempty"`
	HealthStatus     string    `json:"health_status" bson:"health_status"` // healthy, sick, injured, recovering
	Medications      []string  `json:"medications,omitempty" bson:"medications,omitempty"`
	Allergies        []string  `json:"allergies,omitempty" bson:"allergies,omitempty"`
	SpecialNeeds     string    `json:"special_needs,omitempty" bson:"special_needs,omitempty"`
	LastVetVisit     *time.Time `json:"last_vet_visit,omitempty" bson:"last_vet_visit,omitempty"`
	NextVetVisit     *time.Time `json:"next_vet_visit,omitempty" bson:"next_vet_visit,omitempty"`
}

// BehaviorInfo holds behavioral information
type BehaviorInfo struct {
	Temperament      []Temperament `json:"temperament" bson:"temperament"`
	GoodWithKids     bool         `json:"good_with_kids" bson:"good_with_kids"`
	GoodWithDogs     bool         `json:"good_with_dogs" bson:"good_with_dogs"`
	GoodWithCats     bool         `json:"good_with_cats" bson:"good_with_cats"`
	HouseTrained     bool         `json:"house_trained" bson:"house_trained"`
	SpecialNeeds     string       `json:"special_needs,omitempty" bson:"special_needs,omitempty"`
	Notes            string       `json:"notes,omitempty" bson:"notes,omitempty"`
}

// ShelterInfo holds shelter-specific information
type ShelterInfo struct {
	IntakeDate       time.Time             `json:"intake_date" bson:"intake_date"`
	IntakeReason     string                `json:"intake_reason,omitempty" bson:"intake_reason,omitempty"` // stray, surrender, rescue, etc.
	Location         string                `json:"location" bson:"location"` // cage/kennel number or area
	AssignedCaretaker *primitive.ObjectID  `json:"assigned_caretaker,omitempty" bson:"assigned_caretaker,omitempty"` // User ID
	DailyNotes       []DailyNote          `json:"daily_notes,omitempty" bson:"daily_notes,omitempty"`
}

// DailyNote represents a daily observation or note about the animal
type DailyNote struct {
	Date      time.Time            `json:"date" bson:"date"`
	Note      string               `json:"note" bson:"note"`
	CreatedBy primitive.ObjectID   `json:"created_by" bson:"created_by"` // User ID
}

// AdoptionInfo holds adoption-related information
type AdoptionInfo struct {
	AdoptionFee       float64  `json:"adoption_fee" bson:"adoption_fee"`
	Requirements      []string `json:"requirements,omitempty" bson:"requirements,omitempty"`
	AdoptionDate      *time.Time `json:"adoption_date,omitempty" bson:"adoption_date,omitempty"`
	AdopterID         *primitive.ObjectID `json:"adopter_id,omitempty" bson:"adopter_id,omitempty"`
}

// Animal represents an animal in the system
type Animal struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Name         MultilingualName   `json:"name" bson:"name"`
	Category     AnimalCategory     `json:"category" bson:"category"`
	Species      string             `json:"species" bson:"species"` // e.g., "dog", "cat", "parrot"
	Breed        string             `json:"breed,omitempty" bson:"breed,omitempty"`
	Sex          AnimalSex          `json:"sex" bson:"sex"`
	Status       AnimalStatus       `json:"status" bson:"status"`

	// Physical Characteristics
	DateOfBirth  *time.Time         `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	AgeEstimated bool              `json:"age_estimated" bson:"age_estimated"`
	Color        string             `json:"color,omitempty" bson:"color,omitempty"`
	Size         AnimalSize         `json:"size,omitempty" bson:"size,omitempty"`
	Weight       float64            `json:"weight,omitempty" bson:"weight,omitempty"` // in kg

	// Description
	Description  MultilingualName   `json:"description" bson:"description"`

	// Media
	Images       AnimalImages       `json:"images" bson:"images"`

	// Detailed Information
	Medical      MedicalInfo        `json:"medical" bson:"medical"`
	Behavior     BehaviorInfo       `json:"behavior" bson:"behavior"`
	Shelter      ShelterInfo        `json:"shelter" bson:"shelter"`
	Adoption     AdoptionInfo       `json:"adoption" bson:"adoption"`

	// Metadata
	CreatedBy    primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy    primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsAvailableForAdoption checks if the animal is available for adoption
func (a *Animal) IsAvailableForAdoption() bool {
	return a.Status == AnimalStatusAvailable
}

// IsAdopted checks if the animal has been adopted
func (a *Animal) IsAdopted() bool {
	return a.Status == AnimalStatusAdopted
}

// CanBeModified checks if the animal record can be modified
func (a *Animal) CanBeModified() bool {
	return a.Status != AnimalStatusDeceased
}

// GetAge calculates the age of the animal in years
func (a *Animal) GetAge() float64 {
	if a.DateOfBirth == nil {
		return 0
	}
	duration := time.Since(*a.DateOfBirth)
	return duration.Hours() / 24 / 365.25
}

// GetAgeString returns a human-readable age string
func (a *Animal) GetAgeString() string {
	if a.DateOfBirth == nil {
		return "Unknown"
	}

	age := a.GetAge()
	years := int(age)
	months := int((age - float64(years)) * 12)

	if years == 0 {
		if months == 0 {
			weeks := int(time.Since(*a.DateOfBirth).Hours() / 24 / 7)
			if weeks == 0 {
				days := int(time.Since(*a.DateOfBirth).Hours() / 24)
				return fmt.Sprintf("%d days", days)
			}
			return fmt.Sprintf("%d weeks", weeks)
		}
		return fmt.Sprintf("%d months", months)
	}

	if months > 0 {
		return fmt.Sprintf("%d years %d months", years, months)
	}
	return fmt.Sprintf("%d years", years)
}

// AddDailyNote adds a note to the animal's daily notes
func (a *Animal) AddDailyNote(note string, userID primitive.ObjectID) {
	dailyNote := DailyNote{
		Date:      time.Now(),
		Note:      note,
		CreatedBy: userID,
	}
	a.Shelter.DailyNotes = append(a.Shelter.DailyNotes, dailyNote)
}

// MarkAsAdopted marks the animal as adopted
func (a *Animal) MarkAsAdopted(adopterID primitive.ObjectID, adoptionDate time.Time) {
	a.Status = AnimalStatusAdopted
	a.Adoption.AdoptionDate = &adoptionDate
	a.Adoption.AdopterID = &adopterID
}
