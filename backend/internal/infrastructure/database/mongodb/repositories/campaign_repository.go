package repositories

import (
	"context"
	"strconv"
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
	var raw bson.M
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&raw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find campaign")
	}

	return mapCampaignDocument(raw), nil
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

	campaigns, err := r.decodeCampaignCursor(ctx, cursor, "failed to decode campaigns")
	if err != nil {
		return nil, 0, err
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

	campaigns, err := r.decodeCampaignCursor(ctx, cursor, "failed to decode active campaigns")
	if err != nil {
		return nil, err
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

	campaigns, err := r.decodeCampaignCursor(ctx, cursor, "failed to decode featured campaigns")
	if err != nil {
		return nil, err
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

	campaigns, err := r.decodeCampaignCursor(ctx, cursor, "failed to decode public campaigns")
	if err != nil {
		return nil, err
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

	campaigns, err := r.decodeCampaignCursor(ctx, cursor, "failed to decode campaigns")
	if err != nil {
		return nil, err
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
				"_id":         nil,
				"totalGoal":   bson.M{"$sum": "$goal_amount"},
				"totalRaised": bson.M{"$sum": "$current_amount"},
				"totalDonors": bson.M{"$sum": "$donor_count"},
			},
		},
	}

	var amountResult struct {
		TotalGoal   float64 `bson:"totalGoal"`
		TotalRaised float64 `bson:"totalRaised"`
		TotalDonors int64   `bson:"totalDonors"`
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

func (r *campaignRepository) decodeCampaignCursor(ctx context.Context, cursor *mongo.Cursor, message string) ([]*entities.Campaign, error) {
	var rawDocs []bson.M
	if err := cursor.All(ctx, &rawDocs); err != nil {
		return nil, errors.Wrap(err, 500, message)
	}

	campaigns := make([]*entities.Campaign, 0, len(rawDocs))
	for _, raw := range rawDocs {
		campaigns = append(campaigns, mapCampaignDocument(raw))
	}

	return campaigns, nil
}

func mapCampaignDocument(raw bson.M) *entities.Campaign {
	if raw == nil {
		return &entities.Campaign{}
	}

	defaultName := stringFromRaw(raw["name"])
	if defaultName == "" {
		defaultName = "Campaign"
	}

	campaign := &entities.Campaign{
		ID:              objectIDFromRaw(raw["_id"]),
		Name:            decodeMultilingualName(raw["name"], defaultName),
		Description:     entities.MultilingualName{},
		Type:            entities.CampaignType(stringFromRaw(raw["type"])),
		Status:          entities.CampaignStatus(stringFromRaw(raw["status"])),
		GoalAmount:      float64FromRaw(raw["goal_amount"]),
		CurrentAmount:   float64FromRaw(raw["current_amount"]),
		DonorCount:      intFromRaw(raw["donor_count"]),
		DonationCount:   intFromRaw(raw["donation_count"]),
		AverageDonation: float64FromRaw(raw["average_donation"]),
		StartDate:       timeFromRaw(raw["start_date"]),
		EndDate:         timePtrFromRaw(raw["end_date"]),
		ImageURL:        stringFromRaw(raw["image_url"]),
		VideoURL:        stringFromRaw(raw["video_url"]),
		Tags:            stringSliceFromRaw(raw["tags"]),
		Public:          boolFromRaw(raw["public"]),
		Featured:        boolFromRaw(raw["featured"]),
		Manager:         objectIDFromRaw(raw["manager"]),
		ContactEmail:    stringFromRaw(raw["contact_email"]),
		ContactPhone:    stringFromRaw(raw["contact_phone"]),
		Notes:           stringFromRaw(raw["notes"]),
		CreatedBy:       objectIDFromRaw(raw["created_by"]),
		UpdatedBy:       objectIDFromRaw(raw["updated_by"]),
		CreatedAt:       timeFromRaw(raw["created_at"]),
		UpdatedAt:       timeFromRaw(raw["updated_at"]),
	}

	campaign.Description = decodeMultilingualName(raw["description"], campaign.Name.English)

	if campaign.Type == "" {
		campaign.Type = entities.CampaignTypeGeneral
	}
	if campaign.Status == "" {
		campaign.Status = entities.CampaignStatusDraft
	}

	if campaign.GoalAmount == 0 {
		campaign.GoalAmount = float64FromRaw(raw["goal"])
	}
	if campaign.CurrentAmount == 0 {
		campaign.CurrentAmount = float64FromRaw(raw["raised"])
	}

	if campaign.Description.English == "" && campaign.Description.Polish == "" {
		campaign.Description = campaign.Name
	}

	if campaign.CreatedAt.IsZero() {
		campaign.CreatedAt = time.Now()
	}
	if campaign.UpdatedAt.IsZero() {
		campaign.UpdatedAt = campaign.CreatedAt
	}
	if campaign.StartDate.IsZero() {
		campaign.StartDate = campaign.CreatedAt
	}

	return campaign
}

func decodeMultilingualName(value interface{}, fallback string) entities.MultilingualName {
	result := entities.MultilingualName{}

	switch v := value.(type) {
	case nil:
	case string:
		result.English = v
		result.Polish = v
	case bson.M:
		result.English = stringFromRaw(v["en"])
		result.Polish = stringFromRaw(v["pl"])
		result.Latin = stringFromRaw(v["latin"])
	case map[string]interface{}:
		result.English = stringFromRaw(v["en"])
		result.Polish = stringFromRaw(v["pl"])
		result.Latin = stringFromRaw(v["latin"])
	case primitive.D:
		return decodeMultilingualName(v.Map(), fallback)
	default:
		str := stringFromRaw(value)
		if str != "" {
			result.English = str
			result.Polish = str
		}
	}

	if result.English == "" {
		result.English = fallback
	}
	if result.Polish == "" {
		result.Polish = fallback
	}
	if result.Latin == "" {
		result.Latin = fallback
	}

	return result
}

func stringFromRaw(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case primitive.ObjectID:
		return v.Hex()
	default:
		return ""
	}
}

func stringSliceFromRaw(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str := stringFromRaw(item); str != "" {
				result = append(result, str)
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return nil
	}
}

func float64FromRaw(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case primitive.Decimal128:
		if parsed, err := strconv.ParseFloat(v.String(), 64); err == nil {
			return parsed
		}
	}
	return 0
}

func intFromRaw(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func boolFromRaw(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return parsed
		}
	case int:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0
	default:
		return false
	}
	return false
}

func timeFromRaw(value interface{}) time.Time {
	switch v := value.(type) {
	case time.Time:
		return v
	case primitive.DateTime:
		return v.Time()
	case string:
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			return parsed
		}
	default:
		return time.Time{}
	}
	return time.Time{}
}

func timePtrFromRaw(value interface{}) *time.Time {
	t := timeFromRaw(value)
	if t.IsZero() {
		return nil
	}
	return &t
}

func objectIDFromRaw(value interface{}) primitive.ObjectID {
	switch v := value.(type) {
	case primitive.ObjectID:
		return v
	case string:
		if id, err := primitive.ObjectIDFromHex(v); err == nil {
			return id
		}
	}
	return primitive.NilObjectID
}
