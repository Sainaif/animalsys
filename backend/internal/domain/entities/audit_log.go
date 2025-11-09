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
	ActionLogin  AuditAction = "login"
	ActionLogout AuditAction = "logout"
	ActionView   AuditAction = "view"
	ActionExport AuditAction = "export"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID     `bson:"user_id" json:"user_id"`
	Action     AuditAction            `bson:"action" json:"action"`
	EntityType string                 `bson:"entity_type" json:"entity_type"` // "user", "animal", "adoption", etc.
	EntityID   *primitive.ObjectID    `bson:"entity_id,omitempty" json:"entity_id,omitempty"`
	Changes    map[string]interface{} `bson:"changes,omitempty" json:"changes,omitempty"` // old vs new values
	IPAddress  string                 `bson:"ip_address" json:"ip_address"`
	UserAgent  string                 `bson:"user_agent" json:"user_agent"`
	Timestamp  time.Time              `bson:"timestamp" json:"timestamp"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(userID primitive.ObjectID, action AuditAction, entityType string, ipAddress, userAgent string) *AuditLog {
	return &AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Timestamp:  time.Now(),
	}
}

// WithEntityID sets the entity ID
func (a *AuditLog) WithEntityID(id primitive.ObjectID) *AuditLog {
	a.EntityID = &id
	return a
}

// WithChanges sets the changes map
func (a *AuditLog) WithChanges(changes map[string]interface{}) *AuditLog {
	a.Changes = changes
	return a
}
