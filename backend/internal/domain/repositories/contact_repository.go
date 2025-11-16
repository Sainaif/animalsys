package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContactFilter defines query params
type ContactFilter struct {
	Type    string
	Status  string
	OwnerID string
	Search  string
	Limit   int64
	Offset  int64
}

// ContactRepository describes persistence methods
type ContactRepository interface {
	Create(ctx context.Context, contact *entities.Contact) error
	Update(ctx context.Context, contact *entities.Contact) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Contact, error)
	List(ctx context.Context, filter ContactFilter) ([]*entities.Contact, int64, error)
	EnsureIndexes(ctx context.Context) error
}
