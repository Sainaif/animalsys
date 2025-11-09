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

type reportExecutionRepository struct {
	db *mongodb.Database
}

// NewReportExecutionRepository creates a new report execution repository
func NewReportExecutionRepository(db *mongodb.Database) repositories.ReportExecutionRepository {
	return &reportExecutionRepository{db: db}
}

func (r *reportExecutionRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.ReportExecutions)
}

// EnsureIndexes creates necessary indexes for report_executions collection
func (r *reportExecutionRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "report_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "executed_by", Value: 1}}},
		{Keys: bson.D{
			{Key: "report_id", Value: 1},
			{Key: "started_at", Value: -1},
		}},
		{Keys: bson.D{{Key: "started_at", Value: -1}}},
		{Keys: bson.D{{Key: "completed_at", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *reportExecutionRepository) Create(ctx context.Context, execution *entities.ReportExecution) error {
	execution.ID = primitive.NewObjectID()
	execution.CreatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, execution)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create report execution")
	}
	return nil
}

func (r *reportExecutionRepository) Update(ctx context.Context, execution *entities.ReportExecution) error {
	filter := bson.M{"_id": execution.ID}
	update := bson.M{"$set": execution}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update report execution")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *reportExecutionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete report execution")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *reportExecutionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.ReportExecution, error) {
	var execution entities.ReportExecution
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&execution)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find report execution")
	}

	return &execution, nil
}

func (r *reportExecutionRepository) List(ctx context.Context, filter *repositories.ReportExecutionFilter) ([]*entities.ReportExecution, int64, error) {
	query := bson.M{}

	if filter.ReportID != nil {
		query["report_id"] = *filter.ReportID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.ExecutedBy != nil {
		query["executed_by"] = *filter.ExecutedBy
	}

	if filter.StartDate != nil || filter.EndDate != nil {
		dateQuery := bson.M{}
		if filter.StartDate != nil {
			dateQuery["$gte"] = *filter.StartDate
		}
		if filter.EndDate != nil {
			dateQuery["$lte"] = *filter.EndDate
		}
		query["started_at"] = dateQuery
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count report executions")
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
		sortBy = "started_at"
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to list report executions")
	}
	defer cursor.Close(ctx)

	var executions []*entities.ReportExecution
	if err := cursor.All(ctx, &executions); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode report executions")
	}

	return executions, total, nil
}

func (r *reportExecutionRepository) GetByReportID(ctx context.Context, reportID primitive.ObjectID) ([]*entities.ReportExecution, error) {
	query := bson.M{"report_id": reportID}
	findOptions := options.Find().SetSort(bson.D{{Key: "started_at", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get executions by report ID")
	}
	defer cursor.Close(ctx)

	var executions []*entities.ReportExecution
	if err := cursor.All(ctx, &executions); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode report executions")
	}

	return executions, nil
}

func (r *reportExecutionRepository) GetRecentExecutions(ctx context.Context, limit int) ([]*entities.ReportExecution, error) {
	findOptions := options.Find().
		SetSort(bson.D{{Key: "started_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection().Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get recent executions")
	}
	defer cursor.Close(ctx)

	var executions []*entities.ReportExecution
	if err := cursor.All(ctx, &executions); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode report executions")
	}

	return executions, nil
}

func (r *reportExecutionRepository) DeleteOlderThan(ctx context.Context, date time.Time) (int64, error) {
	filter := bson.M{
		"created_at": bson.M{"$lt": date},
	}

	result, err := r.collection().DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, 500, "Failed to delete old report executions")
	}

	return result.DeletedCount, nil
}
