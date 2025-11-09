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

// vaccinationRepository implements the VaccinationRepository interface
type vaccinationRepository struct {
	db *mongodb.Database
}

// NewVaccinationRepository creates a new vaccination repository
func NewVaccinationRepository(db *mongodb.Database) repositories.VaccinationRepository {
	return &vaccinationRepository{db: db}
}

// Create creates a new vaccination record
func (r *vaccinationRepository) Create(ctx context.Context, vaccination *entities.Vaccination) error {
	vaccination.CreatedAt = time.Now()
	vaccination.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Vaccinations)
	result, err := collection.InsertOne(ctx, vaccination)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create vaccination")
	}

	vaccination.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a vaccination by ID
func (r *vaccinationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Vaccination, error) {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	var vaccination entities.Vaccination
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&vaccination)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find vaccination")
	}

	return &vaccination, nil
}

// Update updates an existing vaccination
func (r *vaccinationRepository) Update(ctx context.Context, vaccination *entities.Vaccination) error {
	vaccination.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Vaccinations)
	filter := bson.M{"_id": vaccination.ID}

	result, err := collection.ReplaceOne(ctx, filter, vaccination)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update vaccination")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a vaccination by ID
func (r *vaccinationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete vaccination")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of vaccinations with pagination and filters
func (r *vaccinationRepository) List(ctx context.Context, filter repositories.VaccinationFilter) ([]*entities.Vaccination, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	// Build filter query
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.VaccineType != "" {
		query["vaccine_type"] = filter.VaccineType
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
		query["date_administered"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count vaccinations")
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
	sortField := "date_administered"
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
		return nil, 0, errors.Wrap(err, 500, "failed to query vaccinations")
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err := cursor.All(ctx, &vaccinations); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode vaccinations")
	}

	return vaccinations, total, nil
}

// GetByAnimalID returns all vaccinations for a specific animal
func (r *vaccinationRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Vaccination, error) {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	findOptions := options.Find().SetSort(bson.D{{Key: "date_administered", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"animal_id": animalID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query vaccinations")
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err := cursor.All(ctx, &vaccinations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode vaccinations")
	}

	return vaccinations, nil
}

// GetDueVaccinations returns vaccinations that are due or overdue
func (r *vaccinationRepository) GetDueVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	// Find vaccinations with next_due_date between now and futureDate
	query := bson.M{
		"next_due_date": bson.M{
			"$lte": futureDate,
			"$gte": now,
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "next_due_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query due vaccinations")
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err := cursor.All(ctx, &vaccinations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode due vaccinations")
	}

	return vaccinations, nil
}

// GetByVaccineType returns vaccinations of a specific type for an animal
func (r *vaccinationRepository) GetByVaccineType(ctx context.Context, animalID primitive.ObjectID, vaccineType entities.VaccinationType) ([]*entities.Vaccination, error) {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	query := bson.M{
		"animal_id":    animalID,
		"vaccine_type": vaccineType,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "date_administered", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query vaccinations by type")
	}
	defer cursor.Close(ctx)

	var vaccinations []*entities.Vaccination
	if err := cursor.All(ctx, &vaccinations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode vaccinations")
	}

	return vaccinations, nil
}

// EnsureIndexes creates necessary indexes for the vaccinations collection
func (r *vaccinationRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Vaccinations)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "vaccine_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "date_administered", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "next_due_date", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "animal_id", Value: 1},
				{Key: "vaccine_type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "animal_id", Value: 1},
				{Key: "date_administered", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
