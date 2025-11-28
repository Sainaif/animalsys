package repositories

import (
	"context"
	"strings"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// auditLogRepository implements the AuditLogRepository interface
type auditLogRepository struct {
	db *mongodb.Database
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *mongodb.Database) repositories.AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *auditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	collection := r.db.Collection(mongodb.Collections.AuditLogs)
	result, err := collection.InsertOne(ctx, log)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create audit log")
	}

	log.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an audit log entry by ID
func (r *auditLogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	collection := r.db.Collection(mongodb.Collections.AuditLogs)

	var log entities.AuditLog
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find audit log")
	}

	return &log, nil
}

// List returns a list of audit logs with pagination
func (r *auditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	collection := r.db.Collection(mongodb.Collections.AuditLogs)

	// Build MongoDB filter
	mongoFilter := bson.M{}

	if filter.UserID != nil {
		mongoFilter["user_id"] = *filter.UserID
	}
	if filter.Action != "" {
		mongoFilter["action"] = filter.Action
	}
	if filter.EntityType != "" {
		mongoFilter["entity_type"] = filter.EntityType
	}
	if filter.EntityID != nil {
		mongoFilter["entity_id"] = *filter.EntityID
	}
	if filter.FromDate != nil || filter.ToDate != nil {
		timestampFilter := bson.M{}
		if filter.FromDate != nil {
			timestampFilter["$gte"] = filter.FromDate.Time()
		}
		if filter.ToDate != nil {
			timestampFilter["$lte"] = filter.ToDate.Time()
		}
		mongoFilter["timestamp"] = timestampFilter
	}
	if filter.Search != "" {
		searchRegex := bson.M{"$regex": primitive.Regex{Pattern: filter.Search, Options: "i"}}
		mongoFilter["$or"] = []bson.M{
			{"action": searchRegex},
			{"entity_type": searchRegex},
			{"ip_address": searchRegex},
			{"user_agent": searchRegex},
		}
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count audit logs")
	}

	// Sorting
	sortOptions := bson.D{{Key: "timestamp", Value: -1}} // Default sort
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ":")
		field := parts[0]
		direction := -1 // Default desc
		if len(parts) > 1 && strings.ToLower(parts[1]) == "asc" {
			direction = 1
		}
		// Basic validation
		allowedSortFields := map[string]bool{
			"user_id":     true,
			"action":      true,
			"entity_type": true,
			"timestamp":   true,
		}
		if allowedSortFields[field] {
			sortOptions = bson.D{{Key: field, Value: direction}}
		}
	}

	// Find documents with pagination
	opts := options.Find().
		SetSort(sortOptions).
		SetSkip(filter.Offset).
		SetLimit(filter.Limit)

	cursor, err := collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to list audit logs")
	}
	defer cursor.Close(ctx)

	var logs []*entities.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode audit logs")
	}

	return logs, total, nil
}

// DeleteOlderThan deletes audit logs older than the specified days
func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	collection := r.db.Collection(mongodb.Collections.AuditLogs)

	cutoffDate := time.Now().AddDate(0, 0, -days)
	filter := bson.M{"timestamp": bson.M{"$lt": cutoffDate}}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, 500, "failed to delete old audit logs")
	}

	return result.DeletedCount, nil
}

// EnsureIndexes creates necessary indexes for the audit_logs collection
func (r *auditLogRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.AuditLogs)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "action", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "entity_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "entity_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(90 * 24 * 60 * 60), // 90 days TTL
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
