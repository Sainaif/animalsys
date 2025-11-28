package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorType represents the type of donor
type DonorType string

const (
	DonorTypeIndividual  DonorType = "individual"
	DonorTypeOrganization DonorType = "organization"
	DonorTypeCorporate   DonorType = "corporate"
	DonorTypeFoundation  DonorType = "foundation"
)

// DonorStatus represents the status of a donor
type DonorStatus string

const (
	DonorStatusActive   DonorStatus = "active"
	DonorStatusInactive DonorStatus = "inactive"
	DonorStatusLapsed   DonorStatus = "lapsed" // Not donated in a while
	DonorStatusMajor    DonorStatus = "major"  // Major donor
)

// PreferredContact represents preferred contact method
type PreferredContact string

const (
	PreferredContactEmail PreferredContact = "email"
	PreferredContactPhone PreferredContact = "phone"
	PreferredContactMail  PreferredContact = "mail"
	PreferredContactNone  PreferredContact = "none"
)

// DonorContact represents contact information for a donor
type DonorContact struct {
	Email         string `json:"email,omitempty" bson:"email,omitempty"`
	Phone         string `json:"phone,omitempty" bson:"phone,omitempty"`
	AlternatePhone string `json:"alternate_phone,omitempty" bson:"alternate_phone,omitempty"`
	Website       string `json:"website,omitempty" bson:"website,omitempty"`
}

// DonorAddress represents address for a donor
type DonorAddress struct {
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	State   string `json:"state,omitempty" bson:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty" bson:"zip_code,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`
}

// DonorPreferences represents donor communication preferences
type DonorPreferences struct {
	Newsletter       bool             `json:"newsletter" bson:"newsletter"`
	EventInvites     bool             `json:"event_invites" bson:"event_invites"`
	TaxReceipts      bool             `json:"tax_receipts" bson:"tax_receipts"`
	AnnualReport     bool             `json:"annual_report" bson:"annual_report"`
	PreferredContact PreferredContact `json:"preferred_contact" bson:"preferred_contact"`
	Anonymous        bool             `json:"anonymous" bson:"anonymous"` // Don't publish name
	EmailOptIn       bool             `json:"email_opt_in" bson:"email_opt_in"`
	SmsOptIn         bool             `json:"sms_opt_in" bson:"sms_opt_in"`
}

// Donor represents a donor to the foundation
type Donor struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Type       DonorType   `json:"type" bson:"type"`
	Status     DonorStatus `json:"status" bson:"status"`

	// Individual Donor Information
	FirstName  string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName   string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`

	// Organization Information
	OrganizationName string `json:"organization_name,omitempty" bson:"organization_name,omitempty"`
	ContactPerson    string `json:"contact_person,omitempty" bson:"contact_person,omitempty"`
	TaxID            string `json:"tax_id,omitempty" bson:"tax_id,omitempty"` // EIN for organizations

	// Contact Information
	Contact DonorContact `json:"contact" bson:"contact"`
	Address DonorAddress `json:"address" bson:"address"`

	// Preferences
	Preferences DonorPreferences `json:"preferences" bson:"preferences"`

	// Donation History Summary
	FirstDonationDate  *time.Time `json:"first_donation_date,omitempty" bson:"first_donation_date,omitempty"`
	LastDonationDate   *time.Time `json:"last_donation_date,omitempty" bson:"last_donation_date,omitempty"`
	TotalDonated       float64    `json:"total_donated" bson:"total_donated"`
	DonationCount      int        `json:"donation_count" bson:"donation_count"`
	AverageDonation    float64    `json:"average_donation" bson:"average_donation"`
	LargestDonation    float64    `json:"largest_donation" bson:"largest_donation"`

	// Engagement
	VolunteerHours     int        `json:"volunteer_hours" bson:"volunteer_hours"`
	EventsAttended     int        `json:"events_attended" bson:"events_attended"`
	LastContactDate    *time.Time `json:"last_contact_date,omitempty" bson:"last_contact_date,omitempty"`
	Notes              string     `json:"notes,omitempty" bson:"notes,omitempty"`

	// Additional Information
	Tags               []string `json:"tags,omitempty" bson:"tags,omitempty"` // e.g., "major_donor", "monthly_donor", "legacy"
	Source             string   `json:"source,omitempty" bson:"source,omitempty"` // How they found us
	Interests          []string `json:"interests,omitempty" bson:"interests,omitempty"` // Areas of interest

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetFullName returns the full name for individual donors
func (d *Donor) GetFullName() string {
	if d.Type == DonorTypeIndividual {
		return d.FirstName + " " + d.LastName
	}
	return d.OrganizationName
}

// IsActive checks if the donor is active
func (d *Donor) IsActive() bool {
	return d.Status == DonorStatusActive || d.Status == DonorStatusMajor
}

// IsMajorDonor checks if this is a major donor
func (d *Donor) IsMajorDonor() bool {
	return d.Status == DonorStatusMajor || d.TotalDonated >= 10000 // $10,000 threshold
}

// DaysSinceLastDonation returns days since last donation
func (d *Donor) DaysSinceLastDonation() int {
	if d.LastDonationDate == nil {
		return -1
	}
	duration := time.Since(*d.LastDonationDate)
	return int(duration.Hours() / 24)
}

// UpdateDonationStats updates donation statistics
func (d *Donor) UpdateDonationStats(amount float64) {
	now := time.Now()

	if d.FirstDonationDate == nil {
		d.FirstDonationDate = &now
	}
	d.LastDonationDate = &now

	d.TotalDonated += amount
	d.DonationCount++
	d.AverageDonation = d.TotalDonated / float64(d.DonationCount)

	if amount > d.LargestDonation {
		d.LargestDonation = amount
	}

	// Update status based on total donated
	if d.TotalDonated >= 10000 {
		d.Status = DonorStatusMajor
	} else if d.Status != DonorStatusActive {
		d.Status = DonorStatusActive
	}
}

// NewDonor creates a new donor
func NewDonor(donorType DonorType, createdBy primitive.ObjectID) *Donor {
	now := time.Now()
	return &Donor{
		Type:      donorType,
		Status:    DonorStatusActive,
		Preferences: DonorPreferences{
			Newsletter:       true,
			EventInvites:     true,
			TaxReceipts:      true,
			AnnualReport:     true,
			PreferredContact: PreferredContactEmail,
			Anonymous:        false,
		},
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
