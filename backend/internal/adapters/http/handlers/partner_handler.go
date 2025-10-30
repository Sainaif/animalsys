package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type PartnerHandler struct {
	partnerUseCase *usecases.PartnerUseCase
}

func NewPartnerHandler(partnerUseCase *usecases.PartnerUseCase) *PartnerHandler {
	return &PartnerHandler{
		partnerUseCase: partnerUseCase,
	}
}

func (h *PartnerHandler) Create(c *gin.Context) {
	var req entities.PartnerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	partner, err := h.partnerUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, partner)
}

func (h *PartnerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	partner, err := h.partnerUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, partner)
}

func (h *PartnerHandler) List(c *gin.Context) {
	partnerType := entities.PartnerType(c.Query("type"))
	limit := parseIntQuery(c.Query("limit"), 10)
	offset := parseIntQuery(c.Query("offset"), 0)

	partners, total, err := h.partnerUseCase.List(c.Request.Context(), partnerType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   partners,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *PartnerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.PartnerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	partner, err := h.partnerUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, partner)
}

func (h *PartnerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.partnerUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "partner deleted successfully"})
}

func (h *PartnerHandler) GetActive(c *gin.Context) {
	partners, err := h.partnerUseCase.GetActivePartners(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, partners)
}

func (h *PartnerHandler) AddCollaboration(c *gin.Context) {
	id := c.Param("id")

	var req entities.CollaborationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	addedBy, _ := c.Get("user_id")
	err := h.partnerUseCase.AddCollaboration(c.Request.Context(), id, &req, addedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "collaboration added successfully"})
}

func (h *PartnerHandler) GetStatistics(c *gin.Context) {
	stats, err := h.partnerUseCase.GetPartnerStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
