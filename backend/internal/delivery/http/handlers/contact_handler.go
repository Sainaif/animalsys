package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/contact"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContactHandler handles HTTP endpoints
type ContactHandler struct {
	useCase *contact.UseCase
}

// NewContactHandler creates handler
func NewContactHandler(useCase *contact.UseCase) *ContactHandler {
	return &ContactHandler{useCase: useCase}
}

// List contacts
func (h *ContactHandler) List(c *gin.Context) {
	var filter repositories.ContactFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contacts, total, err := h.useCase.List(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"total":    total,
		"limit":    filter.Limit,
		"offset":   filter.Offset,
	})
}

// Get contact
func (h *ContactHandler) Get(c *gin.Context) {
	contactID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact id"})
		return
	}

	contact, err := h.useCase.Get(c.Request.Context(), contactID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, contact)
}

// Create contact
func (h *ContactHandler) Create(c *gin.Context) {
	var req contact.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newContact, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newContact)
}

// Update contact
func (h *ContactHandler) Update(c *gin.Context) {
	contactID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact id"})
		return
	}

	var req contact.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.useCase.Update(c.Request.Context(), contactID, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete contact
func (h *ContactHandler) Delete(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}
	_ = userID // placeholder for future audit logs

	contactID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact id"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), contactID); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
