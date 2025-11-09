package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerType represents the type of partner organization
type PartnerType string

const (
	PartnerTypeRescue          PartnerType = "rescue"
	PartnerTypeShelter         PartnerType = "shelter"
	PartnerTypeVeterinary      PartnerType = "veterinary"
	PartnerTypeFoster          PartnerType = "foster_network"
	PartnerTypeTransport       PartnerType = "transport"
	PartnerTypeSanctuary       PartnerType = "sanctuary"
	PartnerTypeGovernment      PartnerType = "government"
	PartnerTypeCorporate       PartnerType = "corporate"
	PartnerTypeOther           PartnerType = "other"
)

// PartnerStatus represents the status of the partnership
type PartnerStatus string

const (
	PartnerStatusActive    PartnerStatus = "active"
	PartnerStatusInactive  PartnerStatus = "inactive"
	PartnerStatusSuspended PartnerStatus = "suspended"
	PartnerStatusPending   PartnerStatus = "pending"
)

// Partner represents a partner organization
type Partner struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	LegalName   string             `json:"legal_name,omitempty" bson:"legal_name,omitempty"`
	Type        PartnerType        `json:"type" bson:"type"`
	Status      PartnerStatus      `json:"status" bson:"status"`

	// Contact Information
	ContactInfo ContactDetails `json:"contact_info" bson:"contact_info"`
	Address     AddressInfo    `json:"address,omitempty" bson:"address,omitempty"`

	// Primary Contact Person
	PrimaryContact ContactPerson `json:"primary_contact,omitempty" bson:"primary_contact,omitempty"`

	// Partnership Details
	PartnerSince    time.Time  `json:"partner_since" bson:"partner_since"`
	AgreementNumber string     `json:"agreement_number,omitempty" bson:"agreement_number,omitempty"`
	AgreementDate   *time.Time `json:"agreement_date,omitempty" bson:"agreement_date,omitempty"`
	AgreementExpiry *time.Time `json:"agreement_expiry,omitempty" bson:"agreement_expiry,omitempty"`

	// Services
	ServicesProvided []string `json:"services_provided,omitempty" bson:"services_provided,omitempty"`
	Specializations  []string `json:"specializations,omitempty" bson:"specializations,omitempty"` // e.g., "cats", "dogs", "exotic"

	// Capacity
	MaxCapacity     int  `json:"max_capacity,omitempty" bson:"max_capacity,omitempty"`
	CurrentCapacity int  `json:"current_capacity,omitempty" bson:"current_capacity,omitempty"`
	AcceptsIntakes  bool `json:"accepts_intakes" bson:"accepts_intakes"`

	// Financial
	DiscountPercentage float64 `json:"discount_percentage,omitempty" bson:"discount_percentage,omitempty"` // For veterinary partners
	StandardRate       float64 `json:"standard_rate,omitempty" bson:"standard_rate,omitempty"`

	// Rating and Performance
	Rating      float64 `json:"rating" bson:"rating"` // 0-5
	TotalRatings int    `json:"total_ratings" bson:"total_ratings"`

	// Statistics
	TotalTransfersIn  int `json:"total_transfers_in" bson:"total_transfers_in"`
	TotalTransfersOut int `json:"total_transfers_out" bson:"total_transfers_out"`
	SuccessfulPlacements int `json:"successful_placements" bson:"successful_placements"`

	// Documents
	Documents []string `json:"documents,omitempty" bson:"documents,omitempty"` // Document IDs

	// Notes
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`

	// Social Media
	Website     string            `json:"website,omitempty" bson:"website,omitempty"`
	SocialMedia map[string]string `json:"social_media,omitempty" bson:"social_media,omitempty"`

	// Metadata
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// ContactPerson represents a contact person at a partner organization
type ContactPerson struct {
	Name     string `json:"name" bson:"name"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
	Mobile   string `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

// NewPartner creates a new partner
func NewPartner(name string, partnerType PartnerType, createdBy primitive.ObjectID) *Partner {
	now := time.Now()
	return &Partner{
		ID:                  primitive.NewObjectID(),
		Name:                name,
		Type:                partnerType,
		Status:              PartnerStatusPending,
		PartnerSince:        now,
		AcceptsIntakes:      false,
		Rating:              0,
		TotalRatings:        0,
		TotalTransfersIn:    0,
		TotalTransfersOut:   0,
		SuccessfulPlacements: 0,
		ServicesProvided:    []string{},
		Specializations:     []string{},
		Documents:           []string{},
		Tags:                []string{},
		SocialMedia:         make(map[string]string),
		CreatedBy:           createdBy,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// Activate activates the partnership
func (p *Partner) Activate() {
	p.Status = PartnerStatusActive
	p.UpdatedAt = time.Now()
}

// Suspend suspends the partnership
func (p *Partner) Suspend() {
	p.Status = PartnerStatusSuspended
	p.UpdatedAt = time.Now()
}

// Deactivate deactivates the partnership
func (p *Partner) Deactivate() {
	p.Status = PartnerStatusInactive
	p.UpdatedAt = time.Now()
}

// AddRating adds a rating to the partner
func (p *Partner) AddRating(rating float64) {
	if rating < 0 {
		rating = 0
	}
	if rating > 5 {
		rating = 5
	}

	totalScore := p.Rating * float64(p.TotalRatings)
	p.TotalRatings++
	p.Rating = (totalScore + rating) / float64(p.TotalRatings)
	p.UpdatedAt = time.Now()
}

// IncrementTransfersIn increments the incoming transfers counter
func (p *Partner) IncrementTransfersIn() {
	p.TotalTransfersIn++
	p.UpdatedAt = time.Now()
}

// IncrementTransfersOut increments the outgoing transfers counter
func (p *Partner) IncrementTransfersOut() {
	p.TotalTransfersOut++
	p.UpdatedAt = time.Now()
}

// UpdateCapacity updates the current capacity
func (p *Partner) UpdateCapacity(current int) {
	p.CurrentCapacity = current
	p.UpdatedAt = time.Now()
}

// IsAgreementExpired checks if the partnership agreement has expired
func (p *Partner) IsAgreementExpired() bool {
	if p.AgreementExpiry == nil {
		return false
	}
	return time.Now().After(*p.AgreementExpiry)
}

// HasCapacity checks if the partner has capacity for more animals
func (p *Partner) HasCapacity() bool {
	if p.MaxCapacity == 0 {
		return true // No limit set
	}
	return p.CurrentCapacity < p.MaxCapacity
}
