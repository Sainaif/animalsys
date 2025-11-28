package event

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventUseCaseInterface interface {
	CreateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error
	UpdateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error
	DeleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
	GetEvent(ctx context.Context, eventID primitive.ObjectID) (*entities.Event, error)
	ListEvents(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error)
	GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error)
	GetActiveEvents(ctx context.Context) ([]*entities.Event, error)
	GetPastEvents(ctx context.Context, limit int) ([]*entities.Event, error)
	GetPublicEvents(ctx context.Context) ([]*entities.Event, error)
	GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error)
	GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error)
	GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error)
	ActivateEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
	CompleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
	CancelEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
	AssignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error
	UnassignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error
	UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error
	GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error)
	RegisterForEvent(ctx context.Context, eventID primitive.ObjectID, req entities.EventAttendance, userID primitive.ObjectID) (*entities.EventAttendance, error)
	GetEventRegistrations(ctx context.Context, eventID primitive.ObjectID, limit, offset int64) ([]*entities.EventAttendance, int64, error)
	GetEventStatisticsDetail(ctx context.Context, eventID primitive.ObjectID) (*entities.Event, error)
	PublishEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error
	SendEventReminder(ctx context.Context, eventID primitive.ObjectID) (int, error)
}
