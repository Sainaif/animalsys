package donor

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/delivery/http/dto"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorUseCaseInterface defines the interface for the donor use case
type DonorUseCaseInterface interface {
	GetTopDonors(ctx context.Context, limit int) ([]*dto.TopDonorResponse, error)
	GetRecurringDonors(ctx context.Context) ([]*dto.RecurringDonorResponse, error)
	UpdateCommunicationPreferences(ctx context.Context, donorID primitive.ObjectID, prefs entities.DonorPreferences, userID primitive.ObjectID) error
	CreateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error
	GetDonor(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error)
	UpdateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error
	DeleteDonor(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	ListDonors(ctx context.Context, filter *repositories.DonorFilter) ([]*entities.Donor, int64, error)
	GetMajorDonors(ctx context.Context) ([]*entities.Donor, error)
	GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error)
	GetDonorStatistics(ctx context.Context) (*repositories.DonorStatistics, error)
	UpdateDonorEngagement(ctx context.Context, id primitive.ObjectID, volunteerHours int, eventsAttended int, userID primitive.ObjectID) error
}
