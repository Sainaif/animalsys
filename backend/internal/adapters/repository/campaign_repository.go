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

type campaignRepository struct {
	collection *mongo.Collection
}

func NewCampaignRepository(db *mongo.Database) interfaces.CampaignRepository {
	return &campaignRepository{
		collection: db.Collection("campaigns"),
	}
}

func (r *campaignRepository) Create(ctx context.Context, campaign *entities.Campaign) error {
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, campaign)
	return err
}

func (r *campaignRepository) GetByID(ctx context.Context, id string) (*entities.Campaign, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid campaign ID")
	}

	var campaign entities.Campaign
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&campaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("campaign not found")
		}
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) Update(ctx context.Context, id string, campaign *entities.Campaign) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid campaign ID")
	}

	campaign.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": campaign})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("campaign not found")
	}
	return nil
}

func (r *campaignRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid campaign ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("campaign not found")
	}
	return nil
}

func (r *campaignRepository) List(ctx context.Context, status entities.CampaignStatus, limit, offset int) ([]*entities.Campaign, int64, error) {
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
		SetSort(bson.M{"start_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, 0, err
	}
	return campaigns, total, nil
}

func (r *campaignRepository) GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	now := time.Now()
	filter := bson.M{
		"status":     entities.CampaignStatusActive,
		"start_date": bson.M{"$lte": now},
		"end_date":   bson.M{"$gte": now},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var campaigns []*entities.Campaign
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) UpdateProgress(ctx context.Context, id string, progress float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid campaign ID")
	}

	update := bson.M{
		"$set": bson.M{
			"current_progress": progress,
			"updated_at":       time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("campaign not found")
	}
	return nil
}

func (r *campaignRepository) AddExpense(ctx context.Context, id string, amount float64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid campaign ID")
	}

	update := bson.M{
		"$inc": bson.M{"total_expenses": amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("campaign not found")
	}
	return nil
}
