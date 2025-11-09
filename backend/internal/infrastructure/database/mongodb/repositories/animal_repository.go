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

// animalRepository implements the AnimalRepository interface
type animalRepository struct {
	db *mongodb.Database
}

// NewAnimalRepository creates a new animal repository
func NewAnimalRepository(db *mongodb.Database) repositories.AnimalRepository {
	return &animalRepository{db: db}
}

// Create creates a new animal
func (r *animalRepository) Create(ctx context.Context, animal *entities.Animal) error {
	animal.CreatedAt = time.Now()
	animal.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Animals)
	result, err := collection.InsertOne(ctx, animal)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create animal")
	}

	animal.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an animal by ID
func (r *animalRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Animal, error) {
	collection := r.db.Collection(mongodb.Collections.Animals)

	var animal entities.Animal
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&animal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find animal")
	}

	return &animal, nil
}

// Update updates an existing animal
func (r *animalRepository) Update(ctx context.Context, animal *entities.Animal) error {
	animal.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Animals)
	filter := bson.M{"_id": animal.ID}

	result, err := collection.ReplaceOne(ctx, filter, animal)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update animal")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes an animal by ID
func (r *animalRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Animals)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete animal")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of animals with pagination and filters
func (r *animalRepository) List(ctx context.Context, filter repositories.AnimalFilter) ([]*entities.Animal, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Animals)

	// Build filter query
	query := bson.M{}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.Species != "" {
		query["species"] = filter.Species
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Sex != "" {
		query["sex"] = filter.Sex
	}

	if filter.Size != "" {
		query["size"] = filter.Size
	}

	if filter.AvailableOnly {
		query["status"] = entities.AnimalStatusAvailable
	}

	if filter.GoodWithKids != nil {
		query["behavior.good_with_kids"] = *filter.GoodWithKids
	}

	if filter.GoodWithDogs != nil {
		query["behavior.good_with_dogs"] = *filter.GoodWithDogs
	}

	if filter.GoodWithCats != nil {
		query["behavior.good_with_cats"] = *filter.GoodWithCats
	}

	if filter.AssignedCaretaker != nil {
		query["shelter.assigned_caretaker"] = *filter.AssignedCaretaker
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name.en": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"name.pl": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description.en": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description.pl": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"breed": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Age filters (calculate from date_of_birth)
	if filter.MinAge != nil || filter.MaxAge != nil {
		ageQuery := bson.M{}

		if filter.MinAge != nil {
			// Animals older than minAge should have date_of_birth less than now - minAge
			maxDate := time.Now().AddDate(-int(*filter.MinAge), 0, 0)
			ageQuery["$lte"] = maxDate
		}

		if filter.MaxAge != nil {
			// Animals younger than maxAge should have date_of_birth greater than now - maxAge
			minDate := time.Now().AddDate(-int(*filter.MaxAge), 0, 0)
			ageQuery["$gte"] = minDate
		}

		query["date_of_birth"] = ageQuery
	}

	// Count total documents matching filter
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count animals")
	}

	// Build find options
	findOptions := options.Find()

	// Set pagination
	if filter.Limit > 0 {
		findOptions.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(filter.Offset)
	}

	// Set sorting
	sortField := "created_at"
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
		return nil, 0, errors.Wrap(err, 500, "failed to query animals")
	}
	defer cursor.Close(ctx)

	var animals []*entities.Animal
	if err := cursor.All(ctx, &animals); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode animals")
	}

	return animals, total, nil
}

// AddDailyNote adds a daily note to an animal
func (r *animalRepository) AddDailyNote(ctx context.Context, animalID primitive.ObjectID, note entities.DailyNote) error {
	collection := r.db.Collection(mongodb.Collections.Animals)

	filter := bson.M{"_id": animalID}
	update := bson.M{
		"$push": bson.M{"shelter.daily_notes": note},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "failed to add daily note")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// UpdateImages updates the images for an animal
func (r *animalRepository) UpdateImages(ctx context.Context, animalID primitive.ObjectID, images entities.AnimalImages) error {
	collection := r.db.Collection(mongodb.Collections.Animals)

	filter := bson.M{"_id": animalID}
	update := bson.M{
		"$set": bson.M{
			"images":     images,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update images")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// UpdateStatus updates the status of an animal
func (r *animalRepository) UpdateStatus(ctx context.Context, animalID primitive.ObjectID, status entities.AnimalStatus) error {
	collection := r.db.Collection(mongodb.Collections.Animals)

	filter := bson.M{"_id": animalID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update status")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// GetStatistics returns statistics about animals
func (r *animalRepository) GetStatistics(ctx context.Context) (*repositories.AnimalStatistics, error) {
	collection := r.db.Collection(mongodb.Collections.Animals)

	stats := &repositories.AnimalStatistics{
		ByStatus:   make(map[string]int64),
		ByCategory: make(map[string]int64),
		BySpecies:  make(map[string]int64),
	}

	// Total animals
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count total animals")
	}
	stats.TotalAnimals = total

	// Count by status
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by status")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats.ByStatus[result.ID] = result.Count
	}

	// Available for adoption
	stats.AvailableForAdoption = stats.ByStatus[string(entities.AnimalStatusAvailable)]

	// Count by category
	pipeline = mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by category")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats.ByCategory[result.ID] = result.Count
	}

	// Count by species
	pipeline = mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$species"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by species")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats.BySpecies[result.ID] = result.Count
	}

	// Adopted this month
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	adoptedThisMonth, err := collection.CountDocuments(ctx, bson.M{
		"status": entities.AnimalStatusAdopted,
		"adoption.adoption_date": bson.M{
			"$gte": startOfMonth,
		},
	})
	if err == nil {
		stats.AdoptedThisMonth = adoptedThisMonth
	}

	// Adopted this year
	startOfYear := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	adoptedThisYear, err := collection.CountDocuments(ctx, bson.M{
		"status": entities.AnimalStatusAdopted,
		"adoption.adoption_date": bson.M{
			"$gte": startOfYear,
		},
	})
	if err == nil {
		stats.AdoptedThisYear = adoptedThisYear
	}

	return stats, nil
}

// EnsureIndexes creates necessary indexes for the animals collection
func (r *animalRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Animals)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "species", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "shelter.assigned_caretaker", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "name.en", Value: "text"},
				{Key: "name.pl", Value: "text"},
				{Key: "description.en", Value: "text"},
				{Key: "description.pl", Value: "text"},
				{Key: "breed", Value: "text"},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
