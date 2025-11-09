package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApplicationStatus represents the status of an adoption application
type ApplicationStatus string

const (
	ApplicationStatusPending      ApplicationStatus = "pending"
	ApplicationStatusUnderReview  ApplicationStatus = "under_review"
	ApplicationStatusApproved     ApplicationStatus = "approved"
	ApplicationStatusRejected     ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn    ApplicationStatus = "withdrawn"
	ApplicationStatusCompleted    ApplicationStatus = "completed"
)

// HousingType represents the type of housing
type HousingType string

const (
	HousingTypeHouse     HousingType = "house"
	HousingTypeApartment HousingType = "apartment"
	HousingTypeCondo     HousingType = "condo"
	HousingTypeFarm      HousingType = "farm"
	HousingTypeOther     HousingType = "other"
)

// OwnershipStatus represents housing ownership status
type OwnershipStatus string

const (
	OwnershipOwned  OwnershipStatus = "owned"
	OwnershipRented OwnershipStatus = "rented"
	OwnershipOther  OwnershipStatus = "other"
)

// ApplicantInfo represents information about the adoption applicant
type ApplicantInfo struct {
	FirstName   string    `json:"first_name" bson:"first_name"`
	LastName    string    `json:"last_name" bson:"last_name"`
	Email       string    `json:"email" bson:"email"`
	Phone       string    `json:"phone" bson:"phone"`
	DateOfBirth time.Time `json:"date_of_birth" bson:"date_of_birth"`
	Occupation  string    `json:"occupation,omitempty" bson:"occupation,omitempty"`
	Employer    string    `json:"employer,omitempty" bson:"employer,omitempty"`
}

// AddressInfo represents address details
type AddressInfo struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state,omitempty" bson:"state,omitempty"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
	Country string `json:"country" bson:"country"`
}

// HousingInfo represents housing details
type HousingInfo struct {
	Type              HousingType     `json:"type" bson:"type"`
	Ownership         OwnershipStatus `json:"ownership" bson:"ownership"`
	LandlordName      string          `json:"landlord_name,omitempty" bson:"landlord_name,omitempty"`
	LandlordPhone     string          `json:"landlord_phone,omitempty" bson:"landlord_phone,omitempty"`
	LandlordApproval  bool            `json:"landlord_approval" bson:"landlord_approval"`
	HasYard           bool            `json:"has_yard" bson:"has_yard"`
	YardFenced        bool            `json:"yard_fenced" bson:"yard_fenced"`
	YardSize          string          `json:"yard_size,omitempty" bson:"yard_size,omitempty"`
	AllowsPets        bool            `json:"allows_pets" bson:"allows_pets"`
	PetDeposit        float64         `json:"pet_deposit,omitempty" bson:"pet_deposit,omitempty"`
}

// HouseholdMember represents a member of the household
type HouseholdMember struct {
	Name         string `json:"name" bson:"name"`
	Age          int    `json:"age" bson:"age"`
	Relationship string `json:"relationship" bson:"relationship"`
}

// CurrentPet represents a current pet in the household
type CurrentPet struct {
	Species     string `json:"species" bson:"species"`
	Breed       string `json:"breed,omitempty" bson:"breed,omitempty"`
	Age         int    `json:"age" bson:"age"`
	Name        string `json:"name" bson:"name"`
	Spayed      bool   `json:"spayed" bson:"spayed"`
	Vaccinated  bool   `json:"vaccinated" bson:"vaccinated"`
	VetName     string `json:"vet_name,omitempty" bson:"vet_name,omitempty"`
	VetPhone    string `json:"vet_phone,omitempty" bson:"vet_phone,omitempty"`
}

