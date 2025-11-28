package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *entities.AuditLog) error

	// FindByID finds an audit log entry by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error)

	// List returns a list of audit logs with pagination
	List(ctx context.Context, filter AuditLogFilter) ([]*entities.AuditLog, int64, error)

	// DeleteOlderThan deletes audit logs older than the specified days
	DeleteOlderThan(ctx context.Context, days int) (int64, error)

	// EnsureIndexes creates necessary indexes for the audit_logs collection
	EnsureIndexes(ctx context.Context) error
}

// AuditLogFilter defines filter criteria for listing audit logs
type AuditLogFilter struct {
	UserID     *primitive.ObjectID
	Action     string
	EntityType string
	EntityID   *primitive.ObjectID
	FromDate   *primitive.DateTime
	ToDate     *primitive.DateTime
	Search     string
	Sort       string
	Limit      int64
	Offset     int64
}
