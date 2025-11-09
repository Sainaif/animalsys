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

type partnerRepository struct {
	db *mongodb.Database
}

// NewPartnerRepository creates a new partner repository
func NewPartnerRepository(db *mongodb.Database) repositories.PartnerRepository {
	return &partnerRepository{db: db}
}

func (r *partnerRepository) collection() *mongo.Collection {
	return r.db.DB.Collection("partners")
}

// Create creates a new partner
func (r *partnerRepository) Create(ctx context.Context, partner *entities.Partner) error {
	if partner.ID.IsZero() {
		partner.ID = primitive.NewObjectID()
	}

	_, err := r.collection().InsertOne(ctx, partner)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrConflict
		}
		return err
	}

	return nil
}

// FindByID finds a partner by ID
func (r *partnerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	var partner entities.Partner
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &partner, nil
}

// Update updates a partner
func (r *partnerRepository) Update(ctx context.Context, partner *entities.Partner) error {
	partner.UpdatedAt = time.Now()

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": partner.ID},
		bson.M{"$set": partner},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a partner
func (r *partnerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List lists partners with filtering and pagination
func (r *partnerRepository) List(ctx context.Context, filter *repositories.PartnerFilter) ([]*entities.Partner, int64, error) {
	query := bson.M{}

	// Apply filters
	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.AcceptsIntakes != nil {
		query["accepts_intakes"] = *filter.AcceptsIntakes
	}

	if filter.HasCapacity != nil {
		if *filter.HasCapacity {
			query["$expr"] = bson.M{"$lt": []interface{}{"$current_capacity", "$max_capacity"}}
		}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"contact_name": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options
	opts := options.Find()

	// Sorting
	sortField := "name"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := 1 // ascending by default
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Pagination
	if filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	// Execute query
	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var partners []*entities.Partner
	if err = cursor.All(ctx, &partners); err != nil {
		return nil, 0, err
	}

	return partners, total, nil
}

// GetByType gets partners by type
func (r *partnerRepository) GetByType(ctx context.Context, partnerType entities.PartnerType) ([]*entities.Partner, error) {
	query := bson.M{"type": partnerType}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []*entities.Partner
	if err = cursor.All(ctx, &partners); err != nil {
		return nil, err
	}

	return partners, nil
}

// GetByStatus gets partners by status
func (r *partnerRepository) GetByStatus(ctx context.Context, status entities.PartnerStatus) ([]*entities.Partner, error) {
	query := bson.M{"status": status}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []*entities.Partner
	if err = cursor.All(ctx, &partners); err != nil {
		return nil, err
	}

	return partners, nil
}

// GetActivePartners gets all active partners
func (r *partnerRepository) GetActivePartners(ctx context.Context) ([]*entities.Partner, error) {
	return r.GetByStatus(ctx, entities.PartnerStatusActive)
}

// GetPartnersWithCapacity gets partners with available capacity
func (r *partnerRepository) GetPartnersWithCapacity(ctx context.Context) ([]*entities.Partner, error) {
	query := bson.M{
		"status":          entities.PartnerStatusActive,
		"accepts_intakes": true,
		"$expr":           bson.M{"$lt": []interface{}{"$current_capacity", "$max_capacity"}},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []*entities.Partner
	if err = cursor.All(ctx, &partners); err != nil {
		return nil, err
	}

	return partners, nil
}

// GetPartnerStatistics gets partner statistics
func (r *partnerRepository) GetPartnerStatistics(ctx context.Context) (*repositories.PartnerStatistics, error) {
	stats := &repositories.PartnerStatistics{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
	}

	// Total partners
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalPartners = total

	// Active partners
	activeCount, err := r.collection().CountDocuments(ctx, bson.M{"status": entities.PartnerStatusActive})
	if err != nil {
		return nil, err
	}
	stats.ActivePartners = activeCount

	// Accepting intakes
	acceptingCount, err := r.collection().CountDocuments(ctx, bson.M{
		"status":          entities.PartnerStatusActive,
		"accepts_intakes": true,
	})
	if err != nil {
		return nil, err
	}
	stats.AcceptingIntakes = acceptingCount

	// By type and status
	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.M{
			"byType": []bson.M{
				{"$group": bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}},
			},
			"byStatus": []bson.M{
				{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
			},
			"totals": []bson.M{
				{"$group": bson.M{
					"_id":            nil,
					"totalTransfers": bson.M{"$sum": bson.M{"$add": []interface{}{"$total_transfers_in", "$total_transfers_out"}}},
					"avgRating":      bson.M{"$avg": "$rating"},
				}},
			},
		}}},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		result := results[0]

		// By type
		if byType, ok := result["byType"].([]interface{}); ok {
			for _, item := range byType {
				if m, ok := item.(bson.M); ok {
					partnerType := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByType[partnerType] = int64(count)
				}
			}
		}

		// By status
		if byStatus, ok := result["byStatus"].([]interface{}); ok {
			for _, item := range byStatus {
				if m, ok := item.(bson.M); ok {
					status := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByStatus[status] = int64(count)
				}
			}
		}

		// Totals
		if totals, ok := result["totals"].([]interface{}); ok && len(totals) > 0 {
			if m, ok := totals[0].(bson.M); ok {
				if totalTransfers, ok := m["totalTransfers"].(int32); ok {
					stats.TotalTransfers = int64(totalTransfers)
				}
				if avgRating, ok := m["avgRating"].(float64); ok {
					stats.AverageRating = avgRating
				}
			}
		}
	}

	return stats, nil
}

// EnsureIndexes creates the necessary indexes for the partners collection
func (r *partnerRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "accepts_intakes", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
