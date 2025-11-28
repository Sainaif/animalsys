package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/event"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	eventUseCase event.EventUseCaseInterface
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventUseCase event.EventUseCaseInterface) *EventHandler {
	return &EventHandler{
		eventUseCase: eventUseCase,
	}
}

// CreateEvent creates a new event
// @Summary Create a new event
// @Tags events
// @Accept json
// @Produce json
// @Param event body entities.Event true "Event data"
// @Success 201 {object} entities.Event
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event entities.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.CreateEvent(c.Request.Context(), &event, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvent gets an event by ID
// @Summary Get event by ID
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} entities.Event
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.eventUseCase.GetEvent(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent updates an event
// @Summary Update event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body entities.Event true "Event data"
// @Success 200 {object} entities.Event
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event entities.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.UpdateEvent(c.Request.Context(), &event, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent deletes an event
// @Summary Delete event
// @Tags events
// @Param id path string true "Event ID"
// @Success 204
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.DeleteEvent(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListEvents lists events with filters
// @Summary List events
// @Tags events
// @Produce json
// @Param type query string false "Event type"
// @Param status query string false "Event status"
// @Param public query bool false "Public events"
// @Param featured query bool false "Featured events"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	filter := &repositories.EventFilter{
		Type:      string(entities.EventType(c.Query("type"))),
		Status:    string(entities.EventStatus(c.Query("status"))),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "start_date"),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
	}

	// Parse public
	if publicStr := c.Query("public"); publicStr != "" {
		if public, err := strconv.ParseBool(publicStr); err == nil {
			filter.Public = &public
		}
	}

	// Parse featured
	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filter.Featured = &featured
		}
	}

	// Parse start date
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	// Parse end date
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.ParseInt(limit, 10, 64); err == nil {
			filter.Limit = l
		}
	} else {
		filter.Limit = 20
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.ParseInt(offset, 10, 64); err == nil {
			filter.Offset = o
		}
	}

	events, total, err := h.eventUseCase.ListEvents(c.Request.Context(), filter)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", c.FullPath()).
			Str("type", c.Query("type")).
			Str("status", c.Query("status")).
			Msg("failed to list events")
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetUpcomingEvents gets upcoming events
// @Summary Get upcoming events
// @Tags events
// @Produce json
// @Success 200 {array} entities.Event
// @Router /events/upcoming [get]
func (h *EventHandler) GetUpcomingEvents(c *gin.Context) {
	events, err := h.eventUseCase.GetUpcomingEvents(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetActiveEvents gets active events
// @Summary Get active events
// @Tags events
// @Produce json
// @Success 200 {array} entities.Event
// @Router /events/active [get]
func (h *EventHandler) GetActiveEvents(c *gin.Context) {
	events, err := h.eventUseCase.GetActiveEvents(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetPublicEvents gets public events
// @Summary Get public events
// @Tags events
// @Produce json
// @Success 200 {array} entities.Event
// @Router /events/public [get]
func (h *EventHandler) GetPublicEvents(c *gin.Context) {
	events, err := h.eventUseCase.GetPublicEvents(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetFeaturedEvents gets featured events
// @Summary Get featured events
// @Tags events
// @Produce json
// @Success 200 {array} entities.Event
// @Router /events/featured [get]
func (h *EventHandler) GetFeaturedEvents(c *gin.Context) {
	events, err := h.eventUseCase.GetFeaturedEvents(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetEventsNeedingVolunteers gets events needing volunteers
// @Summary Get events needing volunteers
// @Tags events
// @Produce json
// @Success 200 {array} entities.Event
// @Router /events/needing-volunteers [get]
func (h *EventHandler) GetEventsNeedingVolunteers(c *gin.Context) {
	events, err := h.eventUseCase.GetEventsNeedingVolunteers(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// ActivateEvent activates an event
// @Summary Activate event
// @Tags events
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /events/{id}/activate [post]
func (h *EventHandler) ActivateEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.ActivateEvent(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event activated successfully"})
}

// CompleteEvent completes an event
// @Summary Complete event
// @Tags events
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /events/{id}/complete [post]
func (h *EventHandler) CompleteEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.CompleteEvent(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event completed successfully"})
}

// CancelEvent cancels an event
// @Summary Cancel event
// @Tags events
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /events/{id}/cancel [post]
func (h *EventHandler) CancelEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.CancelEvent(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event cancelled successfully"})
}

// AssignVolunteer assigns a volunteer to an event
// @Summary Assign volunteer to event
// @Tags events
// @Param id path string true "Event ID"
// @Param volunteer_id body map[string]string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /events/{id}/assign-volunteer [post]
func (h *EventHandler) AssignVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req struct {
		VolunteerID string `json:"volunteer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteerID, err := primitive.ObjectIDFromHex(req.VolunteerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.AssignVolunteer(c.Request.Context(), id, volunteerID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer assigned successfully"})
}

// UnassignVolunteer removes a volunteer from an event
// @Summary Unassign volunteer from event
// @Tags events
// @Param id path string true "Event ID"
// @Param volunteer_id body map[string]string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /events/{id}/unassign-volunteer [post]
func (h *EventHandler) UnassignVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req struct {
		VolunteerID string `json:"volunteer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteerID, err := primitive.ObjectIDFromHex(req.VolunteerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.UnassignVolunteer(c.Request.Context(), id, volunteerID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer unassigned successfully"})
}

// GetEventStatistics gets event statistics
// @Summary Get event statistics
// @Tags events
// @Produce json
// @Success 200 {object} repositories.EventStatistics
// @Router /events/statistics [get]
func (h *EventHandler) GetEventStatistics(c *gin.Context) {
	stats, err := h.eventUseCase.GetEventStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPastEvents gets past events
func (h *EventHandler) GetPastEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	events, err := h.eventUseCase.GetPastEvents(c.Request.Context(), limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  len(events),
	})
}

// RegisterForEvent registers a user for an event
func (h *EventHandler) RegisterForEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req entities.EventAttendance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	attendance, err := h.eventUseCase.RegisterForEvent(c.Request.Context(), id, req, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attendance)
}

// GetEventRegistrations gets all registrations for an event
func (h *EventHandler) GetEventRegistrations(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	registrations, total, err := h.eventUseCase.GetEventRegistrations(c.Request.Context(), id, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registrations": registrations,
		"total":         total,
	})
}

// GetEventStatisticsDetail gets detailed statistics for an event
func (h *EventHandler) GetEventStatisticsDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.eventUseCase.GetEventStatisticsDetail(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attendee_count":    event.AttendeeCount,
		"volunteer_count":   event.VolunteerCount,
		"funds_raised":      event.FundsRaised,
		"animals_adopted":   event.AnimalsAdopted,
		"registrations":     event.Registration.CurrentCount,
	})
}

// PublishEvent publishes an event
func (h *EventHandler) PublishEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.eventUseCase.PublishEvent(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event published successfully"})
}

// SendEventReminder sends reminder for an event
func (h *EventHandler) SendEventReminder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	count, err := h.eventUseCase.SendEventReminder(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Reminders sent successfully",
		"recipients_count": count,
	})
}

// GetEventAttendance gets attendance for an event
func (h *EventHandler) GetEventAttendance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	attendance, total, err := h.eventUseCase.GetEventRegistrations(c.Request.Context(), id, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attendance": attendance,
		"total":      total,
	})
}
