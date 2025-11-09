package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdoptionStatus represents the status of an adoption
type AdoptionStatus string

const (
	AdoptionStatusPending   AdoptionStatus = "pending"
	AdoptionStatusCompleted AdoptionStatus = "completed"
	AdoptionStatusReturned  AdoptionStatus = "returned"
	AdoptionStatusCancelled AdoptionStatus = "cancelled"
)

// PaymentStatus represents the payment status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPartial   PaymentStatus = "partial"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusWaived    PaymentStatus = "waived"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// FollowUpSchedule represents a follow-up schedule
type FollowUpSchedule struct {
	ScheduledDate time.Time          `json:"scheduled_date" bson:"scheduled_date"`
	CompletedDate *time.Time         `json:"completed_date,omitempty" bson:"completed_date,omitempty"`
	Type          string             `json:"type" bson:"type"` // phone, email, visit
	Notes         string             `json:"notes,omitempty" bson:"notes,omitempty"`
	CompletedBy   *primitive.ObjectID `json:"completed_by,omitempty" bson:"completed_by,omitempty"`
}

// AdoptionContract represents contract details
type AdoptionContract struct {
	ContractURL      string     `json:"contract_url,omitempty" bson:"contract_url,omitempty"`
	SignedDate       *time.Time `json:"signed_date,omitempty" bson:"signed_date,omitempty"`
	SignedBy         string     `json:"signed_by,omitempty" bson:"signed_by,omitempty"`
	WitnessName      string     `json:"witness_name,omitempty" bson:"witness_name,omitempty"`
	WitnessSignature string     `json:"witness_signature,omitempty" bson:"witness_signature,omitempty"`
	Terms            []string   `json:"terms,omitempty" bson:"terms,omitempty"`
}

// Adoption represents a completed or in-progress adoption
type Adoption struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// References
	ApplicationID primitive.ObjectID  `json:"application_id" bson:"application_id"`
	AnimalID      primitive.ObjectID  `json:"animal_id" bson:"animal_id"`
	AdopterID     primitive.ObjectID  `json:"adopter_id" bson:"adopter_id"` // User ID of adopter

	// Adoption Details
	Status       AdoptionStatus `json:"status" bson:"status"`
	AdoptionDate time.Time      `json:"adoption_date" bson:"adoption_date"`
	TrialPeriod  bool           `json:"trial_period" bson:"trial_period"`
	TrialEndDate *time.Time     `json:"trial_end_date,omitempty" bson:"trial_end_date,omitempty"`

	// Financial Information
	AdoptionFee      float64       `json:"adoption_fee" bson:"adoption_fee"`
	PaymentStatus    PaymentStatus `json:"payment_status" bson:"payment_status"`
	AmountPaid       float64       `json:"amount_paid" bson:"amount_paid"`
	PaymentDate      *time.Time    `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	PaymentMethod    string        `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	ReceiptNumber    string        `json:"receipt_number,omitempty" bson:"receipt_number,omitempty"`

	// Contract Information
	Contract              AdoptionContract `json:"contract" bson:"contract"`
	AgreesToReturnPolicy  bool             `json:"agrees_to_return_policy" bson:"agrees_to_return_policy"`
	AgreesToSpayNeuter    bool             `json:"agrees_to_spay_neuter" bson:"agrees_to_spay_neuter"`
	AgreesToMedicalCare   bool             `json:"agrees_to_medical_care" bson:"agrees_to_medical_care"`
	AgreesToFollowUp      bool             `json:"agrees_to_follow_up" bson:"agrees_to_follow_up"`

	// Follow-up Information
	FollowUpSchedule []FollowUpSchedule `json:"follow_up_schedule,omitempty" bson:"follow_up_schedule,omitempty"`
	NextFollowUpDate *time.Time         `json:"next_follow_up_date,omitempty" bson:"next_follow_up_date,omitempty"`

	// Return Information (if applicable)
	ReturnDate   *time.Time `json:"return_date,omitempty" bson:"return_date,omitempty"`
	ReturnReason string     `json:"return_reason,omitempty" bson:"return_reason,omitempty"`
	ReturnNotes  string     `json:"return_notes,omitempty" bson:"return_notes,omitempty"`

	// Additional Information
	Notes       string   `json:"notes,omitempty" bson:"notes,omitempty"`
	Attachments []string `json:"attachments,omitempty" bson:"attachments,omitempty"`

	// Metadata
	ProcessedBy primitive.ObjectID `json:"processed_by" bson:"processed_by"` // Staff who processed
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy   primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsCompleted checks if the adoption is completed
func (a *Adoption) IsCompleted() bool {
	return a.Status == AdoptionStatusCompleted
}

// IsReturned checks if the animal was returned
func (a *Adoption) IsReturned() bool {
	return a.Status == AdoptionStatusReturned
}

// IsInTrialPeriod checks if adoption is in trial period
func (a *Adoption) IsInTrialPeriod() bool {
	if !a.TrialPeriod || a.TrialEndDate == nil {
		return false
	}
	return time.Now().Before(*a.TrialEndDate)
}

// IsPaid checks if the adoption fee is fully paid
func (a *Adoption) IsPaid() bool {
	return a.PaymentStatus == PaymentStatusPaid || a.PaymentStatus == PaymentStatusWaived
}

// AddFollowUp adds a follow-up to the schedule
func (a *Adoption) AddFollowUp(followUp FollowUpSchedule) {
	a.FollowUpSchedule = append(a.FollowUpSchedule, followUp)
	a.UpdateNextFollowUpDate()
}

// UpdateNextFollowUpDate updates the next follow-up date based on pending follow-ups
func (a *Adoption) UpdateNextFollowUpDate() {
	var nextDate *time.Time
	now := time.Now()

	for i := range a.FollowUpSchedule {
		followUp := &a.FollowUpSchedule[i]
		// Skip completed follow-ups
		if followUp.CompletedDate != nil {
			continue
		}
		// Skip past follow-ups
		if followUp.ScheduledDate.Before(now) {
			continue
		}
		// Find the earliest upcoming follow-up
		if nextDate == nil || followUp.ScheduledDate.Before(*nextDate) {
			nextDate = &followUp.ScheduledDate
		}
	}

	a.NextFollowUpDate = nextDate
}

// NewAdoption creates a new adoption record
func NewAdoption(
	applicationID, animalID, adopterID primitive.ObjectID,
	adoptionFee float64,
	processedBy primitive.ObjectID,
) *Adoption {
	now := time.Now()
	return &Adoption{
		ApplicationID:   applicationID,
		AnimalID:        animalID,
		AdopterID:       adopterID,
		Status:          AdoptionStatusPending,
		AdoptionDate:    now,
		AdoptionFee:     adoptionFee,
		PaymentStatus:   PaymentStatusPending,
		AmountPaid:      0,
		ProcessedBy:     processedBy,
		CreatedBy:       processedBy,
		UpdatedBy:       processedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
