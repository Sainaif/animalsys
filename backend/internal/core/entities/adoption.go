package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdoptionStatus represents the status of an adoption application
type AdoptionStatus string

const (
	AdoptionStatusSubmitted         AdoptionStatus = "submitted"
	AdoptionStatusUnderReview      AdoptionStatus = "under_review"
	AdoptionStatusInterviewScheduled AdoptionStatus = "interview_scheduled"
	AdoptionStatusApproved         AdoptionStatus = "approved"
	AdoptionStatusRejected         AdoptionStatus = "rejected"
	AdoptionStatusCompleted        AdoptionStatus = "completed"
)

// Adoption represents an adoption application
type Adoption struct {
	ID                  primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	AnimalID            string                 `bson:"animal_id" json:"animal_id"`
	AnimalName          string                 `bson:"animal_name" json:"animal_name"`
	ApplicantID         string                 `bson:"applicant_id,omitempty" json:"applicant_id,omitempty"`
	ApplicantName       string                 `bson:"applicant_name" json:"applicant_name"`
	ApplicantEmail      string                 `bson:"applicant_email" json:"applicant_email"`
	ApplicantPhone      string                 `bson:"applicant_phone" json:"applicant_phone"`
	ApplicantAddress    string                 `bson:"applicant_address" json:"applicant_address"`
	ApplicationData     map[string]interface{} `bson:"application_data,omitempty" json:"application_data,omitempty"`
	Status              AdoptionStatus         `bson:"status" json:"status"`
	ApplicationDate     time.Time              `bson:"application_date" json:"application_date"`
	InterviewDate       *time.Time             `bson:"interview_date,omitempty" json:"interview_date,omitempty"`
	ApprovalDate        *time.Time             `bson:"approval_date,omitempty" json:"approval_date,omitempty"`
	RejectionReason     string                 `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	ContractDocumentURL string                 `bson:"contract_document_url,omitempty" json:"contract_document_url,omitempty"`
	AdoptionFee         float64                `bson:"adoption_fee" json:"adoption_fee"`
	FeePaid             bool                   `bson:"fee_paid" json:"fee_paid"`
	Notes               string                 `bson:"notes,omitempty" json:"notes,omitempty"`
	FollowUpDate        *time.Time             `bson:"follow_up_date,omitempty" json:"follow_up_date,omitempty"`
	FollowUpNotes       string                 `bson:"follow_up_notes,omitempty" json:"follow_up_notes,omitempty"`
	CreatedAt           time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time              `bson:"updated_at" json:"updated_at"`
	ProcessedBy         string                 `bson:"processed_by,omitempty" json:"processed_by,omitempty"`
}

// AdoptionCreateRequest represents adoption application submission
type AdoptionCreateRequest struct {
	AnimalID         string                 `json:"animal_id" validate:"required"`
	ApplicantName    string                 `json:"applicant_name" validate:"required,min=2"`
	ApplicantEmail   string                 `json:"applicant_email" validate:"required,email"`
	ApplicantPhone   string                 `json:"applicant_phone" validate:"required"`
	ApplicantAddress string                 `json:"applicant_address" validate:"required"`
	ApplicationData  map[string]interface{} `json:"application_data,omitempty"`
}

// AdoptionUpdateRequest represents adoption update request
type AdoptionUpdateRequest struct {
	Status          AdoptionStatus `json:"status,omitempty"`
	InterviewDate   *time.Time     `json:"interview_date,omitempty"`
	RejectionReason string         `json:"rejection_reason,omitempty"`
	Notes           string         `json:"notes,omitempty"`
	FeePaid         *bool          `json:"fee_paid,omitempty"`
}

// AdoptionFilter represents filters for querying adoptions
type AdoptionFilter struct {
	AnimalID    string
	ApplicantID string
	Status      AdoptionStatus
	StartDate   time.Time
	EndDate     time.Time
	Search      string
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

// NewAdoption creates a new adoption application
func NewAdoption(animalID, animalName, applicantName, applicantEmail, applicantPhone, applicantAddress string, adoptionFee float64) *Adoption {
	now := time.Now()
	return &Adoption{
		ID:              primitive.NewObjectID(),
		AnimalID:        animalID,
		AnimalName:      animalName,
		ApplicantName:   applicantName,
		ApplicantEmail:  applicantEmail,
		ApplicantPhone:  applicantPhone,
		ApplicantAddress: applicantAddress,
		Status:          AdoptionStatusSubmitted,
		ApplicationDate: now,
		AdoptionFee:     adoptionFee,
		FeePaid:         false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Approve approves the adoption
func (a *Adoption) Approve(processedBy string) {
	now := time.Now()
	a.Status = AdoptionStatusApproved
	a.ApprovalDate = &now
	a.ProcessedBy = processedBy
	a.UpdatedAt = now
}

// Reject rejects the adoption
func (a *Adoption) Reject(reason, processedBy string) {
	a.Status = AdoptionStatusRejected
	a.RejectionReason = reason
	a.ProcessedBy = processedBy
	a.UpdatedAt = time.Now()
}

// Complete marks the adoption as completed
func (a *Adoption) Complete() {
	a.Status = AdoptionStatusCompleted
	a.UpdatedAt = time.Now()
}
