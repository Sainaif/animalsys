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

// userRepository implements the UserRepository interface
type userRepository struct {
	db *mongodb.Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongodb.Database) repositories.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Users)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrEmailAlreadyExists
		}
		return errors.Wrap(err, 500, "failed to create user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	collection := r.db.Collection(mongodb.Collections.Users)

	var user entities.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find user")
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	collection := r.db.Collection(mongodb.Collections.Users)

	var user entities.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, 500, "failed to find user")
	}

	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()

	collection := r.db.Collection(mongodb.Collections.Users)
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrEmailAlreadyExists
		}
		return errors.Wrap(err, 500, "failed to update user")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Users)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.Wrap(err, 500, "failed to delete user")
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List returns a list of users with pagination
func (r *userRepository) List(ctx context.Context, filter repositories.UserFilter) ([]*entities.User, int64, error) {
	collection := r.db.Collection(mongodb.Collections.Users)

	// Build MongoDB filter
	mongoFilter := bson.M{}

	if filter.Role != "" {
		mongoFilter["role"] = filter.Role
	}

	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"email": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"first_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"last_name": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to count users")
	}

	// Find documents with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(filter.Offset).
		SetLimit(filter.Limit)

	cursor, err := collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to list users")
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, errors.Wrap(err, 500, "failed to decode users")
	}

	return users, total, nil
}

// UpdateRefreshToken updates the user's refresh token
func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID primitive.ObjectID, token string) error {
	collection := r.db.Collection(mongodb.Collections.Users)

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"refresh_token": token,
			"updated_at":    time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update refresh token")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	collection := r.db.Collection(mongodb.Collections.Users)

	now := time.Now()
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"last_login": now,
			"updated_at": now,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, 500, "failed to update last login")
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	collection := r.db.Collection(mongodb.Collections.Users)

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, errors.Wrap(err, 500, "failed to check user existence")
	}

	return count > 0, nil
}

// EnsureIndexes creates necessary indexes for the users collection
func (r *userRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection(mongodb.Collections.Users)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "role", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.Wrap(err, 500, "failed to create indexes")
	}

	return nil
}
