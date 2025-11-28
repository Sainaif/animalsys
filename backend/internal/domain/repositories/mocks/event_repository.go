package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventRepository struct {
	mock.Mock
}

func (m *EventRepository) Create(ctx context.Context, event *entities.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventRepository) Update(ctx context.Context, event *entities.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *EventRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *EventRepository) List(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Event), args.Get(1).(int64), args.Error(2)
}

func (m *EventRepository) GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetActiveEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetCompletedEvents(ctx context.Context, limit int) ([]*entities.Event, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetPublicEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error) {
	args := m.Called(ctx, organizerID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetEventsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Event, error) {
	args := m.Called(ctx, campaignID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *EventRepository) UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error {
	args := m.Called(ctx, eventID, attendees, volunteers, fundsRaised, animalsAdopted)
	return args.Error(0)
}

func (m *EventRepository) GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(*repositories.EventStatistics), args.Error(1)
}

func (m *EventRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
