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

// CampaignDonor represents a donor with their total contribution to a campaign
type CampaignDonor struct {
	Donor       *entities.Donor `json:"donor"`
	TotalAmount float64         `json:"total_amount"`
}

// CampaignProgress represents the progress of a fundraising campaign
type CampaignProgress struct {
	GoalAmount    float64 `json:"goal_amount"`
	CurrentAmount float64 `json:"current_amount"`
	Percentage    float64 `json:"percentage"`
	DonationCount int     `json:"donation_count"`
}

// CampaignShareable represents a shareable campaign payload
type CampaignShareable struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// CampaignUseCase handles campaign business logic
type CampaignUseCase struct {
	campaignRepo repositories.CampaignRepository
	donationRepo repositories.DonationRepository
	donorRepo    repositories.DonorRepository
	auditLogRepo repositories.AuditLogRepository
}

// NewCampaignUseCase creates a new campaign use case
func NewCampaignUseCase(
	campaignRepo repositories.CampaignRepository,
	donationRepo repositories.DonationRepository,
	donorRepo repositories.DonorRepository,
	auditLogRepo repositories.AuditLogRepository,
) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo: campaignRepo,
		donationRepo: donationRepo,
		donorRepo:    donorRepo,
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

// GetCampaignDonors gets all donors for a campaign
func (uc *CampaignUseCase) GetCampaignDonors(ctx context.Context, id primitive.ObjectID, limit, offset int64) ([]*CampaignDonor, int64, error) {
	// 1. Aggregate donations by donor at the database level
	aggregationResults, totalDonors, err := uc.donationRepo.AggregateDonorsByCampaignID(ctx, id, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if len(aggregationResults) == 0 {
		return []*CampaignDonor{}, 0, nil
	}

	// 2. Get unique donor IDs from aggregation results
	donorIDs := make([]primitive.ObjectID, len(aggregationResults))
	donorTotals := make(map[primitive.ObjectID]float64, len(aggregationResults))
	for i, result := range aggregationResults {
		donorIDs[i] = result.DonorID
		donorTotals[result.DonorID] = result.TotalAmount
	}

	// 3. Fetch all donor details in a single query
	donors, err := uc.donorRepo.FindManyByIDs(ctx, donorIDs)
	if err != nil {
		return nil, 0, err
	}

	// 4. Map donor details to the campaign donor struct
	donorMap := make(map[primitive.ObjectID]*entities.Donor, len(donors))
	for _, donor := range donors {
		donorMap[donor.ID] = donor
	}

	campaignDonors := make([]*CampaignDonor, len(aggregationResults))
	for i, result := range aggregationResults {
		campaignDonors[i] = &CampaignDonor{
			Donor:       donorMap[result.DonorID],
			TotalAmount: donorTotals[result.DonorID],
		}
	}

	return campaignDonors, totalDonors, nil
}

// UpdateCampaignAmount updates campaign raised amount
func (uc *CampaignUseCase) UpdateCampaignAmount(ctx context.Context, id primitive.ObjectID, amount float64, userID primitive.ObjectID) error {
	if amount < 0 {
		return errors.NewBadRequest("Amount cannot be negative")
	}

	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	campaign.CurrentAmount = amount
	campaign.UpdatedBy = userID
	campaign.UpdatedAt = time.Now()

	if err := uc.campaignRepo.Update(ctx, campaign); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "campaign", "Updated current amount", "").
		WithEntityID(id)
	return uc.auditLogRepo.Create(ctx, auditLog)
}

// GetCampaignProgress gets campaign progress
func (uc *CampaignUseCase) GetCampaignProgress(ctx context.Context, id primitive.ObjectID) (*CampaignProgress, error) {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	progress := &CampaignProgress{
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		DonationCount: campaign.DonationCount,
		Percentage:    campaign.GetProgressPercentage(),
	}

	return progress, nil
}

// ShareCampaign shares a campaign
func (uc *CampaignUseCase) ShareCampaign(ctx context.Context, id primitive.ObjectID) (*CampaignShareable, error) {
	campaign, err := uc.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Get base URL from config
	shareable := &CampaignShareable{
		Title:       campaign.Name.English,
		Description: campaign.Description.English,
		URL:         "/campaigns/" + id.Hex(),
	}

	return shareable, nil
}
