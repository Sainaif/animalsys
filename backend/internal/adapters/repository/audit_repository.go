package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type auditLogRepository struct {
	collection *mongo.Collection
}

func NewAuditLogRepository(db *mongo.Database) interfaces.AuditLogRepository {
	return &auditLogRepository{
		collection: db.Collection("audit_logs"),
	}
}

func (r *auditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	log.Timestamp = time.Now()
	_, err := r.collection.InsertOne(ctx, log)
	return err
}

func (r *auditLogRepository) GetByID(ctx context.Context, id string) (*entities.AuditLog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid audit log ID")
	}

	var auditLog entities.AuditLog
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&auditLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}
	return &auditLog, nil
}

func (r *auditLogRepository) List(ctx context.Context, filter *entities.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	mongoFilter := bson.M{}

	if filter.UserID != "" {
		mongoFilter["user_id"] = filter.UserID
	}
	if filter.Action != "" {
		mongoFilter["action"] = filter.Action
	}
	if filter.ResourceType != "" {
		mongoFilter["resource_type"] = filter.ResourceType
	}
	if filter.ResourceID != "" {
		mongoFilter["resource_id"] = filter.ResourceID
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		mongoFilter["timestamp"] = dateFilter
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}
	findOptions.SetSort(bson.M{"timestamp": -1})

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*entities.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *auditLogRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.AuditLog, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"timestamp": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*entities.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *auditLogRepository) GetByResourceID(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entities.AuditLog, int64, error) {
	filter := bson.M{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"timestamp": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*entities.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	filter := bson.M{
		"timestamp": bson.M{"$lt": cutoffDate},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
