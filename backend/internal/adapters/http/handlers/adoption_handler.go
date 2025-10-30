package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type AdoptionHandler struct {
	adoptionUseCase *usecases.AdoptionUseCase
}

func NewAdoptionHandler(adoptionUseCase *usecases.AdoptionUseCase) *AdoptionHandler {
	return &AdoptionHandler{
		adoptionUseCase: adoptionUseCase,
	}
}

// Create godoc
// @Summary Create adoption application
// @Description Create a new adoption application
// @Tags adoptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entities.AdoptionCreateRequest true "Adoption data"
// @Success 201 {object} entities.Adoption
// @Failure 400 {object} ErrorResponse
// @Router /adoptions [post]
func (h *AdoptionHandler) Create(c *gin.Context) {
	var req entities.AdoptionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	adoption, err := h.adoptionUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, adoption)
}

// GetByID godoc
// @Summary Get adoption by ID
// @Description Get adoption details by ID
// @Tags adoptions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Success 200 {object} entities.Adoption
// @Failure 404 {object} ErrorResponse
// @Router /adoptions/{id} [get]
func (h *AdoptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	adoption, err := h.adoptionUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, adoption)
}

// List godoc
// @Summary List adoptions
// @Description Get list of adoptions with filtering and pagination
// @Tags adoptions
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param animal_id query string false "Filter by animal ID"
// @Param applicant_id query string false "Filter by applicant ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /adoptions [get]
func (h *AdoptionHandler) List(c *gin.Context) {
	filter := &entities.AdoptionFilter{
		Status:      entities.AdoptionStatus(c.Query("status")),
		AnimalID:    c.Query("animal_id"),
		ApplicantID: c.Query("applicant_id"),
		Limit:       parseIntQuery(c.Query("limit"), 10),
		Offset:      parseIntQuery(c.Query("offset"), 0),
	}

	adoptions, total, err := h.adoptionUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   adoptions,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Update godoc
// @Summary Update adoption
// @Description Update adoption details
// @Tags adoptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Param request body entities.AdoptionUpdateRequest true "Update data"
// @Success 200 {object} entities.Adoption
// @Failure 400 {object} ErrorResponse
// @Router /adoptions/{id} [put]
func (h *AdoptionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.AdoptionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	adoption, err := h.adoptionUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, adoption)
}

// Delete godoc
// @Summary Delete adoption
// @Description Delete adoption by ID
// @Tags adoptions
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /adoptions/{id} [delete]
func (h *AdoptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.adoptionUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "adoption deleted successfully",
	})
}

// Approve godoc
// @Summary Approve adoption
// @Description Approve an adoption application
// @Tags adoptions
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /adoptions/{id}/approve [post]
func (h *AdoptionHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	processedBy, _ := c.Get("user_id")
	err := h.adoptionUseCase.Approve(c.Request.Context(), id, processedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "adoption approved successfully",
	})
}

// Reject godoc
// @Summary Reject adoption
// @Description Reject an adoption application
// @Tags adoptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Param request body map[string]string true "Rejection reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /adoptions/{id}/reject [post]
func (h *AdoptionHandler) Reject(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	processedBy, _ := c.Get("user_id")
	err := h.adoptionUseCase.Reject(c.Request.Context(), id, req.Reason, processedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "adoption rejected successfully",
	})
}

// Complete godoc
// @Summary Complete adoption
// @Description Mark adoption as completed
// @Tags adoptions
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /adoptions/{id}/complete [post]
func (h *AdoptionHandler) Complete(c *gin.Context) {
	id := c.Param("id")

	completedBy, _ := c.Get("user_id")
	err := h.adoptionUseCase.Complete(c.Request.Context(), id, completedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "adoption completed successfully",
	})
}
