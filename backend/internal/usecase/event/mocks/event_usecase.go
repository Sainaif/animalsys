package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventUseCase struct {
	mock.Mock
}

func (m *EventUseCase) CreateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error {
	args := m.Called(ctx, event, userID)
	return args.Error(0)
}

func (m *EventUseCase) UpdateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error {
	args := m.Called(ctx, event, userID)
	return args.Error(0)
}

func (m *EventUseCase) DeleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *EventUseCase) GetEvent(ctx context.Context, eventID primitive.ObjectID) (*entities.Event, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *EventUseCase) ListEvents(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Event), args.Get(1).(int64), args.Error(2)
}

func (m *EventUseCase) GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetActiveEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetPastEvents(ctx context.Context, limit int) ([]*entities.Event, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetPublicEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error) {
	args := m.Called(ctx, organizerID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventUseCase) ActivateEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *EventUseCase) CompleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *EventUseCase) CancelEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *EventUseCase) AssignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, volunteerID, userID)
	return args.Error(0)
}

func (m *EventUseCase) UnassignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, volunteerID, userID)
	return args.Error(0)
}

func (m *EventUseCase) UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error {
	args := m.Called(ctx, eventID, attendees, volunteers, fundsRaised, animalsAdopted)
	return args.Error(0)
}

func (m *EventUseCase) GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(*repositories.EventStatistics), args.Error(1)
}

func (m *EventUseCase) RegisterForEvent(ctx context.Context, eventID primitive.ObjectID, req entities.EventAttendance, userID primitive.ObjectID) (*entities.EventAttendance, error) {
	args := m.Called(ctx, eventID, req, userID)
	return args.Get(0).(*entities.EventAttendance), args.Error(1)
}

func (m *EventUseCase) GetEventRegistrations(ctx context.Context, eventID primitive.ObjectID, limit, offset int64) ([]*entities.EventAttendance, int64, error) {
	args := m.Called(ctx, eventID, limit, offset)
	return args.Get(0).([]*entities.EventAttendance), args.Get(1).(int64), args.Error(2)
}

func (m *EventUseCase) GetEventStatisticsDetail(ctx context.Context, eventID primitive.ObjectID) (*entities.Event, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *EventUseCase) PublishEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *EventUseCase) SendEventReminder(ctx context.Context, eventID primitive.ObjectID) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}
