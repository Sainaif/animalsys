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

type notificationRepository struct {
	db *mongodb.Database
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *mongodb.Database) repositories.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Notifications)
}

// EnsureIndexes creates necessary indexes for notifications collection
func (r *notificationRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "priority", Value: 1}}},
		{Keys: bson.D{{Key: "read", Value: 1}}},
		{Keys: bson.D{{Key: "dismissed", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "related_type", Value: 1}}},
		{Keys: bson.D{{Key: "related_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "read", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "dismissed", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "group_key", Value: 1},
		}},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *notificationRepository) Create(ctx context.Context, notification *entities.Notification) error {
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, notification)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to create notification")
	}
	return nil
}

func (r *notificationRepository) Update(ctx context.Context, notification *entities.Notification) error {
	notification.UpdatedAt = time.Now()

	filter := bson.M{"_id": notification.ID}
	update := bson.M{"$set": notification}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to update notification")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to delete notification")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Notification, error) {
	var notification entities.Notification
	filter := bson.M{"_id": id}

	err := r.collection().FindOne(ctx, filter).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find notification")
	}

	return &notification, nil
}

func (r *notificationRepository) preferencesCollection() *mongo.Collection {
	return r.db.DB.Collection("notification_preferences")
}

func (r *notificationRepository) FindPreferencesByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.NotificationPreferences, error) {
	var preferences entities.NotificationPreferences
	err := r.preferencesCollection().FindOne(ctx, bson.M{"user_id": userID}).Decode(&preferences)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil, not an error, if not found
		}
		return nil, err
	}
	return &preferences, nil
}

func (r *notificationRepository) UpsertPreferences(ctx context.Context, preferences *entities.NotificationPreferences) error {
	preferences.UpdatedAt = time.Now()

	filter := bson.M{"user_id": preferences.UserID}
	update := bson.M{"$set": preferences}
	opts := options.Update().SetUpsert(true)

	_, err := r.preferencesCollection().UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *notificationRepository) List(ctx context.Context, filter *repositories.NotificationFilter) ([]*entities.Notification, int64, error) {
	query := bson.M{}

	if filter.UserID != nil {
		query["user_id"] = *filter.UserID
	}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.Priority != "" {
		query["priority"] = filter.Priority
	}

	if filter.Read != nil {
		query["read"] = *filter.Read
	}

	if filter.Dismissed != nil {
		query["dismissed"] = *filter.Dismissed
	}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.RelatedType != "" {
		query["related_type"] = filter.RelatedType
	}

	if filter.RelatedID != nil {
		query["related_id"] = *filter.RelatedID
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to count notifications")
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
		return nil, 0, errors.Wrap(err, 500, "Failed to list notifications")
	}
	defer cursor.Close(ctx)

	var notifications []*entities.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, 0, errors.Wrap(err, 500, "Failed to decode notifications")
	}

	return notifications, total, nil
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	query := bson.M{
		"user_id":   userID,
		"dismissed": false,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get notifications by user")
	}
	defer cursor.Close(ctx)

	var notifications []*entities.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode notifications")
	}

	return notifications, nil
}

func (r *notificationRepository) GetUnreadByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	query := bson.M{
		"user_id":   userID,
		"read":      false,
		"dismissed": false,
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection().Find(ctx, query, findOptions)
	if err != nil {
		return nil, errors.Wrap(err, 500, "Failed to get unread notifications")
	}
	defer cursor.Close(ctx)

	var notifications []*entities.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, errors.Wrap(err, 500, "Failed to decode notifications")
	}

	return notifications, nil
}

func (r *notificationRepository) CountUnreadByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	query := bson.M{
		"user_id":   userID,
		"read":      false,
		"dismissed": false,
	}

	count, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return 0, errors.Wrap(err, 500, "Failed to count unread notifications")
	}

	return count, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"read":       true,
			"read_at":    now,
			"updated_at": now,
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to mark notification as read")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) MarkAsUnread(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"read":       false,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{
			"read_at": "",
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to mark notification as unread")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{
		"user_id": userID,
		"read":    false,
	}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"read":       true,
			"read_at":    now,
			"updated_at": now,
		},
	}

	_, err := r.collection().UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to mark all notifications as read")
	}

	return nil
}

func (r *notificationRepository) Dismiss(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"dismissed":  true,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "Failed to dismiss notification")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *notificationRepository) DeleteExpired(ctx context.Context) (int64, error) {
	filter := bson.M{
		"expires_at": bson.M{
			"$ne":  nil,
			"$lte": time.Now(),
		},
	}

	result, err := r.collection().DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, 500, "Failed to delete expired notifications")
	}

	return result.DeletedCount, nil
}

func (r *notificationRepository) FindByGroupKey(ctx context.Context, userID primitive.ObjectID, groupKey string) (*entities.Notification, error) {
	var notification entities.Notification
	filter := bson.M{
		"user_id":   userID,
		"group_key": groupKey,
	}
	findOptions := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := r.collection().FindOne(ctx, filter, findOptions).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "Failed to find notification by group key")
	}

	return &notification, nil
}
