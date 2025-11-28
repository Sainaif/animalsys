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

type volunteerRepository struct {
	db *mongodb.Database
}

// NewVolunteerRepository creates a new volunteer repository
func NewVolunteerRepository(db *mongodb.Database) repositories.VolunteerRepository {
	return &volunteerRepository{db: db}
}

func (r *volunteerRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Volunteers)
}

// EnsureIndexes creates necessary indexes for volunteers collection
func (r *volunteerRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "skills.name", Value: 1}}},
		{Keys: bson.D{{Key: "total_hours", Value: -1}}},
		{Keys: bson.D{{Key: "rating", Value: -1}}},
		{Keys: bson.D{{Key: "application_date", Value: -1}}},
		{Keys: bson.D{{Key: "last_activity_date", Value: -1}}},
		{Keys: bson.D{
			{Key: "first_name", Value: 1},
			{Key: "last_name", Value: 1},
		}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *volunteerRepository) Create(ctx context.Context, volunteer *entities.Volunteer) error {
	volunteer.ID = primitive.NewObjectID()
	volunteer.CreatedAt = time.Now()
	volunteer.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, volunteer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewConflict("Volunteer with this email already exists")
		}
		return errors.Wrap(err, 500, "Failed to create volunteer")
	}
	return nil
}

func (r *volunteerRepository) Update(ctx context.Context, volunteer *entities.Volunteer) error {
	volunteer.UpdatedAt = time.Now()

	filter := bson.M{"_id": volunteer.ID}
	update := bson.M{"$set": volunteer}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewConflict("Volunteer with this email already exists")
		}
		return errors.Wrap(err, 500, "Failed to update volunteer")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete volunteer")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Volunteer, error) {
	var volunteer entities.Volunteer
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&volunteer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find volunteer")
	}

	return &volunteer, nil
}

func (r *volunteerRepository) FindByEmail(ctx context.Context, email string) (*entities.Volunteer, error) {
	var volunteer entities.Volunteer
	filter := bson.M{"email": email}

	err := r.collection().FindOne(ctx, filter).Decode(&volunteer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find volunteer")
	}

	return &volunteer, nil
}

func (r *volunteerRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Volunteer, error) {
	var volunteer entities.Volunteer
	filter := bson.M{"user_id": userID}

	err := r.collection().FindOne(ctx, filter).Decode(&volunteer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find volunteer")
	}

	return &volunteer, nil
}

func (r *volunteerRepository) List(ctx context.Context, filter *repositories.VolunteerFilter) ([]*entities.Volunteer, int64, error) {
	query := bson.M{}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if len(filter.Skills) > 0 {
		query["skills.name"] = bson.M{"$in": filter.Skills}
	}

	if len(filter.Roles) > 0 {
		query["roles"] = bson.M{"$in": filter.Roles}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"first_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"last_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"email": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if filter.HasUserID != nil {
		if *filter.HasUserID {
			query["user_id"] = bson.M{"$ne": nil}
		} else {
			query["user_id"] = nil
		}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count volunteers")
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
		return nil, 0, errors.Wrap(err, 500, "Failed to list volunteers")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, total, nil
}

func (r *volunteerRepository) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	query := bson.M{"status": entities.VolunteerStatusActive}
	findOptions := options.Find().SetSort(bson.D{{Key: "last_name", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get active volunteers")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, nil
}

func (r *volunteerRepository) GetVolunteersBySkill(ctx context.Context, skill string) ([]*entities.Volunteer, error) {
	query := bson.M{
		"skills.name": skill,
		"status":      entities.VolunteerStatusActive,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "skills.level", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get volunteers by skill")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, nil
}

func (r *volunteerRepository) GetVolunteersNeedingBackgroundCheck(ctx context.Context) ([]*entities.Volunteer, error) {
	query := bson.M{
		"$or": []bson.M{
			{"background_check": nil},
			{"background_check.status": bson.M{"$ne": "passed"}},
			{"background_check.expiration_date": bson.M{"$lt": time.Now()}},
		},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get volunteers needing background check")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, nil
}

func (r *volunteerRepository) GetVolunteersWithExpiredCertifications(ctx context.Context) ([]*entities.Volunteer, error) {
	query := bson.M{
		"certifications.expiration_date": bson.M{"$lt": time.Now()},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get volunteers with expired certifications")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, nil
}

func (r *volunteerRepository) GetTopVolunteers(ctx context.Context, limit int) ([]*entities.Volunteer, error) {
	findOptions := options.Find().
		SetSort(bson.D{{Key: "total_hours", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection().Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get top volunteers")
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err := cursor.All(ctx, &volunteers); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteers")
	}

	return volunteers, nil
}

func (r *volunteerRepository) UpdateHours(ctx context.Context, volunteerID primitive.ObjectID, hours float64, notes string) error {
	filter := bson.M{"_id": volunteerID}
	update := bson.M{
		"$inc": bson.M{"total_hours": hours},
		"$set": bson.M{
			"last_activity_date": time.Now(),
			"updated_at":         time.Now(),
		},
	}

	if notes != "" {
		update["$push"] = bson.M{"notes": notes}
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update volunteer hours")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerRepository) IncrementEventsAttended(ctx context.Context, volunteerID primitive.ObjectID) error {
	filter := bson.M{"_id": volunteerID}
	update := bson.M{
		"$inc": bson.M{"events_attended": 1},
		"$set": bson.M{
			"last_activity_date": time.Now(),
			"updated_at":         time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to increment events attended")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerRepository) GetVolunteerStatistics(ctx context.Context) (*repositories.VolunteerStatistics, error) {
	stats := &repositories.VolunteerStatistics{
		ByStatus: make(map[string]int64),
	}

	// Total volunteers
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to count volunteers")
	}
	stats.TotalVolunteers = total

	// Active volunteers
	activeCount, _ := r.collection().CountDocuments(ctx, bson.M{"status": entities.VolunteerStatusActive})
	stats.ActiveVolunteers = activeCount

	// Inactive volunteers
	inactiveCount, _ := r.collection().CountDocuments(ctx, bson.M{"status": entities.VolunteerStatusInactive})
	stats.InactiveVolunteers = inactiveCount

	// Aggregate statistics
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":           nil,
				"total_hours":   bson.M{"$sum": "$total_hours"},
				"avg_hours":     bson.M{"$avg": "$total_hours"},
				"total_events":  bson.M{"$sum": "$events_attended"},
				"avg_rating":    bson.M{"$avg": "$rating"},
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
			TotalHours   float64 `bson:"total_hours"`
			AverageHours float64 `bson:"avg_hours"`
			TotalEvents  int64   `bson:"total_events"`
			AverageRating float64 `bson:"avg_rating"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.TotalHours = result.TotalHours
			stats.AverageHours = result.AverageHours
			stats.TotalEvents = result.TotalEvents
			stats.AverageRating = result.AverageRating
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

	// Top volunteers
	topVolunteers, err := r.GetTopVolunteers(ctx, 10)
	if err == nil {
		stats.TopVolunteers = topVolunteers
	}

	return stats, nil
}
