package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
	ActionRead   AuditAction = "read"
	ActionLogin  AuditAction = "login"
	ActionLogout AuditAction = "logout"
	ActionExport AuditAction = "export"
	ActionImport AuditAction = "import"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID         string                 `bson:"user_id" json:"user_id"`
	UserEmail      string                 `bson:"user_email" json:"user_email"`
	UserRole       string                 `bson:"user_role" json:"user_role"`
	Action         AuditAction            `bson:"action" json:"action"`
	ResourceType   string                 `bson:"resource_type" json:"resource_type"`
	ResourceID     string                 `bson:"resource_id,omitempty" json:"resource_id,omitempty"`
	Description    string                 `bson:"description" json:"description"`
	IPAddress      string                 `bson:"ip_address" json:"ip_address"`
	UserAgent      string                 `bson:"user_agent" json:"user_agent"`
	Before         map[string]interface{} `bson:"before,omitempty" json:"before,omitempty"`
	After          map[string]interface{} `bson:"after,omitempty" json:"after,omitempty"`
	Changes        map[string]interface{} `bson:"changes,omitempty" json:"changes,omitempty"`
	Status         string                 `bson:"status" json:"status"`
	ErrorMessage   string                 `bson:"error_message,omitempty" json:"error_message,omitempty"`
	Timestamp      time.Time              `bson:"timestamp" json:"timestamp"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(userID, userEmail, userRole string, action AuditAction, resourceType, resourceID, description string) *AuditLog {
	return &AuditLog{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		UserEmail:    userEmail,
		UserRole:     userRole,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
		Status:       "success",
		Timestamp:    time.Now(),
	}
}

// SetIPAndUserAgent sets IP address and user agent
func (a *AuditLog) SetIPAndUserAgent(ip, userAgent string) {
	a.IPAddress = ip
	a.UserAgent = userAgent
}

// SetBefore sets the before state
func (a *AuditLog) SetBefore(before map[string]interface{}) {
	a.Before = before
}

// SetAfter sets the after state
func (a *AuditLog) SetAfter(after map[string]interface{}) {
	a.After = after
}

// SetChanges sets the changes
func (a *AuditLog) SetChanges(changes map[string]interface{}) {
	a.Changes = changes
}

// SetError marks the audit log as failed with an error message
func (a *AuditLog) SetError(errorMsg string) {
	a.Status = "failed"
	a.ErrorMessage = errorMsg
}

// AuditLogFilter represents filters for querying audit logs
type AuditLogFilter struct {
	UserID       string
	Action       AuditAction
	ResourceType string
	ResourceID   string
	StartDate    time.Time
	EndDate      time.Time
	Status       string
	Limit        int
	Offset       int
}
