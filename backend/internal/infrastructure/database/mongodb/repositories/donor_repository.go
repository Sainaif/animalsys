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

// donorRepository implements the DonorRepository interface
type donorRepository struct {
	db *mongodb.Database
}

// NewDonorRepository creates a new donor repository
func NewDonorRepository(db *mongodb.Database) repositories.DonorRepository {
	return &donorRepository{db: db}
}

// Create creates a new donor
func (r *donorRepository) Create(ctx context.Context, donor *entities.Donor) error {
	donor.CreatedAt = time.Now()
	donor.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Donors)
	result, err := collection.InsertOne(ctx, donor)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create donor")
	}

	donor.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a donor by ID
func (r *donorRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	var donor entities.Donor
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&donor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find donor")
	}

	return &donor, nil
}

// FindManyByIDs finds multiple donors by their IDs
func (r *donorRepository) FindManyByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*entities.Donor, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	filter := bson.M{"_id": bson.M{"$in": ids}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query donors")
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err := cursor.All(ctx, &donors); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode donors")
	}

	return donors, nil
}

// Update updates an existing donor
func (r *donorRepository) Update(ctx context.Context, donor *entities.Donor) error {
	donor.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Donors)
	filter := bson.M{"_id": donor.ID}

	result, err := collection.ReplaceOne(ctx, filter, donor)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update donor")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a donor by ID
func (r *donorRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Donors)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete donor")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of donors with pagination and filters
func (r *donorRepository) List(ctx context.Context, filter *repositories.DonorFilter) ([]*entities.Donor, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	// Build filter query
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"first_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"last_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"organization_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"contact.email": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if filter.MinTotalDonated != nil {
		query["total_donated"] = bson.M{"$gte": *filter.MinTotalDonated}
	}

	if filter.MaxTotalDonated != nil {
		if existing, ok := query["total_donated"].(bson.M); ok {
			existing["$lte"] = *filter.MaxTotalDonated
		} else {
			query["total_donated"] = bson.M{"$lte": *filter.MaxTotalDonated}
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
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
		query["created_at"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count donors")
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
		return nil, 0, errors.Wrap(err, 500, "failed to query donors")
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err := cursor.All(ctx, &donors); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode donors")
	}

	return donors, total, nil
}

// FindByEmail finds a donor by email
func (r *donorRepository) FindByEmail(ctx context.Context, email string) (*entities.Donor, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	var donor entities.Donor
	err := collection.FindOne(ctx, bson.M{"contact.email": email}).Decode(&donor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil, not error if not found
		}
		return nil, errors.Wrap(err, 500, "failed to find donor by email")
	}

	return &donor, nil
}

// GetMajorDonors returns all major donors
func (r *donorRepository) GetMajorDonors(ctx context.Context) ([]*entities.Donor, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	query := bson.M{
		"$or": []bson.M{
			{"status": entities.DonorStatusMajor},
			{"total_donated": bson.M{"$gte": 10000}},
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "total_donated", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query major donors")
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err := cursor.All(ctx, &donors); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode major donors")
	}

	return donors, nil
}

// GetLapsedDonors returns donors who haven't donated in N days
func (r *donorRepository) GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	cutoffDate := time.Now().AddDate(0, 0, -days)

	query := bson.M{
		"last_donation_date": bson.M{
			"$lt": cutoffDate,
			"$ne": nil,
		},
		"status": bson.M{"$ne": entities.DonorStatusInactive},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "last_donation_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query lapsed donors")
	}
	defer cursor.Close(ctx)

	var donors []*entities.Donor
	if err := cursor.All(ctx, &donors); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode lapsed donors")
	}

	return donors, nil
}

// GetDonorStatistics returns donor statistics
func (r *donorRepository) GetDonorStatistics(ctx context.Context) (*repositories.DonorStatistics, error) {
	collection := r.db.Collection(mongodb.Collections.Donors)

	stats := &repositories.DonorStatistics{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
	}

	// Total donors
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count total donors")
	}
	stats.TotalDonors = total

	// By type
	typePipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}},
	}
	typeCursor, err := collection.Aggregate(ctx, typePipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by type")
	}
	defer typeCursor.Close(ctx)

	var typeResults []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err := typeCursor.All(ctx, &typeResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode type results")
	}

	for _, result := range typeResults {
		stats.ByType[result.ID] = result.Count
		if result.ID == string(entities.DonorTypeIndividual) {
			stats.IndividualDonors = result.Count
		} else if result.ID == string(entities.DonorTypeOrganization) || result.ID == string(entities.DonorTypeCorporate) {
			stats.OrganizationDonors += result.Count
		}
	}

	// By status
	statusPipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}},
	}
	statusCursor, err := collection.Aggregate(ctx, statusPipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by status")
	}
	defer statusCursor.Close(ctx)

	var statusResults []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err := statusCursor.All(ctx, &statusResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode status results")
	}

	for _, result := range statusResults {
		stats.ByStatus[result.ID] = result.Count
		switch result.ID {
		case string(entities.DonorStatusActive):
			stats.ActiveDonors = result.Count
		case string(entities.DonorStatusInactive):
			stats.InactiveDonors = result.Count
		case string(entities.DonorStatusMajor):
			stats.MajorDonors = result.Count
		case string(entities.DonorStatusLapsed):
			stats.LapsedDonors = result.Count
		}
	}

	// Lifetime value statistics
	valuePipeline := []bson.M{
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$total_donated"},
			"avg":   bson.M{"$avg": "$total_donated"},
		}},
	}
	valueCursor, err := collection.Aggregate(ctx, valuePipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate lifetime value")
	}
	defer valueCursor.Close(ctx)

	var valueResults []struct {
		Total float64 `bson:"total"`
		Avg   float64 `bson:"avg"`
	}
	if err := valueCursor.All(ctx, &valueResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode value results")
	}

	if len(valueResults) > 0 {
		stats.TotalLifetimeValue = valueResults[0].Total
		stats.AverageLifetimeValue = valueResults[0].Avg
	}

	return stats, nil
}

// EnsureIndexes creates necessary indexes for the donors collection
func (r *donorRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Donors)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "contact.email", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "total_donated", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "last_donation_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "organization_name", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
