package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/partner"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerHandler handles partner-related HTTP requests
type PartnerHandler struct {
	partnerUseCase *partner.PartnerUseCase
}

// NewPartnerHandler creates a new partner handler
func NewPartnerHandler(partnerUseCase *partner.PartnerUseCase) *PartnerHandler {
	return &PartnerHandler{
		partnerUseCase: partnerUseCase,
	}
}

// CreatePartner creates a new partner
func (h *PartnerHandler) CreatePartner(c *gin.Context) {
	var partnerReq entities.Partner
	if err := c.ShouldBindJSON(&partnerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.CreatePartner(c.Request.Context(), &partnerReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, partnerReq)
}

// GetPartner gets a partner by ID
func (h *PartnerHandler) GetPartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	partner, err := h.partnerUseCase.GetPartnerByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partner)
}

// UpdatePartner updates a partner
func (h *PartnerHandler) UpdatePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var partnerReq entities.Partner
	if err := c.ShouldBindJSON(&partnerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	partnerReq.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.UpdatePartner(c.Request.Context(), &partnerReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partnerReq)
}

// DeletePartner deletes a partner
func (h *PartnerHandler) DeletePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.DeletePartner(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deleted successfully"})
}

// ListPartners lists partners with filtering
func (h *PartnerHandler) ListPartners(c *gin.Context) {
	filter := &repositories.PartnerFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.Status = c.Query("status")
	filter.Search = c.Query("search")

	if acceptsIntakesStr := c.Query("accepts_intakes"); acceptsIntakesStr != "" {
		acceptsIntakes := acceptsIntakesStr == "true"
		filter.AcceptsIntakes = &acceptsIntakes
	}

	if hasCapacityStr := c.Query("has_capacity"); hasCapacityStr != "" {
		hasCapacity := hasCapacityStr == "true"
		filter.HasCapacity = &hasCapacity
	}

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "name")
	filter.SortOrder = c.DefaultQuery("sort_order", "asc")

	partners, total, err := h.partnerUseCase.ListPartners(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   partners,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetPartnersByType gets partners by type
func (h *PartnerHandler) GetPartnersByType(c *gin.Context) {
	partnerType := entities.PartnerType(c.Param("type"))

	partners, err := h.partnerUseCase.GetPartnersByType(c.Request.Context(), partnerType)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partners)
}

// GetPartnersByStatus gets partners by status
func (h *PartnerHandler) GetPartnersByStatus(c *gin.Context) {
	status := entities.PartnerStatus(c.Param("status"))

	partners, err := h.partnerUseCase.GetPartnersByStatus(c.Request.Context(), status)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partners)
}

// GetActivePartners gets all active partners
func (h *PartnerHandler) GetActivePartners(c *gin.Context) {
	partners, err := h.partnerUseCase.GetActivePartners(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partners)
}

// GetPartnersWithCapacity gets partners with available capacity
func (h *PartnerHandler) GetPartnersWithCapacity(c *gin.Context) {
	partners, err := h.partnerUseCase.GetPartnersWithCapacity(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, partners)
}

// GetPartnerStatistics gets partner statistics
func (h *PartnerHandler) GetPartnerStatistics(c *gin.Context) {
	stats, err := h.partnerUseCase.GetPartnerStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ActivatePartner activates a partner
func (h *PartnerHandler) ActivatePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.ActivatePartner(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner activated successfully"})
}

// SuspendPartner suspends a partner
func (h *PartnerHandler) SuspendPartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.SuspendPartner(c.Request.Context(), id, userID, req.Reason); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner suspended successfully"})
}

// DeactivatePartner deactivates a partner
func (h *PartnerHandler) DeactivatePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.DeactivatePartner(c.Request.Context(), id, userID, req.Reason); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deactivated successfully"})
}

// AddRating adds a rating to a partner
func (h *PartnerHandler) AddRating(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		Rating float64 `json:"rating" binding:"required,min=0,max=5"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.AddRating(c.Request.Context(), id, userID, req.Rating); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating added successfully"})
}

// UpdateCapacity updates the current capacity of a partner
func (h *PartnerHandler) UpdateCapacity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		CurrentCapacity int `json:"current_capacity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.UpdateCapacity(c.Request.Context(), id, userID, req.CurrentCapacity); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Capacity updated successfully"})
}

// SetAgreementExpiry sets the agreement expiry date for a partner
func (h *PartnerHandler) SetAgreementExpiry(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		ExpiryDate string `json:"expiry_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiryDate, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.SetAgreementExpiry(c.Request.Context(), id, userID, expiryDate); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agreement expiry date set successfully"})
}


// GetPartnersAcceptingIntakes gets all partners currently accepting intakes
func (h *PartnerHandler) GetPartnersAcceptingIntakes(c *gin.Context) {
	filter := &repositories.PartnerFilter{}

	// Set the primary filter for this endpoint
	acceptsIntakes := true
	filter.AcceptsIntakes = &acceptsIntakes

	// Check for optional capacity filter
	if hasCapacityStr := c.Query("has_capacity"); hasCapacityStr != "" {
		hasCapacity := hasCapacityStr == "true"
		filter.HasCapacity = &hasCapacity
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "name")
	filter.SortOrder = c.DefaultQuery("sort_order", "asc")

	partners, total, err := h.partnerUseCase.ListPartners(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   partners,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}


// RatePartner rates a partner
func (h *PartnerHandler) RatePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		Rating float64 `json:"rating" binding:"required,min=0,max=5"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.AddRating(c.Request.Context(), id, userID, req.Rating); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating added successfully"})
}

// GetPartnerStatisticsDetail gets detailed statistics for a partner
func (h *PartnerHandler) GetPartnerStatisticsDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	partner, err := h.partnerUseCase.GetPartnerByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// For simplicity, we are not fetching detailed transfer stats here.
	// A more advanced implementation could fetch transfer data.
	response := gin.H{
		"partner_id":            partner.ID,
		"name":                  partner.Name,
		"status":                partner.Status,
		"accepts_intakes":       partner.AcceptsIntakes,
		"current_capacity":      partner.CurrentCapacity,
		"max_capacity":          partner.MaxCapacity,
		"average_rating":        partner.Rating,
		"total_ratings":         partner.TotalRatings,
		"total_transfers_in":    partner.TotalTransfersIn,
		"total_transfers_out":   partner.TotalTransfersOut,
		"successful_placements": partner.SuccessfulPlacements,
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePartnerCapacity updates partner capacity
func (h *PartnerHandler) UpdatePartnerCapacity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req struct {
		CurrentCapacity int `json:"current_capacity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.partnerUseCase.UpdateCapacity(c.Request.Context(), id, userID, req.CurrentCapacity); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Capacity updated successfully"})
}

