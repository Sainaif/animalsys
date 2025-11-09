package repositories

import (
	"context"
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

type reportRepository struct {
	db *mongodb.Database
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *mongodb.Database) repositories.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Reports)
}

// EnsureIndexes creates necessary indexes for reports collection
func (r *reportRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "format", Value: 1}}},
		{Keys: bson.D{{Key: "active", Value: 1}}},
		{Keys: bson.D{{Key: "is_public", Value: 1}}},
		{Keys: bson.D{{Key: "created_by", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
		{Keys: bson.D{
			{Key: "schedule.enabled", Value: 1},
			{Key: "schedule.next_run_at", Value: 1},
		}},
		{Keys: bson.D{{Key: "last_execution_at", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *reportRepository) Create(ctx context.Context, report *entities.Report) error {
	report.ID = primitive.NewObjectID()
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, report)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create report")
	}
	return nil
}

func (r *reportRepository) Update(ctx context.Context, report *entities.Report) error {
	report.UpdatedAt = time.Now()

	filter := bson.M{"_id": report.ID}
	update := bson.M{"$set": report}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update report")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *reportRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete report")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *reportRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Report, error) {
	var report entities.Report
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find report")
	}

	return &report, nil
}

func (r *reportRepository) List(ctx context.Context, filter *repositories.ReportFilter) ([]*entities.Report, int64, error) {
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Format != "" {
		query["format"] = filter.Format
	}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.Active != nil {
		query["active"] = *filter.Active
	}

	if filter.IsPublic != nil {
		query["is_public"] = *filter.IsPublic
	}

	if filter.CreatedBy != nil {
		query["created_by"] = *filter.CreatedBy
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count reports")
	}

	// Find with pagination
	findOptions := options.Find()
	findOptions.SetLimit(filter.Limit)
	findOptions.SetSkip(filter.Offset)

	// Sort
	sortOrder := 1
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to list reports")
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode reports")
	}

	return reports, total, nil
}

func (r *reportRepository) GetActiveReports(ctx context.Context) ([]*entities.Report, error) {
	query := bson.M{"active": true}
	findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get active reports")
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode reports")
	}

	return reports, nil
}

func (r *reportRepository) GetPublicReports(ctx context.Context) ([]*entities.Report, error) {
	query := bson.M{
		"active":    true,
		"is_public": true,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get public reports")
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode reports")
	}

	return reports, nil
}

func (r *reportRepository) GetScheduledReports(ctx context.Context) ([]*entities.Report, error) {
	query := bson.M{
		"active":           true,
		"schedule.enabled": true,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "schedule.next_run_at", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get scheduled reports")
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode reports")
	}

	return reports, nil
}

func (r *reportRepository) GetReportsForExecution(ctx context.Context) ([]*entities.Report, error) {
	now := time.Now()
	query := bson.M{
		"active":                true,
		"schedule.enabled":      true,
		"schedule.next_run_at":  bson.M{"$lte": now},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get reports for execution")
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode reports")
	}

	return reports, nil
}

func (r *reportRepository) IncrementExecutionCount(ctx context.Context, id primitive.ObjectID, status entities.ReportStatus) error {
	filter := bson.M{"_id": id}
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"execution_count": 1},
		"$set": bson.M{
			"last_execution_at":     now,
			"last_execution_status": status,
			"updated_at":            now,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to increment execution count")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}
