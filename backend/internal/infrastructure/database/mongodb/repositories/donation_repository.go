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

// donationRepository implements the DonationRepository interface
type donationRepository struct {
	db *mongodb.Database
}

// NewDonationRepository creates a new donation repository
func NewDonationRepository(db *mongodb.Database) repositories.DonationRepository {
	return &donationRepository{db: db}
}

// Create creates a new donation
func (r *donationRepository) Create(ctx context.Context, donation *entities.Donation) error {
	donation.CreatedAt = time.Now()
	donation.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Donations)
	result, err := collection.InsertOne(ctx, donation)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create donation")
	}

	donation.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a donation by ID
func (r *donationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	var donation entities.Donation
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&donation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find donation")
	}

	return &donation, nil
}

// Update updates an existing donation
func (r *donationRepository) Update(ctx context.Context, donation *entities.Donation) error {
	donation.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Donations)
	filter := bson.M{"_id": donation.ID}

	result, err := collection.ReplaceOne(ctx, filter, donation)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update donation")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a donation by ID
func (r *donationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Donations)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete donation")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of donations with pagination and filters
func (r *donationRepository) List(ctx context.Context, filter *repositories.DonationFilter) ([]*entities.Donation, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	// Build filter query
	query := bson.M{}

	if filter.DonorID != nil {
		query["donor_id"] = *filter.DonorID
	}

	if filter.CampaignID != nil {
		query["campaign_id"] = *filter.CampaignID
	}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.MinAmount != nil {
		query["amount"] = bson.M{"$gte": *filter.MinAmount}
	}

	if filter.MaxAmount != nil {
		if existing, ok := query["amount"].(bson.M); ok {
			existing["$lte"] = *filter.MaxAmount
		} else {
			query["amount"] = bson.M{"$lte": *filter.MaxAmount}
		}
	}

	if filter.PaymentMethod != "" {
		query["payment.method"] = filter.PaymentMethod
	}

	if filter.Designation != "" {
		query["designation"] = filter.Designation
	}

	if filter.IsRecurring != nil {
		query["is_recurring"] = *filter.IsRecurring
	}

	if filter.Anonymous != nil {
		query["anonymous"] = *filter.Anonymous
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
		query["donation_date"] = dateFilter
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count donations")
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
	sortField := "donation_date"
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
		return nil, 0, errors.Wrap(err, 500, "failed to query donations")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode donations")
	}

	return donations, total, nil
}

// GetByDonorID returns all donations for a specific donor
func (r *donationRepository) GetByDonorID(ctx context.Context, donorID primitive.ObjectID) ([]*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	findOptions := options.Find().SetSort(bson.D{{Key: "donation_date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"donor_id": donorID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query donations")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode donations")
	}

	return donations, nil
}

// GetByCampaignID returns all donations for a specific campaign
func (r *donationRepository) GetByCampaignID(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	findOptions := options.Find().SetSort(bson.D{{Key: "donation_date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"campaign_id": campaignID}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query donations")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode donations")
	}

	return donations, nil
}

// GetRecurringDonations returns all active recurring donations
func (r *donationRepository) GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	query := bson.M{
		"is_recurring":              true,
		"recurring_info.active":     true,
		"status":                    entities.DonationStatusCompleted,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "recurring_info.next_billing_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query recurring donations")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode recurring donations")
	}

	return donations, nil
}

// GetPendingThankYous returns donations without thank you sent
func (r *donationRepository) GetPendingThankYous(ctx context.Context) ([]*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	query := bson.M{
		"thank_you_sent": false,
		"status":         entities.DonationStatusCompleted,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "donation_date", Value: 1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query pending thank yous")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode pending thank yous")
	}

	return donations, nil
}

// GetDonationsByDateRange returns donations within a date range
func (r *donationRepository) GetDonationsByDateRange(ctx context.Context, from, to time.Time) ([]*entities.Donation, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	query := bson.M{
		"donation_date": bson.M{
			"$gte": from,
			"$lte": to,
		},
		"status": entities.DonationStatusCompleted,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "donation_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to query donations by date range")
	}
	defer cursor.Close(ctx)

	var donations []*entities.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode donations")
	}

	return donations, nil
}

