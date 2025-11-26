package event

import (
	"context"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockEventRepository struct {
	events map[primitive.ObjectID]*entities.Event
}

func (m *mockEventRepository) Create(ctx context.Context, event *entities.Event) error {
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockEventRepository) Update(ctx context.Context, event *entities.Event) error {
	m.events[event.ID] = event
	return nil
}

func (m *mockEventRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Event, error) {
	event, ok := m.events[id]
	if !ok {
		return nil, nil
	}
	return event, nil
}

func (m *mockEventRepository) List(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error) {
	return nil, 0, nil
}

func (m *mockEventRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.events, id)
	return nil
}

func (m *mockEventRepository) GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetActiveEvents(ctx context.Context) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetCompletedEvents(ctx context.Context, limit int) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetPublicEvents(ctx context.Context) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetEventsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error) {
	return nil, nil
}

func (m *mockEventRepository) UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error {
	event, ok := m.events[eventID]
	if !ok {
		return nil
	}
	event.UpdateStatistics(attendees, volunteers, fundsRaised, animalsAdopted)
	return nil
}

func (m *mockEventRepository) GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error) {
	return nil, nil
}

func (m *mockEventRepository) EnsureIndexes(ctx context.Context) error {
	return nil
}

func TestEventUseCase_UpdateStatistics(t *testing.T) {
	repo := &mockEventRepository{
		events: make(map[primitive.ObjectID]*entities.Event),
	}
	useCase := NewEventUseCase(repo, nil, nil, nil)

	event := &entities.Event{
		ID:   primitive.NewObjectID(),
		Type: entities.EventTypeShopping,
	}
	repo.events[event.ID] = event

	err := useCase.UpdateEventStatistics(context.Background(), event.ID, 0, 0, 100, 0)
	assert.NoError(t, err)

	updatedEvent, err := repo.FindByID(context.Background(), event.ID)
	assert.NoError(t, err)
	assert.Equal(t, float64(100), updatedEvent.Budget)
}
