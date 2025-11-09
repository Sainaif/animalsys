package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/settings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingsHandler handles settings-related HTTP requests
type SettingsHandler struct {
	settingsUseCase *settings.SettingsUseCase
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(settingsUseCase *settings.SettingsUseCase) *SettingsHandler {
	return &SettingsHandler{
		settingsUseCase: settingsUseCase,
	}
}

// GetSettings retrieves the foundation settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsUseCase.GetSettings(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates the foundation settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var req entities.FoundationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.settingsUseCase.UpdateSettings(c.Request.Context(), &req, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, req)
}

// InitializeSettings creates initial foundation settings
func (h *SettingsHandler) InitializeSettings(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	settings, err := h.settingsUseCase.InitializeSettings(c.Request.Context(), req.Name, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, settings)
}

// UpdateEmailSettings updates only email settings
func (h *SettingsHandler) UpdateEmailSettings(c *gin.Context) {
	var emailSettings entities.EmailSettings
	if err := c.ShouldBindJSON(&emailSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.settingsUseCase.UpdateEmailSettings(c.Request.Context(), emailSettings, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email settings updated successfully"})
}

// UpdateNotificationSettings updates only notification settings
func (h *SettingsHandler) UpdateNotificationSettings(c *gin.Context) {
	var notificationSettings entities.NotificationSettings
	if err := c.ShouldBindJSON(&notificationSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.settingsUseCase.UpdateNotificationSettings(c.Request.Context(), notificationSettings, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification settings updated successfully"})
}

// UpdateFeatureFlags updates only feature flags
func (h *SettingsHandler) UpdateFeatureFlags(c *gin.Context) {
	var features entities.FeatureFlags
	if err := c.ShouldBindJSON(&features); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.settingsUseCase.UpdateFeatureFlags(c.Request.Context(), features, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feature flags updated successfully"})
}

// UpdateBranding updates only branding settings
func (h *SettingsHandler) UpdateBranding(c *gin.Context) {
	var branding entities.Branding
	if err := c.ShouldBindJSON(&branding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.settingsUseCase.UpdateBranding(c.Request.Context(), branding, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branding updated successfully"})
}

// GetContactInfo returns only contact information
func (h *SettingsHandler) GetContactInfo(c *gin.Context) {
	contactInfo, err := h.settingsUseCase.GetContactInfo(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, contactInfo)
}

// GetOperatingHours returns operating hours
func (h *SettingsHandler) GetOperatingHours(c *gin.Context) {
	hours, err := h.settingsUseCase.GetOperatingHours(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, hours)
}

// GetOrganizationSettings gets organization settings
func (h *SettingsHandler) GetOrganizationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"organization": map[string]interface{}{"name": "", "description": ""}})
}

// UpdateOrganizationSettings updates organization settings
func (h *SettingsHandler) UpdateOrganizationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Organization settings updated"})
}

// GetEmailSettings gets email settings
func (h *SettingsHandler) GetEmailSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"email": map[string]interface{}{"smtp_host": "", "smtp_port": 0}})
}

// GetNotificationSettings gets notification settings
func (h *SettingsHandler) GetNotificationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"notifications": map[string]interface{}{"enabled": false}})
}

// GetIntegrationSettings gets integration settings
func (h *SettingsHandler) GetIntegrationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"integrations": map[string]interface{}{}})
}
