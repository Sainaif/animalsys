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

type communicationRepository struct {
	db *mongodb.Database
}

// NewCommunicationRepository creates a new communication repository
func NewCommunicationRepository(db *mongodb.Database) repositories.CommunicationRepository {
	return &communicationRepository{db: db}
}

func (r *communicationRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Communications)
}

// EnsureIndexes creates necessary indexes for communications collection
func (r *communicationRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "recipient_type", Value: 1}}},
		{Keys: bson.D{{Key: "recipient_id", Value: 1}}},
		{Keys: bson.D{{Key: "sender_id", Value: 1}}},
		{Keys: bson.D{{Key: "template_id", Value: 1}}},
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}},
		{Keys: bson.D{{Key: "batch_id", Value: 1}}},
		{Keys: bson.D{{Key: "related_type", Value: 1}}},
		{Keys: bson.D{{Key: "related_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "recipient_id", Value: 1},
			{Key: "status", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "scheduled_at", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "next_retry_at", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "campaign_id", Value: 1},
			{Key: "status", Value: 1},
		}},
		{Keys: bson.D{{Key: "scheduled_at", Value: 1}}},
		{Keys: bson.D{{Key: "sent_at", Value: -1}}},
		{Keys: bson.D{{Key: "delivered_at", Value: -1}}},
		{Keys: bson.D{{Key: "next_retry_at", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *communicationRepository) Create(ctx context.Context, communication *entities.Communication) error {
	communication.ID = primitive.NewObjectID()
	communication.CreatedAt = time.Now()
	communication.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, communication)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create communication")
	}
	return nil
}

func (r *communicationRepository) Update(ctx context.Context, communication *entities.Communication) error {
	communication.UpdatedAt = time.Now()

	filter := bson.M{"_id": communication.ID}
	update := bson.M{"$set": communication}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update communication")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete communication")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Communication, error) {
	var communication entities.Communication
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&communication)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find communication")
	}

	return &communication, nil
}

func (r *communicationRepository) List(ctx context.Context, filter *repositories.CommunicationFilter) ([]*entities.Communication, int64, error) {
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.RecipientType != "" {
		query["recipient_type"] = filter.RecipientType
	}

	if filter.RecipientID != nil {
		query["recipient_id"] = *filter.RecipientID
	}

	if filter.SenderID != nil {
		query["sender_id"] = *filter.SenderID
	}

	if filter.TemplateID != nil {
		query["template_id"] = *filter.TemplateID
	}

	if filter.CampaignID != nil {
		query["campaign_id"] = *filter.CampaignID
	}

	if filter.BatchID != "" {
		query["batch_id"] = filter.BatchID
	}

	if filter.RelatedType != "" {
		query["related_type"] = filter.RelatedType
	}

	if filter.RelatedID != nil {
		query["related_id"] = *filter.RelatedID
	}

	// Date range filter
	if filter.StartDate != nil || filter.EndDate != nil {
		dateQuery := bson.M{}
		if filter.StartDate != nil {
			dateQuery["$gte"] = *filter.StartDate
		}
		if filter.EndDate != nil {
			dateQuery["$lte"] = *filter.EndDate
		}
		query["sent_at"] = dateQuery
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count communications")
	}

	// Find with pagination
	findOptions := options.Find()
	findOptions.SetLimit(filter.Limit)
	findOptions.SetSkip(filter.Offset)

	// Sort
	sortOrder := 1
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to list communications")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, total, nil
}

