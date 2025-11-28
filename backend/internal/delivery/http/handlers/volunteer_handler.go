package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/volunteer"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VolunteerHandler handles volunteer-related HTTP requests
type VolunteerHandler struct {
	volunteerUseCase volunteer.VolunteerUseCase
}

// NewVolunteerHandler creates a new volunteer handler
func NewVolunteerHandler(volunteerUseCase volunteer.VolunteerUseCase) *VolunteerHandler {
	return &VolunteerHandler{
		volunteerUseCase: volunteerUseCase,
	}
}

// CreateVolunteer creates a new volunteer
// @Summary Create a new volunteer
// @Tags volunteers
// @Accept json
// @Produce json
// @Param volunteer body entities.Volunteer true "Volunteer data"
// @Success 201 {object} entities.Volunteer
// @Router /volunteers [post]
func (h *VolunteerHandler) CreateVolunteer(c *gin.Context) {
	var volunteer entities.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.CreateVolunteer(c.Request.Context(), &volunteer, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, volunteer)
}

// GetVolunteer gets a volunteer by ID
// @Summary Get volunteer by ID
// @Tags volunteers
// @Produce json
// @Param id path string true "Volunteer ID"
// @Success 200 {object} entities.Volunteer
// @Router /volunteers/{id} [get]
func (h *VolunteerHandler) GetVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	volunteer, err := h.volunteerUseCase.GetVolunteer(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, volunteer)
}

// UpdateVolunteer updates a volunteer
// @Summary Update volunteer
// @Tags volunteers
// @Accept json
// @Produce json
// @Param id path string true "Volunteer ID"
// @Param volunteer body entities.Volunteer true "Volunteer data"
// @Success 200 {object} entities.Volunteer
// @Router /volunteers/{id} [put]
func (h *VolunteerHandler) UpdateVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var volunteer entities.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteer.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.UpdateVolunteer(c.Request.Context(), &volunteer, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, volunteer)
}

// DeleteVolunteer deletes a volunteer
// @Summary Delete volunteer
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Success 204
// @Router /volunteers/{id} [delete]
func (h *VolunteerHandler) DeleteVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.DeleteVolunteer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListVolunteers lists volunteers with filters
// @Summary List volunteers
// @Tags volunteers
// @Produce json
// @Param status query string false "Volunteer status"
// @Param skill query string false "Skill"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers [get]
func (h *VolunteerHandler) ListVolunteers(c *gin.Context) {
	filter := &repositories.VolunteerFilter{
		Status:    string(entities.VolunteerStatus(c.Query("status"))),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "last_name"),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
	}

	// Parse skills
	if skillStr := c.Query("skill"); skillStr != "" {
		filter.Skills = []string{skillStr}
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

	volunteers, total, err := h.volunteerUseCase.ListVolunteers(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"volunteers": volunteers,
		"total":      total,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
	})
}

// GetActiveVolunteers gets all active volunteers
// @Summary Get active volunteers
// @Tags volunteers
// @Produce json
// @Success 200 {array} entities.Volunteer
// @Router /volunteers/active [get]
func (h *VolunteerHandler) GetActiveVolunteers(c *gin.Context) {
	volunteers, err := h.volunteerUseCase.GetActiveVolunteers(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"volunteers": volunteers})
}

// GetVolunteersBySkill gets volunteers with a specific skill
// @Summary Get volunteers by skill
// @Tags volunteers
// @Produce json
// @Param skill query string true "Skill name"
// @Success 200 {array} entities.Volunteer
// @Router /volunteers/by-skill [get]
func (h *VolunteerHandler) GetVolunteersBySkill(c *gin.Context) {
	skill := c.Query("skill")
	if skill == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill parameter is required"})
		return
	}

	volunteers, err := h.volunteerUseCase.GetVolunteersBySkill(c.Request.Context(), skill)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"volunteers": volunteers})
}

// GetVolunteersNeedingBackgroundCheck gets volunteers needing background check
// @Summary Get volunteers needing background check
// @Tags volunteers
// @Produce json
// @Success 200 {array} entities.Volunteer
// @Router /volunteers/needing-background-check [get]
func (h *VolunteerHandler) GetVolunteersNeedingBackgroundCheck(c *gin.Context) {
	volunteers, err := h.volunteerUseCase.GetVolunteersNeedingBackgroundCheck(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"volunteers": volunteers})
}

