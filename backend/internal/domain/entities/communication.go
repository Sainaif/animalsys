package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommunicationStatus represents the status of a communication
type CommunicationStatus string

const (
	CommunicationStatusPending   CommunicationStatus = "pending"
	CommunicationStatusSending   CommunicationStatus = "sending"
	CommunicationStatusSent      CommunicationStatus = "sent"
	CommunicationStatusDelivered CommunicationStatus = "delivered"
	CommunicationStatusFailed    CommunicationStatus = "failed"
	CommunicationStatusBounced   CommunicationStatus = "bounced"
)

// RecipientType represents the type of recipient
type RecipientType string

const (
	RecipientTypeUser      RecipientType = "user"
	RecipientTypeDonor     RecipientType = "donor"
	RecipientTypeVolunteer RecipientType = "volunteer"
	RecipientTypeContact   RecipientType = "contact"
	RecipientTypeExternal  RecipientType = "external"
)

// Communication represents a sent communication (email, SMS, etc.)
type Communication struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Template Reference
	TemplateID *primitive.ObjectID `json:"template_id,omitempty" bson:"template_id,omitempty"`

	// Communication Details
	Type     TemplateType     `json:"type" bson:"type"`
	Category TemplateCategory `json:"category" bson:"category"`
	Status   CommunicationStatus `json:"status" bson:"status"`

	// Recipient Information
	RecipientType  RecipientType       `json:"recipient_type" bson:"recipient_type"`
	RecipientID    *primitive.ObjectID `json:"recipient_id,omitempty" bson:"recipient_id,omitempty"`
	RecipientEmail string              `json:"recipient_email,omitempty" bson:"recipient_email,omitempty"`
	RecipientPhone string              `json:"recipient_phone,omitempty" bson:"recipient_phone,omitempty"`
	RecipientName  string              `json:"recipient_name,omitempty" bson:"recipient_name,omitempty"`

	// Sender Information
	SenderID    primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	FromEmail   string             `json:"from_email,omitempty" bson:"from_email,omitempty"`
	FromName    string             `json:"from_name,omitempty" bson:"from_name,omitempty"`

	// Content
	Subject  string `json:"subject,omitempty" bson:"subject,omitempty"`
	Body     string `json:"body" bson:"body"`
	HTMLBody string `json:"html_body,omitempty" bson:"html_body,omitempty"`

	// Email specific
	CC      []string `json:"cc,omitempty" bson:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty" bson:"bcc,omitempty"`
	ReplyTo string   `json:"reply_to,omitempty" bson:"reply_to,omitempty"`

	// Attachments
	Attachments []CommunicationAttachment `json:"attachments,omitempty" bson:"attachments,omitempty"`

	// Tracking
	SentAt       *time.Time `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" bson:"delivered_at,omitempty"`
	OpenedAt     *time.Time `json:"opened_at,omitempty" bson:"opened_at,omitempty"`
	ClickedAt    *time.Time `json:"clicked_at,omitempty" bson:"clicked_at,omitempty"`
	OpenCount    int        `json:"open_count" bson:"open_count"`
	ClickCount   int        `json:"click_count" bson:"click_count"`

	// Error Handling
	ErrorMessage string     `json:"error_message,omitempty" bson:"error_message,omitempty"`
	RetryCount   int        `json:"retry_count" bson:"retry_count"`
	MaxRetries   int        `json:"max_retries" bson:"max_retries"`
	NextRetryAt  *time.Time `json:"next_retry_at,omitempty" bson:"next_retry_at,omitempty"`

	// Related Resources
	RelatedType string              `json:"related_type,omitempty" bson:"related_type,omitempty"`       // adoption, donation, event, etc.
	RelatedID   *primitive.ObjectID `json:"related_id,omitempty" bson:"related_id,omitempty"`

	// Campaign/Batch Information
	CampaignID *primitive.ObjectID `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"`
	BatchID    string              `json:"batch_id,omitempty" bson:"batch_id,omitempty"`

	// Additional Data
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CommunicationAttachment represents an attachment
type CommunicationAttachment struct {
	Filename    string `json:"filename" bson:"filename"`
	URL         string `json:"url" bson:"url"`
	ContentType string `json:"content_type" bson:"content_type"`
	Size        int64  `json:"size" bson:"size"`
}

// MarkAsSent marks the communication as sent
func (c *Communication) MarkAsSent() {
	now := time.Now()
	c.Status = CommunicationStatusSent
	c.SentAt = &now
	c.UpdatedAt = now
}

// MarkAsDelivered marks the communication as delivered
func (c *Communication) MarkAsDelivered() {
	now := time.Now()
	c.Status = CommunicationStatusDelivered
	c.DeliveredAt = &now
	c.UpdatedAt = now
}

// MarkAsFailed marks the communication as failed
func (c *Communication) MarkAsFailed(errorMessage string) {
	now := time.Now()
	c.Status = CommunicationStatusFailed
	c.ErrorMessage = errorMessage
	c.RetryCount++
	c.UpdatedAt = now

	// Schedule retry if under max retries
	if c.RetryCount < c.MaxRetries {
		nextRetry := now.Add(time.Duration(c.RetryCount*5) * time.Minute)
		c.NextRetryAt = &nextRetry
		c.Status = CommunicationStatusPending
	}
}

// MarkAsOpened marks the communication as opened
func (c *Communication) MarkAsOpened() {
	now := time.Now()
	if c.OpenedAt == nil {
		c.OpenedAt = &now
	}
	c.OpenCount++
	c.UpdatedAt = now
}

// MarkAsClicked marks the communication as clicked
func (c *Communication) MarkAsClicked() {
	now := time.Now()
	if c.ClickedAt == nil {
		c.ClickedAt = &now
	}
	c.ClickCount++
	c.UpdatedAt = now
}

// IsDelivered checks if the communication was delivered
func (c *Communication) IsDelivered() bool {
	return c.Status == CommunicationStatusDelivered || c.Status == CommunicationStatusSent
}

// CanRetry checks if the communication can be retried
func (c *Communication) CanRetry() bool {
	return c.Status == CommunicationStatusFailed && c.RetryCount < c.MaxRetries
}

// NewCommunication creates a new communication
func NewCommunication(
	commType TemplateType,
	category TemplateCategory,
	recipientEmail string,
	subject string,
	body string,
	senderID primitive.ObjectID,
) *Communication {
	now := time.Now()
	return &Communication{
		Type:           commType,
		Category:       category,
		Status:         CommunicationStatusPending,
		RecipientType:  RecipientTypeExternal,
		RecipientEmail: recipientEmail,
		SenderID:       senderID,
		Subject:        subject,
		Body:           body,
		OpenCount:      0,
		ClickCount:     0,
		RetryCount:     0,
		MaxRetries:     3,
		CC:             []string{},
		BCC:            []string{},
		Attachments:    []CommunicationAttachment{},
		Metadata:       make(map[string]interface{}),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
