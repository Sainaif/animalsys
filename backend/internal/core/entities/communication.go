package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommunicationType represents the type of communication
type CommunicationType string

const (
	CommunicationTypeEmail      CommunicationType = "email"
	CommunicationTypeSMS        CommunicationType = "sms"
	CommunicationTypeNewsletter CommunicationType = "newsletter"
)

// CommunicationStatus represents the status of communication
type CommunicationStatus string

const (
	CommunicationStatusDraft     CommunicationStatus = "draft"
	CommunicationStatusScheduled CommunicationStatus = "scheduled"
	CommunicationStatusSent      CommunicationStatus = "sent"
	CommunicationStatusFailed    CommunicationStatus = "failed"
)

// Communication represents a communication (email/SMS/newsletter)
type Communication struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Type          CommunicationType   `bson:"type" json:"type"`
	Status        CommunicationStatus `bson:"status" json:"status"`
	Subject       string              `bson:"subject,omitempty" json:"subject,omitempty"`
	Message       string              `bson:"message" json:"message"`
	TemplateID    string              `bson:"template_id,omitempty" json:"template_id,omitempty"`
	RecipientType string              `bson:"recipient_type" json:"recipient_type"` // all, volunteers, donors, adopters, etc.
	Recipients    []string            `bson:"recipients,omitempty" json:"recipients,omitempty"` // specific email addresses
	RecipientCount int                `bson:"recipient_count" json:"recipient_count"`
	SentCount     int                 `bson:"sent_count" json:"sent_count"`
	FailedCount   int                 `bson:"failed_count" json:"failed_count"`
	ScheduledFor  *time.Time          `bson:"scheduled_for,omitempty" json:"scheduled_for,omitempty"`
	SentAt        *time.Time          `bson:"sent_at,omitempty" json:"sent_at,omitempty"`
	CreatedAt     time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updated_at"`
	CreatedBy     string              `bson:"created_by" json:"created_by"`
}

// CommunicationTemplate represents a communication template
type CommunicationTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        CommunicationType  `bson:"type" json:"type"`
	Subject     string             `bson:"subject,omitempty" json:"subject,omitempty"`
	Body        string             `bson:"body" json:"body"`
	Variables   []string           `bson:"variables,omitempty" json:"variables,omitempty"` // e.g., {name}, {amount}
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
}

// ContactOptIn represents user opt-in preferences
type ContactOptIn struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email           string             `bson:"email" json:"email"`
	EmailOptIn      bool               `bson:"email_opt_in" json:"email_opt_in"`
	SMSOptIn        bool               `bson:"sms_opt_in" json:"sms_opt_in"`
	NewsletterOptIn bool               `bson:"newsletter_opt_in" json:"newsletter_opt_in"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// CommunicationCreateRequest represents communication creation
type CommunicationCreateRequest struct {
	Type          CommunicationType `json:"type" validate:"required,oneof=email sms newsletter"`
	Subject       string            `json:"subject,omitempty"`
	Message       string            `json:"message" validate:"required"`
	TemplateID    string            `json:"template_id,omitempty"`
	RecipientType string            `json:"recipient_type" validate:"required"`
	Recipients    []string          `json:"recipients,omitempty"`
	ScheduledFor  *time.Time        `json:"scheduled_for,omitempty"`
}

// CommunicationTemplateCreateRequest represents template creation
type CommunicationTemplateCreateRequest struct {
	Name        string            `json:"name" validate:"required"`
	Type        CommunicationType `json:"type" validate:"required,oneof=email sms newsletter"`
	Subject     string            `json:"subject,omitempty"`
	Body        string            `json:"body" validate:"required"`
	Variables   []string          `json:"variables,omitempty"`
	Description string            `json:"description,omitempty"`
}

// NewCommunication creates a new communication
func NewCommunication(commType CommunicationType, subject, message, recipientType string, createdBy string) *Communication {
	now := time.Now()
	return &Communication{
		ID:            primitive.NewObjectID(),
		Type:          commType,
		Status:        CommunicationStatusDraft,
		Subject:       subject,
		Message:       message,
		RecipientType: recipientType,
		SentCount:     0,
		FailedCount:   0,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     createdBy,
	}
}

// MarkAsSent marks communication as sent
func (c *Communication) MarkAsSent(sentCount, failedCount int) {
	now := time.Now()
	c.Status = CommunicationStatusSent
	c.SentCount = sentCount
	c.FailedCount = failedCount
	c.SentAt = &now
	c.UpdatedAt = now
}

// Schedule schedules the communication
func (c *Communication) Schedule(scheduledFor time.Time) {
	c.Status = CommunicationStatusScheduled
	c.ScheduledFor = &scheduledFor
	c.UpdatedAt = time.Now()
}