// GetDonationStatistics returns donation statistics
func (r *donationRepository) GetDonationStatistics(ctx context.Context) (*repositories.DonationStatistics, error) {
	collection := r.db.Collection(mongodb.Collections.Donations)

	stats := &repositories.DonationStatistics{
		ByType:          make(map[string]int64),
		ByStatus:        make(map[string]int64),
		ByPaymentMethod: make(map[string]int64),
	}

	// Total completed donations
	totalQuery := bson.M{"status": entities.DonationStatusCompleted}
	total, err := collection.CountDocuments(ctx, totalQuery)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to count total donations")
	}
	stats.TotalDonations = total

	// Aggregation pipeline for amount statistics
	amountPipeline := []bson.M{
		{"$match": bson.M{"status": entities.DonationStatusCompleted}},
		{"$group": bson.M{
			"_id":      nil,
			"total":    bson.M{"$sum": "$amount"},
			"avg":      bson.M{"$avg": "$amount"},
			"max":      bson.M{"$max": "$amount"},
			"min":      bson.M{"$min": "$amount"},
		}},
	}
	amountCursor, err := collection.Aggregate(ctx, amountPipeline)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to aggregate amounts")
	}
	defer amountCursor.Close(ctx)

	var amountResults []struct {
		Total float64 `bson:"total"`
		Avg   float64 `bson:"avg"`
		Max   float64 `bson:"max"`
		Min   float64 `bson:"min"`
	}
	if err := amountCursor.All(ctx, &amountResults); err != nil {
		return nil, errors.Wrap(err, 500, "failed to decode amount results")
	}

	if len(amountResults) > 0 {
		stats.TotalAmount = amountResults[0].Total
		stats.AverageDonation = amountResults[0].Avg
		stats.LargestDonation = amountResults[0].Max
		stats.SmallestDonation = amountResults[0].Min
	}

	// Time-based statistics
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// Today
	todayQuery := bson.M{
		"donation_date": bson.M{"$gte": startOfDay},
		"status":        entities.DonationStatusCompleted,
	}
	todayCount, _ := collection.CountDocuments(ctx, todayQuery)
	stats.DonationsToday = todayCount

	todayPipeline := []bson.M{
		{"$match": todayQuery},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	todayCursor, _ := collection.Aggregate(ctx, todayPipeline)
	if todayCursor != nil {
		defer todayCursor.Close(ctx)
		var results []struct{ Total float64 `bson:"total"` }
		if todayCursor.All(ctx, &results) == nil && len(results) > 0 {
			stats.AmountToday = results[0].Total
		}
	}

	// This week
	weekQuery := bson.M{
		"donation_date": bson.M{"$gte": startOfWeek},
		"status":        entities.DonationStatusCompleted,
	}
	weekCount, _ := collection.CountDocuments(ctx, weekQuery)
	stats.DonationsThisWeek = weekCount

	weekPipeline := []bson.M{
		{"$match": weekQuery},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	weekCursor, _ := collection.Aggregate(ctx, weekPipeline)
	if weekCursor != nil {
		defer weekCursor.Close(ctx)
		var results []struct{ Total float64 `bson:"total"` }
		if weekCursor.All(ctx, &results) == nil && len(results) > 0 {
			stats.AmountThisWeek = results[0].Total
		}
	}

	// This month
	monthQuery := bson.M{
		"donation_date": bson.M{"$gte": startOfMonth},
		"status":        entities.DonationStatusCompleted,
	}
	monthCount, _ := collection.CountDocuments(ctx, monthQuery)
	stats.DonationsThisMonth = monthCount

	monthPipeline := []bson.M{
		{"$match": monthQuery},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	monthCursor, _ := collection.Aggregate(ctx, monthPipeline)
	if monthCursor != nil {
		defer monthCursor.Close(ctx)
		var results []struct{ Total float64 `bson:"total"` }
		if monthCursor.All(ctx, &results) == nil && len(results) > 0 {
			stats.AmountThisMonth = results[0].Total
		}
	}

	// This year
	yearQuery := bson.M{
		"donation_date": bson.M{"$gte": startOfYear},
		"status":        entities.DonationStatusCompleted,
	}
	yearCount, _ := collection.CountDocuments(ctx, yearQuery)
	stats.DonationsThisYear = yearCount

	yearPipeline := []bson.M{
		{"$match": yearQuery},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	yearCursor, _ := collection.Aggregate(ctx, yearPipeline)
	if yearCursor != nil {
		defer yearCursor.Close(ctx)
		var results []struct{ Total float64 `bson:"total"` }
		if yearCursor.All(ctx, &results) == nil && len(results) > 0 {
			stats.AmountThisYear = results[0].Total
		}
	}

	// Recurring donations
	recurringQuery := bson.M{
		"is_recurring": true,
		"recurring_info.active": true,
	}
	recurringCount, _ := collection.CountDocuments(ctx, recurringQuery)
	stats.RecurringDonations = recurringCount

	recurringPipeline := []bson.M{
		{"$match": recurringQuery},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}
	recurringCursor, _ := collection.Aggregate(ctx, recurringPipeline)
	if recurringCursor != nil {
		defer recurringCursor.Close(ctx)
		var results []struct{ Total float64 `bson:"total"` }
		if recurringCursor.All(ctx, &results) == nil && len(results) > 0 {
			stats.RecurringAmount = results[0].Total
		}
	}

	// By type
	typePipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}},
	}
	typeCursor, _ := collection.Aggregate(ctx, typePipeline)
	if typeCursor != nil {
		defer typeCursor.Close(ctx)
		var typeResults []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if typeCursor.All(ctx, &typeResults) == nil {
			for _, result := range typeResults {
				stats.ByType[result.ID] = result.Count
			}
		}
	}

	// By status
	statusPipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}},
	}
	statusCursor, _ := collection.Aggregate(ctx, statusPipeline)
	if statusCursor != nil {
		defer statusCursor.Close(ctx)
		var statusResults []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if statusCursor.All(ctx, &statusResults) == nil {
			for _, result := range statusResults {
				stats.ByStatus[result.ID] = result.Count
			}
		}
	}

	// By payment method
	paymentPipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$payment.method",
			"count": bson.M{"$sum": 1},
		}},
	}
	paymentCursor, _ := collection.Aggregate(ctx, paymentPipeline)
	if paymentCursor != nil {
		defer paymentCursor.Close(ctx)
		var paymentResults []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if paymentCursor.All(ctx, &paymentResults) == nil {
			for _, result := range paymentResults {
				stats.ByPaymentMethod[result.ID] = result.Count
			}
		}
	}

	return stats, nil
}

// EnsureIndexes creates necessary indexes for the donations collection
func (r *donationRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Donations)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "donor_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "campaign_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "donation_date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "amount", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "payment.method", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_recurring", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "thank_you_sent", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "donor_id", Value: 1},
				{Key: "donation_date", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "campaign_id", Value: 1},
				{Key: "donation_date", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "donation_date", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
