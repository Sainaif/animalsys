package campaign

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ICampaignUseCase defines the interface for campaign business logic
type ICampaignUseCase interface {
	CreateCampaign(ctx context.Context, campaign *entities.Campaign, userID primitive.ObjectID) error
	GetCampaign(ctx context.Context, id primitive.ObjectID) (*entities.Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *entities.Campaign, userID primitive.ObjectID) error
	DeleteCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	ListCampaigns(ctx context.Context, filter *repositories.CampaignFilter) ([]*entities.Campaign, int64, error)
	GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error)
	GetFeaturedCampaigns(ctx context.Context) ([]*entities.Campaign, error)
	GetPublicCampaigns(ctx context.Context) ([]*entities.Campaign, error)
	GetCampaignsByManager(ctx context.Context, managerID primitive.ObjectID) ([]*entities.Campaign, error)
	ActivateCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	PauseCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	CompleteCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	CancelCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	GetCampaignStatistics(ctx context.Context) (*repositories.CampaignStatistics, error)
	GetCampaignDonors(ctx context.Context, id primitive.ObjectID, limit, offset int64) ([]*CampaignDonor, int64, error)
	UpdateCampaignAmount(ctx context.Context, id primitive.ObjectID, amount float64, userID primitive.ObjectID) error
	GetCampaignProgress(ctx context.Context, id primitive.ObjectID) (*CampaignProgress, error)
	ShareCampaign(ctx context.Context, id primitive.ObjectID) (*CampaignShareable, error)
}
