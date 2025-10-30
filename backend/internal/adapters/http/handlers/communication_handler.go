package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type CommunicationHandler struct {
	communicationUseCase *usecases.CommunicationUseCase
}

func NewCommunicationHandler(communicationUseCase *usecases.CommunicationUseCase) *CommunicationHandler {
	return &CommunicationHandler{
		communicationUseCase: communicationUseCase,
	}
}

// Communication endpoints

func (h *CommunicationHandler) CreateCommunication(c *gin.Context) {
	var req entities.CommunicationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	communication, err := h.communicationUseCase.CreateCommunication(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, communication)
}

func (h *CommunicationHandler) GetCommunicationByID(c *gin.Context) {
	id := c.Param("id")

	communication, err := h.communicationUseCase.GetCommunicationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, communication)
}

func (h *CommunicationHandler) ListCommunications(c *gin.Context) {
	status := entities.CommunicationStatus(c.Query("status"))
	limit := parseIntQuery(c.Query("limit"), 10)
	offset := parseIntQuery(c.Query("offset"), 0)

	communications, total, err := h.communicationUseCase.ListCommunications(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   communications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CommunicationHandler) UpdateCommunication(c *gin.Context) {
	id := c.Param("id")

	var req entities.CommunicationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	communication, err := h.communicationUseCase.UpdateCommunication(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, communication)
}

func (h *CommunicationHandler) DeleteCommunication(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.communicationUseCase.DeleteCommunication(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "communication deleted successfully"})
}

func (h *CommunicationHandler) SendCommunication(c *gin.Context) {
	id := c.Param("id")

	sentBy, _ := c.Get("user_id")
	err := h.communicationUseCase.SendCommunication(c.Request.Context(), id, sentBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "communication sent successfully"})
}

func (h *CommunicationHandler) GetScheduled(c *gin.Context) {
	communications, err := h.communicationUseCase.GetScheduledCommunications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, communications)
}

func (h *CommunicationHandler) GetStatistics(c *gin.Context) {
	stats, err := h.communicationUseCase.GetCommunicationStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Template endpoints

func (h *CommunicationHandler) CreateTemplate(c *gin.Context) {
	var req entities.CommunicationTemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	template, err := h.communicationUseCase.CreateTemplate(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *CommunicationHandler) GetTemplateByID(c *gin.Context) {
	id := c.Param("id")

	template, err := h.communicationUseCase.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *CommunicationHandler) ListTemplates(c *gin.Context) {
	commType := entities.CommunicationType(c.Query("type"))

	templates, err := h.communicationUseCase.ListTemplates(c.Request.Context(), commType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *CommunicationHandler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")

	var req entities.CommunicationTemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	template, err := h.communicationUseCase.UpdateTemplate(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *CommunicationHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.communicationUseCase.DeleteTemplate(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted successfully"})
}

func (h *CommunicationHandler) CreateFromTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var req struct {
		Recipients []string          `json:"recipients" binding:"required"`
		Variables  map[string]string `json:"variables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	communication, err := h.communicationUseCase.CreateFromTemplate(c.Request.Context(), templateID, req.Recipients, req.Variables, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, communication)
}
