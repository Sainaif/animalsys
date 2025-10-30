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

type volunteerRepository struct {
	collection *mongo.Collection
}

func NewVolunteerRepository(db *mongo.Database) interfaces.VolunteerRepository {
	return &volunteerRepository{
		collection: db.Collection("volunteers"),
	}
}

func (r *volunteerRepository) Create(ctx context.Context, volunteer *entities.Volunteer) error {
	volunteer.CreatedAt = time.Now()
	volunteer.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, volunteer)
	return err
}

func (r *volunteerRepository) GetByID(ctx context.Context, id string) (*entities.Volunteer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid volunteer ID")
	}

	var volunteer entities.Volunteer
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&volunteer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("volunteer not found")
		}
		return nil, err
	}
	return &volunteer, nil
}

func (r *volunteerRepository) GetByUserID(ctx context.Context, userID string) (*entities.Volunteer, error) {
	var volunteer entities.Volunteer
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&volunteer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("volunteer not found")
		}
		return nil, err
	}
	return &volunteer, nil
}

func (r *volunteerRepository) Update(ctx context.Context, id string, volunteer *entities.Volunteer) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid volunteer ID")
	}

	volunteer.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": volunteer})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("volunteer not found")
	}
	return nil
}

func (r *volunteerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid volunteer ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("volunteer not found")
	}
	return nil
}

func (r *volunteerRepository) List(ctx context.Context, status entities.VolunteerStatus, limit, offset int) ([]*entities.Volunteer, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"registration_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err = cursor.All(ctx, &volunteers); err != nil {
		return nil, 0, err
	}
	return volunteers, total, nil
}

func (r *volunteerRepository) AddTraining(ctx context.Context, id string, training entities.Training) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid volunteer ID")
	}

	update := bson.M{
		"$push": bson.M{"trainings": training},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("volunteer not found")
	}
	return nil
}

func (r *volunteerRepository) LogHours(ctx context.Context, id string, hours float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid volunteer ID")
	}

	update := bson.M{
		"$inc": bson.M{
			"total_hours":       hours,
			"hours_this_month": hours,
		},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("volunteer not found")
	}
	return nil
}

func (r *volunteerRepository) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": entities.VolunteerStatusActive})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var volunteers []*entities.Volunteer
	if err = cursor.All(ctx, &volunteers); err != nil {
		return nil, err
	}
	return volunteers, nil
}

func (r *volunteerRepository) ResetMonthlyHours(ctx context.Context) error {
	update := bson.M{
		"$set": bson.M{
			"hours_this_month": 0,
			"updated_at":       time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, bson.M{}, update)
	return err
}

type volunteerHourRepository struct {
	collection *mongo.Collection
}

func NewVolunteerHourRepository(db *mongo.Database) interfaces.VolunteerHourRepository {
	return &volunteerHourRepository{
		collection: db.Collection("volunteer_hours"),
	}
}

func (r *volunteerHourRepository) Create(ctx context.Context, entry *entities.VolunteerHourEntry) error {
	entry.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

func (r *volunteerHourRepository) GetByID(ctx context.Context, id string) (*entities.VolunteerHourEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid entry ID")
	}

	var entry entities.VolunteerHourEntry
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *volunteerHourRepository) GetByVolunteerID(ctx context.Context, volunteerID string, limit, offset int) ([]*entities.VolunteerHourEntry, int64, error) {
	filter := bson.M{"volunteer_id": volunteerID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var entries []*entities.VolunteerHourEntry
	if err = cursor.All(ctx, &entries); err != nil {
		return nil, 0, err
	}
	return entries, total, nil
}

func (r *volunteerHourRepository) Update(ctx context.Context, id string, entry *entities.VolunteerHourEntry) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid entry ID")
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": entry})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("entry not found")
	}
	return nil
}

func (r *volunteerHourRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid entry ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("entry not found")
	}
	return nil
}

func (r *volunteerHourRepository) GetTotalHours(ctx context.Context, volunteerID string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"volunteer_id": volunteerID}}},
		{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$hours"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}

func (r *volunteerHourRepository) GetHoursByDateRange(ctx context.Context, volunteerID string, startDate, endDate string) (float64, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"volunteer_id": volunteerID,
			"date":         bson.M{"$gte": start, "$lte": end},
		}}},
		{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$hours"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}
