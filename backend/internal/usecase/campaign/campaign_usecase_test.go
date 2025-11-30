package campaign

import (
	"context"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCampaignUseCase_GetCampaignDonors(t *testing.T) {
	mockCampaignRepo := new(mocks.CampaignRepository)
	mockDonationRepo := new(mocks.DonationRepository)
	mockDonorRepo := new(mocks.DonorRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewCampaignUseCase(mockCampaignRepo, mockDonationRepo, mockDonorRepo, mockAuditLogRepo)

	campaignID := primitive.NewObjectID()
	donor1ID := primitive.NewObjectID()
	donor2ID := primitive.NewObjectID()

	aggregationResults := []*repositories.DonorAggregationResult{
		{DonorID: donor2ID, TotalAmount: 200.0},
		{DonorID: donor1ID, TotalAmount: 150.0},
	}
	donorsList := []*entities.Donor{
		{ID: donor1ID},
		{ID: donor2ID},
	}
	donorIDs := []primitive.ObjectID{donor2ID, donor1ID}

	mockDonationRepo.On("AggregateDonorsByCampaignID", mock.Anything, campaignID, int64(20), int64(0)).Return(aggregationResults, int64(2), nil)
	mockDonorRepo.On("FindManyByIDs", mock.Anything, donorIDs).Return(donorsList, nil)

	donors, total, err := uc.GetCampaignDonors(context.Background(), campaignID, 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, donors, 2)
	assert.Equal(t, donor2ID, donors[0].Donor.ID)
	assert.Equal(t, 200.0, donors[0].TotalAmount)
	assert.Equal(t, donor1ID, donors[1].Donor.ID)
	assert.Equal(t, 150.0, donors[1].TotalAmount)

	mockDonationRepo.AssertExpectations(t)
	mockDonorRepo.AssertExpectations(t)
}

func TestCampaignUseCase_UpdateCampaignAmount(t *testing.T) {
	mockCampaignRepo := new(mocks.CampaignRepository)
	mockDonationRepo := new(mocks.DonationRepository)
	mockDonorRepo := new(mocks.DonorRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewCampaignUseCase(mockCampaignRepo, mockDonationRepo, mockDonorRepo, mockAuditLogRepo)

	campaignID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	campaign := &entities.Campaign{ID: campaignID, GoalAmount: 1000}

	mockCampaignRepo.On("FindByID", mock.Anything, campaignID).Return(campaign, nil)
	mockCampaignRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Campaign")).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil)

	err := uc.UpdateCampaignAmount(context.Background(), campaignID, 500.0, userID)

	assert.NoError(t, err)
	assert.Equal(t, 500.0, campaign.CurrentAmount)
	assert.Equal(t, userID, campaign.UpdatedBy)

	mockCampaignRepo.AssertExpectations(t)
	mockAuditLogRepo.AssertExpectations(t)
}

func TestCampaignUseCase_GetCampaignProgress(t *testing.T) {
	mockCampaignRepo := new(mocks.CampaignRepository)
	mockDonationRepo := new(mocks.DonationRepository)
	mockDonorRepo := new(mocks.DonorRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewCampaignUseCase(mockCampaignRepo, mockDonationRepo, mockDonorRepo, mockAuditLogRepo)

	campaignID := primitive.NewObjectID()
	campaign := &entities.Campaign{
		ID:            campaignID,
		GoalAmount:    1000,
		CurrentAmount: 250,
		DonationCount: 10,
	}

	mockCampaignRepo.On("FindByID", mock.Anything, campaignID).Return(campaign, nil)

	progress, err := uc.GetCampaignProgress(context.Background(), campaignID)

	assert.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, 1000.0, progress.GoalAmount)
	assert.Equal(t, 250.0, progress.CurrentAmount)
	assert.Equal(t, 10, progress.DonationCount)
	assert.Equal(t, 25.0, progress.Percentage)

	mockCampaignRepo.AssertExpectations(t)
}

func TestCampaignUseCase_ShareCampaign(t *testing.T) {
	mockCampaignRepo := new(mocks.CampaignRepository)
	mockDonationRepo := new(mocks.DonationRepository)
	mockDonorRepo := new(mocks.DonorRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewCampaignUseCase(mockCampaignRepo, mockDonationRepo, mockDonorRepo, mockAuditLogRepo)

	campaignID := primitive.NewObjectID()
	campaign := &entities.Campaign{
		ID:          campaignID,
		Name:        entities.MultilingualName{English: "Test Campaign"},
		Description: entities.MultilingualName{English: "Test Description"},
	}

	mockCampaignRepo.On("FindByID", mock.Anything, campaignID).Return(campaign, nil)

	shareable, err := uc.ShareCampaign(context.Background(), campaignID)

	assert.NoError(t, err)
	assert.NotNil(t, shareable)
	assert.Equal(t, "Test Campaign", shareable.Title)
	assert.Equal(t, "Test Description", shareable.Description)
	assert.Equal(t, "/campaigns/"+campaignID.Hex(), shareable.URL)

	mockCampaignRepo.AssertExpectations(t)
}
