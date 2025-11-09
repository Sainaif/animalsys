package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/donation"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonationHandler handles donation-related HTTP requests
type DonationHandler struct {
	donationUseCase *donation.DonationUseCase
}

// NewDonationHandler creates a new donation handler
func NewDonationHandler(donationUseCase *donation.DonationUseCase) *DonationHandler {
	return &DonationHandler{
		donationUseCase: donationUseCase,
	}
}

// CreateDonation creates a new donation
// @Summary Create a new donation
// @Tags donations
// @Accept json
// @Produce json
// @Param donation body entities.Donation true "Donation data"
// @Success 201 {object} entities.Donation
// @Router /donations [post]
func (h *DonationHandler) CreateDonation(c *gin.Context) {
	var donation entities.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.CreateDonation(c.Request.Context(), &donation, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, donation)
}

// ProcessDonation processes a donation
// @Summary Process a donation
// @Tags donations
// @Param id path string true "Donation ID"
// @Success 200 {object} map[string]interface{}
// @Router /donations/{id}/process [post]
func (h *DonationHandler) ProcessDonation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.ProcessDonation(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation processed successfully"})
}

// RefundDonation refunds a donation
// @Summary Refund a donation
// @Tags donations
// @Param id path string true "Donation ID"
// @Success 200 {object} map[string]interface{}
// @Router /donations/{id}/refund [post]
func (h *DonationHandler) RefundDonation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.RefundDonation(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation refunded successfully"})
}

// GetDonation gets a donation by ID
// @Summary Get donation by ID
// @Tags donations
// @Produce json
// @Param id path string true "Donation ID"
// @Success 200 {object} entities.Donation
// @Router /donations/{id} [get]
func (h *DonationHandler) GetDonation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	donation, err := h.donationUseCase.GetDonation(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, donation)
}

// UpdateDonation updates a donation
// @Summary Update donation
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Param donation body entities.Donation true "Donation data"
// @Success 200 {object} entities.Donation
// @Router /donations/{id} [put]
func (h *DonationHandler) UpdateDonation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	var donation entities.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.UpdateDonation(c.Request.Context(), &donation, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, donation)
}

// DeleteDonation deletes a donation
// @Summary Delete donation
// @Tags donations
// @Param id path string true "Donation ID"
// @Success 204
// @Router /donations/{id} [delete]
func (h *DonationHandler) DeleteDonation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.DeleteDonation(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListDonations lists donations with filters
// @Summary List donations
// @Tags donations
// @Produce json
// @Param donor_id query string false "Donor ID"
// @Param campaign_id query string false "Campaign ID"
// @Param type query string false "Donation type"
// @Param status query string false "Donation status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /donations [get]
func (h *DonationHandler) ListDonations(c *gin.Context) {
	filter := &repositories.DonationFilter{
		Type:          string(entities.DonationType(c.Query("type"))),
		Status:        string(entities.DonationStatus(c.Query("status"))),
		PaymentMethod: string(entities.PaymentMethodType(c.Query("payment_method"))),
		SortBy:        c.DefaultQuery("sort_by", "created_at"),
		SortOrder:     c.DefaultQuery("sort_order", "desc"),
	}

	// Parse donor ID
	if donorIDStr := c.Query("donor_id"); donorIDStr != "" {
		if donorID, err := primitive.ObjectIDFromHex(donorIDStr); err == nil {
			filter.DonorID = &donorID
		}
	}

	// Parse campaign ID
	if campaignIDStr := c.Query("campaign_id"); campaignIDStr != "" {
		if campaignID, err := primitive.ObjectIDFromHex(campaignIDStr); err == nil {
			filter.CampaignID = &campaignID
		}
	}

	// Parse is_recurring
	if recurringStr := c.Query("is_recurring"); recurringStr != "" {
		if recurring, err := strconv.ParseBool(recurringStr); err == nil {
			filter.IsRecurring = &recurring
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

	// Parse amount range
	if minStr := c.Query("amount_min"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			filter.MinAmount = &min
		}
	}
	if maxStr := c.Query("amount_max"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filter.MaxAmount = &max
		}
	}

	donations, total, err := h.donationUseCase.ListDonations(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"donations": donations,
		"total":     total,
		"limit":     filter.Limit,
		"offset":    filter.Offset,
	})
}

// GetDonationsByDonor gets all donations for a donor
// @Summary Get donations by donor
// @Tags donations
// @Produce json
// @Param id path string true "Donor ID"
// @Success 200 {array} entities.Donation
// @Router /donors/{id}/donations [get]
func (h *DonationHandler) GetDonationsByDonor(c *gin.Context) {
	donorIDParam := c.Param("id")
	donorID, err := primitive.ObjectIDFromHex(donorIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	donations, err := h.donationUseCase.GetDonationsByDonor(c.Request.Context(), donorID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}

// GetDonationsByCampaign gets all donations for a campaign
// @Summary Get donations by campaign
// @Tags donations
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {array} entities.Donation
// @Router /campaigns/{id}/donations [get]
func (h *DonationHandler) GetDonationsByCampaign(c *gin.Context) {
	campaignIDParam := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	donations, err := h.donationUseCase.GetDonationsByCampaign(c.Request.Context(), campaignID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}

// GetRecurringDonations gets all active recurring donations
// @Summary Get recurring donations
// @Tags donations
// @Produce json
// @Success 200 {array} entities.Donation
// @Router /donations/recurring [get]
func (h *DonationHandler) GetRecurringDonations(c *gin.Context) {
	donations, err := h.donationUseCase.GetRecurringDonations(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}

// GetPendingThankYous gets donations needing thank you notes
// @Summary Get pending thank you notes
// @Tags donations
// @Produce json
// @Success 200 {array} entities.Donation
// @Router /donations/pending-thank-yous [get]
func (h *DonationHandler) GetPendingThankYous(c *gin.Context) {
	donations, err := h.donationUseCase.GetPendingThankYous(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}

// SendThankYou marks thank you as sent
// @Summary Send thank you note
// @Tags donations
// @Param id path string true "Donation ID"
// @Success 200 {object} map[string]interface{}
// @Router /donations/{id}/thank-you [post]
func (h *DonationHandler) SendThankYou(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.SendThankYou(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thank you note sent successfully"})
}

// GenerateTaxReceipt generates a tax receipt
// @Summary Generate tax receipt
// @Tags donations
// @Accept json
// @Produce json
// @Param id path string true "Donation ID"
// @Param receipt body map[string]string true "Receipt data"
// @Success 200 {object} map[string]interface{}
// @Router /donations/{id}/tax-receipt [post]
func (h *DonationHandler) GenerateTaxReceipt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	var req struct {
		ReceiptURL string `json:"receipt_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.donationUseCase.GenerateTaxReceipt(c.Request.Context(), id, req.ReceiptURL, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tax receipt generated successfully"})
}

// GetDonationStatistics gets donation statistics
// @Summary Get donation statistics
// @Tags donations
// @Produce json
// @Success 200 {object} repositories.DonationStatistics
// @Router /donations/statistics [get]
func (h *DonationHandler) GetDonationStatistics(c *gin.Context) {
	stats, err := h.donationUseCase.GetDonationStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDonationsByDateRange gets donations within a date range
// @Summary Get donations by date range
// @Tags donations
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} entities.Donation
// @Router /donations/date-range [get]
func (h *DonationHandler) GetDonationsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	donations, err := h.donationUseCase.GetDonationsByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}

// GetDonationReceipt gets a donation receipt
func (h *DonationHandler) GetDonationReceipt(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"receipt": map[string]interface{}{}})
}

