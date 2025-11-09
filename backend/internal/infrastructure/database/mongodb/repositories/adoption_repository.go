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

// adoptionRepository implements the AdoptionRepository interface
type adoptionRepository struct {
	db *mongodb.Database
}

// NewAdoptionRepository creates a new adoption repository
func NewAdoptionRepository(db *mongodb.Database) repositories.AdoptionRepository {
	return &adoptionRepository{db: db}
}

// Create creates a new adoption record
func (r *adoptionRepository) Create(ctx context.Context, adoption *entities.Adoption) error {
	adoption.CreatedAt = time.Now()
	adoption.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Adoptions)
	result, err := collection.InsertOne(ctx, adoption)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create adoption")
	}

	adoption.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an adoption by ID
func (r *adoptionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Adoption, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	var adoption entities.Adoption
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&adoption)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find adoption")
	}

	return &adoption, nil
}

// Update updates an existing adoption
func (r *adoptionRepository) Update(ctx context.Context, adoption *entities.Adoption) error {
	adoption.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Adoptions)
	filter := bson.M{"_id": adoption.ID}

	result, err := collection.ReplaceOne(ctx, filter, adoption)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update adoption")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes an adoption by ID
func (r *adoptionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete adoption")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of adoptions with pagination and filters
func (r *adoptionRepository) List(ctx context.Context, filter repositories.AdoptionFilter) ([]*entities.Adoption, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	// Build filter query
	query := bson.M{}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
	}

	if filter.AdopterID != nil {
		query["adopter_id"] = *filter.AdopterID
	}

	if filter.ApplicationID != nil {
		query["application_id"] = *filter.ApplicationID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.PaymentStatus != "" {
		query["payment_status"] = filter.PaymentStatus
	}

	if filter.TrialPeriod != nil {
		query["trial_period"] = *filter.TrialPeriod
	}

	if filter.ProcessedBy != nil {
		query["processed_by"] = *filter.ProcessedBy
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
		query["adoption_date"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count adoptions")
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
	sortField := "adoption_date"
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
		return nil, 0, errors.Wrap(err, 500, "failed to query adoptions")
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err := cursor.All(ctx, &adoptions); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode adoptions")
	}

	return adoptions, total, nil
}

// GetByAnimalID returns adoption for a specific animal
func (r *adoptionRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) (*entities.Adoption, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	var adoption entities.Adoption
	err := collection.FindOne(
		ctx,
		bson.M{"animal_id": animalID},
		options.FindOne().SetSort(bson.D{{Key: "adoption_date", Value: -1}}),
	).Decode(&adoption)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no adoption found
		}
		return nil, errors.Wrap(err, 500, "failed to query adoption")
	}

	return &adoption, nil
}

// GetByAdopterID returns all adoptions for a specific adopter
func (r *adoptionRepository) GetByAdopterID(ctx context.Context, adopterID primitive.ObjectID) ([]*entities.Adoption, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	findOptions := options.Find().SetSort(bson.D{{Key: "adoption_date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"adopter_id": adopterID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query adoptions")
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err := cursor.All(ctx, &adoptions); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode adoptions")
	}

	return adoptions, nil
}

// GetByApplicationID returns adoption by application ID
func (r *adoptionRepository) GetByApplicationID(ctx context.Context, applicationID primitive.ObjectID) (*entities.Adoption, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	var adoption entities.Adoption
	err := collection.FindOne(ctx, bson.M{"application_id": applicationID}).Decode(&adoption)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no adoption found
		}
		return nil, errors.Wrap(err, 500, "failed to find adoption")
	}

	return &adoption, nil
}

// GetPendingFollowUps returns adoptions with pending follow-ups
func (r *adoptionRepository) GetPendingFollowUps(ctx context.Context, days int) ([]*entities.Adoption, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	query := bson.M{
		"next_follow_up_date": bson.M{
			"$lte": futureDate,
			"$gte": now,
		},
		"status": bson.M{
			"$in": []entities.AdoptionStatus{
				entities.AdoptionStatusPending,
				entities.AdoptionStatusCompleted,
			},
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "next_follow_up_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query pending follow-ups")
	}
	defer cursor.Close(ctx)

	var adoptions []*entities.Adoption
	if err := cursor.All(ctx, &adoptions); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode pending follow-ups")
	}

	return adoptions, nil
}

// GetAdoptionStatistics returns adoption statistics
func (r *adoptionRepository) GetAdoptionStatistics(ctx context.Context) (*repositories.AdoptionStatistics, error) {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	stats := &repositories.AdoptionStatistics{
		ByStatus:        make(map[string]int64),
		ByPaymentStatus: make(map[string]int64),
	}

	// Total adoptions
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count total adoptions")
	}
	stats.TotalAdoptions = total

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
		case string(entities.AdoptionStatusCompleted):
			stats.CompletedAdoptions = result.Count
		case string(entities.AdoptionStatusPending):
			stats.PendingAdoptions = result.Count
		case string(entities.AdoptionStatusReturned):
			stats.ReturnedAnimals = result.Count
		}
	}

	// By payment status
	paymentPipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$payment_status",
			"count": bson.M{"$sum": 1},
		}},
	}
	paymentCursor, err := collection.Aggregate(ctx, paymentPipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate by payment status")
	}
	defer paymentCursor.Close(ctx)

	var paymentResults []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err := paymentCursor.All(ctx, &paymentResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode payment results")
	}

	for _, result := range paymentResults {
		stats.ByPaymentStatus[result.ID] = result.Count
	}

	// Financial stats
	feePipeline := []bson.M{
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$adoption_fee"},
			"avg":   bson.M{"$avg": "$adoption_fee"},
		}},
	}
	feeCursor, err := collection.Aggregate(ctx, feePipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate fees")
	}
	defer feeCursor.Close(ctx)

	var feeResults []struct {
		Total float64 `bson:"total"`
		Avg   float64 `bson:"avg"`
	}
	if err := feeCursor.All(ctx, &feeResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode fee results")
	}

	if len(feeResults) > 0 {
		stats.TotalAdoptionFees = feeResults[0].Total
		stats.AverageAdoptionFee = feeResults[0].Avg
	}

	// Adoptions this month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthCount, err := collection.CountDocuments(ctx, bson.M{
		"adoption_date": bson.M{"$gte": startOfMonth},
	})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count month adoptions")
	}
	stats.AdoptionsThisMonth = monthCount

	// Adoptions this year
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	yearCount, err := collection.CountDocuments(ctx, bson.M{
		"adoption_date": bson.M{"$gte": startOfYear},
	})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count year adoptions")
	}
	stats.AdoptionsThisYear = yearCount

	// Pending follow-ups
	followUpCount, err := collection.CountDocuments(ctx, bson.M{
		"next_follow_up_date": bson.M{
			"$ne":  nil,
			"$lte": now,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count pending follow-ups")
	}
	stats.PendingFollowUps = followUpCount

	// Return rate
	if stats.TotalAdoptions > 0 {
		stats.ReturnRate = (float64(stats.ReturnedAnimals) / float64(stats.TotalAdoptions)) * 100
	}

	return stats, nil
}

// EnsureIndexes creates necessary indexes for the adoptions collection
func (r *adoptionRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Adoptions)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "adopter_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "application_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "payment_status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "adoption_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "next_follow_up_date", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "animal_id", Value: 1},
				{Key: "adoption_date", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "adopter_id", Value: 1},
				{Key: "adoption_date", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "next_follow_up_date", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