func (r *communicationRepository) GetPending(ctx context.Context) ([]*entities.Communication, error) {
	now := time.Now()
	query := bson.M{
		"status": entities.CommunicationStatusPending,
		"$or": []bson.M{
			{"scheduled_at": nil},
			{"scheduled_at": bson.M{"$lte": now}},
		},
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "scheduled_at", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get pending communications")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, nil
}

func (r *communicationRepository) GetForRetry(ctx context.Context) ([]*entities.Communication, error) {
	now := time.Now()
	query := bson.M{
		"status":        entities.CommunicationStatusPending,
		"retry_count":   bson.M{"$gt": 0, "$lt": 3},
		"next_retry_at": bson.M{"$lte": now},
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "next_retry_at", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get communications for retry")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, nil
}

func (r *communicationRepository) GetByRecipient(ctx context.Context, recipientType entities.RecipientType, recipientID primitive.ObjectID) ([]*entities.Communication, error) {
	query := bson.M{
		"recipient_type": recipientType,
		"recipient_id":   recipientID,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "sent_at", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get communications by recipient")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, nil
}

func (r *communicationRepository) GetByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Communication, error) {
	query := bson.M{"campaign_id": campaignID}
	findOptions := options.Find().SetSort(bson.D{{Key: "sent_at", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get communications by campaign")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, nil
}

func (r *communicationRepository) GetByBatch(ctx context.Context, batchID string) ([]*entities.Communication, error) {
	query := bson.M{"batch_id": batchID}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get communications by batch")
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err := cursor.All(ctx, &communications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communications")
	}

	return communications, nil
}

func (r *communicationRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status entities.CommunicationStatus) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	// Set sent_at or delivered_at based on status
	if status == entities.CommunicationStatusSent {
		update["$set"].(bson.M)["sent_at"] = time.Now()
	} else if status == entities.CommunicationStatusDelivered {
		update["$set"].(bson.M)["delivered_at"] = time.Now()
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update communication status")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationRepository) MarkAsOpened(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"open_count": 1},
		"$set": bson.M{
			"opened_at":  time.Now(),
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to mark communication as opened")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationRepository) MarkAsClicked(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"click_count": 1},
		"$set": bson.M{
			"clicked_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to mark communication as clicked")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationRepository) GetStatistics(ctx context.Context, startDate, endDate time.Time) (*repositories.CommunicationStatistics, error) {
	stats := &repositories.CommunicationStatistics{
		ByType:     make(map[string]int64),
		ByCategory: make(map[string]int64),
		ByStatus:   make(map[string]int64),
	}

	query := bson.M{
		"sent_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	// Count by status
	pipeline := []bson.M{
		{"$match": query},
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int64  `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByStatus[result.ID] = result.Count

				// Update totals
				switch result.ID {
				case string(entities.CommunicationStatusSent):
					stats.TotalSent += result.Count
				case string(entities.CommunicationStatusDelivered):
					stats.TotalSent += result.Count
					stats.TotalDelivered += result.Count
				case string(entities.CommunicationStatusFailed):
					stats.TotalFailed += result.Count
				case string(entities.CommunicationStatusBounced):
					stats.TotalBounced += result.Count
				}
			}
		}
	}

	// Count opened and clicked
	openedCount, _ := r.collection().CountDocuments(ctx, bson.M{
		"sent_at":    bson.M{"$gte": startDate, "$lte": endDate},
		"open_count": bson.M{"$gt": 0},
	})
	stats.TotalOpened = openedCount

	clickedCount, _ := r.collection().CountDocuments(ctx, bson.M{
		"sent_at":     bson.M{"$gte": startDate, "$lte": endDate},
		"click_count": bson.M{"$gt": 0},
	})
	stats.TotalClicked = clickedCount

	// Calculate rates
	if stats.TotalSent > 0 {
		stats.DeliveryRate = float64(stats.TotalDelivered) / float64(stats.TotalSent) * 100
		stats.OpenRate = float64(stats.TotalOpened) / float64(stats.TotalSent) * 100
		stats.ClickRate = float64(stats.TotalClicked) / float64(stats.TotalSent) * 100
	}

	// Count by type
	typePipeline := []bson.M{
		{"$match": query},
		{
			"$group": bson.M{
				"_id":   "$type",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err = r.collection().Aggregate(ctx, typePipeline)
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int64  `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByType[result.ID] = result.Count
			}
		}
	}

	// Count by category
	categoryPipeline := []bson.M{
		{"$match": query},
		{
			"$group": bson.M{
				"_id":   "$category",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err = r.collection().Aggregate(ctx, categoryPipeline)
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int64  `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByCategory[result.ID] = result.Count
			}
		}
	}

	return stats, nil
}
