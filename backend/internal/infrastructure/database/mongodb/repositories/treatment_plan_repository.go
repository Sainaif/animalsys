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

type treatmentPlanRepository struct {
	collection *mongo.Collection
}

// NewTreatmentPlanRepository creates a new treatment plan repository
func NewTreatmentPlanRepository(db *mongodb.Database) repositories.TreatmentPlanRepository {
	return &treatmentPlanRepository{
		collection: db.Collection("treatment_plans"),
	}
}

func (r *treatmentPlanRepository) Create(ctx context.Context, plan *entities.TreatmentPlan) error {
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	if plan.Goals == nil {
		plan.Goals = []string{}
	}
	if plan.Medications == nil {
		plan.Medications = []primitive.ObjectID{}
	}
	if plan.Procedures == nil {
		plan.Procedures = []entities.PlannedProcedure{}
	}
	if plan.FollowUpSchedule == nil {
		plan.FollowUpSchedule = []entities.TreatmentFollowUpSchedule{}
	}
	if plan.Progress == nil {
		plan.Progress = []entities.ProgressNote{}
	}

	result, err := r.collection.InsertOne(ctx, plan)
	if err != nil {
		return errors.ErrInternalServer
	}

	plan.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *treatmentPlanRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.TreatmentPlan, error) {
	var plan entities.TreatmentPlan
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}

	return &plan, nil
}

func (r *treatmentPlanRepository) Update(ctx context.Context, plan *entities.TreatmentPlan) error {
	plan.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": plan.ID},
		plan,
	)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *treatmentPlanRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *treatmentPlanRepository) List(ctx context.Context, filter repositories.TreatmentPlanFilter) ([]*entities.TreatmentPlan, int64, error) {
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

	if filter.CreatedBy != nil {
		query["created_by"] = *filter.CreatedBy
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

	var plans []*entities.TreatmentPlan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, 0, errors.ErrInternalServer
	}

	return plans, total, nil
}

func (r *treatmentPlanRepository) FindByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"animal_id": animalID},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var plans []*entities.TreatmentPlan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, errors.ErrInternalServer
	}

	return plans, nil
}

func (r *treatmentPlanRepository) FindActiveByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"animal_id": animalID,
			"status":    entities.TreatmentPlanStatusActive,
		},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var plans []*entities.TreatmentPlan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, errors.ErrInternalServer
	}

	return plans, nil
}

func (r *treatmentPlanRepository) FindByCondition(ctx context.Context, conditionID primitive.ObjectID) ([]*entities.TreatmentPlan, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"condition_id": conditionID},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var plans []*entities.TreatmentPlan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, errors.ErrInternalServer
	}

	return plans, nil
}

func (r *treatmentPlanRepository) AddProgressNote(ctx context.Context, planID primitive.ObjectID, note entities.ProgressNote) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": planID},
		bson.M{
			"$push": bson.M{"progress": note},
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

func (r *treatmentPlanRepository) UpdateProcedureStatus(ctx context.Context, planID primitive.ObjectID, procedureIndex int, status string, completedDate *time.Time, performedBy *primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"procedures." + string(rune(procedureIndex)) + ".status":     status,
			"updated_at": time.Now(),
		},
	}

	if completedDate != nil {
		update["$set"].(bson.M)["procedures."+string(rune(procedureIndex))+".completed_date"] = completedDate
	}

	if performedBy != nil {
		update["$set"].(bson.M)["procedures."+string(rune(procedureIndex))+".performed_by"] = performedBy
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": planID}, update)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *treatmentPlanRepository) UpdateFollowUpStatus(ctx context.Context, planID primitive.ObjectID, followUpIndex int, completedDate *time.Time, vetVisitID *primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if completedDate != nil {
		update["$set"].(bson.M)["follow_up_schedule."+string(rune(followUpIndex))+".completed_date"] = completedDate
	}

	if vetVisitID != nil {
		update["$set"].(bson.M)["follow_up_schedule."+string(rune(followUpIndex))+".vet_visit_id"] = vetVisitID
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": planID}, update)
	if err != nil {
		return errors.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *treatmentPlanRepository) EnsureIndexes(ctx context.Context) error {
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
			Keys: bson.D{{Key: "created_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "start_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}, {Key: "status", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
