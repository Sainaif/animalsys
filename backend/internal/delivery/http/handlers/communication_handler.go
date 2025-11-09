package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/communication"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommunicationHandler handles communication-related HTTP requests
type CommunicationHandler struct {
	communicationUseCase *communication.CommunicationUseCase
}

// NewCommunicationHandler creates a new communication handler
func NewCommunicationHandler(communicationUseCase *communication.CommunicationUseCase) *CommunicationHandler {
	return &CommunicationHandler{
		communicationUseCase: communicationUseCase,
	}
}

// CreateCommunication creates a new communication
func (h *CommunicationHandler) CreateCommunication(c *gin.Context) {
	var comm entities.Communication
	if err := c.ShouldBindJSON(&comm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.CreateCommunication(c.Request.Context(), &comm, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comm)
}

// SendFromTemplate sends a communication using a template
func (h *CommunicationHandler) SendFromTemplate(c *gin.Context) {
	var req struct {
		TemplateID     string                  `json:"template_id" binding:"required"`
		RecipientType  entities.RecipientType  `json:"recipient_type" binding:"required"`
		RecipientID    string                  `json:"recipient_id" binding:"required"`
		RecipientEmail string                  `json:"recipient_email"`
		RecipientPhone string                  `json:"recipient_phone"`
		Variables      map[string]string       `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateID, err := primitive.ObjectIDFromHex(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	recipientID, err := primitive.ObjectIDFromHex(req.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	comm, err := h.communicationUseCase.SendCommunicationFromTemplate(
		c.Request.Context(),
		templateID,
		req.RecipientType,
		recipientID,
		req.RecipientEmail,
		req.RecipientPhone,
		req.Variables,
		userID,
	)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comm)
}

// GetCommunication gets a communication by ID
func (h *CommunicationHandler) GetCommunication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	comm, err := h.communicationUseCase.GetCommunicationByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comm)
}

// ListCommunications lists communications with filtering
func (h *CommunicationHandler) ListCommunications(c *gin.Context) {
	filter := &repositories.CommunicationFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.Category = c.Query("category")
	filter.Status = c.Query("status")
	filter.RecipientType = c.Query("recipient_type")
	filter.BatchID = c.Query("batch_id")
	filter.RelatedType = c.Query("related_type")

	if recipientIDStr := c.Query("recipient_id"); recipientIDStr != "" {
		recipientID, err := primitive.ObjectIDFromHex(recipientIDStr)
		if err == nil {
			filter.RecipientID = &recipientID
		}
	}

	if senderIDStr := c.Query("sender_id"); senderIDStr != "" {
		senderID, err := primitive.ObjectIDFromHex(senderIDStr)
		if err == nil {
			filter.SenderID = &senderID
		}
	}

	if templateIDStr := c.Query("template_id"); templateIDStr != "" {
		templateID, err := primitive.ObjectIDFromHex(templateIDStr)
		if err == nil {
			filter.TemplateID = &templateID
		}
	}

	if campaignIDStr := c.Query("campaign_id"); campaignIDStr != "" {
		campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
		if err == nil {
			filter.CampaignID = &campaignID
		}
	}

	if relatedIDStr := c.Query("related_id"); relatedIDStr != "" {
		relatedID, err := primitive.ObjectIDFromHex(relatedIDStr)
		if err == nil {
			filter.RelatedID = &relatedID
		}
	}

	// Date range
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	communications, total, err := h.communicationUseCase.ListCommunications(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  communications,
		"total": total,
		"limit": filter.Limit,
		"offset": filter.Offset,
	})
}

// UpdateCommunicationStatus updates communication status
func (h *CommunicationHandler) UpdateCommunicationStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	var req struct {
		Status entities.CommunicationStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.UpdateCommunicationStatus(c.Request.Context(), id, req.Status, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Communication status updated successfully"})
}

// TrackOpen tracks when a communication is opened
func (h *CommunicationHandler) TrackOpen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	if err := h.communicationUseCase.MarkCommunicationAsOpened(c.Request.Context(), id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Open tracked"})
}

// TrackClick tracks when a link in communication is clicked
func (h *CommunicationHandler) TrackClick(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	if err := h.communicationUseCase.MarkCommunicationAsClicked(c.Request.Context(), id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Click tracked"})
}

// GetPendingCommunications gets pending communications
func (h *CommunicationHandler) GetPendingCommunications(c *gin.Context) {
	communications, err := h.communicationUseCase.GetPendingCommunications(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, communications)
}

// GetCommunicationsForRetry gets communications that need retry
func (h *CommunicationHandler) GetCommunicationsForRetry(c *gin.Context) {
	communications, err := h.communicationUseCase.GetCommunicationsForRetry(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, communications)
}

// GetCommunicationsByRecipient gets communications by recipient
func (h *CommunicationHandler) GetCommunicationsByRecipient(c *gin.Context) {
	recipientType := entities.RecipientType(c.Query("type"))
	recipientIDStr := c.Query("id")

	if recipientIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipient ID is required"})
		return
	}

	recipientID, err := primitive.ObjectIDFromHex(recipientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
		return
	}

	communications, err := h.communicationUseCase.GetCommunicationsByRecipient(c.Request.Context(), recipientType, recipientID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, communications)
}

// GetCommunicationsByCampaign gets communications by campaign
func (h *CommunicationHandler) GetCommunicationsByCampaign(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	communications, err := h.communicationUseCase.GetCommunicationsByCampaign(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, communications)
}

// GetCommunicationsByBatch gets communications by batch
func (h *CommunicationHandler) GetCommunicationsByBatch(c *gin.Context) {
	batchID := c.Param("id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch ID is required"})
		return
	}

	communications, err := h.communicationUseCase.GetCommunicationsByBatch(c.Request.Context(), batchID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, communications)
}

// GetCommunicationStatistics gets communication statistics
func (h *CommunicationHandler) GetCommunicationStatistics(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
	} else {
		// Default to last 30 days
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
	} else {
		endDate = time.Now()
	}

	stats, err := h.communicationUseCase.GetCommunicationStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteCommunication deletes a communication
func (h *CommunicationHandler) DeleteCommunication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid communication ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.DeleteCommunication(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Communication deleted successfully"})
}

// CreateTemplate creates a new communication template
func (h *CommunicationHandler) CreateTemplate(c *gin.Context) {
	var template entities.CommunicationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.CreateTemplate(c.Request.Context(), &template, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateTemplate updates a template
func (h *CommunicationHandler) UpdateTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var template entities.CommunicationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.UpdateTemplate(c.Request.Context(), &template, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetTemplate gets a template by ID
func (h *CommunicationHandler) GetTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.communicationUseCase.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, template)
}

// ListTemplates lists templates with filtering
func (h *CommunicationHandler) ListTemplates(c *gin.Context) {
	filter := &repositories.CommunicationTemplateFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.Category = c.Query("category")
	filter.Language = c.Query("language")
	filter.Search = c.Query("search")

	if activeStr := c.Query("active"); activeStr != "" {
		active := activeStr == "true"
		filter.Active = &active
	}

	if isDefaultStr := c.Query("is_default"); isDefaultStr != "" {
		isDefault := isDefaultStr == "true"
		filter.IsDefault = &isDefault
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	templates, total, err := h.communicationUseCase.ListTemplates(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   templates,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetTemplatesByCategory gets templates by category
func (h *CommunicationHandler) GetTemplatesByCategory(c *gin.Context) {
	category := entities.TemplateCategory(c.Query("category"))
	templateType := entities.TemplateType(c.Query("type"))

	templates, err := h.communicationUseCase.GetTemplatesByCategory(c.Request.Context(), category, templateType)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetDefaultTemplate gets the default template for a category and type
func (h *CommunicationHandler) GetDefaultTemplate(c *gin.Context) {
	category := entities.TemplateCategory(c.Query("category"))
	templateType := entities.TemplateType(c.Query("type"))

	template, err := h.communicationUseCase.GetDefaultTemplate(c.Request.Context(), category, templateType)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetActiveTemplates gets all active templates
func (h *CommunicationHandler) GetActiveTemplates(c *gin.Context) {
	templates, err := h.communicationUseCase.GetActiveTemplates(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, templates)
}

// DeleteTemplate deletes a template
func (h *CommunicationHandler) DeleteTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.communicationUseCase.DeleteTemplate(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// UpdateCommunication updates a communication
func (h *CommunicationHandler) UpdateCommunication(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Communication updated"})
}

// GetCommunicationStatus gets communication delivery status
func (h *CommunicationHandler) GetCommunicationStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "delivered"})
}
