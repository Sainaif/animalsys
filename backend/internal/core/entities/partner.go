package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerType represents the type of partner
type PartnerType string

const (
	PartnerTypeVeterinarian PartnerType = "veterinarian"
	PartnerTypePetStore     PartnerType = "pet_store"
	PartnerTypeSponsor      PartnerType = "sponsor"
	PartnerTypeOrganization PartnerType = "organization"
	PartnerTypeVolunteer    PartnerType = "volunteer_org"
)

// Partner represents a partner organization or individual
type Partner struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name                string               `bson:"name" json:"name"`
	Type                PartnerType          `bson:"type" json:"type"`
	ContactPerson       string               `bson:"contact_person,omitempty" json:"contact_person,omitempty"`
	Email               string               `bson:"email" json:"email"`
	Phone               string               `bson:"phone,omitempty" json:"phone,omitempty"`
	Address             string               `bson:"address,omitempty" json:"address,omitempty"`
	Website             string               `bson:"website,omitempty" json:"website,omitempty"`
	Description         string               `bson:"description,omitempty" json:"description,omitempty"`
	Services            []string             `bson:"services,omitempty" json:"services,omitempty"`
	AgreementStartDate  *time.Time           `bson:"agreement_start_date,omitempty" json:"agreement_start_date,omitempty"`
	AgreementEndDate    *time.Time           `bson:"agreement_end_date,omitempty" json:"agreement_end_date,omitempty"`
	AgreementDocumentURL string              `bson:"agreement_document_url,omitempty" json:"agreement_document_url,omitempty"`
	Benefits            string               `bson:"benefits,omitempty" json:"benefits,omitempty"`
	CollaborationHistory []Collaboration     `bson:"collaboration_history,omitempty" json:"collaboration_history,omitempty"`
	Active              bool                 `bson:"active" json:"active"`
	Notes               string               `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt           time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time            `bson:"updated_at" json:"updated_at"`
	CreatedBy           string               `bson:"created_by,omitempty" json:"created_by,omitempty"`
}

// Collaboration represents a collaboration activity
type Collaboration struct {
	Date        time.Time `bson:"date" json:"date"`
	Type        string    `bson:"type" json:"type"` // donation, service, event, etc.
	Description string    `bson:"description" json:"description"`
	Value       float64   `bson:"value,omitempty" json:"value,omitempty"`
}

// PartnerCreateRequest represents partner creation
type PartnerCreateRequest struct {
	Name               string      `json:"name" validate:"required"`
	Type               PartnerType `json:"type" validate:"required,oneof=veterinarian pet_store sponsor organization volunteer_org"`
	ContactPerson      string      `json:"contact_person,omitempty"`
	Email              string      `json:"email" validate:"required,email"`
	Phone              string      `json:"phone,omitempty"`
	Address            string      `json:"address,omitempty"`
	Website            string      `json:"website,omitempty"`
	Description        string      `json:"description,omitempty"`
	Services           []string    `json:"services,omitempty"`
	AgreementStartDate *time.Time  `json:"agreement_start_date,omitempty"`
	AgreementEndDate   *time.Time  `json:"agreement_end_date,omitempty"`
	Benefits           string      `json:"benefits,omitempty"`
	Notes              string      `json:"notes,omitempty"`
}

// PartnerUpdateRequest represents partner update
type PartnerUpdateRequest struct {
	Name               string      `json:"name,omitempty"`
	ContactPerson      string      `json:"contact_person,omitempty"`
	Email              string      `json:"email,omitempty" validate:"omitempty,email"`
	Phone              string      `json:"phone,omitempty"`
	Address            string      `json:"address,omitempty"`
	Website            string      `json:"website,omitempty"`
	Description        string      `json:"description,omitempty"`
	Services           []string    `json:"services,omitempty"`
	AgreementStartDate *time.Time  `json:"agreement_start_date,omitempty"`
	AgreementEndDate   *time.Time  `json:"agreement_end_date,omitempty"`
	Benefits           string      `json:"benefits,omitempty"`
	Active             *bool       `json:"active,omitempty"`
	Notes              string      `json:"notes,omitempty"`
}

// NewPartner creates a new partner
func NewPartner(name string, partnerType PartnerType, email, createdBy string) *Partner {
	now := time.Now()
	return &Partner{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Type:      partnerType,
		Email:     email,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
	}
}

// AddCollaboration adds a collaboration record
func (p *Partner) AddCollaboration(collaboration Collaboration) {
	p.CollaborationHistory = append(p.CollaborationHistory, collaboration)
	p.UpdatedAt = time.Now()
}
