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

type eventAttendanceRepository struct {
	db *mongodb.Database
}

// NewEventAttendanceRepository creates a new event attendance repository
func NewEventAttendanceRepository(db *mongodb.Database) repositories.EventAttendanceRepository {
	return &eventAttendanceRepository{db: db}
}

func (r *eventAttendanceRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.EventAttendances)
}

// EnsureIndexes creates necessary indexes for event_attendances collection
func (r *eventAttendanceRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "event_id", Value: 1}}},
		{Keys: bson.D{{Key: "volunteer_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "donor_id", Value: 1}}},
		{Keys: bson.D{{Key: "attendee_type", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "payment_status", Value: 1}}},
		{Keys: bson.D{
			{Key: "event_id", Value: 1},
			{Key: "status", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "event_id", Value: 1},
			{Key: "volunteer_id", Value: 1},
		}},
		{Keys: bson.D{{Key: "registration_date", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *eventAttendanceRepository) Create(ctx context.Context, attendance *entities.EventAttendance) error {
	attendance.ID = primitive.NewObjectID()
	attendance.CreatedAt = time.Now()
	attendance.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, attendance)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create event attendance")
	}
	return nil
}

func (r *eventAttendanceRepository) Update(ctx context.Context, attendance *entities.EventAttendance) error {
	attendance.UpdatedAt = time.Now()

	filter := bson.M{"_id": attendance.ID}
	update := bson.M{"$set": attendance}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update event attendance")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *eventAttendanceRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete event attendance")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *eventAttendanceRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.EventAttendance, error) {
	var attendance entities.EventAttendance
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&attendance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find event attendance")
	}

	return &attendance, nil
}

func (r *eventAttendanceRepository) List(ctx context.Context, filter *repositories.EventAttendanceFilter) ([]*entities.EventAttendance, int64, error) {
	query := bson.M{}

	if filter.EventID != nil {
		query["event_id"] = *filter.EventID
	}

	if filter.VolunteerID != nil {
		query["volunteer_id"] = *filter.VolunteerID
	}

	if filter.UserID != nil {
		query["user_id"] = *filter.UserID
	}

	if filter.DonorID != nil {
		query["donor_id"] = *filter.DonorID
	}

	if filter.AttendeeType != "" {
		query["attendee_type"] = filter.AttendeeType
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count event attendances")
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
		return nil, 0, errors.Wrap(err, 500, "Failed to list event attendances")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, total, nil
}

func (r *eventAttendanceRepository) GetAttendanceByEvent(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	query := bson.M{"event_id": eventID}
	findOptions := options.Find().SetSort(bson.D{{Key: "registration_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get attendance by event")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) GetAttendanceByVolunteer(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	query := bson.M{"volunteer_id": volunteerID}
	findOptions := options.Find().SetSort(bson.D{{Key: "registration_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get attendance by volunteer")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) GetAttendanceByDonor(ctx context.Context, donorID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	query := bson.M{"donor_id": donorID}
	findOptions := options.Find().SetSort(bson.D{{Key: "registration_date", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get attendance by donor")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) GetConfirmedAttendees(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	query := bson.M{
		"event_id": eventID,
		"status": bson.M{"$in": []string{
			string(entities.AttendanceStatusConfirmed),
			string(entities.AttendanceStatusAttended),
		}},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get confirmed attendees")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) GetNoShows(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	query := bson.M{
		"event_id": eventID,
		"status":   entities.AttendanceStatusNoShow,
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get no-shows")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) GetPendingPayments(ctx context.Context) ([]*entities.EventAttendance, error) {
	query := bson.M{
		"payment_status": "pending",
		"registration_fee": bson.M{"$gt": 0},
	}

	cursor, err := r.collection().Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get pending payments")
	}
	defer cursor.Close(ctx)

	var attendances []*entities.EventAttendance
	if err := cursor.All(ctx, &attendances); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode event attendances")
	}

	return attendances, nil
}

func (r *eventAttendanceRepository) CountAttendeesByEvent(ctx context.Context, eventID primitive.ObjectID) (int64, error) {
	query := bson.M{
		"event_id": eventID,
		"status": bson.M{"$in": []string{
			string(entities.AttendanceStatusConfirmed),
			string(entities.AttendanceStatusAttended),
		}},
	}

	count, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return 0, errors.Wrap(err, 500, "Failed to count attendees")
	}

	return count, nil
}
