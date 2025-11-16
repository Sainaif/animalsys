package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContactType represents type of contact
type ContactType string

const (
	ContactTypeAdopter   ContactType = "adopter"
	ContactTypeDonor     ContactType = "donor"
	ContactTypeVolunteer ContactType = "volunteer"
	ContactTypePartner   ContactType = "partner"
	ContactTypeVendor    ContactType = "vendor"
	ContactTypeOther     ContactType = "other"
)

// ContactStatus represents lifecycle state
type ContactStatus string

const (
	ContactStatusActive   ContactStatus = "active"
	ContactStatusInactive ContactStatus = "inactive"
	ContactStatusProspect ContactStatus = "prospect"
	ContactStatusArchived ContactStatus = "archived"
)

// Contact represents CRM contact record
type Contact struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	FirstName        string              `bson:"first_name" json:"first_name" validate:"required"`
	LastName         string              `bson:"last_name" json:"last_name" validate:"required"`
	Organization     string              `bson:"organization,omitempty" json:"organization,omitempty"`
	Email            string              `bson:"email,omitempty" json:"email,omitempty"`
	Phone            string              `bson:"phone,omitempty" json:"phone,omitempty"`
	Type             ContactType         `bson:"type" json:"type" validate:"required"`
	Status           ContactStatus       `bson:"status" json:"status" validate:"required"`
	Tags             []string            `bson:"tags,omitempty" json:"tags,omitempty"`
	OwnerID          *primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
	OwnerName        string              `bson:"owner_name,omitempty" json:"owner_name,omitempty"`
	PreferredChannel string              `bson:"preferred_channel,omitempty" json:"preferred_channel,omitempty"`
	LastContactedAt  *time.Time          `bson:"last_contacted_at,omitempty" json:"last_contacted_at,omitempty"`
	NextFollowUpAt   *time.Time          `bson:"next_follow_up_at,omitempty" json:"next_follow_up_at,omitempty"`
	Notes            string              `bson:"notes,omitempty" json:"notes,omitempty"`
	Address          *AddressInfo        `bson:"address,omitempty" json:"address,omitempty"`
	Activities       []ContactActivity   `bson:"activities,omitempty" json:"activities,omitempty"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updated_at"`
}

// ContactActivity logs interactions
type ContactActivity struct {
	ID            primitive.ObjectID `bson:"id" json:"id"`
	Type          string             `bson:"type" json:"type"`
	Subject       string             `bson:"subject" json:"subject"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedByName string             `bson:"created_by_name,omitempty" json:"created_by_name,omitempty"`
	Outcome       string             `bson:"outcome,omitempty" json:"outcome,omitempty"`
}

// NewContact creates contact with defaults
func NewContact(firstName, lastName string, contactType ContactType, status ContactStatus) *Contact {
	return &Contact{
		ID:        primitive.NewObjectID(),
		FirstName: firstName,
		LastName:  lastName,
		Type:      contactType,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      []string{},
	}
}
