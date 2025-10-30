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

type reportRepository struct {
	collection *mongo.Collection
}

func NewReportRepository(db *mongo.Database) interfaces.ReportRepository {
	return &reportRepository{
		collection: db.Collection("reports"),
	}
}

func (r *reportRepository) Create(ctx context.Context, report *entities.Report) error {
	report.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, report)
	return err
}

func (r *reportRepository) GetByID(ctx context.Context, id string) (*entities.Report, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid report ID")
	}

	var report entities.Report
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("report not found")
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) Update(ctx context.Context, id string, report *entities.Report) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid report ID")
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": report})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("report not found")
	}
	return nil
}

func (r *reportRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid report ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("report not found")
	}
	return nil
}

func (r *reportRepository) List(ctx context.Context, reportType entities.ReportType, limit, offset int) ([]*entities.Report, int64, error) {
	filter := bson.M{}
	if reportType != "" {
		filter["type"] = reportType
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (r *reportRepository) GetByCreatedBy(ctx context.Context, userID string, limit, offset int) ([]*entities.Report, int64, error) {
	filter := bson.M{"created_by": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []*entities.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (r *reportRepository) DeleteExpired(ctx context.Context) (int64, error) {
	filter := bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *reportRepository) MarkAsCompleted(ctx context.Context, id, fileURL string, fileSize int64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid report ID")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       entities.ReportStatusCompleted,
			"file_url":     fileURL,
			"file_size":    fileSize,
			"generated_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("report not found")
	}
	return nil
}

func (r *reportRepository) MarkAsFailed(ctx context.Context, id, errorMsg string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid report ID")
	}

	update := bson.M{
		"$set": bson.M{
			"status": entities.ReportStatusFailed,
			"error":  errorMsg,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("report not found")
	}
	return nil
}
