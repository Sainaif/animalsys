package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of users with pagination
	List(ctx context.Context, filter UserFilter) ([]*entities.User, int64, error)

	// UpdateRefreshToken updates the user's refresh token
	UpdateRefreshToken(ctx context.Context, userID primitive.ObjectID, token string) error

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// EnsureIndexes creates necessary indexes for the users collection
	EnsureIndexes(ctx context.Context) error
}

// UserFilter defines filter criteria for listing users
type UserFilter struct {
	Role   string
	Status string
	Search string // search in email, first name, last name
	Limit  int64
	Offset int64
}