// GetTopVolunteers gets top volunteers by hours
// @Summary Get top volunteers
// @Tags volunteers
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Success 200 {array} entities.Volunteer
// @Router /volunteers/top [get]
func (h *VolunteerHandler) GetTopVolunteers(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	volunteers, err := h.volunteerUseCase.GetTopVolunteers(c.Request.Context(), limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"volunteers": volunteers})
}

// ApproveVolunteer approves a volunteer application
// @Summary Approve volunteer
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/approve [post]
func (h *VolunteerHandler) ApproveVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.ApproveVolunteer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer approved successfully"})
}

// SuspendVolunteer suspends a volunteer
// @Summary Suspend volunteer
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Param suspension body map[string]string true "Suspension reason"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/suspend [post]
func (h *VolunteerHandler) SuspendVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.SuspendVolunteer(c.Request.Context(), id, req.Reason, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer suspended successfully"})
}


// AddCommendation adds a commendation to a volunteer
// @Summary Add commendation
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Param commendation body map[string]string true "Commendation note"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/commendation [post]
func (h *VolunteerHandler) AddCommendation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var req struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.AddCommendation(c.Request.Context(), id, req.Note, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Commendation added successfully"})
}

// AddWarning adds a warning to a volunteer
// @Summary Add warning
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Param warning body map[string]string true "Warning reason"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/warning [post]
func (h *VolunteerHandler) AddWarning(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.AddWarning(c.Request.Context(), id, req.Reason, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warning added successfully"})
}

// AddCertification adds a certification to a volunteer
// @Summary Add certification
// @Tags volunteers
// @Param id path string true "Volunteer ID"
// @Param certification body entities.Certification true "Certification data"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/certification [post]
func (h *VolunteerHandler) AddCertification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var cert entities.Certification
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.AddCertification(c.Request.Context(), id, cert, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certification added successfully"})
}

// GetVolunteerStatistics gets volunteer statistics
// @Summary Get volunteer statistics
// @Tags volunteers
// @Produce json
// @Success 200 {object} repositories.VolunteerStatistics
// @Router /volunteers/statistics [get]
func (h *VolunteerHandler) GetVolunteerStatistics(c *gin.Context) {
	stats, err := h.volunteerUseCase.GetVolunteerStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// LogVolunteerHours logs volunteer hours
// @Summary Log volunteer hours
// @Tags volunteers
// @Accept json
// @Produce json
// @Param id path string true "Volunteer ID"
// @Param log body struct{ Hours float64 `json:"hours"`; Notes string `json:"notes"` } true "Log data"
// @Success 200 {object} entities.Volunteer
// @Router /volunteers/{id}/log-hours [post]
func (h *VolunteerHandler) LogVolunteerHours(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var req struct {
		Hours float64 `json:"hours" binding:"required"`
		Notes string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	volunteer, err := h.volunteerUseCase.LogHours(c.Request.Context(), id, req.Hours, req.Notes, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, volunteer)
}

// GetVolunteerHours gets volunteer hours
// @Summary Get volunteer hours
// @Tags volunteers
// @Produce json
// @Param id path string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/hours [get]
func (h *VolunteerHandler) GetVolunteerHours(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	hours, err := h.volunteerUseCase.GetVolunteerHours(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_hours": hours})
}

// ActivateVolunteer activates a volunteer
// @Summary Activate a volunteer
// @Tags volunteers
// @Produce json
// @Param id path string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/activate [post]
func (h *VolunteerHandler) ActivateVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.ActivateVolunteer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer activated"})
}

// DeactivateVolunteer deactivates a volunteer
// @Summary Deactivate a volunteer
// @Tags volunteers
// @Produce json
// @Param id path string true "Volunteer ID"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/{id}/deactivate [post]
func (h *VolunteerHandler) DeactivateVolunteer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.volunteerUseCase.DeactivateVolunteer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer deactivated"})
}

// GetVolunteersByRole gets volunteers by role
// @Summary Get volunteers by role
// @Tags volunteers
// @Produce json
// @Param role path string true "Volunteer Role"
// @Success 200 {object} map[string]interface{}
// @Router /volunteers/role/{role} [get]
func (h *VolunteerHandler) GetVolunteersByRole(c *gin.Context) {
	roleParam := c.Param("role")
	role := entities.VolunteerRole(roleParam)

	if !role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer role"})
		return
	}

	volunteers, total, err := h.volunteerUseCase.GetVolunteersByRole(c.Request.Context(), role)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"volunteers": volunteers, "total": total})
}
