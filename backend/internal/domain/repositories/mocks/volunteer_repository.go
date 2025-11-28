package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VolunteerRepository struct {
	mock.Mock
}

func (m *VolunteerRepository) Create(ctx context.Context, volunteer *entities.Volunteer) error {
	args := m.Called(ctx, volunteer)
	return args.Error(0)
}

func (m *VolunteerRepository) Update(ctx context.Context, volunteer *entities.Volunteer) error {
	args := m.Called(ctx, volunteer)
	return args.Error(0)
}

func (m *VolunteerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *VolunteerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Volunteer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) FindByEmail(ctx context.Context, email string) (*entities.Volunteer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Volunteer, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) List(ctx context.Context, filter *repositories.VolunteerFilter) ([]*entities.Volunteer, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Volunteer), args.Get(1).(int64), args.Error(2)
}

func (m *VolunteerRepository) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) GetVolunteersBySkill(ctx context.Context, skill string) ([]*entities.Volunteer, error) {
	args := m.Called(ctx, skill)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) GetVolunteersNeedingBackgroundCheck(ctx context.Context) ([]*entities.Volunteer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) GetVolunteersWithExpiredCertifications(ctx context.Context) ([]*entities.Volunteer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) GetTopVolunteers(ctx context.Context, limit int) ([]*entities.Volunteer, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Volunteer), args.Error(1)
}

func (m *VolunteerRepository) UpdateHours(ctx context.Context, volunteerID primitive.ObjectID, hours float64, notes string) error {
	args := m.Called(ctx, volunteerID, hours, notes)
	return args.Error(0)
}

func (m *VolunteerRepository) IncrementEventsAttended(ctx context.Context, volunteerID primitive.ObjectID) error {
	args := m.Called(ctx, volunteerID)
	return args.Error(0)
}

func (m *VolunteerRepository) GetVolunteerStatistics(ctx context.Context) (*repositories.VolunteerStatistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.VolunteerStatistics), args.Error(1)
}

func (m *VolunteerRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
