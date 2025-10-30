package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type VolunteerHandler struct {
	volunteerUseCase *usecases.VolunteerUseCase
}

func NewVolunteerHandler(volunteerUseCase *usecases.VolunteerUseCase) *VolunteerHandler {
	return &VolunteerHandler{
		volunteerUseCase: volunteerUseCase,
	}
}

// Create godoc
// @Summary Create volunteer
// @Description Create a new volunteer record
// @Tags volunteers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entities.VolunteerCreateRequest true "Volunteer data"
// @Success 201 {object} entities.Volunteer
// @Failure 400 {object} ErrorResponse
// @Router /volunteers [post]
func (h *VolunteerHandler) Create(c *gin.Context) {
	var req entities.VolunteerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	volunteer, err := h.volunteerUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, volunteer)
}

// GetByID godoc
// @Summary Get volunteer by ID
// @Description Get volunteer details by ID
// @Tags volunteers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Success 200 {object} entities.Volunteer
// @Failure 404 {object} ErrorResponse
// @Router /volunteers/{id} [get]
func (h *VolunteerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	volunteer, err := h.volunteerUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volunteer)
}

// List godoc
// @Summary List volunteers
// @Description Get list of volunteers with filtering and pagination
// @Tags volunteers
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param search query string false "Search in name/email"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /volunteers [get]
func (h *VolunteerHandler) List(c *gin.Context) {
	filter := &entities.VolunteerFilter{
		Status: entities.VolunteerStatus(c.Query("status")),
		Search: c.Query("search"),
		Limit:  parseIntQuery(c.Query("limit"), 10),
		Offset: parseIntQuery(c.Query("offset"), 0),
	}

	volunteers, total, err := h.volunteerUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   volunteers,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Update godoc
// @Summary Update volunteer
// @Description Update volunteer details
// @Tags volunteers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Param request body entities.VolunteerUpdateRequest true "Update data"
// @Success 200 {object} entities.Volunteer
// @Failure 400 {object} ErrorResponse
// @Router /volunteers/{id} [put]
func (h *VolunteerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.VolunteerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	volunteer, err := h.volunteerUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volunteer)
}

// Delete godoc
// @Summary Delete volunteer
// @Description Delete volunteer by ID
// @Tags volunteers
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /volunteers/{id} [delete]
func (h *VolunteerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.volunteerUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "volunteer deleted successfully",
	})
}

// AddTraining godoc
// @Summary Add training
// @Description Add training record to volunteer
// @Tags volunteers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Param request body entities.Training true "Training data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /volunteers/{id}/trainings [post]
func (h *VolunteerHandler) AddTraining(c *gin.Context) {
	id := c.Param("id")

	var training entities.Training
	if err := c.ShouldBindJSON(&training); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	addedBy, _ := c.Get("user_id")
	err := h.volunteerUseCase.AddTraining(c.Request.Context(), id, training, addedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training added successfully",
	})
}

// LogHours godoc
// @Summary Log volunteer hours
// @Description Log hours worked by volunteer
// @Tags volunteers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entities.VolunteerHourCreateRequest true "Hours data"
// @Success 201 {object} entities.VolunteerHour
// @Failure 400 {object} ErrorResponse
// @Router /volunteers/hours [post]
func (h *VolunteerHandler) LogHours(c *gin.Context) {
	var req entities.VolunteerHourCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	loggedBy, _ := c.Get("user_id")
	hour, err := h.volunteerUseCase.LogHours(c.Request.Context(), &req, loggedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hour)
}

// GetVolunteerHours godoc
// @Summary Get volunteer hours
// @Description Get hours logged by a volunteer
// @Tags volunteers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /volunteers/{id}/hours [get]
func (h *VolunteerHandler) GetVolunteerHours(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit := parseIntQuery(c.Query("limit"), 10)
	offset := parseIntQuery(c.Query("offset"), 0)

	hours, total, err := h.volunteerUseCase.GetVolunteerHours(c.Request.Context(), id, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   hours,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetActive godoc
// @Summary Get active volunteers
// @Description Get list of active volunteers
// @Tags volunteers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entities.Volunteer
// @Failure 500 {object} ErrorResponse
// @Router /volunteers/active [get]
func (h *VolunteerHandler) GetActive(c *gin.Context) {
	volunteers, err := h.volunteerUseCase.GetActiveVolunteers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volunteers)
}

// GetStatistics godoc
// @Summary Get volunteer statistics
// @Description Get statistics for a volunteer
// @Tags volunteers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Volunteer ID"
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /volunteers/{id}/statistics [get]
func (h *VolunteerHandler) GetStatistics(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.volunteerUseCase.GetVolunteerStatistics(c.Request.Context(), id, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
