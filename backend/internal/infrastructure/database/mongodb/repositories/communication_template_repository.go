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

type communicationTemplateRepository struct {
	db *mongodb.Database
}

// NewCommunicationTemplateRepository creates a new communication template repository
func NewCommunicationTemplateRepository(db *mongodb.Database) repositories.CommunicationTemplateRepository {
	return &communicationTemplateRepository{db: db}
}

func (r *communicationTemplateRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Templates)
}

// EnsureIndexes creates necessary indexes for communication_templates collection
func (r *communicationTemplateRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "active", Value: 1}}},
		{Keys: bson.D{{Key: "is_default", Value: 1}}},
		{Keys: bson.D{{Key: "language", Value: 1}}},
		{Keys: bson.D{
			{Key: "type", Value: 1},
			{Key: "category", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "type", Value: 1},
			{Key: "active", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "category", Value: 1},
			{Key: "is_default", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "category", Value: 1},
			{Key: "type", Value: 1},
			{Key: "is_default", Value: 1},
		}},
		{Keys: bson.D{{Key: "usage_count", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *communicationTemplateRepository) Create(ctx context.Context, template *entities.CommunicationTemplate) error {
	template.ID = primitive.NewObjectID()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, template)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewConflict("Template with this name already exists")
		}
		return errors.Wrap(err, 500, "Failed to create communication template")
	}
	return nil
}

func (r *communicationTemplateRepository) Update(ctx context.Context, template *entities.CommunicationTemplate) error {
	template.UpdatedAt = time.Now()

	filter := bson.M{"_id": template.ID}
	update := bson.M{"$set": template}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewConflict("Template with this name already exists")
		}
		return errors.Wrap(err, 500, "Failed to update communication template")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete communication template")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *communicationTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.CommunicationTemplate, error) {
	var template entities.CommunicationTemplate
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find communication template")
	}

	return &template, nil
}

func (r *communicationTemplateRepository) List(ctx context.Context, filter *repositories.CommunicationTemplateFilter) ([]*entities.CommunicationTemplate, int64, error) {
	query := bson.M{}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.Active != nil {
		query["active"] = *filter.Active
	}

	if filter.IsDefault != nil {
		query["is_default"] = *filter.IsDefault
	}

	if filter.Language != "" {
		query["language"] = filter.Language
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"subject": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count communication templates")
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
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to list communication templates")
	}
	defer cursor.Close(ctx)

	var templates []*entities.CommunicationTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode communication templates")
	}

	return templates, total, nil
}

func (r *communicationTemplateRepository) GetByCategory(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) ([]*entities.CommunicationTemplate, error) {
	query := bson.M{
		"category": category,
		"type":     templateType,
		"active":   true,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get templates by category")
	}
	defer cursor.Close(ctx)

	var templates []*entities.CommunicationTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communication templates")
	}

	return templates, nil
}

func (r *communicationTemplateRepository) GetDefault(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) (*entities.CommunicationTemplate, error) {
	var template entities.CommunicationTemplate
	query := bson.M{
		"category":   category,
		"type":       templateType,
		"is_default": true,
		"active":     true,
	}

	err := r.collection().FindOne(ctx, query).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find default template")
	}

	return &template, nil
}

func (r *communicationTemplateRepository) GetActiveTemplates(ctx context.Context) ([]*entities.CommunicationTemplate, error) {
	query := bson.M{"active": true}
	findOptions := options.Find().SetSort(bson.D{{Key: "category", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get active templates")
	}
	defer cursor.Close(ctx)

	var templates []*entities.CommunicationTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode communication templates")
	}

	return templates, nil
}

func (r *communicationTemplateRepository) IncrementUsage(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"usage_count": 1},
		"$set": bson.M{
			"last_used_at": time.Now(),
			"updated_at":   time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to increment template usage")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}
