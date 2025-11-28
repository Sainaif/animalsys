package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VolunteerUseCase struct {
	mock.Mock
}

func (m *VolunteerUseCase) LogHours(ctx context.Context, volunteerID primitive.ObjectID, hours float64, notes string, userID primitive.ObjectID) (*entities.Volunteer, error) {
	args := m.Called(ctx, volunteerID, hours, notes, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}

func (m *VolunteerUseCase) GetVolunteerHours(ctx context.Context, volunteerID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, volunteerID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *VolunteerUseCase) ActivateVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, volunteerID, userID)
	return args.Error(0)
}

func (m *VolunteerUseCase) DeactivateVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, volunteerID, userID)
	return args.Error(0)
}

func (m *VolunteerUseCase) GetVolunteersByRole(ctx context.Context, role entities.VolunteerRole) ([]*entities.Volunteer, int64, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Volunteer), args.Get(1).(int64), args.Error(2)
}

func (m *VolunteerUseCase) CreateVolunteer(ctx context.Context, volunteer *entities.Volunteer, userID primitive.ObjectID) error {
	args := m.Called(ctx, volunteer, userID)
	return args.Error(0)
}
func (m *VolunteerUseCase) GetVolunteer(ctx context.Context, volunteerID primitive.ObjectID) (*entities.Volunteer, error) {
	args := m.Called(ctx, volunteerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}
func (m *VolunteerUseCase) UpdateVolunteer(ctx context.Context, volunteer *entities.Volunteer, userID primitive.ObjectID) error {
	args := m.Called(ctx, volunteer, userID)
	return args.Error(0)
}
func (m *VolunteerUseCase) GetVolunteerByEmail(ctx context.Context, email string) (*entities.Volunteer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Volunteer), args.Error(1)
}
func (m *VolunteerUseCase) DeleteVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) ListVolunteers(ctx context.Context, filter *repositories.VolunteerFilter) ([]*entities.Volunteer, int64, error) {
	panic("unimplemented")
}
func (m *VolunteerUseCase) GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error) {
	panic("unimplemented")
}
func (m *VolunteerUseCase) GetVolunteersBySkill(ctx context.Context, skill string) ([]*entities.Volunteer, error) {
	panic("unimplemented")
}
func (m *VolunteerUseCase) GetVolunteersNeedingBackgroundCheck(ctx context.Context) ([]*entities.Volunteer, error) {
	panic("unimplemented")
}
func (m *VolunteerUseCase) GetTopVolunteers(ctx context.Context, limit int) ([]*entities.Volunteer, error) {
	panic("unimplemented")
}
func (m *VolunteerUseCase) ApproveVolunteer(ctx context.Context, volunteerID primitive.ObjectID, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) SuspendVolunteer(ctx context.Context, volunteerID primitive.ObjectID, reason string, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) AddCommendation(ctx context.Context, volunteerID primitive.ObjectID, note string, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) AddWarning(ctx context.Context, volunteerID primitive.ObjectID, reason string, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) AddCertification(ctx context.Context, volunteerID primitive.ObjectID, cert entities.Certification, userID primitive.ObjectID) error {
	panic("unimplemented")
}
func (m *VolunteerUseCase) GetVolunteerStatistics(ctx context.Context) (*repositories.VolunteerStatistics, error) {
	panic("unimplemented")
}
