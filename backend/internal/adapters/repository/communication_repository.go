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

type communicationRepository struct {
	collection *mongo.Collection
}

func NewCommunicationRepository(db *mongo.Database) interfaces.CommunicationRepository {
	return &communicationRepository{
		collection: db.Collection("communications"),
	}
}

func (r *communicationRepository) Create(ctx context.Context, communication *entities.Communication) error {
	communication.CreatedAt = time.Now()
	communication.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, communication)
	return err
}

func (r *communicationRepository) GetByID(ctx context.Context, id string) (*entities.Communication, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid communication ID")
	}

	var communication entities.Communication
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&communication)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("communication not found")
		}
		return nil, err
	}
	return &communication, nil
}

func (r *communicationRepository) Update(ctx context.Context, id string, communication *entities.Communication) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid communication ID")
	}

	communication.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": communication})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("communication not found")
	}
	return nil
}

func (r *communicationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid communication ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("communication not found")
	}
	return nil
}

func (r *communicationRepository) List(ctx context.Context, status entities.CommunicationStatus, limit, offset int) ([]*entities.Communication, int64, error) {
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
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err = cursor.All(ctx, &communications); err != nil {
		return nil, 0, err
	}
	return communications, total, nil
}

func (r *communicationRepository) GetScheduled(ctx context.Context) ([]*entities.Communication, error) {
	filter := bson.M{
		"status":        entities.CommunicationStatusScheduled,
		"scheduled_for": bson.M{"$lte": time.Now()},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"scheduled_for": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var communications []*entities.Communication
	if err = cursor.All(ctx, &communications); err != nil {
		return nil, err
	}
	return communications, nil
}

func (r *communicationRepository) MarkAsSent(ctx context.Context, id string, sentCount, failedCount int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid communication ID")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       entities.CommunicationStatusSent,
			"sent_count":   sentCount,
			"failed_count": failedCount,
			"sent_at":      now,
			"updated_at":   now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("communication not found")
	}
	return nil
}

type communicationTemplateRepository struct {
	collection *mongo.Collection
}

func NewCommunicationTemplateRepository(db *mongo.Database) interfaces.CommunicationTemplateRepository {
	return &communicationTemplateRepository{
		collection: db.Collection("communication_templates"),
	}
}

func (r *communicationTemplateRepository) Create(ctx context.Context, template *entities.CommunicationTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, template)
	return err
}

func (r *communicationTemplateRepository) GetByID(ctx context.Context, id string) (*entities.CommunicationTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid template ID")
	}

	var template entities.CommunicationTemplate
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("template not found")
		}
		return nil, err
	}
	return &template, nil
}

func (r *communicationTemplateRepository) Update(ctx context.Context, id string, template *entities.CommunicationTemplate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid template ID")
	}

	template.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": template})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("template not found")
	}
	return nil
}

func (r *communicationTemplateRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid template ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("template not found")
	}
	return nil
}

func (r *communicationTemplateRepository) List(ctx context.Context, commType entities.CommunicationType) ([]*entities.CommunicationTemplate, error) {
	filter := bson.M{}
	if commType != "" {
		filter["type"] = commType
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*entities.CommunicationTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}
