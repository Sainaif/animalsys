package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnimalRepository defines the interface for animal data access
type AnimalRepository interface {
	// Create creates a new animal
	Create(ctx context.Context, animal *entities.Animal) error

	// FindByID finds an animal by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Animal, error)

	// Update updates an existing animal
	Update(ctx context.Context, animal *entities.Animal) error

	// Delete deletes an animal by ID (soft delete by marking as deleted)
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of animals with pagination and filters
	List(ctx context.Context, filter AnimalFilter) ([]*entities.Animal, int64, error)

	// AddDailyNote adds a daily note to an animal
	AddDailyNote(ctx context.Context, animalID primitive.ObjectID, note entities.DailyNote) error

	// UpdateImages updates the images for an animal
	UpdateImages(ctx context.Context, animalID primitive.ObjectID, images entities.AnimalImages) error

	// UpdateStatus updates the status of an animal
	UpdateStatus(ctx context.Context, animalID primitive.ObjectID, status entities.AnimalStatus) error

	// GetStatistics returns statistics about animals
	GetStatistics(ctx context.Context) (*AnimalStatistics, error)

	// EnsureIndexes creates necessary indexes for the animals collection
	EnsureIndexes(ctx context.Context) error
}

// AnimalFilter defines filter criteria for listing animals
type AnimalFilter struct {
	Category         string   // Filter by category
	Species          string   // Filter by species
	Status           string   // Filter by status
	Sex              string   // Filter by sex
	Size             string   // Filter by size
	AvailableOnly    bool     // Show only available animals
	GoodWithKids     *bool    // Filter by good_with_kids
	GoodWithDogs     *bool    // Filter by good_with_dogs
	GoodWithCats     *bool    // Filter by good_with_cats
	Search           string   // Search in name and description
	AssignedCaretaker *primitive.ObjectID // Filter by assigned caretaker
	MinAge           *float64 // Minimum age in years
	MaxAge           *float64 // Maximum age in years
	Limit            int64    // Limit results
	Offset           int64    // Offset for pagination
	SortBy           string   // Field to sort by (e.g., "created_at", "name")
	SortOrder        string   // Sort order ("asc" or "desc")
}

// AnimalStatistics holds statistical information about animals
type AnimalStatistics struct {
	TotalAnimals       int64                      `json:"total_animals"`
	ByStatus           map[string]int64           `json:"by_status"`
	ByCategory         map[string]int64           `json:"by_category"`
	BySpecies          map[string]int64           `json:"by_species"`
	AvailableForAdoption int64                    `json:"available_for_adoption"`
	AdoptedThisMonth   int64                      `json:"adopted_this_month"`
	AdoptedThisYear    int64                      `json:"adopted_this_year"`
}
