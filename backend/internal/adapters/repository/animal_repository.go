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

type animalRepository struct {
	collection *mongo.Collection
}

// NewAnimalRepository creates a new animal repository
func NewAnimalRepository(db *mongo.Database) interfaces.AnimalRepository {
	return &animalRepository{
		collection: db.Collection("animals"),
	}
}

func (r *animalRepository) Create(ctx context.Context, animal *entities.Animal) error {
	animal.CreatedAt = time.Now()
	animal.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, animal)
	return err
}

func (r *animalRepository) GetByID(ctx context.Context, id string) (*entities.Animal, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid animal ID")
	}

	var animal entities.Animal
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&animal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("animal not found")
		}
		return nil, err
	}

	return &animal, nil
}

func (r *animalRepository) Update(ctx context.Context, id string, animal *entities.Animal) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid animal ID")
	}

	animal.UpdatedAt = time.Now()

	update := bson.M{
		"$set": animal,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("animal not found")
	}

	return nil
}

func (r *animalRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid animal ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("animal not found")
	}

	return nil
}

func (r *animalRepository) List(ctx context.Context, filter *entities.AnimalFilter) ([]*entities.Animal, int64, error) {
	// Build filter
	mongoFilter := bson.M{}

	if filter.Species != "" {
		mongoFilter["species"] = filter.Species
	}

	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}

	if filter.Gender != "" {
		mongoFilter["gender"] = filter.Gender
	}

	if filter.Size != "" {
		mongoFilter["size"] = filter.Size
	}

	if filter.Neutered != nil {
		mongoFilter["neutered"] = *filter.Neutered
	}

	if filter.Vaccinated != nil {
		mongoFilter["vaccinated"] = *filter.Vaccinated
	}

	// Age filter
	if filter.MinAge > 0 || filter.MaxAge > 0 {
		ageFilter := bson.M{}
		if filter.MinAge > 0 {
			ageFilter["$gte"] = filter.MinAge
		}
		if filter.MaxAge > 0 {
			ageFilter["$lte"] = filter.MaxAge
		}
		mongoFilter["age_years"] = ageFilter
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"breed": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Build options
	findOptions := options.Find()

	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}

	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}

	// Sort
	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}

	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}

	findOptions.SetSort(bson.M{sortField: sortOrder})

	// Find
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var animals []*entities.Animal
	if err = cursor.All(ctx, &animals); err != nil {
		return nil, 0, err
	}

	return animals, total, nil
}

func (r *animalRepository) AddMedicalRecord(ctx context.Context, id string, record entities.MedicalRecord) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid animal ID")
	}

	update := bson.M{
		"$push": bson.M{"medical_history": record},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("animal not found")
	}

	return nil
}

func (r *animalRepository) AddPhoto(ctx context.Context, id string, photoURL string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid animal ID")
	}

	update := bson.M{
		"$push": bson.M{"photos": photoURL},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("animal not found")
	}

	return nil
}

func (r *animalRepository) UpdateStatus(ctx context.Context, id string, status entities.AnimalStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid animal ID")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("animal not found")
	}

	return nil
}

func (r *animalRepository) GetAvailableForAdoption(ctx context.Context, limit, offset int) ([]*entities.Animal, int64, error) {
	filter := bson.M{"status": entities.AnimalStatusAvailable}

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

	var animals []*entities.Animal
	if err = cursor.All(ctx, &animals); err != nil {
		return nil, 0, err
	}

	return animals, total, nil
}

func (r *animalRepository) GetByStatus(ctx context.Context, status entities.AnimalStatus, limit, offset int) ([]*entities.Animal, int64, error) {
	filter := bson.M{"status": status}

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

	var animals []*entities.Animal
	if err = cursor.All(ctx, &animals); err != nil {
		return nil, 0, err
	}

	return animals, total, nil
}
