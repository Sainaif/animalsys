package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerRepository defines the interface for partner data access
type PartnerRepository interface {
	Create(ctx context.Context, partner *entities.Partner) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error)
	Update(ctx context.Context, partner *entities.Partner) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter *PartnerFilter) ([]*entities.Partner, int64, error)
	GetByType(ctx context.Context, partnerType entities.PartnerType) ([]*entities.Partner, error)
	GetByStatus(ctx context.Context, status entities.PartnerStatus) ([]*entities.Partner, error)
	GetActivePartners(ctx context.Context) ([]*entities.Partner, error)
	GetPartnersWithCapacity(ctx context.Context) ([]*entities.Partner, error)
	GetPartnerStatistics(ctx context.Context) (*PartnerStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// PartnerFilter defines filter criteria for listing partners
type PartnerFilter struct {
	Type            string
	Status          string
	AcceptsIntakes  *bool
	HasCapacity     *bool
	Search          string
	Tags            []string
	Limit           int64
	Offset          int64
	SortBy          string
	SortOrder       string
}

// PartnerStatistics represents partner statistics
type PartnerStatistics struct {
	TotalPartners       int64            `json:"total_partners"`
	ByType              map[string]int64 `json:"by_type"`
	ByStatus            map[string]int64 `json:"by_status"`
	ActivePartners      int64            `json:"active_partners"`
	AcceptingIntakes    int64            `json:"accepting_intakes"`
	TotalTransfers      int64            `json:"total_transfers"`
	AverageRating       float64          `json:"average_rating"`
}
