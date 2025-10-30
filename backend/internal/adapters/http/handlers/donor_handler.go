package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type DonorHandler struct {
	donorUseCase *usecases.DonorUseCase
}

func NewDonorHandler(donorUseCase *usecases.DonorUseCase) *DonorHandler {
	return &DonorHandler{
		donorUseCase: donorUseCase,
	}
}

// Donor endpoints

func (h *DonorHandler) CreateDonor(c *gin.Context) {
	var req entities.DonorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	donor, err := h.donorUseCase.CreateDonor(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, donor)
}

func (h *DonorHandler) GetDonorByID(c *gin.Context) {
	id := c.Param("id")

	donor, err := h.donorUseCase.GetDonorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donor)
}

func (h *DonorHandler) ListDonors(c *gin.Context) {
	filter := &entities.DonorFilter{
		Type:   entities.DonorType(c.Query("type")),
		Search: c.Query("search"),
		Limit:  parseIntQuery(c.Query("limit"), 10),
		Offset: parseIntQuery(c.Query("offset"), 0),
	}

	donors, total, err := h.donorUseCase.ListDonors(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   donors,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *DonorHandler) UpdateDonor(c *gin.Context) {
	id := c.Param("id")

	var req entities.DonorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	donor, err := h.donorUseCase.UpdateDonor(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donor)
}

func (h *DonorHandler) DeleteDonor(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.donorUseCase.DeleteDonor(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "donor deleted successfully"})
}

func (h *DonorHandler) GetTopDonors(c *gin.Context) {
	limit := parseIntQuery(c.Query("limit"), 10)

	donors, err := h.donorUseCase.GetTopDonors(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donors)
}

// Donation endpoints

func (h *DonorHandler) RecordDonation(c *gin.Context) {
	var req entities.DonationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	recordedBy, _ := c.Get("user_id")
	donation, err := h.donorUseCase.RecordDonation(c.Request.Context(), &req, recordedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, donation)
}

func (h *DonorHandler) GetDonationByID(c *gin.Context) {
	id := c.Param("id")

	donation, err := h.donorUseCase.GetDonationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (h *DonorHandler) ListDonations(c *gin.Context) {
	filter := &entities.DonationFilter{
		DonorID:   c.Query("donor_id"),
		Type:      entities.DonationType(c.Query("type")),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Limit:     parseIntQuery(c.Query("limit"), 10),
		Offset:    parseIntQuery(c.Query("offset"), 0),
	}

	donations, total, err := h.donorUseCase.ListDonations(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   donations,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *DonorHandler) UpdateDonation(c *gin.Context) {
	id := c.Param("id")

	var req entities.DonationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	donation, err := h.donorUseCase.UpdateDonation(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (h *DonorHandler) DeleteDonation(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.donorUseCase.DeleteDonation(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "donation deleted successfully"})
}

func (h *DonorHandler) GetDonationStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_date and end_date are required"})
		return
	}

	stats, err := h.donorUseCase.GetDonationStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
