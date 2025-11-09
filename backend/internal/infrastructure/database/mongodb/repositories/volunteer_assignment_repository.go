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

type volunteerAssignmentRepository struct {
	db *mongodb.Database
}

// NewVolunteerAssignmentRepository creates a new volunteer assignment repository
func NewVolunteerAssignmentRepository(db *mongodb.Database) repositories.VolunteerAssignmentRepository {
	return &volunteerAssignmentRepository{db: db}
}

func (r *volunteerAssignmentRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.VolunteerAssignments)
}

// EnsureIndexes creates necessary indexes for volunteer_assignments collection
func (r *volunteerAssignmentRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "volunteer_id", Value: 1}}},
		{Keys: bson.D{{Key: "event_id", Value: 1}}},
		{Keys: bson.D{{Key: "animal_id", Value: 1}}},
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "start_date", Value: 1}}},
		{Keys: bson.D{
			{Key: "volunteer_id", Value: 1},
			{Key: "status", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "volunteer_id", Value: 1},
			{Key: "start_date", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "start_date", Value: 1},
		}},
		{Keys: bson.D{{Key: "assigned_date", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *volunteerAssignmentRepository) Create(ctx context.Context, assignment *entities.VolunteerAssignment) error {
	assignment.ID = primitive.NewObjectID()
	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, assignment)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create volunteer assignment")
	}
	return nil
}

func (r *volunteerAssignmentRepository) Update(ctx context.Context, assignment *entities.VolunteerAssignment) error {
	assignment.UpdatedAt = time.Now()

	filter := bson.M{"_id": assignment.ID}
	update := bson.M{"$set": assignment}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update volunteer assignment")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerAssignmentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete volunteer assignment")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *volunteerAssignmentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.VolunteerAssignment, error) {
	var assignment entities.VolunteerAssignment
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&assignment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find volunteer assignment")
	}

	return &assignment, nil
}

func (r *volunteerAssignmentRepository) List(ctx context.Context, filter *repositories.VolunteerAssignmentFilter) ([]*entities.VolunteerAssignment, int64, error) {
	query := bson.M{}

	if filter.VolunteerID != nil {
		query["volunteer_id"] = *filter.VolunteerID
	}

	if filter.EventID != nil {
		query["event_id"] = *filter.EventID
	}

	if filter.AnimalID != nil {
		query["animal_id"] = *filter.AnimalID
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

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count volunteer assignments")
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
		return nil, 0, errors.Wrap(err, 500, "Failed to list volunteer assignments")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, total, nil
}

func (r *volunteerAssignmentRepository) GetAssignmentsByVolunteer(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{"volunteer_id": volunteerID}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get assignments by volunteer")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetAssignmentsByEvent(ctx context.Context, eventID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{"event_id": eventID}
	findOptions := options.Find().SetSort(bson.D{{Key: "assigned_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get assignments by event")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetUpcomingAssignments(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{
		"volunteer_id": volunteerID,
		"start_date":   bson.M{"$gt": time.Now()},
		"status": bson.M{"$in": []string{
			string(entities.AssignmentStatusAssigned),
			string(entities.AssignmentStatusConfirmed),
		}},
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get upcoming assignments")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetActiveAssignments(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{
		"volunteer_id": volunteerID,
		"status": bson.M{"$in": []string{
			string(entities.AssignmentStatusConfirmed),
			string(entities.AssignmentStatusInProgress),
		}},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get active assignments")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetAssignmentsNeedingReminder(ctx context.Context) ([]*entities.VolunteerAssignment, error) {
	// Assignments that need reminders: 24 hours before start, not yet sent
	reminderTime := time.Now().Add(24 * time.Hour)
	query := bson.M{
		"reminder_sent": false,
		"start_date":    bson.M{"$lte": reminderTime, "$gt": time.Now()},
		"status": bson.M{"$in": []string{
			string(entities.AssignmentStatusAssigned),
			string(entities.AssignmentStatusConfirmed),
		}},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get assignments needing reminder")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetCompletedAssignmentsByVolunteer(ctx context.Context, volunteerID primitive.ObjectID, startDate, endDate time.Time) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{
		"volunteer_id": volunteerID,
		"status":       entities.AssignmentStatusCompleted,
		"start_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get completed assignments")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetAssignmentsByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{"animal_id": animalID}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get assignments by animal")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}

func (r *volunteerAssignmentRepository) GetAssignmentsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.VolunteerAssignment, error) {
	query := bson.M{"campaign_id": campaignID}
	findOptions := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get assignments by campaign")
	}
	defer cursor.Close(ctx)

	var assignments []*entities.VolunteerAssignment
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode volunteer assignments")
	}

	return assignments, nil
}
