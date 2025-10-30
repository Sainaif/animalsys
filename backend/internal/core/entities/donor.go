package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorType represents the type of donor
type DonorType string

const (
	DonorTypeIndividual   DonorType = "individual"
	DonorTypeOrganization DonorType = "organization"
)

// DonationType represents the type of donation
type DonationType string

const (
	DonationTypeOneTime   DonationType = "one_time"
	DonationTypeRecurring DonationType = "recurring"
)

// Donor represents a donor
type Donor struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type            DonorType          `bson:"type" json:"type"`
	FirstName       string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName        string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	OrganizationName string            `bson:"organization_name,omitempty" json:"organization_name,omitempty"`
	Email           string             `bson:"email" json:"email"`
	Phone           string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address         string             `bson:"address,omitempty" json:"address,omitempty"`
	TotalDonated    float64            `bson:"total_donated" json:"total_donated"`
	LastDonation    *time.Time         `bson:"last_donation,omitempty" json:"last_donation,omitempty"`
	DonationCount   int                `bson:"donation_count" json:"donation_count"`
	IsMajorDonor    bool               `bson:"is_major_donor" json:"is_major_donor"`
	TaxID           string             `bson:"tax_id,omitempty" json:"tax_id,omitempty"`
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// Donation represents a donation transaction
type Donation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DonorID         string             `bson:"donor_id" json:"donor_id"`
	DonorName       string             `bson:"donor_name" json:"donor_name"`
	Type            DonationType       `bson:"type" json:"type"`
	Amount          float64            `bson:"amount" json:"amount"`
	Currency        string             `bson:"currency" json:"currency"`
	Date            time.Time          `bson:"date" json:"date"`
	PaymentMethod   string             `bson:"payment_method,omitempty" json:"payment_method,omitempty"`
	Purpose         string             `bson:"purpose,omitempty" json:"purpose,omitempty"`
	Campaign        string             `bson:"campaign,omitempty" json:"campaign,omitempty"`
	TaxReceiptSent  bool               `bson:"tax_receipt_sent" json:"tax_receipt_sent"`
	TaxReceiptURL   string             `bson:"tax_receipt_url,omitempty" json:"tax_receipt_url,omitempty"`
	ThankYouSent    bool               `bson:"thank_you_sent" json:"thank_you_sent"`
	RecurringPlan   string             `bson:"recurring_plan,omitempty" json:"recurring_plan,omitempty"` // monthly, quarterly, yearly
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	RecordedBy      string             `bson:"recorded_by" json:"recorded_by"`
}

// DonorCreateRequest represents donor creation
type DonorCreateRequest struct {
	Type             DonorType `json:"type" validate:"required,oneof=individual organization"`
	FirstName        string    `json:"first_name,omitempty"`
	LastName         string    `json:"last_name,omitempty"`
	OrganizationName string    `json:"organization_name,omitempty"`
	Email            string    `json:"email" validate:"required,email"`
	Phone            string    `json:"phone,omitempty"`
	Address          string    `json:"address,omitempty"`
	TaxID            string    `json:"tax_id,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

// DonationCreateRequest represents donation creation
type DonationCreateRequest struct {
	DonorID       string       `json:"donor_id" validate:"required"`
	Type          DonationType `json:"type" validate:"required,oneof=one_time recurring"`
	Amount        float64      `json:"amount" validate:"required,gt=0"`
	Currency      string       `json:"currency" validate:"required,len=3"`
	Date          time.Time    `json:"date" validate:"required"`
	PaymentMethod string       `json:"payment_method,omitempty"`
	Purpose       string       `json:"purpose,omitempty"`
	Campaign      string       `json:"campaign,omitempty"`
	RecurringPlan string       `json:"recurring_plan,omitempty"`
	Notes         string       `json:"notes,omitempty"`
}

// NewDonor creates a new donor
func NewDonor(donorType DonorType, email string) *Donor {
	now := time.Now()
	return &Donor{
		ID:            primitive.NewObjectID(),
		Type:          donorType,
		Email:         email,
		TotalDonated:  0,
		DonationCount: 0,
		IsMajorDonor:  false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewDonation creates a new donation
func NewDonation(donorID, donorName string, donationType DonationType, amount float64, currency string, date time.Time, recordedBy string) *Donation {
	now := time.Now()
	return &Donation{
		ID:             primitive.NewObjectID(),
		DonorID:        donorID,
		DonorName:      donorName,
		Type:           donationType,
		Amount:         amount,
		Currency:       currency,
		Date:           date,
		TaxReceiptSent: false,
		ThankYouSent:   false,
		CreatedAt:      now,
		RecordedBy:     recordedBy,
	}
}

// AddDonation adds a donation to donor's record
func (d *Donor) AddDonation(amount float64) {
	now := time.Now()
	d.TotalDonated += amount
	d.DonationCount++
	d.LastDonation = &now
	d.UpdatedAt = now

	// Check if major donor (e.g., donated more than 5000)
	if d.TotalDonated >= 5000 {
		d.IsMajorDonor = true
	}
}
