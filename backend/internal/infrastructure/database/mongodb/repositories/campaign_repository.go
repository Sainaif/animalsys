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

type campaignRepository struct {
	db *mongodb.Database
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(db *mongodb.Database) repositories.CampaignRepository {
	return &campaignRepository{
		db: db,
	}
}

// Create creates a new campaign
func (r *campaignRepository) Create(ctx context.Context, campaign *entities.Campaign) error {
	campaign.ID = primitive.NewObjectID()
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Campaigns)
	_, err := collection.InsertOne(ctx, campaign)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create campaign")
	}

	return nil
}

// FindByID finds a campaign by ID
func (r *campaignRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Campaign, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	var campaign entities.Campaign
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&campaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find campaign")
	}

	return &campaign, nil
}

// Update updates a campaign
func (r *campaignRepository) Update(ctx context.Context, campaign *entities.Campaign) error {
	campaign.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Campaigns)
	result, err := collection.ReplaceOne(
		ctx,
		bson.M{"_id": campaign.ID},
		campaign,
	)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update campaign")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a campaign
func (r *campaignRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete campaign")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List lists campaigns with filters
func (r *campaignRepository) List(ctx context.Context, filter *repositories.CampaignFilter) ([]*entities.Campaign, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	query := r.buildQuery(filter)

	// Get total count
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count campaigns")
	}

	// Build find options
	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}

	// Sorting
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
		return nil, 0, errors.Wrap(err, 500, "failed to list campaigns")
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode campaigns")
	}

	return campaigns, total, nil
}

// buildQuery builds MongoDB query from filter
func (r *campaignRepository) buildQuery(filter *repositories.CampaignFilter) bson.M {
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Public != nil {
		query["public"] = *filter.Public
	}

	if filter.Featured != nil {
		query["featured"] = *filter.Featured
	}

	if filter.ManagerID != nil && !filter.ManagerID.IsZero() {
		query["manager"] = *filter.ManagerID
	}

	if filter.StartDateFrom != nil {
		query["start_date"] = bson.M{"$gte": *filter.StartDateFrom}
	}

	if filter.StartDateTo != nil {
		if existing, ok := query["start_date"].(bson.M); ok {
			existing["$lte"] = *filter.StartDateTo
		} else {
			query["start_date"] = bson.M{"$lte": *filter.StartDateTo}
		}
	}

	if filter.EndDateFrom != nil {
		query["end_date"] = bson.M{"$gte": *filter.EndDateFrom}
	}

	if filter.EndDateTo != nil {
		if existing, ok := query["end_date"].(bson.M); ok {
			existing["$lte"] = *filter.EndDateTo
		} else {
			query["end_date"] = bson.M{"$lte": *filter.EndDateTo}
		}
	}

	if filter.GoalAmountMin > 0 {
		query["goal_amount"] = bson.M{"$gte": filter.GoalAmountMin}
	}

	if filter.GoalAmountMax > 0 {
		if existing, ok := query["goal_amount"].(bson.M); ok {
			existing["$lte"] = filter.GoalAmountMax
		} else {
			query["goal_amount"] = bson.M{"$lte": filter.GoalAmountMax}
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name.en": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"name.pl": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description.en": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description.pl": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	return query
}

// GetActiveCampaigns returns all active campaigns within their date range
func (r *campaignRepository) GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	now := time.Now()

	query := bson.M{
		"status":     entities.CampaignStatusActive,
		"start_date": bson.M{"$lte": now},
		"$or": []bson.M{
			{"end_date": nil},
			{"end_date": bson.M{"$gte": now}},
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to get active campaigns")
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode active campaigns")
	}

	return campaigns, nil
}

// GetFeaturedCampaigns returns all featured campaigns
func (r *campaignRepository) GetFeaturedCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	now := time.Now()

	query := bson.M{
		"featured":   true,
		"status":     entities.CampaignStatusActive,
		"start_date": bson.M{"$lte": now},
		"$or": []bson.M{
			{"end_date": nil},
			{"end_date": bson.M{"$gte": now}},
		},
	}

	findOptions := options.Find().SetSort(bson.D{
		{Key: "current_amount", Value: -1}, // Show campaigns with most funding first
		{Key: "start_date", Value: -1},
	})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to get featured campaigns")
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode featured campaigns")
	}

	return campaigns, nil
}

