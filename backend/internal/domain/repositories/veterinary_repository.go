package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VeterinaryVisitRepository defines the interface for veterinary visit data access
type VeterinaryVisitRepository interface {
	// Create creates a new veterinary visit
	Create(ctx context.Context, visit *entities.VeterinaryVisit) error

	// FindByID finds a veterinary visit by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.VeterinaryVisit, error)

	// Update updates an existing veterinary visit
	Update(ctx context.Context, visit *entities.VeterinaryVisit) error

	// Delete deletes a veterinary visit by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of veterinary visits with pagination and filters
	List(ctx context.Context, filter VeterinaryVisitFilter) ([]*entities.VeterinaryVisit, int64, error)

	// GetByAnimalID returns all visits for a specific animal
	GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VeterinaryVisit, error)

	// GetUpcomingVisits returns scheduled visits for the future
	GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error)

	// EnsureIndexes creates necessary indexes for the veterinary_visits collection
	EnsureIndexes(ctx context.Context) error
}

// VeterinaryVisitFilter defines filter criteria for listing veterinary visits
type VeterinaryVisitFilter struct {
	AnimalID      *primitive.ObjectID
	VisitType     string
	Status        string
	FromDate      *time.Time
	ToDate        *time.Time
	VeterinarianName string
	Limit         int64
	Offset        int64
	SortBy        string // Field to sort by
	SortOrder     string // "asc" or "desc"
}

// VaccinationRepository defines the interface for vaccination data access
type VaccinationRepository interface {
	// Create creates a new vaccination record
	Create(ctx context.Context, vaccination *entities.Vaccination) error

	// FindByID finds a vaccination by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Vaccination, error)

	// Update updates an existing vaccination
	Update(ctx context.Context, vaccination *entities.Vaccination) error

	// Delete deletes a vaccination by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of vaccinations with pagination and filters
	List(ctx context.Context, filter VaccinationFilter) ([]*entities.Vaccination, int64, error)

	// GetByAnimalID returns all vaccinations for a specific animal
	GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Vaccination, error)

	// GetDueVaccinations returns vaccinations that are due or overdue
	GetDueVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error)

	// GetByVaccineType returns vaccinations of a specific type for an animal
	GetByVaccineType(ctx context.Context, animalID primitive.ObjectID, vaccineType entities.VaccinationType) ([]*entities.Vaccination, error)

	// EnsureIndexes creates necessary indexes for the vaccinations collection
	EnsureIndexes(ctx context.Context) error
}

// VaccinationFilter defines filter criteria for listing vaccinations
type VaccinationFilter struct {
	AnimalID     *primitive.ObjectID
	VaccineType  string
	Status       string // current, due, overdue, expired
	FromDate     *time.Time
	ToDate       *time.Time
	Limit        int64
	Offset       int64
	SortBy       string // Field to sort by
	SortOrder    string // "asc" or "desc"
}
