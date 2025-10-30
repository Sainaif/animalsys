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

type adoptionRepository struct {
	collection *mongo.Collection
}

func NewAdoptionRepository(db *mongo.Database) interfaces.AdoptionRepository {
	return &adoptionRepository{
		collection: db.Collection("adoptions"),
	}
}

func (r *adoptionRepository) Create(ctx context.Context, adoption *entities.Adoption) error {
	adoption.CreatedAt = time.Now()
	adoption.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, adoption)
	return err
}

func (r *adoptionRepository) GetByID(ctx context.Context, id string) (*entities.Adoption, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid adoption ID")
	}

	var adoption entities.Adoption
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&adoption)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("adoption not found")
		}
		return nil, err
	}
	return &adoption, nil
}

func (r *adoptionRepository) Update(ctx context.Context, id string, adoption *entities.Adoption) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid adoption ID")
	}

	adoption.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": adoption})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("adoption not found")
	}
	return nil
}

func (r *adoptionRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid adoption ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("adoption not found")
	}
	return nil
}

func (r *adoptionRepository) List(ctx context.Context, filter *entities.AdoptionFilter) ([]*entities.Adoption, int64, error) {
	mongoFilter := bson.M{}

	if filter.AnimalID != "" {
		mongoFilter["animal_id"] = filter.AnimalID
	}
	if filter.ApplicantID != "" {
		mongoFilter["applicant_id"] = filter.ApplicantID
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		mongoFilter["application_date"] = dateFilter
	}
	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"applicant_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"applicant_email": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"animal_name": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}
	sortField := "application_date"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	findOptions.SetSort(bson.M{sortField: sortOrder})

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err = cursor.All(ctx, &adoptions); err != nil {
		return nil, 0, err
	}
	return adoptions, total, nil
}

func (r *adoptionRepository) GetByAnimalID(ctx context.Context, animalID string) ([]*entities.Adoption, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"animal_id": animalID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err = cursor.All(ctx, &adoptions); err != nil {
		return nil, err
	}
	return adoptions, nil
}

func (r *adoptionRepository) GetByApplicantID(ctx context.Context, applicantID string) ([]*entities.Adoption, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"applicant_id": applicantID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err = cursor.All(ctx, &adoptions); err != nil {
		return nil, err
	}
	return adoptions, nil
}

func (r *adoptionRepository) UpdateStatus(ctx context.Context, id string, status entities.AdoptionStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid adoption ID")
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
		return errors.New("adoption not found")
	}
	return nil
}

func (r *adoptionRepository) GetPendingApplications(ctx context.Context, limit, offset int) ([]*entities.Adoption, int64, error) {
	filter := bson.M{"status": bson.M{"$in": []entities.AdoptionStatus{
		entities.AdoptionStatusSubmitted,
		entities.AdoptionStatusUnderReview,
	}}}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"application_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err = cursor.All(ctx, &adoptions); err != nil {
		return nil, 0, err
	}
	return adoptions, total, nil
}
