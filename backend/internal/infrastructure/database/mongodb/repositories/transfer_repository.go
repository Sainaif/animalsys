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

type transferRepository struct {
	db *mongodb.Database
}

// NewTransferRepository creates a new transfer repository
func NewTransferRepository(db *mongodb.Database) repositories.TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) collection() *mongo.Collection {
	return r.db.DB.Collection("transfers")
}

// Create creates a new transfer
func (r *transferRepository) Create(ctx context.Context, transfer *entities.Transfer) error {
	if transfer.ID.IsZero() {
		transfer.ID = primitive.NewObjectID()
	}

	_, err := r.collection().InsertOne(ctx, transfer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrConflict
		}
		return err
	}

	return nil
}

// FindByID finds a transfer by ID
func (r *transferRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Transfer, error) {
	var transfer entities.Transfer
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&transfer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &transfer, nil
}

// Update updates a transfer
func (r *transferRepository) Update(ctx context.Context, transfer *entities.Transfer) error {
	transfer.UpdatedAt = time.Now()

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": transfer.ID},
		bson.M{"$set": transfer},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a transfer
func (r *transferRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List lists transfers with filtering and pagination
func (r *transferRepository) List(ctx context.Context, filter *repositories.TransferFilter) ([]*entities.Transfer, int64, error) {
	query := bson.M{}

	// Apply filters
	if filter.Direction != "" {
		query["direction"] = filter.Direction
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Reason != "" {
		query["reason"] = filter.Reason
	}

	if filter.AnimalID != nil {
		query["animal_id"] = filter.AnimalID
	}

	if filter.PartnerID != nil {
		query["partner_id"] = filter.PartnerID
	}

	if filter.RequestedBy != nil {
		query["requested_by"] = filter.RequestedBy
	}

	if filter.ApprovedBy != nil {
		query["approved_by"] = filter.ApprovedBy
	}

	if filter.ScheduledAfter != nil {
		if query["scheduled_date"] == nil {
			query["scheduled_date"] = bson.M{}
		}
		query["scheduled_date"].(bson.M)["$gte"] = filter.ScheduledAfter
	}

	if filter.ScheduledBefore != nil {
		if query["scheduled_date"] == nil {
			query["scheduled_date"] = bson.M{}
		}
		query["scheduled_date"].(bson.M)["$lte"] = filter.ScheduledBefore
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"notes": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"transfer_reference": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options
	opts := options.Find()

	// Sorting
	sortField := "scheduled_date"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1 // descending by default
	if filter.SortOrder == "asc" {
		sortOrder = 1
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

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// GetByAnimal gets transfers for a specific animal
func (r *transferRepository) GetByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Transfer, error) {
	query := bson.M{"animal_id": animalID}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetByPartner gets transfers for a specific partner
func (r *transferRepository) GetByPartner(ctx context.Context, partnerID primitive.ObjectID) ([]*entities.Transfer, error) {
	query := bson.M{"partner_id": partnerID}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetByStatus gets transfers by status
func (r *transferRepository) GetByStatus(ctx context.Context, status entities.TransferStatus) ([]*entities.Transfer, error) {
	query := bson.M{"status": status}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetPendingTransfers gets all pending transfers
func (r *transferRepository) GetPendingTransfers(ctx context.Context) ([]*entities.Transfer, error) {
	return r.GetByStatus(ctx, entities.TransferStatusPending)
}

// GetUpcomingTransfers gets transfers scheduled within the next N days
func (r *transferRepository) GetUpcomingTransfers(ctx context.Context, days int) ([]*entities.Transfer, error) {
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	query := bson.M{
		"status": entities.TransferStatusApproved,
		"scheduled_date": bson.M{
			"$gte": now,
			"$lte": futureDate,
		},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetOverdueTransfers gets transfers that are overdue
func (r *transferRepository) GetOverdueTransfers(ctx context.Context) ([]*entities.Transfer, error) {
	now := time.Now()

	query := bson.M{
		"status": bson.M{"$in": []entities.TransferStatus{
			entities.TransferStatusApproved,
			entities.TransferStatusInTransit,
		}},
		"scheduled_date": bson.M{"$lt": now},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetRequiringFollowUp gets transfers requiring follow-up
func (r *transferRepository) GetRequiringFollowUp(ctx context.Context) ([]*entities.Transfer, error) {
	query := bson.M{"requires_follow_up": true}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []*entities.Transfer
	if err = cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetTransferStatistics gets transfer statistics
func (r *transferRepository) GetTransferStatistics(ctx context.Context) (*repositories.TransferStatistics, error) {
	stats := &repositories.TransferStatistics{
		ByDirection: make(map[string]int64),
		ByStatus:    make(map[string]int64),
		ByReason:    make(map[string]int64),
	}

	// Total transfers
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalTransfers = total

	// Pending transfers
	pendingCount, err := r.collection().CountDocuments(ctx, bson.M{"status": entities.TransferStatusPending})
	if err != nil {
		return nil, err
	}
	stats.PendingTransfers = pendingCount

	// In-transit transfers
	inTransitCount, err := r.collection().CountDocuments(ctx, bson.M{"status": entities.TransferStatusInTransit})
	if err != nil {
		return nil, err
	}
	stats.InTransitTransfers = inTransitCount

	// Completed this month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	completedThisMonth, err := r.collection().CountDocuments(ctx, bson.M{
		"status":       entities.TransferStatusCompleted,
		"completed_at": bson.M{"$gte": startOfMonth},
	})
	if err != nil {
		return nil, err
	}
	stats.CompletedThisMonth = completedThisMonth

	// By direction, status, and reason using aggregation
	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.M{
			"byDirection": []bson.M{
				{"$group": bson.M{"_id": "$direction", "count": bson.M{"$sum": 1}}},
			},
			"byStatus": []bson.M{
				{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
			},
			"byReason": []bson.M{
				{"$group": bson.M{"_id": "$reason", "count": bson.M{"$sum": 1}}},
			},
			"avgDuration": []bson.M{
				{"$match": bson.M{
					"status":       entities.TransferStatusCompleted,
					"started_at":   bson.M{"$exists": true},
					"completed_at": bson.M{"$exists": true},
				}},
				{"$project": bson.M{
					"duration": bson.M{"$subtract": []string{"$completed_at", "$started_at"}},
				}},
				{"$group": bson.M{
					"_id":         nil,
					"avgDuration": bson.M{"$avg": "$duration"},
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

		// By direction
		if byDirection, ok := result["byDirection"].([]interface{}); ok {
			for _, item := range byDirection {
				if m, ok := item.(bson.M); ok {
					direction := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByDirection[direction] = int64(count)
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

		// By reason
		if byReason, ok := result["byReason"].([]interface{}); ok {
			for _, item := range byReason {
				if m, ok := item.(bson.M); ok {
					reason := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByReason[reason] = int64(count)
				}
			}
		}

		// Average duration (in hours)
		if avgDuration, ok := result["avgDuration"].([]interface{}); ok && len(avgDuration) > 0 {
			if m, ok := avgDuration[0].(bson.M); ok {
				if duration, ok := m["avgDuration"].(int64); ok {
					// Convert milliseconds to hours
					stats.AverageDuration = float64(duration) / (1000 * 60 * 60)
				}
			}
		}
	}

	return stats, nil
}

// EnsureIndexes creates the necessary indexes for the transfers collection
func (r *transferRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "animal_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "partner_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "direction", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "scheduled_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "requested_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "requires_follow_up", Value: 1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
