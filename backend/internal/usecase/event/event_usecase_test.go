package event

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

func TestEventUseCase_GetPastEvents(t *testing.T) {
	mockEventRepo := new(mocks.EventRepository)
	uc := NewEventUseCase(mockEventRepo, nil, nil, nil)

	expectedEvents := []*entities.Event{{ID: primitive.NewObjectID()}}
	mockEventRepo.On("GetCompletedEvents", mock.Anything, 20).Return(expectedEvents, nil)

	events, err := uc.GetPastEvents(context.Background(), 20)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvents, events)
	mockEventRepo.AssertExpectations(t)
}

func TestEventUseCase_RegisterForEvent(t *testing.T) {
	mockEventRepo := new(mocks.EventRepository)
	mockAttendanceRepo := new(mocks.EventAttendanceRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewEventUseCase(mockEventRepo, mockAttendanceRepo, nil, mockAuditLogRepo)

	eventID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	event := &entities.Event{
		ID:     eventID,
		Status: entities.EventStatusScheduled,
		Registration: entities.EventRegistration{
			Required:     true,
			MaxAttendees: 10,
			CurrentCount: 5,
		},
	}
	req := entities.EventAttendance{
		AttendeeType: entities.AttendeeTypePublic,
	}

	mockEventRepo.On("FindByID", mock.Anything, eventID).Return(event, nil)
	mockAttendanceRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.EventAttendance")).Return(nil)
	mockEventRepo.On("Update", mock.Anything, event).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil)

	attendance, err := uc.RegisterForEvent(context.Background(), eventID, req, userID)

	assert.NoError(t, err)
	assert.NotNil(t, attendance)
	assert.Equal(t, eventID, attendance.EventID)
	mockEventRepo.AssertExpectations(t)
	mockAttendanceRepo.AssertExpectations(t)
}

func TestEventUseCase_RegisterForEvent_FullEvent(t *testing.T) {
	mockEventRepo := new(mocks.EventRepository)
	uc := NewEventUseCase(mockEventRepo, nil, nil, nil)

	eventID := primitive.NewObjectID()
	event := &entities.Event{
		ID:     eventID,
		Status: entities.EventStatusScheduled,
		Registration: entities.EventRegistration{
			Required:     true,
			MaxAttendees: 10,
			CurrentCount: 10,
		},
	}
	req := entities.EventAttendance{}

	mockEventRepo.On("FindByID", mock.Anything, eventID).Return(event, nil)

	_, err := uc.RegisterForEvent(context.Background(), eventID, req, primitive.NewObjectID())

	assert.Error(t, err)
	assert.Equal(t, "Event is full", err.Error())
}

func TestEventUseCase_GetEventRegistrations(t *testing.T) {
	mockAttendanceRepo := new(mocks.EventAttendanceRepository)
	uc := NewEventUseCase(nil, mockAttendanceRepo, nil, nil)

	eventID := primitive.NewObjectID()
	expectedRegistrations := []*entities.EventAttendance{{ID: primitive.NewObjectID()}}
	filter := &repositories.EventAttendanceFilter{
		EventID: &eventID,
		Limit:   20,
		Offset:  0,
	}

	mockAttendanceRepo.On("List", mock.Anything, filter).Return(expectedRegistrations, int64(len(expectedRegistrations)), nil)

	registrations, total, err := uc.GetEventRegistrations(context.Background(), eventID, 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedRegistrations, registrations)
	assert.Equal(t, int64(1), total)
	mockAttendanceRepo.AssertExpectations(t)
}

func TestEventUseCase_GetEventStatisticsDetail(t *testing.T) {
	mockEventRepo := new(mocks.EventRepository)
	uc := NewEventUseCase(mockEventRepo, nil, nil, nil)

	eventID := primitive.NewObjectID()
	expectedEvent := &entities.Event{
		ID:            eventID,
		AttendeeCount: 100,
		VolunteerCount: 10,
		FundsRaised:   1000.0,
		AnimalsAdopted: 5,
		Registration: entities.EventRegistration{
			CurrentCount: 50,
		},
	}
	mockEventRepo.On("FindByID", mock.Anything, eventID).Return(expectedEvent, nil)

	event, err := uc.GetEventStatisticsDetail(context.Background(), eventID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, event)
	mockEventRepo.AssertExpectations(t)
}

func TestEventUseCase_PublishEvent(t *testing.T) {
	mockEventRepo := new(mocks.EventRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewEventUseCase(mockEventRepo, nil, nil, mockAuditLogRepo)

	eventID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	event := &entities.Event{
		ID:     eventID,
		Status: entities.EventStatusDraft,
	}

	mockEventRepo.On("FindByID", mock.Anything, eventID).Return(event, nil)
	mockEventRepo.On("Update", mock.Anything, event).Return(nil)
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil)

	err := uc.PublishEvent(context.Background(), eventID, userID)

	assert.NoError(t, err)
	assert.Equal(t, entities.EventStatusScheduled, event.Status)
	mockEventRepo.AssertExpectations(t)
}

func TestEventUseCase_SendEventReminder(t *testing.T) {
	mockAttendanceRepo := new(mocks.EventAttendanceRepository)
	uc := NewEventUseCase(nil, mockAttendanceRepo, nil, nil)

	eventID := primitive.NewObjectID()
	expectedAttendees := []*entities.EventAttendance{{ID: primitive.NewObjectID()}, {ID: primitive.NewObjectID()}}
	mockAttendanceRepo.On("GetAttendanceByEvent", mock.Anything, eventID).Return(expectedAttendees, nil)

	count, err := uc.SendEventReminder(context.Background(), eventID)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	mockAttendanceRepo.AssertExpectations(t)
}
