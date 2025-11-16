package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	domainRepo "github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type contactRepository struct {
	db *mongodb.Database
}

// NewContactRepository creates Mongo implementation
func NewContactRepository(db *mongodb.Database) domainRepo.ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Contacts)
}

func (r *contactRepository) Create(ctx context.Context, contact *entities.Contact) error {
	if contact.ID.IsZero() {
		contact.ID = primitive.NewObjectID()
	}
	now := time.Now()
	contact.CreatedAt = now
	contact.UpdatedAt = now
	_, err := r.collection().InsertOne(ctx, contact)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}
	return nil
}

func (r *contactRepository) Update(ctx context.Context, contact *entities.Contact) error {
	contact.UpdatedAt = time.Now()
	filter := bson.M{"_id": contact.ID}
	update := bson.M{"$set": contact}
	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.NewNotFound("contact not found")
	}
	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	if result.DeletedCount == 0 {
		return errors.NewNotFound("contact not found")
	}
	return nil
}

func (r *contactRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Contact, error) {
	var contact entities.Contact
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&contact)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("contact not found")
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}
	return &contact, nil
}

func (r *contactRepository) List(ctx context.Context, filter domainRepo.ContactFilter) ([]*entities.Contact, int64, error) {
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.OwnerID != "" {
		if ownerID, err := primitive.ObjectIDFromHex(filter.OwnerID); err == nil {
			query["owner_id"] = ownerID
		}
	}
	if filter.Search != "" {
		regex := primitive.Regex{Pattern: filter.Search, Options: "i"}
		query["$or"] = []bson.M{
			{"first_name": regex},
			{"last_name": regex},
			{"email": regex},
			{"organization": regex},
		}
	}

	opts := options.Find().
		SetSkip(filter.Offset).
		SetLimit(filter.Limit).
		SetSort(bson.M{"updated_at": -1})

	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list contacts: %w", err)
	}
	defer cursor.Close(ctx)

	var contacts []*entities.Contact
	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, 0, fmt.Errorf("failed to decode contacts: %w", err)
	}

	count, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count contacts: %w", err)
	}

	return contacts, count, nil
}

// EnsureIndexes creates indexes
func (r *contactRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "owner_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}
	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
