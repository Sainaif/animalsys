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

type documentRepository struct {
	collection *mongo.Collection
}

func NewDocumentRepository(db *mongo.Database) interfaces.DocumentRepository {
	return &documentRepository{
		collection: db.Collection("documents"),
	}
}

func (r *documentRepository) Create(ctx context.Context, document *entities.Document) error {
	document.UploadedAt = time.Now()
	document.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, document)
	return err
}

func (r *documentRepository) GetByID(ctx context.Context, id string) (*entities.Document, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid document ID")
	}

	var document entities.Document
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) Update(ctx context.Context, id string, document *entities.Document) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid document ID")
	}

	document.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": document})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("document not found")
	}
	return nil
}

func (r *documentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid document ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("document not found")
	}
	return nil
}

func (r *documentRepository) List(ctx context.Context, filter *entities.DocumentFilter) ([]*entities.Document, int64, error) {
	mongoFilter := bson.M{}

	if filter.FileType != "" {
		mongoFilter["file_type"] = filter.FileType
	}
	if filter.RelatedEntityType != "" {
		mongoFilter["related_entity_type"] = filter.RelatedEntityType
	}
	if filter.RelatedEntityID != "" {
		mongoFilter["related_entity_id"] = filter.RelatedEntityID
	}
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}
	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"file_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"original_file_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		findOptions.SetSkip(int64(filter.Offset))
	}
	sortField := "uploaded_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	findOptions.SetSort(bson.M{sortField: sortOrder})

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, 0, err
	}
	return documents, total, nil
}

func (r *documentRepository) GetByRelatedEntity(ctx context.Context, entityType, entityID string) ([]*entities.Document, error) {
	filter := bson.M{
		"related_entity_type": entityType,
		"related_entity_id":   entityID,
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"uploaded_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	return documents, nil
}

func (r *documentRepository) GetExpiringSoon(ctx context.Context, days int) ([]*entities.Document, error) {
	expiryDate := time.Now().AddDate(0, 0, days)

	filter := bson.M{
		"expiry_date": bson.M{
			"$lte": expiryDate,
			"$gte": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"expiry_date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*entities.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	return documents, nil
}
