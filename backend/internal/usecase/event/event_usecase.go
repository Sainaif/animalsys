package event

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventUseCase handles event-related business logic
type EventUseCase struct {
	eventRepo      repositories.EventRepository
	attendanceRepo repositories.EventAttendanceRepository
	volunteerRepo  repositories.VolunteerRepository
	auditLogRepo   repositories.AuditLogRepository
}

// NewEventUseCase creates a new event use case
func NewEventUseCase(
	eventRepo repositories.EventRepository,
	attendanceRepo repositories.EventAttendanceRepository,
	volunteerRepo repositories.VolunteerRepository,
	auditLogRepo repositories.AuditLogRepository,
) *EventUseCase {
	return &EventUseCase{
		eventRepo:      eventRepo,
		attendanceRepo: attendanceRepo,
		volunteerRepo:  volunteerRepo,
		auditLogRepo:   auditLogRepo,
	}
}

// CreateEvent creates a new event
func (uc *EventUseCase) CreateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error {
	// Validate event
	if err := uc.validateEvent(event); err != nil {
		return err
	}

	// Create event
	if err := uc.eventRepo.Create(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "event", "", "").
			WithEntityID(event.ID))

	return nil
}

// UpdateEvent updates an event
func (uc *EventUseCase) UpdateEvent(ctx context.Context, event *entities.Event, userID primitive.ObjectID) error {
	// Check if event exists
	existingEvent, err := uc.eventRepo.FindByID(ctx, event.ID)
	if err != nil {
		return err
	}

	// Validate event
	if err := uc.validateEvent(event); err != nil {
		return err
	}

	// Update event
	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "", "").
			WithEntityID(event.ID).
			WithChanges(map[string]interface{}{
				"old_status": existingEvent.Status,
				"new_status": event.Status,
			}))

	return nil
}

// DeleteEvent deletes an event
func (uc *EventUseCase) DeleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if event exists
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Check if event can be deleted (only draft or cancelled events)
	if event.Status != entities.EventStatusDraft && event.Status != entities.EventStatusCancelled {
		return errors.NewBadRequest("Only draft or cancelled events can be deleted")
	}

	// Delete event
	if err := uc.eventRepo.Delete(ctx, eventID); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "event", "", "").
			WithEntityID(eventID))

	return nil
}

// GetEvent gets an event by ID
func (uc *EventUseCase) GetEvent(ctx context.Context, eventID primitive.ObjectID) (*entities.Event, error) {
	return uc.eventRepo.FindByID(ctx, eventID)
}

// ListEvents lists events with filters
func (uc *EventUseCase) ListEvents(ctx context.Context, filter *repositories.EventFilter) ([]*entities.Event, int64, error) {
	return uc.eventRepo.List(ctx, filter)
}

// GetUpcomingEvents gets upcoming events
func (uc *EventUseCase) GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error) {
	return uc.eventRepo.GetUpcomingEvents(ctx)
}

// GetActiveEvents gets active events
func (uc *EventUseCase) GetActiveEvents(ctx context.Context) ([]*entities.Event, error) {
	return uc.eventRepo.GetActiveEvents(ctx)
}

// GetPublicEvents gets public events
func (uc *EventUseCase) GetPublicEvents(ctx context.Context) ([]*entities.Event, error) {
	return uc.eventRepo.GetPublicEvents(ctx)
}

// GetFeaturedEvents gets featured events
func (uc *EventUseCase) GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error) {
	return uc.eventRepo.GetFeaturedEvents(ctx)
}

// GetEventsByOrganizer gets events by organizer
func (uc *EventUseCase) GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error) {
	return uc.eventRepo.GetEventsByOrganizer(ctx, organizerID)
}

// GetEventsNeedingVolunteers gets events needing volunteers
func (uc *EventUseCase) GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error) {
	return uc.eventRepo.GetEventsNeedingVolunteers(ctx)
}

