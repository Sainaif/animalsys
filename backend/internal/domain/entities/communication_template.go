package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TemplateType represents the type of communication template
type TemplateType string

const (
	TemplateTypeEmail TemplateType = "email"
	TemplateTypeSMS   TemplateType = "sms"
	TemplateTypePush  TemplateType = "push"
)

// TemplateCategory represents the category of template
type TemplateCategory string

const (
	TemplateCategoryAdoption     TemplateCategory = "adoption"
	TemplateCategoryDonation     TemplateCategory = "donation"
	TemplateCategoryEvent        TemplateCategory = "event"
	TemplateCategoryVolunteer    TemplateCategory = "volunteer"
	TemplateCategoryVeterinary   TemplateCategory = "veterinary"
	TemplateCategoryGeneral      TemplateCategory = "general"
	TemplateCategoryMarketing    TemplateCategory = "marketing"
	TemplateCategoryNotification TemplateCategory = "notification"
)

// CommunicationTemplate represents a template for emails, SMS, etc.
type CommunicationTemplate struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Name        string           `json:"name" bson:"name"`
	Description string           `json:"description,omitempty" bson:"description,omitempty"`
	Type        TemplateType     `json:"type" bson:"type"`
	Category    TemplateCategory `json:"category" bson:"category"`

	// Template Content
	Subject      string `json:"subject,omitempty" bson:"subject,omitempty"`           // For email
	Body         string `json:"body" bson:"body"`                                     // Template body with placeholders
	HTMLBody     string `json:"html_body,omitempty" bson:"html_body,omitempty"`       // HTML version for email
	SMSBody      string `json:"sms_body,omitempty" bson:"sms_body,omitempty"`         // SMS text

	// Template Variables
	Variables    []string          `json:"variables,omitempty" bson:"variables,omitempty"`       // List of available variables
	SampleData   map[string]string `json:"sample_data,omitempty" bson:"sample_data,omitempty"`   // Sample data for preview

	// Settings
	Active       bool   `json:"active" bson:"active"`
	IsDefault    bool   `json:"is_default" bson:"is_default"`                         // Default template for this category
	Language     string `json:"language" bson:"language"`                             // en, pl, etc.

	// Email specific
	FromEmail    string   `json:"from_email,omitempty" bson:"from_email,omitempty"`
	FromName     string   `json:"from_name,omitempty" bson:"from_name,omitempty"`
	ReplyTo      string   `json:"reply_to,omitempty" bson:"reply_to,omitempty"`
	CC           []string `json:"cc,omitempty" bson:"cc,omitempty"`
	BCC          []string `json:"bcc,omitempty" bson:"bcc,omitempty"`

	// Usage Statistics
	UsageCount   int       `json:"usage_count" bson:"usage_count"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty" bson:"last_used_at,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// IncrementUsage increments the usage counter
func (t *CommunicationTemplate) IncrementUsage() {
	t.UsageCount++
	now := time.Now()
	t.LastUsedAt = &now
}

// RenderBody renders the template body with provided data
func (t *CommunicationTemplate) RenderBody(data map[string]string) string {
	body := t.Body
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		body = replaceAll(body, placeholder, value)
	}
	return body
}

// RenderSubject renders the template subject with provided data
func (t *CommunicationTemplate) RenderSubject(data map[string]string) string {
	subject := t.Subject
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		subject = replaceAll(subject, placeholder, value)
	}
	return subject
}

// Validate validates the template
func (t *CommunicationTemplate) Validate() error {
	if t.Name == "" {
		return NewValidationError("Template name is required")
	}
	if t.Type == "" {
		return NewValidationError("Template type is required")
	}
	if t.Category == "" {
		return NewValidationError("Template category is required")
	}
	if t.Body == "" && t.HTMLBody == "" && t.SMSBody == "" {
		return NewValidationError("Template must have at least one body type")
	}
	if t.Type == TemplateTypeEmail && t.Subject == "" {
		return NewValidationError("Email template must have a subject")
	}
	return nil
}

// NewCommunicationTemplate creates a new communication template
func NewCommunicationTemplate(
	name string,
	templateType TemplateType,
	category TemplateCategory,
	body string,
	createdBy primitive.ObjectID,
) *CommunicationTemplate {
	now := time.Now()
	return &CommunicationTemplate{
		Name:       name,
		Type:       templateType,
		Category:   category,
		Body:       body,
		Active:     true,
		IsDefault:  false,
		Language:   "en",
		UsageCount: 0,
		Variables:  []string{},
		SampleData: make(map[string]string),
		CC:         []string{},
		BCC:        []string{},
		CreatedBy:  createdBy,
		UpdatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Helper function to replace all occurrences
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
