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

type eventRepository struct {
	db *mongodb.Database
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *mongodb.Database) repositories.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Events)
}

// EnsureIndexes creates necessary indexes for events collection
func (r *eventRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "start_date", Value: 1}}},
		{Keys: bson.D{{Key: "public", Value: 1}}},
		{Keys: bson.D{{Key: "featured", Value: 1}}},
		{Keys: bson.D{{Key: "organizer", Value: 1}}},
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "start_date", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "public", Value: 1},
			{Key: "featured", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "public", Value: 1},
		}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *eventRepository) Create(ctx context.Context, event *entities.Event) error {
	event.ID = primitive.NewObjectID()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, event)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create event")
	}
	return nil
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) error {
	event.UpdatedAt = time.Now()

	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update event")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete event")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *eventRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Event, error) {
	var event entities.Event
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find event")
	}

	return &event, nil
}

func (r *eventRepository) List(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error) {
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

	if filter.StartDate != nil {
		query["start_date"] = bson.M{"$gte": *filter.StartDate}
	}

	if filter.EndDate != nil {
		if query["start_date"] != nil {
			query["start_date"].(bson.M)["$lte"] = *filter.EndDate
		} else {
			query["start_date"] = bson.M{"$lte": *filter.EndDate}
		}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name.english": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"name.polish": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description.english": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count events")
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
	findOptions.SetSort(bson.D{{Key: filter.SortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to list events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, total, nil
}

func (r *eventRepository) GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error) {
	query := bson.M{
		"status": entities.EventStatusScheduled,
		"start_date": bson.M{"$gt": time.Now()},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get upcoming events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetActiveEvents(ctx context.Context) ([]*entities.Event, error) {
	now := time.Now()
	query := bson.M{
		"status": entities.EventStatusActive,
		"start_date": bson.M{"$lte": now},
		"$or": []bson.M{
			{"end_date": nil},
			{"end_date": bson.M{"$gte": now}},
		},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get active events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetCompletedEvents(ctx context.Context, limit int) ([]*entities.Event, error) {
	query := bson.M{"status": entities.EventStatusCompleted}
	findOptions := options.Find().
		SetSort(bson.D{{Key: "start_date", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get completed events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetPublicEvents(ctx context.Context) ([]*entities.Event, error) {
	query := bson.M{"public": true}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get public events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error) {
	query := bson.M{
		"featured": true,
		"public": true,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get featured events")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error) {
	query := bson.M{"organizer": organizerID}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get events by organizer")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetEventsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Event, error) {
	query := bson.M{"campaign_id": campaignID}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get events by campaign")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error) {
	// Events where assigned_volunteers count < required_volunteers
	query := bson.M{
		"$expr": bson.M{
			"$lt": []interface{}{
				bson.M{"$size": "$assigned_volunteers"},
				"$required_volunteers",
			},
		},
		"status": bson.M{"$in": []string{
			string(entities.EventStatusScheduled),
			string(entities.EventStatusActive),
		}},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get events needing volunteers")
	}
	defer cursor.Close(ctx)

	var events []*entities.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode events")
	}

	return events, nil
}

func (r *eventRepository) UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error {
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$set": bson.M{
			"attendee_count":   attendees,
			"volunteer_count":  volunteers,
			"funds_raised":     fundsRaised,
			"animals_adopted":  animalsAdopted,
			"updated_at":       time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update event statistics")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *eventRepository) GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error) {
	stats := &repositories.EventStatistics{
		ByType:   make(map[string]int64),
		ByStatus: make(map[string]int64),
	}

	// Total events
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to count events")
	}
	stats.TotalEvents = total

	// By status
	upcomingCount, _ := r.collection().CountDocuments(ctx, bson.M{
		"status": entities.EventStatusScheduled,
		"start_date": bson.M{"$gt": time.Now()},
	})
	stats.UpcomingEvents = upcomingCount

	activeCount, _ := r.collection().CountDocuments(ctx, bson.M{"status": entities.EventStatusActive})
	stats.ActiveEvents = activeCount

	completedCount, _ := r.collection().CountDocuments(ctx, bson.M{"status": entities.EventStatusCompleted})
	stats.CompletedEvents = completedCount

	// Aggregate statistics
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":                   nil,
				"total_attendees":       bson.M{"$sum": "$attendee_count"},
				"total_volunteers":      bson.M{"$sum": "$volunteer_count"},
				"total_funds_raised":    bson.M{"$sum": "$funds_raised"},
				"total_animals_adopted": bson.M{"$sum": "$animals_adopted"},
				"avg_attendees":         bson.M{"$avg": "$attendee_count"},
			},
		},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err != nil {
		return stats, nil // Return partial stats
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalAttendees      int64   `bson:"total_attendees"`
			TotalVolunteers     int64   `bson:"total_volunteers"`
			TotalFundsRaised    float64 `bson:"total_funds_raised"`
			TotalAnimalsAdopted int64   `bson:"total_animals_adopted"`
			AverageAttendees    float64 `bson:"avg_attendees"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.TotalAttendees = result.TotalAttendees
			stats.TotalVolunteers = result.TotalVolunteers
			stats.TotalFundsRaised = result.TotalFundsRaised
			stats.TotalAnimalsAdopted = result.TotalAnimalsAdopted
			stats.AverageAttendees = result.AverageAttendees
		}
	}

	// By type
	typePipeline := []bson.M{
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

	// By status
	statusPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err = r.collection().Aggregate(ctx, statusPipeline)
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int64  `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByStatus[result.ID] = result.Count
			}
		}
	}

	return stats, nil
}
