package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventAttendanceRepository struct {
	mock.Mock
}

func (m *EventAttendanceRepository) Create(ctx context.Context, attendance *entities.EventAttendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *EventAttendanceRepository) Update(ctx context.Context, attendance *entities.EventAttendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *EventAttendanceRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *EventAttendanceRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.EventAttendance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) List(ctx context.Context, filter *repositories.EventAttendanceFilter) ([]*entities.EventAttendance, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.EventAttendance), args.Get(1).(int64), args.Error(2)
}

func (m *EventAttendanceRepository) GetAttendanceByEvent(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) GetAttendanceByVolunteer(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx, volunteerID)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) GetAttendanceByDonor(ctx context.Context, donorID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx, donorID)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) GetConfirmedAttendees(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) GetNoShows(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) GetPendingPayments(ctx context.Context) ([]*entities.EventAttendance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.EventAttendance), args.Error(1)
}

func (m *EventAttendanceRepository) CountAttendeesByEvent(ctx context.Context, eventID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *EventAttendanceRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
