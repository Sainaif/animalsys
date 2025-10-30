package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type DocumentHandler struct {
	documentUseCase *usecases.DocumentUseCase
}

func NewDocumentHandler(documentUseCase *usecases.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{
		documentUseCase: documentUseCase,
	}
}

func (h *DocumentHandler) Create(c *gin.Context) {
	var req entities.DocumentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	uploadedBy, _ := c.Get("user_id")
	document, err := h.documentUseCase.Create(c.Request.Context(), &req, uploadedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (h *DocumentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	document, err := h.documentUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *DocumentHandler) List(c *gin.Context) {
	filter := &entities.DocumentFilter{
		EntityType: c.Query("entity_type"),
		EntityID:   c.Query("entity_id"),
		Type:       entities.DocumentType(c.Query("type")),
		Search:     c.Query("search"),
		Limit:      parseIntQuery(c.Query("limit"), 10),
		Offset:     parseIntQuery(c.Query("offset"), 0),
	}

	documents, total, err := h.documentUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   documents,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *DocumentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.DocumentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	document, err := h.documentUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.documentUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted successfully"})
}

func (h *DocumentHandler) GetExpiringSoon(c *gin.Context) {
	days := parseIntQuery(c.Query("days"), 30)

	documents, err := h.documentUseCase.GetExpiringSoon(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}
