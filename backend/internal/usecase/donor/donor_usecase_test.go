package donor

import (
	"context"
	"testing"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockDonorRepository struct {
	mock.Mock
}

func (m *mockDonorRepository) Create(ctx context.Context, donor *entities.Donor) error {
	panic("implement me")
}
func (m *mockDonorRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Donor), args.Error(1)
}
func (m *mockDonorRepository) Update(ctx context.Context, donor *entities.Donor) error {
	args := m.Called(ctx, donor)
	return args.Error(0)
}
func (m *mockDonorRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}
func (m *mockDonorRepository) List(ctx context.Context, filter *repositories.DonorFilter) ([]*entities.Donor, int64, error) {
	panic("implement me")
}
func (m *mockDonorRepository) FindByEmail(ctx context.Context, email string) (*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorRepository) GetMajorDonors(ctx context.Context) ([]*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorRepository) GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorRepository) GetDonorStatistics(ctx context.Context) (*repositories.DonorStatistics, error) {
	panic("implement me")
}
func (m *mockDonorRepository) EnsureIndexes(ctx context.Context) error {
	panic("implement me")
}
func (m *mockDonorRepository) FindDonorsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*entities.Donor, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*entities.Donor), args.Error(1)
}

type mockDonationRepository struct {
	mock.Mock
}

func (m *mockDonationRepository) GetTopDonorsByTotalDonated(ctx context.Context, limit int) ([]*repositories.TopDonorResult, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*repositories.TopDonorResult), args.Error(1)
}

func (m *mockDonationRepository) GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Donation), args.Error(1)
}
func (m *mockDonationRepository) Create(ctx context.Context, donation *entities.Donation) error {
	panic("implement me")
}
func (m *mockDonationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donation, error) {
	panic("implement me")
}
func (m *mockDonationRepository) Update(ctx context.Context, donation *entities.Donation) error {
	panic("implement me")
}
func (m *mockDonationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	panic("implement me")
}
func (m *mockDonationRepository) List(ctx context.Context, filter *repositories.DonationFilter) ([]*entities.Donation, int64, error) {
	panic("implement me")
}
func (m *mockDonationRepository) GetByDonorID(ctx context.Context, donorID primitive.ObjectID) ([]*entities.Donation, error) {
	panic("implement me")
}
func (m *mockDonationRepository) GetByCampaignID(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Donation, error) {
	panic("implement me")
}
func (m *mockDonationRepository) GetPendingThankYous(ctx context.Context) ([]*entities.Donation, error) {
	panic("implement me")
}
func (m *mockDonationRepository) GetDonationStatistics(ctx context.Context) (*repositories.DonationStatistics, error) {
	panic("implement me")
}
func (m *mockDonationRepository) GetDonationsByDateRange(ctx context.Context, from time.Time, to time.Time) ([]*entities.Donation, error) {
	panic("implement me")
}
func (m *mockDonationRepository) EnsureIndexes(ctx context.Context) error {
	panic("implement me")
}

type mockAuditLogRepository struct {
	mock.Mock
}

func (m *mockAuditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *mockAuditLogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	panic("implement me")
}
func (m *mockAuditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	panic("implement me")
}
func (m *mockAuditLogRepository) EnsureIndexes(ctx context.Context) error {
	panic("implement me")
}
func (m *mockAuditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	panic("implement me")
}
func TestDonorUseCase_GetTopDonors(t *testing.T) {
	mockDonorRepo := new(mockDonorRepository)
	mockDonationRepo := new(mockDonationRepository)
	mockAuditLogRepo := new(mockAuditLogRepository)
	uc := NewDonorUseCase(mockDonorRepo, mockDonationRepo, mockAuditLogRepo)

	ctx := context.Background()
	limit := 5

	topDonorsResult := []*repositories.TopDonorResult{
		{DonorID: primitive.NewObjectID(), TotalAmount: 1000},
		{DonorID: primitive.NewObjectID(), TotalAmount: 500},
	}
	var donorIDs []primitive.ObjectID
	for _, result := range topDonorsResult {
		donorIDs = append(donorIDs, result.DonorID)
	}

	donors := []*entities.Donor{
		{ID: donorIDs[0]},
		{ID: donorIDs[1]},
	}

	mockDonationRepo.On("GetTopDonorsByTotalDonated", ctx, limit).Return(topDonorsResult, nil)
	mockDonorRepo.On("FindDonorsByIDs", ctx, donorIDs).Return(donors, nil)

	response, err := uc.GetTopDonors(ctx, limit)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, float64(1000), response[0].TotalDonated)
	assert.Equal(t, float64(500), response[1].TotalDonated)
	mockDonationRepo.AssertExpectations(t)
	mockDonorRepo.AssertExpectations(t)
}

func TestDonorUseCase_GetRecurringDonors(t *testing.T) {
	mockDonorRepo := new(mockDonorRepository)
	mockDonationRepo := new(mockDonationRepository)
	mockAuditLogRepo := new(mockAuditLogRepository)
	uc := NewDonorUseCase(mockDonorRepo, mockDonationRepo, mockAuditLogRepo)

	ctx := context.Background()
	donorID1 := primitive.NewObjectID()
	donorID2 := primitive.NewObjectID()

	donations := []*entities.Donation{
		{DonorID: donorID1, Amount: 100},
		{DonorID: donorID2, Amount: 200},
		{DonorID: donorID1, Amount: 50},
	}

	donors := []*entities.Donor{
		{ID: donorID1},
		{ID: donorID2},
	}

	mockDonationRepo.On("GetRecurringDonations", ctx).Return(donations, nil)
	mockDonorRepo.On("FindDonorsByIDs", ctx, []primitive.ObjectID{donorID1, donorID2}).Return(donors, nil)

	response, err := uc.GetRecurringDonors(ctx)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	for _, r := range response {
		if r.ID == donorID1 {
			assert.Equal(t, 2, r.RecurringCount)
			assert.Equal(t, float64(150), r.RecurringAmount)
		} else if r.ID == donorID2 {
			assert.Equal(t, 1, r.RecurringCount)
			assert.Equal(t, float64(200), r.RecurringAmount)
		}
	}
	mockDonationRepo.AssertExpectations(t)
	mockDonorRepo.AssertExpectations(t)
}

func TestDonorUseCase_UpdateCommunicationPreferences(t *testing.T) {
	mockDonorRepo := new(mockDonorRepository)
	mockDonationRepo := new(mockDonationRepository)
	mockAuditLogRepo := new(mockAuditLogRepository)
	uc := NewDonorUseCase(mockDonorRepo, mockDonationRepo, mockAuditLogRepo)

	ctx := context.Background()
	donorID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	prefs := entities.DonorPreferences{EmailOptIn: true, SmsOptIn: false}

	mockDonorRepo.On("FindByID", ctx, donorID).Return(&entities.Donor{ID: donorID}, nil)
	mockDonorRepo.On("Update", ctx, mock.AnythingOfType("*entities.Donor")).Return(nil)
	mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*entities.AuditLog")).Return(nil)

	err := uc.UpdateCommunicationPreferences(ctx, donorID, prefs, userID)
	assert.NoError(t, err)
	mockDonorRepo.AssertExpectations(t)
}
