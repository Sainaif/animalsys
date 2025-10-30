package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type scheduleRepository struct {
	collection *mongo.Collection
}

func NewScheduleRepository(db *mongo.Database) interfaces.ScheduleRepository {
	return &scheduleRepository{
		collection: db.Collection("schedules"),
	}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *entities.Schedule) error {
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, schedule)
	return err
}

func (r *scheduleRepository) GetByID(ctx context.Context, id string) (*entities.Schedule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid schedule ID")
	}

	var schedule entities.Schedule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Update(ctx context.Context, id string, schedule *entities.Schedule) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid schedule ID")
	}

	schedule.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": schedule})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("schedule not found")
	}
	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid schedule ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("schedule not found")
	}
	return nil
}

func (r *scheduleRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.Schedule, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"shift_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var schedules []*entities.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, 0, err
	}
	return schedules, total, nil
}

func (r *scheduleRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Schedule, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	filter := bson.M{
		"shift_date": bson.M{"$gte": start, "$lte": end},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"shift_date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*entities.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *scheduleRepository) GetByStatus(ctx context.Context, status entities.ShiftStatus) ([]*entities.Schedule, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*entities.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *scheduleRepository) UpdateStatus(ctx context.Context, id string, status entities.ShiftStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid schedule ID")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("schedule not found")
	}
	return nil
}

func (r *scheduleRepository) GetUpcomingShifts(ctx context.Context, userID string, limit int) ([]*entities.Schedule, error) {
	filter := bson.M{
		"user_id":    userID,
		"shift_date": bson.M{"$gte": time.Now()},
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"shift_date": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*entities.Schedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}
