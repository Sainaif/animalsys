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

type medicalConditionRepository struct {
	collection *mongo.Collection
}

// NewMedicalConditionRepository creates a new medical condition repository
func NewMedicalConditionRepository(db *mongodb.Database) repositories.MedicalConditionRepository {
	return &medicalConditionRepository{
		collection: db.Collection("medical_conditions"),
	}
}

func (r *medicalConditionRepository) Create(ctx context.Context, condition *entities.MedicalCondition) error {
	condition.CreatedAt = time.Now()
	condition.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, condition)
	if err != nil {
		return errors.ErrInternalServer
	}

	condition.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *medicalConditionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.MedicalCondition, error) {
	var condition entities.MedicalCondition
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&condition)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}

	return &condition, nil
}

func (r *medicalConditionRepository) Update(ctx context.Context, condition *entities.MedicalCondition) error {
	condition.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": condition.ID},
		condition,
	)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *medicalConditionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *medicalConditionRepository) List(ctx context.Context, filter repositories.MedicalConditionFilter) ([]*entities.MedicalCondition, int64, error) {
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Severity != "" {
		query["severity"] = filter.Severity
	}

	if filter.IsChronic != nil {
		query["is_chronic"] = *filter.IsChronic
	}

	if filter.FromDate != nil || filter.ToDate != nil {
		dateFilter := bson.M{}
		if filter.FromDate != nil {
			dateFilter["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			dateFilter["$lte"] = *filter.ToDate
		}
		query["diagnosis_date"] = dateFilter
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.ErrInternalServer
	}

	// Set default limit if not specified
	limit := filter.Limit
	if limit == 0 {
		limit = 50
	}

	// Find options
	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(filter.Offset).
		SetSort(bson.D{{Key: "diagnosis_date", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var conditions []*entities.MedicalCondition
	if err = cursor.All(ctx, &conditions); err != nil {
		return nil, 0, errors.ErrInternalServer
	}

	return conditions, total, nil
}

func (r *medicalConditionRepository) FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"animal_id": animalID},
		options.Find().SetSort(bson.D{{Key: "diagnosis_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var conditions []*entities.MedicalCondition
	if err = cursor.All(ctx, &conditions); err != nil {
		return nil, errors.ErrInternalServer
	}

	return conditions, nil
}

func (r *medicalConditionRepository) FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.MedicalCondition, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"animal_id": animalID,
			"status": bson.M{
				"$in": []entities.ConditionStatus{
					entities.ConditionStatusActive,
					entities.ConditionStatusTreating,
					entities.ConditionStatusMonitored,
				},
			},
		},
		options.Find().SetSort(bson.D{{Key: "diagnosis_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var conditions []*entities.MedicalCondition
	if err = cursor.All(ctx, &conditions); err != nil {
		return nil, errors.ErrInternalServer
	}

	return conditions, nil
}

func (r *medicalConditionRepository) FindChronicConditions(ctx context.Context) ([]*entities.MedicalCondition, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"is_chronic": true},
		options.Find().SetSort(bson.D{{Key: "animal_id", Value: 1}, {Key: "diagnosis_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var conditions []*entities.MedicalCondition
	if err = cursor.All(ctx, &conditions); err != nil {
		return nil, errors.ErrInternalServer
	}

	return conditions, nil
}

func (r *medicalConditionRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "diagnosis_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "is_chronic", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}, {Key: "status", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
