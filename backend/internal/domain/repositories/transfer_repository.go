package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransferRepository defines the interface for transfer data access
type TransferRepository interface {
	Create(ctx context.Context, transfer *entities.Transfer) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Transfer, error)
	Update(ctx context.Context, transfer *entities.Transfer) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter *TransferFilter) ([]*entities.Transfer, int64, error)
	GetByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Transfer, error)
	GetByPartner(ctx context.Context, partnerID primitive.ObjectID) ([]*entities.Transfer, error)
	GetByStatus(ctx context.Context, status entities.TransferStatus) ([]*entities.Transfer, error)
	GetPendingTransfers(ctx context.Context) ([]*entities.Transfer, error)
	GetUpcomingTransfers(ctx context.Context, days int) ([]*entities.Transfer, error)
	GetOverdueTransfers(ctx context.Context) ([]*entities.Transfer, error)
	GetRequiringFollowUp(ctx context.Context) ([]*entities.Transfer, error)
	GetTransferStatistics(ctx context.Context) (*TransferStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// TransferFilter defines filter criteria for listing transfers
type TransferFilter struct {
	Direction       string
	Status          string
	Reason          string
	AnimalID        *primitive.ObjectID
	PartnerID       *primitive.ObjectID
	RequestedBy     *primitive.ObjectID
	ApprovedBy      *primitive.ObjectID
	ScheduledAfter  *time.Time
	ScheduledBefore *time.Time
	Search          string
	Limit           int64
	Offset          int64
	SortBy          string
	SortOrder       string
}

// TransferStatistics represents transfer statistics
type TransferStatistics struct {
	TotalTransfers     int64            `json:"total_transfers"`
	ByDirection        map[string]int64 `json:"by_direction"`
	ByStatus           map[string]int64 `json:"by_status"`
	ByReason           map[string]int64 `json:"by_reason"`
	PendingTransfers   int64            `json:"pending_transfers"`
	InTransitTransfers int64            `json:"in_transit_transfers"`
	CompletedThisMonth int64            `json:"completed_this_month"`
	AverageDuration    float64          `json:"average_duration"` // in hours
}
