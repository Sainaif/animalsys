package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerRepository is a mock type for the PartnerRepository type
type PartnerRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, partner
func (m *PartnerRepository) Create(ctx context.Context, partner *entities.Partner) error {
	args := m.Called(ctx, partner)
	return args.Error(0)
}

// FindByID provides a mock function with given fields: ctx, id
func (m *PartnerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Partner), args.Error(1)
}

// Update provides a mock function with given fields: ctx, partner
func (m *PartnerRepository) Update(ctx context.Context, partner *entities.Partner) error {
	args := m.Called(ctx, partner)
	return args.Error(0)
}

// Delete provides a mock function with given fields: ctx, id
func (m *PartnerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List provides a mock function with given fields: ctx, filter
func (m *PartnerRepository) List(ctx context.Context, filter *repositories.PartnerFilter) ([]*entities.Partner, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Partner), args.Get(1).(int64), args.Error(2)
}

// GetByType provides a mock function with given fields: ctx, partnerType
func (m *PartnerRepository) GetByType(ctx context.Context, partnerType entities.PartnerType) ([]*entities.Partner, error) {
	args := m.Called(ctx, partnerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Partner), args.Error(1)
}

// GetByStatus provides a mock function with given fields: ctx, status
func (m *PartnerRepository) GetByStatus(ctx context.Context, status entities.PartnerStatus) ([]*entities.Partner, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Partner), args.Error(1)
}

// GetActivePartners provides a mock function with given fields: ctx
func (m *PartnerRepository) GetActivePartners(ctx context.Context) ([]*entities.Partner, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Partner), args.Error(1)
}

// GetPartnersWithCapacity provides a mock function with given fields: ctx
func (m *PartnerRepository) GetPartnersWithCapacity(ctx context.Context) ([]*entities.Partner, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Partner), args.Error(1)
}

// GetPartnerStatistics provides a mock function with given fields: ctx
func (m *PartnerRepository) GetPartnerStatistics(ctx context.Context) (*repositories.PartnerStatistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.PartnerStatistics), args.Error(1)
}

// EnsureIndexes provides a mock function with given fields: ctx
func (m *PartnerRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
