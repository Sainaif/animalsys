package campaign

import (
	"context"
	"strings"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CampaignUseCase handles campaign business logic
type CampaignUseCase struct {
	campaignRepo repositories.CampaignRepository
	auditLogRepo repositories.AuditLogRepository
}

// NewCampaignUseCase creates a new campaign use case
func NewCampaignUseCase(
	campaignRepo repositories.CampaignRepository,
	auditLogRepo repositories.AuditLogRepository,
) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo: campaignRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateCampaign creates a new campaign
func (uc *CampaignUseCase) CreateCampaign(ctx context.Context, campaign *entities.Campaign, userID primitive.ObjectID) error {
	// Validate campaign
	if err := uc.validateCampaign(campaign); err != nil {
		return err
	}

	campaign.CreatedBy = userID
	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Create(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionCreate, "campaign", "", "").
		WithEntityID(campaign.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetCampaign gets a campaign by ID
func (uc *CampaignUseCase) GetCampaign(ctx context.Context, id primitive.ObjectID) (*entities.Campaign, error) {
	return uc.campaignRepo.FindByID(ctx, id)
}

// UpdateCampaign updates a campaign
func (uc *CampaignUseCase) UpdateCampaign(ctx context.Context, campaign *entities.Campaign, userID primitive.ObjectID) error {
	// Validate campaign
	if err := uc.validateCampaign(campaign); err != nil {
		return err
	}

	// Check if campaign exists
	_, err := uc.campaignRepo.FindByID(ctx, campaign.ID)
	if err != nil {
		return err
	}

	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "", "").
		WithEntityID(campaign.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// DeleteCampaign deletes a campaign
func (uc *CampaignUseCase) DeleteCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if campaign exists
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if campaign has donations (business rule: can't delete campaigns with donations)
	if campaign.DonationCount > 0 {
		return errors.NewBadRequest("Cannot delete campaign with donations. Consider marking as cancelled instead.")
	}

	if err := uc.campaignRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionDelete, "campaign", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListCampaigns lists campaigns with filters
func (uc *CampaignUseCase) ListCampaigns(ctx context.Context, filter *repositories.CampaignFilter) ([]*entities.Campaign, int64, error) {
	return uc.campaignRepo.List(ctx, filter)
}

// GetActiveCampaigns gets all active campaigns
func (uc *CampaignUseCase) GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	return uc.campaignRepo.GetActiveCampaigns(ctx)
}

// GetFeaturedCampaigns gets all featured campaigns
func (uc *CampaignUseCase) GetFeaturedCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	return uc.campaignRepo.GetFeaturedCampaigns(ctx)
}

// GetPublicCampaigns gets all public campaigns
func (uc *CampaignUseCase) GetPublicCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	return uc.campaignRepo.GetPublicCampaigns(ctx)
}

// GetCampaignsByManager gets campaigns managed by a specific user
func (uc *CampaignUseCase) GetCampaignsByManager(ctx context.Context, managerID primitive.ObjectID) ([]*entities.Campaign, error) {
	return uc.campaignRepo.GetCampaignsByManager(ctx, managerID)
}

// ActivateCampaign activates a campaign
func (uc *CampaignUseCase) ActivateCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if campaign.Status == entities.CampaignStatusActive {
		return errors.NewBadRequest("Campaign is already active")
	}

	if campaign.Status == entities.CampaignStatusCompleted {
		return errors.NewBadRequest("Cannot activate a completed campaign")
	}

	campaign.Status = entities.CampaignStatusActive
	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// PauseCampaign pauses a campaign
func (uc *CampaignUseCase) PauseCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if campaign.Status != entities.CampaignStatusActive {
		return errors.NewBadRequest("Only active campaigns can be paused")
	}

	campaign.Status = entities.CampaignStatusPaused
	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// CompleteCampaign completes a campaign
func (uc *CampaignUseCase) CompleteCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	campaign.Status = entities.CampaignStatusCompleted
	now := time.Now()
	campaign.EndDate = &now
	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// CancelCampaign cancels a campaign
func (uc *CampaignUseCase) CancelCampaign(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if campaign.Status == entities.CampaignStatusCompleted {
		return errors.NewBadRequest("Cannot cancel a completed campaign")
	}

	campaign.Status = entities.CampaignStatusCancelled
	campaign.UpdatedBy = userID

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetCampaignStatistics gets campaign statistics
func (uc *CampaignUseCase) GetCampaignStatistics(ctx context.Context) (*repositories.CampaignStatistics, error) {
	return uc.campaignRepo.GetCampaignStatistics(ctx)
}

// validateCampaign validates campaign data
func (uc *CampaignUseCase) validateCampaign(campaign *entities.Campaign) error {
	if strings.TrimSpace(campaign.Name.English) == "" {
		return errors.NewBadRequest("Campaign name (English) is required")
	}

	if campaign.Type == "" {
		return errors.NewBadRequest("Campaign type is required")
	}

	if campaign.GoalAmount <= 0 {
		return errors.NewBadRequest("Goal amount must be greater than zero")
	}

	if campaign.Manager.IsZero() {
		return errors.NewBadRequest("Campaign manager is required")
	}

	// Validate dates
	if campaign.EndDate != nil {
		if campaign.EndDate.Before(campaign.StartDate) {
			return errors.NewBadRequest("End date cannot be before start date")
		}
	}

	return nil
}
