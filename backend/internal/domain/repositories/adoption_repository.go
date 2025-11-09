package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdoptionApplicationRepository defines the interface for adoption application data access
type AdoptionApplicationRepository interface {
	// Create creates a new adoption application
	Create(ctx context.Context, application *entities.AdoptionApplication) error

	// FindByID finds an adoption application by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AdoptionApplication, error)

	// Update updates an existing adoption application
	Update(ctx context.Context, application *entities.AdoptionApplication) error

	// Delete deletes an adoption application by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of adoption applications with pagination and filters
	List(ctx context.Context, filter AdoptionApplicationFilter) ([]*entities.AdoptionApplication, int64, error)

	// GetByAnimalID returns all applications for a specific animal
	GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.AdoptionApplication, error)

	// GetByApplicantEmail returns applications by applicant email
	GetByApplicantEmail(ctx context.Context, email string) ([]*entities.AdoptionApplication, error)

	// GetPendingApplications returns all pending applications
	GetPendingApplications(ctx context.Context) ([]*entities.AdoptionApplication, error)

	// GetApplicationsByStatus returns applications by status
	GetApplicationsByStatus(ctx context.Context, status entities.ApplicationStatus) ([]*entities.AdoptionApplication, error)

	// EnsureIndexes creates necessary indexes for the applications collection
	EnsureIndexes(ctx context.Context) error
}

// AdoptionApplicationFilter defines filter criteria for listing adoption applications
type AdoptionApplicationFilter struct {
	AnimalID        *primitive.ObjectID
	Status          string
	ApplicantEmail  string
	ApplicantName   string
	FromDate        *time.Time
	ToDate          *time.Time
	ReviewedBy      *primitive.ObjectID
	Limit           int64
	Offset          int64
	SortBy          string // Field to sort by
	SortOrder       string // "asc" or "desc"
}

// AdoptionRepository defines the interface for adoption data access
type AdoptionRepository interface {
	// Create creates a new adoption record
	Create(ctx context.Context, adoption *entities.Adoption) error

	// FindByID finds an adoption by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Adoption, error)

	// Update updates an existing adoption
	Update(ctx context.Context, adoption *entities.Adoption) error

	// Delete deletes an adoption by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of adoptions with pagination and filters
	List(ctx context.Context, filter AdoptionFilter) ([]*entities.Adoption, int64, error)

	// GetByAnimalID returns adoption for a specific animal
	GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) (*entities.Adoption, error)

	// GetByAdopterID returns all adoptions for a specific adopter
	GetByAdopterID(ctx context.Context, adopterID primitive.ObjectID) ([]*entities.Adoption, error)

	// GetByApplicationID returns adoption by application ID
	GetByApplicationID(ctx context.Context, applicationID primitive.ObjectID) (*entities.Adoption, error)

	// GetPendingFollowUps returns adoptions with pending follow-ups
	GetPendingFollowUps(ctx context.Context, days int) ([]*entities.Adoption, error)

	// GetAdoptionStatistics returns adoption statistics
	GetAdoptionStatistics(ctx context.Context) (*AdoptionStatistics, error)

	// EnsureIndexes creates necessary indexes for the adoptions collection
	EnsureIndexes(ctx context.Context) error
}

// AdoptionFilter defines filter criteria for listing adoptions
type AdoptionFilter struct {
	AnimalID      *primitive.ObjectID
	AdopterID     *primitive.ObjectID
	ApplicationID *primitive.ObjectID
	Status        string
	PaymentStatus string
	FromDate      *time.Time
	ToDate        *time.Time
	TrialPeriod   *bool
	ProcessedBy   *primitive.ObjectID
	Limit         int64
	Offset        int64
	SortBy        string // Field to sort by
	SortOrder     string // "asc" or "desc"
}

// AdoptionStatistics represents adoption statistics
type AdoptionStatistics struct {
	TotalAdoptions       int64                       `json:"total_adoptions" bson:"total_adoptions"`
	CompletedAdoptions   int64                       `json:"completed_adoptions" bson:"completed_adoptions"`
	PendingAdoptions     int64                       `json:"pending_adoptions" bson:"pending_adoptions"`
	ReturnedAnimals      int64                       `json:"returned_animals" bson:"returned_animals"`
	TotalAdoptionFees    float64                     `json:"total_adoption_fees" bson:"total_adoption_fees"`
	AverageAdoptionFee   float64                     `json:"average_adoption_fee" bson:"average_adoption_fee"`
	AdoptionsThisMonth   int64                       `json:"adoptions_this_month" bson:"adoptions_this_month"`
	AdoptionsThisYear    int64                       `json:"adoptions_this_year" bson:"adoptions_this_year"`
	ByStatus             map[string]int64            `json:"by_status" bson:"by_status"`
	ByPaymentStatus      map[string]int64            `json:"by_payment_status" bson:"by_payment_status"`
	PendingFollowUps     int64                       `json:"pending_follow_ups" bson:"pending_follow_ups"`
	ReturnRate           float64                     `json:"return_rate" bson:"return_rate"` // percentage
}
