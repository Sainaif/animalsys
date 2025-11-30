package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/campaign"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CampaignHandler handles campaign-related HTTP requests
type CampaignHandler struct {
	campaignUseCase campaign.ICampaignUseCase
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(campaignUseCase campaign.ICampaignUseCase) *CampaignHandler {
	return &CampaignHandler{
		campaignUseCase: campaignUseCase,
	}
}

// CreateCampaign creates a new campaign
// @Summary Create a new campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body entities.Campaign true "Campaign data"
// @Success 201 {object} entities.Campaign
// @Router /campaigns [post]
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var campaign entities.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		HandleError(c, errors.NewBadRequest(err.Error()))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.CreateCampaign(c.Request.Context(), &campaign, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// GetCampaign gets a campaign by ID
// @Summary Get campaign by ID
// @Tags campaigns
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} entities.Campaign
// @Router /campaigns/{id} [get]
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	campaign, err := h.campaignUseCase.GetCampaign(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// UpdateCampaign updates a campaign
// @Summary Update campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param campaign body entities.Campaign true "Campaign data"
// @Success 200 {object} entities.Campaign
// @Router /campaigns/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	var campaign entities.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		HandleError(c, errors.NewBadRequest(err.Error()))
		return
	}

	campaign.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.UpdateCampaign(c.Request.Context(), &campaign, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign deletes a campaign
// @Summary Delete campaign
// @Tags campaigns
// @Param id path string true "Campaign ID"
// @Success 204
// @Router /campaigns/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.DeleteCampaign(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCampaigns lists campaigns with filters
// @Summary List campaigns
// @Tags campaigns
// @Produce json
// @Param type query string false "Campaign type"
// @Param status query string false "Campaign status"
// @Param public query bool false "Public campaigns"
// @Param featured query bool false "Featured campaigns"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns [get]
func (h *CampaignHandler) ListCampaigns(c *gin.Context) {
	filter := &repositories.CampaignFilter{
		Type:      string(entities.CampaignType(c.Query("type"))),
		Status:    string(entities.CampaignStatus(c.Query("status"))),
		Search:    c.Query("search"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
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

	// Parse manager ID
	if managerIDStr := c.Query("manager_id"); managerIDStr != "" {
		if managerID, err := primitive.ObjectIDFromHex(managerIDStr); err == nil {
			filter.ManagerID = &managerID
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

	// Parse goal amount range
	if minStr := c.Query("goal_amount_min"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.GoalAmountMin = min
		}
	}
	if maxStr := c.Query("goal_amount_max"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.GoalAmountMax = max
		}
	}

	campaigns, total, err := h.campaignUseCase.ListCampaigns(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"total":     total,
		"limit":     filter.Limit,
		"offset":    filter.Offset,
	})
}

// GetActiveCampaigns gets all active campaigns
// @Summary Get active campaigns
// @Tags campaigns
// @Produce json
// @Success 200 {array} entities.Campaign
// @Router /campaigns/active [get]
func (h *CampaignHandler) GetActiveCampaigns(c *gin.Context) {
	campaigns, err := h.campaignUseCase.GetActiveCampaigns(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// GetFeaturedCampaigns gets all featured campaigns
// @Summary Get featured campaigns
// @Tags campaigns
// @Produce json
// @Success 200 {array} entities.Campaign
// @Router /campaigns/featured [get]
func (h *CampaignHandler) GetFeaturedCampaigns(c *gin.Context) {
	campaigns, err := h.campaignUseCase.GetFeaturedCampaigns(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// GetPublicCampaigns gets all public campaigns
// @Summary Get public campaigns
// @Tags campaigns
// @Produce json
// @Success 200 {array} entities.Campaign
// @Router /campaigns/public [get]
func (h *CampaignHandler) GetPublicCampaigns(c *gin.Context) {
	campaigns, err := h.campaignUseCase.GetPublicCampaigns(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// GetCampaignsByManager gets campaigns managed by a user
// @Summary Get campaigns by manager
// @Tags campaigns
// @Produce json
// @Param id path string true "Manager ID"
// @Success 200 {array} entities.Campaign
// @Router /users/{id}/campaigns [get]
func (h *CampaignHandler) GetCampaignsByManager(c *gin.Context) {
	managerID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid manager ID"))
		return
	}

	campaigns, err := h.campaignUseCase.GetCampaignsByManager(c.Request.Context(), managerID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// ActivateCampaign activates a campaign
// @Summary Activate campaign
// @Tags campaigns
// @Param id path string true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/activate [post]
func (h *CampaignHandler) ActivateCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.ActivateCampaign(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign activated successfully"})
}

// PauseCampaign pauses a campaign
// @Summary Pause campaign
// @Tags campaigns
// @Param id path string true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/pause [post]
func (h *CampaignHandler) PauseCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.PauseCampaign(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign paused successfully"})
}

// CompleteCampaign completes a campaign
// @Summary Complete campaign
// @Tags campaigns
// @Param id path string true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/complete [post]
func (h *CampaignHandler) CompleteCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.CompleteCampaign(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign completed successfully"})
}

// CancelCampaign cancels a campaign
// @Summary Cancel campaign
// @Tags campaigns
// @Param id path string true "Campaign ID"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/cancel [post]
func (h *CampaignHandler) CancelCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.CancelCampaign(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign cancelled successfully"})
}

// GetCampaignStatistics gets campaign statistics
// @Summary Get campaign statistics
// @Tags campaigns
// @Produce json
// @Success 200 {object} repositories.CampaignStatistics
// @Router /campaigns/statistics [get]
func (h *CampaignHandler) GetCampaignStatistics(c *gin.Context) {
	stats, err := h.campaignUseCase.GetCampaignStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCampaignDonors gets all donors for a campaign
// @Summary Get campaign donors
// @Tags campaigns
// @Produce json
// @Param id path string true "Campaign ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/donors [get]
func (h *CampaignHandler) GetCampaignDonors(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	donors, total, err := h.campaignUseCase.GetCampaignDonors(c.Request.Context(), id, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"donors": donors,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateCampaignAmount updates campaign raised amount
// @Summary Update campaign amount
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param amount body float64 true "New amount"
// @Success 200 {object} map[string]interface{}
// @Router /campaigns/{id}/update-amount [post]
func (h *CampaignHandler) UpdateCampaignAmount(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleError(c, errors.NewBadRequest("Invalid request payload: "+err.Error()))
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.campaignUseCase.UpdateCampaignAmount(c.Request.Context(), id, payload.Amount, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign amount updated successfully"})
}

// GetCampaignProgress gets campaign progress
// @Summary Get campaign progress
// @Tags campaigns
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} campaign.CampaignProgress
// @Router /campaigns/{id}/progress [get]
func (h *CampaignHandler) GetCampaignProgress(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	progress, err := h.campaignUseCase.GetCampaignProgress(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, progress)
}

// ShareCampaign shares a campaign
// @Summary Share campaign
// @Tags campaigns
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} campaign.CampaignShareable
// @Router /campaigns/{id}/share [post]
func (h *CampaignHandler) ShareCampaign(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleError(c, errors.NewBadRequest("Invalid campaign ID"))
		return
	}

	shareable, err := h.campaignUseCase.ShareCampaign(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, shareable)
}
