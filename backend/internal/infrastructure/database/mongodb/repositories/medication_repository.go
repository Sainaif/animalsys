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

type medicationRepository struct {
	collection *mongo.Collection
}

// NewMedicationRepository creates a new medication repository
func NewMedicationRepository(db *mongodb.Database) repositories.MedicationRepository {
	return &medicationRepository{
		collection: db.Collection("medications"),
	}
}

func (r *medicationRepository) Create(ctx context.Context, medication *entities.Medication) error {
	medication.CreatedAt = time.Now()
	medication.UpdatedAt = time.Now()

	if medication.AdministrationLogs == nil {
		medication.AdministrationLogs = []entities.AdministrationLog{}
	}

	result, err := r.collection.InsertOne(ctx, medication)
	if err != nil {
		return errors.ErrInternalServer
	}

	medication.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *medicationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Medication, error) {
	var medication entities.Medication
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&medication)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}

	return &medication, nil
}

func (r *medicationRepository) Update(ctx context.Context, medication *entities.Medication) error {
	medication.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": medication.ID},
		medication,
	)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *medicationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *medicationRepository) List(ctx context.Context, filter repositories.MedicationFilter) ([]*entities.Medication, int64, error) {
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.ConditionID != nil {
		query["condition_id"] = *filter.ConditionID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.FromDate != nil || filter.ToDate != nil {
		dateFilter := bson.M{}
		if filter.FromDate != nil {
			dateFilter["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			dateFilter["$lte"] = *filter.ToDate
		}
		query["start_date"] = dateFilter
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
		SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, 0, errors.ErrInternalServer
	}

	return medications, total, nil
}

func (r *medicationRepository) FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"animal_id": animalID},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, errors.ErrInternalServer
	}

	return medications, nil
}

func (r *medicationRepository) FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Medication, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"animal_id": animalID,
			"status":    entities.MedicationStatusActive,
		},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, errors.ErrInternalServer
	}

	return medications, nil
}

func (r *medicationRepository) FindByCondition(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.Medication, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"condition_id": conditionID},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, errors.ErrInternalServer
	}

	return medications, nil
}

func (r *medicationRepository) FindDueForRefill(ctx context.Context) ([]*entities.Medication, error) {
	now := time.Now()

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"status": entities.MedicationStatusActive,
			"$or": []bson.M{
				{"next_refill_due": bson.M{"$lte": now}},
				{"refills_remaining": bson.M{"$lte": 1}},
			},
		},
		options.Find().SetSort(bson.D{{Key: "next_refill_due", Value: 1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, errors.ErrInternalServer
	}

	return medications, nil
}

func (r *medicationRepository) FindExpiringSoon(ctx context.Context, days int) ([]*entities.Medication, error) {
	expiryDate := time.Now().AddDate(0, 0, days)

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"status": entities.MedicationStatusActive,
			"end_date": bson.M{
				"$lte": expiryDate,
				"$gte": time.Now(),
			},
		},
		options.Find().SetSort(bson.D{{Key: "end_date", Value: 1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var medications []*entities.Medication
	if err = cursor.All(ctx, &medications); err != nil {
		return nil, errors.ErrInternalServer
	}

	return medications, nil
}

func (r *medicationRepository) AddAdministrationLog(ctx context.Context, medicationID primitive.ObjectID, log entities.AdministrationLog) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": medicationID},
		bson.M{
			"$push": bson.M{"administration_logs": log},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *medicationRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "condition_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "start_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "end_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "next_refill_due", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}, {Key: "status", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
