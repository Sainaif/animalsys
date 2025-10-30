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

type partnerRepository struct {
	collection *mongo.Collection
}

func NewPartnerRepository(db *mongo.Database) interfaces.PartnerRepository {
	return &partnerRepository{
		collection: db.Collection("partners"),
	}
}

func (r *partnerRepository) Create(ctx context.Context, partner *entities.Partner) error {
	partner.CreatedAt = time.Now()
	partner.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, partner)
	return err
}

func (r *partnerRepository) GetByID(ctx context.Context, id string) (*entities.Partner, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid partner ID")
	}

	var partner entities.Partner
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}
	return &partner, nil
}

func (r *partnerRepository) Update(ctx context.Context, id string, partner *entities.Partner) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid partner ID")
	}

	partner.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": partner})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("partner not found")
	}
	return nil
}

func (r *partnerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid partner ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("partner not found")
	}
	return nil
}

func (r *partnerRepository) List(ctx context.Context, partnerType entities.PartnerType, limit, offset int) ([]*entities.Partner, int64, error) {
	filter := bson.M{}
	if partnerType != "" {
		filter["type"] = partnerType
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"name": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
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

func (r *partnerRepository) GetActivePartners(ctx context.Context) ([]*entities.Partner, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"active": true})
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

func (r *partnerRepository) AddCollaboration(ctx context.Context, id string, collaboration entities.Collaboration) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid partner ID")
	}

	update := bson.M{
		"$push": bson.M{"collaboration_history": collaboration},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("partner not found")
	}
	return nil
}
