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

// adoptionApplicationRepository implements the AdoptionApplicationRepository interface
type adoptionApplicationRepository struct {
	db *mongodb.Database
}

// NewAdoptionApplicationRepository creates a new adoption application repository
func NewAdoptionApplicationRepository(db *mongodb.Database) repositories.AdoptionApplicationRepository {
	return &adoptionApplicationRepository{db: db}
}

// Create creates a new adoption application
func (r *adoptionApplicationRepository) Create(ctx context.Context, application *entities.AdoptionApplication) error {
	application.CreatedAt = time.Now()
	application.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)
	result, err := collection.InsertOne(ctx, application)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create adoption application")
	}

	application.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an adoption application by ID
func (r *adoptionApplicationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AdoptionApplication, error) {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	var application entities.AdoptionApplication
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&application)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find adoption application")
	}

	return &application, nil
}

// Update updates an existing adoption application
func (r *adoptionApplicationRepository) Update(ctx context.Context, application *entities.AdoptionApplication) error {
	application.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)
	filter := bson.M{"_id": application.ID}

	result, err := collection.ReplaceOne(ctx, filter, application)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update adoption application")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes an adoption application by ID
func (r *adoptionApplicationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete adoption application")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of adoption applications with pagination and filters
func (r *adoptionApplicationRepository) List(ctx context.Context, filter repositories.AdoptionApplicationFilter) ([]*entities.AdoptionApplication, int64, error) {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	// Build filter query
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.ApplicantEmail != "" {
		query["applicant.email"] = bson.M{"$regex": filter.ApplicantEmail, "$options": "i"}
	}

	if filter.ApplicantName != "" {
		query["$or"] = []bson.M{
			{"applicant.first_name": bson.M{"$regex": filter.ApplicantName, "$options": "i"}},
			{"applicant.last_name": bson.M{"$regex": filter.ApplicantName, "$options": "i"}},
		}
	}

	if filter.ReviewedBy != nil {
		query["reviewed_by"] = *filter.ReviewedBy
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
		query["application_date"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count adoption applications")
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
	sortField := "application_date"
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
		return nil, 0, errors.Wrap(err, 500, "failed to query adoption applications")
	}
	defer cursor.Close(ctx)

	var applications []*entities.AdoptionApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode adoption applications")
	}

	return applications, total, nil
}

// GetByAnimalID returns all applications for a specific animal
func (r *adoptionApplicationRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.AdoptionApplication, error) {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	findOptions := options.Find().SetSort(bson.D{{Key: "application_date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"animal_id": animalID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query adoption applications")
	}
	defer cursor.Close(ctx)

	var applications []*entities.AdoptionApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode adoption applications")
	}

	return applications, nil
}

// GetByApplicantEmail returns applications by applicant email
func (r *adoptionApplicationRepository) GetByApplicantEmail(ctx context.Context, email string) ([]*entities.AdoptionApplication, error) {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	query := bson.M{"applicant.email": email}
	findOptions := options.Find().SetSort(bson.D{{Key: "application_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query adoption applications")
	}
	defer cursor.Close(ctx)

	var applications []*entities.AdoptionApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode adoption applications")
	}

	return applications, nil
}

// GetPendingApplications returns all pending applications
func (r *adoptionApplicationRepository) GetPendingApplications(ctx context.Context) ([]*entities.AdoptionApplication, error) {
	return r.GetApplicationsByStatus(ctx, entities.ApplicationStatusPending)
}

// GetApplicationsByStatus returns applications by status
func (r *adoptionApplicationRepository) GetApplicationsByStatus(ctx context.Context, status entities.ApplicationStatus) ([]*entities.AdoptionApplication, error) {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	query := bson.M{"status": status}
	findOptions := options.Find().SetSort(bson.D{{Key: "application_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query adoption applications")
	}
	defer cursor.Close(ctx)

	var applications []*entities.AdoptionApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode adoption applications")
	}

	return applications, nil
}

// EnsureIndexes creates necessary indexes for the adoption_applications collection
func (r *adoptionApplicationRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.AdoptionApplications)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "applicant.email", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "application_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "reviewed_by", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "animal_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "applicant.email", Value: 1},
				{Key: "application_date", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
