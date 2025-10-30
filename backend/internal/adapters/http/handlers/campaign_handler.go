package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type CampaignHandler struct {
	campaignUseCase *usecases.CampaignUseCase
}

func NewCampaignHandler(campaignUseCase *usecases.CampaignUseCase) *CampaignHandler {
	return &CampaignHandler{
		campaignUseCase: campaignUseCase,
	}
}

func (h *CampaignHandler) Create(c *gin.Context) {
	var req entities.CampaignCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	campaign, err := h.campaignUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

func (h *CampaignHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	campaign, err := h.campaignUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func (h *CampaignHandler) List(c *gin.Context) {
	filter := &entities.CampaignFilter{
		Type:   entities.CampaignType(c.Query("type")),
		Status: entities.CampaignStatus(c.Query("status")),
		Limit:  parseIntQuery(c.Query("limit"), 10),
		Offset: parseIntQuery(c.Query("offset"), 0),
	}

	campaigns, total, err := h.campaignUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   campaigns,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *CampaignHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.CampaignUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	campaign, err := h.campaignUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func (h *CampaignHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.campaignUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "campaign deleted successfully"})
}

func (h *CampaignHandler) UpdateProgress(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		AmountRaised      float64 `json:"amount_raised"`
		ParticipantsCount int     `json:"participants_count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.campaignUseCase.UpdateProgress(c.Request.Context(), id, req.AmountRaised, req.ParticipantsCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress updated successfully"})
}

func (h *CampaignHandler) GetActive(c *gin.Context) {
	campaigns, err := h.campaignUseCase.GetActiveCampaigns(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (h *CampaignHandler) GetStatistics(c *gin.Context) {
	id := c.Param("id")

	stats, err := h.campaignUseCase.GetCampaignStatistics(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *CampaignHandler) GetAllStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_date and end_date are required"})
		return
	}

	stats, err := h.campaignUseCase.GetAllCampaignsStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
