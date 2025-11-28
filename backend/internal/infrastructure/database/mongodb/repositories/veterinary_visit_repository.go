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

// veterinaryVisitRepository implements the VeterinaryVisitRepository interface
type veterinaryVisitRepository struct {
	db *mongodb.Database
}

// NewVeterinaryVisitRepository creates a new veterinary visit repository
func NewVeterinaryVisitRepository(db *mongodb.Database) repositories.VeterinaryVisitRepository {
	return &veterinaryVisitRepository{db: db}
}

// Create creates a new veterinary visit
func (r *veterinaryVisitRepository) Create(ctx context.Context, visit *entities.VeterinaryVisit) error {
	visit.CreatedAt = time.Now()
	visit.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)
	result, err := collection.InsertOne(ctx, visit)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create veterinary visit")
	}

	visit.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a veterinary visit by ID
func (r *veterinaryVisitRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	var visit entities.VeterinaryVisit
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&visit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find veterinary visit")
	}

	return &visit, nil
}

// Update updates an existing veterinary visit
func (r *veterinaryVisitRepository) Update(ctx context.Context, visit *entities.VeterinaryVisit) error {
	visit.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)
	filter := bson.M{"_id": visit.ID}

	result, err := collection.ReplaceOne(ctx, filter, visit)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update veterinary visit")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *veterinaryVisitRepository) ListCombined(ctx context.Context, filter repositories.CombinedFilter) ([]*entities.VeterinaryRecord, int64, error) {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	matchStage := bson.D{}
	if filter.AnimalID != nil {
		matchStage = append(matchStage, bson.E{Key: "animal_id", Value: *filter.AnimalID})
	}
	if filter.FromDate != nil || filter.ToDate != nil {
		dateFilter := bson.M{}
		if filter.FromDate != nil {
			dateFilter["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			dateFilter["$lte"] = *filter.ToDate
		}
		matchStage = append(matchStage, bson.E{Key: "date", Value: dateFilter})
	}

	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "record_type", Value: "visit"},
				{Key: "date", Value: "$visit_date"},
				{Key: "visit", Value: "$$ROOT"},
			}},
		},
		bson.D{
			{Key: "$unionWith", Value: bson.D{
				{Key: "coll", Value: mongodb.Collections.Vaccinations},
				{Key: "pipeline", Value: mongo.Pipeline{
					bson.D{
						{Key: "$project", Value: bson.D{
							{Key: "record_type", Value: "vaccination"},
							{Key: "date", Value: "$date_administered"},
							{Key: "vaccination", Value: "$$ROOT"},
						}},
					},
				}},
			}},
		},
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: filter.SortBy, Value: sortOrder}}}},
		bson.D{{Key: "$skip", Value: filter.Offset}},
		bson.D{{Key: "$limit", Value: filter.Limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to aggregate records")
	}
	defer cursor.Close(ctx)

	var records []*entities.VeterinaryRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode records")
	}

	// Get total count
	countPipeline := mongo.Pipeline{
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "date", Value: "$visit_date"},
				{Key: "animal_id", Value: "$animal_id"},
			}},
		},
		bson.D{
			{Key: "$unionWith", Value: bson.D{
				{Key: "coll", Value: mongodb.Collections.Vaccinations},
				{Key: "pipeline", Value: mongo.Pipeline{
					bson.D{
						{Key: "$project", Value: bson.D{
							{Key: "date", Value: "$date_administered"},
							{Key: "animal_id", Value: "$animal_id"},
						}},
					},
				}},
			}},
		},
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$count", Value: "total"}},
	}

	countCursor, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count records")
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode count")
	}

	var total int64
	if len(countResult) > 0 {
		total = int64(countResult[0]["total"].(int32))
	}

	return records, total, nil
}

// Delete deletes a veterinary visit by ID
func (r *veterinaryVisitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete veterinary visit")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of veterinary visits with pagination and filters
func (r *veterinaryVisitRepository) List(ctx context.Context, filter repositories.VeterinaryVisitFilter) ([]*entities.VeterinaryVisit, int64, error) {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	// Build filter query
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.VisitType != "" {
		query["visit_type"] = filter.VisitType
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.VeterinarianName != "" {
		query["veterinarian_name"] = bson.M{"$regex": filter.VeterinarianName, "$options": "i"}
	}

	// Date range filter
	if filter.FromDate != nil || filter.ToDate != nil {
		dateFilter := bson.M{}
		if filter.FromDate != nil {
			dateFilter["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			dateFilter["$lte"] = *filter.ToDate
		}
		query["visit_date"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count veterinary visits")
	}

	// Build find options
	findOptions := options.Find()

	if filter.Limit > 0 {
		findOptions.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(filter.Offset)
	}

	// Set sorting
	sortField := "visit_date"
	sortOrder := -1 // descending by default

	if filter.SortBy != "" {
		sortField = filter.SortBy
	}

	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Execute query
	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to query veterinary visits")
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err := cursor.All(ctx, &visits); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode veterinary visits")
	}

	return visits, total, nil
}

// GetByAnimalID returns all visits for a specific animal
func (r *veterinaryVisitRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VeterinaryVisit, error) {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	findOptions := options.Find().SetSort(bson.D{{Key: "visit_date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"animal_id": animalID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query veterinary visits")
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err := cursor.All(ctx, &visits); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode veterinary visits")
	}

	return visits, nil
}

// GetUpcomingVisits returns scheduled visits for the future
func (r *veterinaryVisitRepository) GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	query := bson.M{
		"status": entities.VisitStatusScheduled,
		"scheduled_date": bson.M{
			"$gte": now,
			"$lte": futureDate,
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query upcoming visits")
	}
	defer cursor.Close(ctx)

	var visits []*entities.VeterinaryVisit
	if err := cursor.All(ctx, &visits); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode upcoming visits")
	}

	return visits, nil
}

// EnsureIndexes creates necessary indexes for the veterinary_visits collection
func (r *veterinaryVisitRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.VeterinaryVisits)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "visit_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "visit_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "scheduled_date", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "animal_id", Value: 1},
				{Key: "visit_date", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