// ActivateEvent activates an event
func (uc *EventUseCase) ActivateEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != entities.EventStatusScheduled {
		return errors.NewBadRequest("Only scheduled events can be activated")
	}

	event.Status = entities.EventStatusActive
	event.UpdatedBy = userID

	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "activated", "").
			WithEntityID(eventID))

	return nil
}

// CompleteEvent completes an event
func (uc *EventUseCase) CompleteEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != entities.EventStatusActive {
		return errors.NewBadRequest("Only active events can be completed")
	}

	event.Status = entities.EventStatusCompleted
	event.UpdatedBy = userID

	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "completed", "").
			WithEntityID(eventID))

	return nil
}

// CancelEvent cancels an event
func (uc *EventUseCase) CancelEvent(ctx context.Context, eventID primitive.ObjectID, userID primitive.ObjectID) error {
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status == entities.EventStatusCompleted || event.Status == entities.EventStatusCancelled {
		return errors.NewBadRequest("Cannot cancel completed or already cancelled events")
	}

	event.Status = entities.EventStatusCancelled
	event.UpdatedBy = userID

	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "cancelled", "").
			WithEntityID(eventID))

	return nil
}

// AssignVolunteer assigns a volunteer to an event
func (uc *EventUseCase) AssignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error {
	// Get event
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Check if volunteer exists
	volunteer, err := uc.volunteerRepo.FindByID(ctx, volunteerID)
	if err != nil {
		return errors.NewNotFound("Volunteer not found")
	}

	// Check if volunteer is active
	if volunteer.Status != entities.VolunteerStatusActive {
		return errors.NewBadRequest("Only active volunteers can be assigned")
	}

	// Check if volunteer already assigned
	for _, vID := range event.AssignedVolunteers {
		if vID == volunteerID {
			return errors.NewBadRequest("Volunteer already assigned to this event")
		}
	}

	// Add volunteer to event
	event.AssignedVolunteers = append(event.AssignedVolunteers, volunteerID)
	event.UpdatedBy = userID

	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "assigned_volunteer", "").
			WithEntityID(eventID).
			WithChanges(map[string]interface{}{
				"volunteer_id": volunteerID,
			}))

	return nil
}

// UnassignVolunteer removes a volunteer from an event
func (uc *EventUseCase) UnassignVolunteer(ctx context.Context, eventID, volunteerID, userID primitive.ObjectID) error {
	event, err := uc.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Remove volunteer from event
	newAssigned := []primitive.ObjectID{}
	found := false
	for _, vID := range event.AssignedVolunteers {
		if vID != volunteerID {
			newAssigned = append(newAssigned, vID)
		} else {
			found = true
		}
	}

	if !found {
		return errors.NewBadRequest("Volunteer not assigned to this event")
	}

	event.AssignedVolunteers = newAssigned
	event.UpdatedBy = userID

	if err := uc.eventRepo.Update(ctx, event); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "event", "unassigned_volunteer", "").
			WithEntityID(eventID).
			WithChanges(map[string]interface{}{
				"volunteer_id": volunteerID,
			}))

	return nil
}

// UpdateEventStatistics updates event statistics
func (uc *EventUseCase) UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error {
	return uc.eventRepo.UpdateEventStatistics(ctx, eventID, attendees, volunteers, fundsRaised, animalsAdopted)
}

// GetEventStatistics gets event statistics
func (uc *EventUseCase) GetEventStatistics(ctx context.Context) (*repositories.EventStatistics, error) {
	return uc.eventRepo.GetEventStatistics(ctx)
}

// validateEvent validates event data
func (uc *EventUseCase) validateEvent(event *entities.Event) error {
	if event.Name.English == "" {
		return errors.NewBadRequest("Event name is required")
	}

	if event.Type == "" {
		return errors.NewBadRequest("Event type is required")
	}

	if event.StartDate.IsZero() {
		return errors.NewBadRequest("Event start date is required")
	}

	if event.Duration <= 0 {
		return errors.NewBadRequest("Event duration must be greater than 0")
	}

	return nil
}