// GetPublicCampaigns returns all public campaigns
func (r *campaignRepository) GetPublicCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	now := time.Now()

	query := bson.M{
		"public":     true,
		"status":     entities.CampaignStatusActive,
		"start_date": bson.M{"$lte": now},
		"$or": []bson.M{
			{"end_date": nil},
			{"end_date": bson.M{"$gte": now}},
		},
	}

	findOptions := options.Find().SetSort(bson.D{
		{Key: "featured", Value: -1}, // Featured first
		{Key: "start_date", Value: -1},
	})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to get public campaigns")
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode public campaigns")
	}

	return campaigns, nil
}

// GetCampaignsByManager returns all campaigns managed by a specific user
func (r *campaignRepository) GetCampaignsByManager(ctx context.Context, managerID primitive.ObjectID) ([]*entities.Campaign, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	query := bson.M{"manager": managerID}

	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to get campaigns by manager")
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode campaigns")
	}

	return campaigns, nil
}

// UpdateCampaignStats updates campaign statistics after a donation
func (r *campaignRepository) UpdateCampaignStats(ctx context.Context, campaignID primitive.ObjectID, donationAmount float64, isNewDonor bool) error {
	// Get the campaign
	campaign, err := r.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	// Update stats using entity method
	campaign.UpdateStats(donationAmount, isNewDonor)

	// Save the updated campaign
	return r.Update(ctx, campaign)
}

// EnsureIndexes creates necessary indexes for the campaigns collection
func (r *campaignRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "public", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "featured", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "start_date", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "end_date", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "manager", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "start_date", Value: 1},
				{Key: "end_date", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "public", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "featured", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "tags", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create campaign indexes")
	}

	return nil
}

// GetCampaignStatistics returns statistics for campaigns
func (r *campaignRepository) GetCampaignStatistics(ctx context.Context) (*repositories.CampaignStatistics, error) {
	collection := r.db.Collection(mongodb.Collections.Campaigns)
	stats := &repositories.CampaignStatistics{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
	}

	// Total campaigns
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count campaigns")
	}
	stats.TotalCampaigns = total

	// Active campaigns
	activeCount, err := collection.CountDocuments(ctx, bson.M{"status": entities.CampaignStatusActive})
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count active campaigns")
	}
	stats.ActiveCampaigns = activeCount

	// Aggregate by type
	typePipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$type",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	typeCursor, err := collection.Aggregate(ctx, typePipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate campaigns by type")
	}
	defer typeCursor.Close(ctx)

	for typeCursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := typeCursor.Decode(&result); err != nil {
			continue
		}
		stats.ByType[result.ID] = result.Count
	}

	// Aggregate by status
	statusPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	statusCursor, err := collection.Aggregate(ctx, statusPipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate campaigns by status")
	}
	defer statusCursor.Close(ctx)

	for statusCursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := statusCursor.Decode(&result); err != nil {
			continue
		}
		stats.ByStatus[result.ID] = result.Count
	}

	// Aggregate amounts
	amountPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":          nil,
				"totalGoal":    bson.M{"$sum": "$goal_amount"},
				"totalRaised":  bson.M{"$sum": "$current_amount"},
				"totalDonors":  bson.M{"$sum": "$donor_count"},
			},
		},
	}

	var amountResult struct {
		TotalGoal    float64 `bson:"totalGoal"`
		TotalRaised  float64 `bson:"totalRaised"`
		TotalDonors  int64   `bson:"totalDonors"`
	}

	amountCursor, err := collection.Aggregate(ctx, amountPipeline)
	if err == nil {
		defer amountCursor.Close(ctx)
		if amountCursor.Next(ctx) {
			if err := amountCursor.Decode(&amountResult); err == nil {
				stats.TotalGoalAmount = amountResult.TotalGoal
				stats.TotalRaisedAmount = amountResult.TotalRaised
				stats.TotalDonors = amountResult.TotalDonors
			}
		}
	}

	// Calculate average progress
	if stats.TotalGoalAmount > 0 {
		stats.AverageProgress = (stats.TotalRaisedAmount / stats.TotalGoalAmount) * 100
	}

	return stats, nil
}