// Reference represents a personal reference
type Reference struct {
	Name         string `json:"name" bson:"name"`
	Relationship string `json:"relationship" bson:"relationship"`
	Phone        string `json:"phone" bson:"phone"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
	Contacted    bool   `json:"contacted" bson:"contacted"`
	Notes        string `json:"notes,omitempty" bson:"notes,omitempty"`
}

// AdoptionApplication represents an adoption application
type AdoptionApplication struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Animal Information
	AnimalID primitive.ObjectID `json:"animal_id" bson:"animal_id"`

	// Applicant Information
	Applicant ApplicantInfo `json:"applicant" bson:"applicant"`
	Address   AddressInfo   `json:"address" bson:"address"`

	// Housing Information
	Housing HousingInfo `json:"housing" bson:"housing"`

	// Household Information
	HouseholdSize    int               `json:"household_size" bson:"household_size"`
	HouseholdMembers []HouseholdMember `json:"household_members,omitempty" bson:"household_members,omitempty"`
	HasChildren      bool              `json:"has_children" bson:"has_children"`
	ChildrenAges     []int             `json:"children_ages,omitempty" bson:"children_ages,omitempty"`

	// Pet Information
	CurrentPets       []CurrentPet `json:"current_pets,omitempty" bson:"current_pets,omitempty"`
	PreviousPets      string       `json:"previous_pets,omitempty" bson:"previous_pets,omitempty"`
	PetExperience     string       `json:"pet_experience,omitempty" bson:"pet_experience,omitempty"`
	SurrenderedPets   bool         `json:"surrendered_pets" bson:"surrendered_pets"`
	SurrenderReason   string       `json:"surrender_reason,omitempty" bson:"surrender_reason,omitempty"`

	// Adoption Details
	ReasonForAdoption string   `json:"reason_for_adoption" bson:"reason_for_adoption"`
	PetLocation       string   `json:"pet_location" bson:"pet_location"` // Where will pet stay (indoor/outdoor)
	AloneTime         int      `json:"alone_time" bson:"alone_time"`     // Hours per day pet will be alone
	ActivityLevel     string   `json:"activity_level,omitempty" bson:"activity_level,omitempty"`
	PreparedFor       []string `json:"prepared_for,omitempty" bson:"prepared_for,omitempty"` // Costs, training, etc

	// References
	References []Reference `json:"references,omitempty" bson:"references,omitempty"`

	// Veterinarian Information
	HasVeterinarian bool   `json:"has_veterinarian" bson:"has_veterinarian"`
	VetName         string `json:"vet_name,omitempty" bson:"vet_name,omitempty"`
	VetPhone        string `json:"vet_phone,omitempty" bson:"vet_phone,omitempty"`
	VetAddress      string `json:"vet_address,omitempty" bson:"vet_address,omitempty"`

	// Agreement and Consent
	AgreesToHomeVisit     bool `json:"agrees_to_home_visit" bson:"agrees_to_home_visit"`
	AgreesToFollowUp      bool `json:"agrees_to_follow_up" bson:"agrees_to_follow_up"`
	AgreesToReturnPolicy  bool `json:"agrees_to_return_policy" bson:"agrees_to_return_policy"`
	UnderstandsCommitment bool `json:"understands_commitment" bson:"understands_commitment"`

	// Application Status
	Status            ApplicationStatus `json:"status" bson:"status"`
	ApplicationDate   time.Time         `json:"application_date" bson:"application_date"`
	ReviewDate        *time.Time        `json:"review_date,omitempty" bson:"review_date,omitempty"`
	ApprovalDate      *time.Time        `json:"approval_date,omitempty" bson:"approval_date,omitempty"`
	RejectionDate     *time.Time        `json:"rejection_date,omitempty" bson:"rejection_date,omitempty"`
	RejectionReason   string            `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`

	// Review Information
	ReviewedBy      *primitive.ObjectID `json:"reviewed_by,omitempty" bson:"reviewed_by,omitempty"`
	ReviewNotes     string              `json:"review_notes,omitempty" bson:"review_notes,omitempty"`
	HomeVisitDate   *time.Time          `json:"home_visit_date,omitempty" bson:"home_visit_date,omitempty"`
	HomeVisitNotes  string              `json:"home_visit_notes,omitempty" bson:"home_visit_notes,omitempty"`
	InterviewDate   *time.Time          `json:"interview_date,omitempty" bson:"interview_date,omitempty"`
	InterviewNotes  string              `json:"interview_notes,omitempty" bson:"interview_notes,omitempty"`

	// Additional Information
	AdditionalInfo string   `json:"additional_info,omitempty" bson:"additional_info,omitempty"`
	Attachments    []string `json:"attachments,omitempty" bson:"attachments,omitempty"` // URLs to uploaded documents

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsApproved checks if the application is approved
func (a *AdoptionApplication) IsApproved() bool {
	return a.Status == ApplicationStatusApproved
}

// IsRejected checks if the application is rejected
func (a *AdoptionApplication) IsRejected() bool {
	return a.Status == ApplicationStatusRejected
}

// IsPending checks if the application is pending
func (a *AdoptionApplication) IsPending() bool {
	return a.Status == ApplicationStatusPending || a.Status == ApplicationStatusUnderReview
}

// CanBeReviewed checks if the application can be reviewed
func (a *AdoptionApplication) CanBeReviewed() bool {
	return a.Status == ApplicationStatusPending || a.Status == ApplicationStatusUnderReview
}

// NewAdoptionApplication creates a new adoption application
func NewAdoptionApplication(animalID primitive.ObjectID, applicant ApplicantInfo, createdBy primitive.ObjectID) *AdoptionApplication {
	now := time.Now()
	return &AdoptionApplication{
		AnimalID:        animalID,
		Applicant:       applicant,
		Status:          ApplicationStatusPending,
		ApplicationDate: now,
		CreatedBy:       createdBy,
		UpdatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
