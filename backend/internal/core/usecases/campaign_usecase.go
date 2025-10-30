package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type CampaignUseCase struct {
	campaignRepo interfaces.CampaignRepository
	auditRepo    interfaces.AuditLogRepository
}

func NewCampaignUseCase(
	campaignRepo interfaces.CampaignRepository,
	auditRepo interfaces.AuditLogRepository,
) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo: campaignRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *CampaignUseCase) Create(ctx context.Context, req *entities.CampaignCreateRequest, createdBy string) (*entities.Campaign, error) {
	campaign := entities.NewCampaign(
		req.Name,
		req.Type,
		req.StartDate,
		req.EndDate,
		req.GoalAmount,
	)

	campaign.Description = req.Description
	campaign.TargetAudience = req.TargetAudience
	campaign.Budget = req.Budget
	campaign.CreatedBy = createdBy

	if err := uc.campaignRepo.Create(ctx, campaign); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "campaign", campaign.ID.Hex(), "Campaign created: "+req.Name)
	uc.auditRepo.Create(ctx, auditLog)

	return campaign, nil
}

func (uc *CampaignUseCase) GetByID(ctx context.Context, id string) (*entities.Campaign, error) {
	return uc.campaignRepo.GetByID(ctx, id)
}

func (uc *CampaignUseCase) List(ctx context.Context, filter *entities.CampaignFilter) ([]*entities.Campaign, int64, error) {
	return uc.campaignRepo.List(ctx, filter)
}

func (uc *CampaignUseCase) Update(ctx context.Context, id string, req *entities.CampaignUpdateRequest, updatedBy string) (*entities.Campaign, error) {
	campaign, err := uc.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		campaign.Name = req.Name
	}
	if req.Description != "" {
		campaign.Description = req.Description
	}
	if req.Type != "" {
		campaign.Type = req.Type
	}
	if !req.StartDate.IsZero() {
		campaign.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		campaign.EndDate = req.EndDate
	}
	if req.GoalAmount > 0 {
		campaign.GoalAmount = req.GoalAmount
	}
	if req.TargetAudience != "" {
		campaign.TargetAudience = req.TargetAudience
	}
	if req.Budget > 0 {
		campaign.Budget = req.Budget
	}
	if req.Status != "" {
		campaign.Status = req.Status
	}
	campaign.UpdatedBy = updatedBy

	if err := uc.campaignRepo.Update(ctx, id, campaign); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "campaign", id, "Campaign updated")
	uc.auditRepo.Create(ctx, auditLog)

	return campaign, nil
}

func (uc *CampaignUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.campaignRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "campaign", id, "Campaign deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *CampaignUseCase) UpdateProgress(ctx context.Context, id string, amountRaised float64, participantsCount int) error {
	campaign, err := uc.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	campaign.UpdateProgress(amountRaised, participantsCount)

	if err := uc.campaignRepo.Update(ctx, id, campaign); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog("system", "", "", entities.ActionUpdate, "campaign", id, "Campaign progress updated")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *CampaignUseCase) GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error) {
	return uc.campaignRepo.GetActiveCampaigns(ctx)
}

func (uc *CampaignUseCase) GetByDateRange(ctx context.Context, startDate, endDate string, limit, offset int) ([]*entities.Campaign, int64, error) {
	return uc.campaignRepo.GetByDateRange(ctx, startDate, endDate, limit, offset)
}

func (uc *CampaignUseCase) GetCampaignStatistics(ctx context.Context, id string) (map[string]interface{}, error) {
	campaign, err := uc.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"campaign_id":         id,
		"campaign_name":       campaign.Name,
		"type":                campaign.Type,
		"status":              campaign.Status,
		"goal_amount":         campaign.GoalAmount,
		"amount_raised":       campaign.AmountRaised,
		"participants_count":  campaign.ParticipantsCount,
		"progress_percentage": campaign.ProgressPercentage,
		"budget":              campaign.Budget,
		"start_date":          campaign.StartDate,
		"end_date":            campaign.EndDate,
	}

	return stats, nil
}

func (uc *CampaignUseCase) GetAllCampaignsStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	campaigns, _, err := uc.campaignRepo.GetByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return nil, err
	}

	totalGoal := 0.0
	totalRaised := 0.0
	totalParticipants := 0
	byType := make(map[string]int)
	byStatus := make(map[string]int)

	for _, campaign := range campaigns {
		totalGoal += campaign.GoalAmount
		totalRaised += campaign.AmountRaised
		totalParticipants += campaign.ParticipantsCount
		byType[string(campaign.Type)]++
		byStatus[string(campaign.Status)]++
	}

	overallProgress := 0.0
	if totalGoal > 0 {
		overallProgress = (totalRaised / totalGoal) * 100
	}

	stats := map[string]interface{}{
		"period": map[string]string{
			"start": startDate,
			"end":   endDate,
		},
		"total_campaigns":      len(campaigns),
		"total_goal":           totalGoal,
		"total_raised":         totalRaised,
		"overall_progress":     overallProgress,
		"total_participants":   totalParticipants,
		"by_type":              byType,
		"by_status":            byStatus,
	}

	return stats, nil
}
