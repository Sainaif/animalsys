package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/donor"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorHandler handles donor-related HTTP requests
type DonorHandler struct {
	donorUseCase donor.DonorUseCaseInterface
}

// NewDonorHandler creates a new donor handler
func NewDonorHandler(donorUseCase donor.DonorUseCaseInterface) *DonorHandler {
	return &DonorHandler{
		donorUseCase: donorUseCase,
	}
}

// CreateDonor creates a new donor
// @Summary Create a new donor
// @Tags donors
// @Accept json
// @Produce json
// @Param donor body entities.Donor true "Donor data"
// @Success 201 {object} entities.Donor
// @Router /donors [post]
func (h *DonorHandler) CreateDonor(c *gin.Context) {
	var donor entities.Donor
	if err := c.ShouldBindJSON(&donor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donorUseCase.CreateDonor(c.Request.Context(), &donor, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, donor)
}

// GetDonor gets a donor by ID
// @Summary Get donor by ID
// @Tags donors
// @Produce json
// @Param id path string true "Donor ID"
// @Success 200 {object} entities.Donor
// @Router /donors/{id} [get]
func (h *DonorHandler) GetDonor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	donor, err := h.donorUseCase.GetDonor(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, donor)
}

// UpdateDonor updates a donor
// @Summary Update donor
// @Tags donors
// @Accept json
// @Produce json
// @Param id path string true "Donor ID"
// @Param donor body entities.Donor true "Donor data"
// @Success 200 {object} entities.Donor
// @Router /donors/{id} [put]
func (h *DonorHandler) UpdateDonor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	var donor entities.Donor
	if err := c.ShouldBindJSON(&donor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donor.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donorUseCase.UpdateDonor(c.Request.Context(), &donor, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, donor)
}

// DeleteDonor deletes a donor
// @Summary Delete donor
// @Tags donors
// @Param id path string true "Donor ID"
// @Success 204
// @Router /donors/{id} [delete]
func (h *DonorHandler) DeleteDonor(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donorUseCase.DeleteDonor(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListDonors lists donors with filters
// @Summary List donors
// @Tags donors
// @Produce json
// @Param type query string false "Donor type"
// @Param status query string false "Donor status"
// @Param search query string false "Search query"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /donors [get]
func (h *DonorHandler) ListDonors(c *gin.Context) {
	filter := &repositories.DonorFilter{
		Type:      string(entities.DonorType(c.Query("type"))),
		Status:    string(entities.DonorStatus(c.Query("status"))),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
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

	// Parse donation amount range
	if minStr := c.Query("total_donated_min"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.MinTotalDonated = &min
		}
	}
	if maxStr := c.Query("total_donated_max"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.MaxTotalDonated = &max
		}
	}

	donors, total, err := h.donorUseCase.ListDonors(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"donors": donors,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetMajorDonors gets all major donors
// @Summary Get major donors
// @Tags donors
// @Produce json
// @Success 200 {array} entities.Donor
// @Router /donors/major [get]
func (h *DonorHandler) GetMajorDonors(c *gin.Context) {
	donors, err := h.donorUseCase.GetMajorDonors(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donors": donors})
}

// GetLapsedDonors gets lapsed donors
// @Summary Get lapsed donors
// @Tags donors
// @Produce json
// @Param days query int false "Days since last donation (default: 365)"
// @Success 200 {array} entities.Donor
// @Router /donors/lapsed [get]
func (h *DonorHandler) GetLapsedDonors(c *gin.Context) {
	days := 365
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	donors, err := h.donorUseCase.GetLapsedDonors(c.Request.Context(), days)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donors": donors, "days": days})
}

// GetDonorStatistics gets donor statistics
// @Summary Get donor statistics
// @Tags donors
// @Produce json
// @Success 200 {object} repositories.DonorStatistics
// @Router /donors/statistics [get]
func (h *DonorHandler) GetDonorStatistics(c *gin.Context) {
	stats, err := h.donorUseCase.GetDonorStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateDonorEngagement updates donor engagement metrics
// @Summary Update donor engagement
// @Tags donors
// @Accept json
// @Produce json
// @Param id path string true "Donor ID"
// @Param engagement body map[string]int true "Engagement data"
// @Success 200 {object} map[string]interface{}
// @Router /donors/{id}/engagement [post]
func (h *DonorHandler) UpdateDonorEngagement(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	var req struct {
		VolunteerHours int `json:"volunteer_hours"`
		EventsAttended int `json:"events_attended"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donorUseCase.UpdateDonorEngagement(c.Request.Context(), id, req.VolunteerHours, req.EventsAttended, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Engagement updated successfully"})
}

// GetTopDonors gets the top donors by donation amount
func (h *DonorHandler) GetTopDonors(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	if limit > 100 { // Max limit
		limit = 100
	}

	donors, err := h.donorUseCase.GetTopDonors(c.Request.Context(), limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"donors": donors,
		"limit":  limit,
		"total":  len(donors),
	})
}

// GetRecurringDonors gets all recurring donors
func (h *DonorHandler) GetRecurringDonors(c *gin.Context) {
	donors, err := h.donorUseCase.GetRecurringDonors(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"donors": donors,
		"total":  len(donors),
	})
}

// HandleError handles application errors
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

// UpdateCommunicationPreferences updates donor communication preferences
func (h *DonorHandler) UpdateCommunicationPreferences(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	var prefs entities.DonorPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donorUseCase.UpdateCommunicationPreferences(c.Request.Context(), id, prefs, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated"})
}
