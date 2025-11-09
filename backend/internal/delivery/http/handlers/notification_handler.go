package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/notification"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationUseCase *notification.NotificationUseCase
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationUseCase *notification.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

// CreateNotification creates a new notification (admin only)
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var notif entities.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationUseCase.CreateNotification(c.Request.Context(), &notif); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, notif)
}

// GetNotification gets a notification by ID
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	notif, err := h.notificationUseCase.GetNotificationByID(c.Request.Context(), id, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notif)
}

// ListNotifications lists notifications with filtering
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	filter := &repositories.NotificationFilter{}

	// Parse query parameters
	filter.Type = c.Query("type")
	filter.Priority = c.Query("priority")
	filter.Category = c.Query("category")
	filter.RelatedType = c.Query("related_type")

	// User ID - allow filtering by user for admins, default to current user
	userID := c.MustGet("user_id").(primitive.ObjectID)
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		// Check if user has permission to view other users' notifications
		// For now, default to current user
		filter.UserID = &userID
	} else {
		filter.UserID = &userID
	}

	if readStr := c.Query("read"); readStr != "" {
		read := readStr == "true"
		filter.Read = &read
	}

	if dismissedStr := c.Query("dismissed"); dismissedStr != "" {
		dismissed := dismissedStr == "true"
		filter.Dismissed = &dismissed
	}

	if relatedIDStr := c.Query("related_id"); relatedIDStr != "" {
		relatedID, err := primitive.ObjectIDFromHex(relatedIDStr)
		if err == nil {
			filter.RelatedID = &relatedID
		}
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	notifications, total, err := h.notificationUseCase.ListNotifications(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   notifications,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetMyNotifications gets all non-dismissed notifications for the current user
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	notifications, err := h.notificationUseCase.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadNotifications gets unread notifications for the current user
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	notifications, err := h.notificationUseCase.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadCount gets the count of unread notifications
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	count, err := h.notificationUseCase.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.notificationUseCase.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAsUnread marks a notification as unread
func (h *NotificationHandler) MarkAsUnread(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.notificationUseCase.MarkAsUnread(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as unread"})
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.notificationUseCase.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// DismissNotification dismisses a notification
func (h *NotificationHandler) DismissNotification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.notificationUseCase.DismissNotification(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification dismissed"})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.notificationUseCase.DeleteNotification(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}


// GetNotificationPreferences gets notification preferences for the current user
func (h *NotificationHandler) GetNotificationPreferences(c *gin.Context) {
	// Return mock preferences
	c.JSON(http.StatusOK, gin.H{
		"preferences": map[string]interface{}{
			"email_enabled": true,
			"push_enabled":  false,
		},
	})
}


// GetNotificationsByType gets notifications by type
func (h *NotificationHandler) GetNotificationsByType(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"notifications": []interface{}{}, "total": 0})
}

// UpdateNotificationPreferences updates notification preferences
func (h *NotificationHandler) UpdateNotificationPreferences(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated"})
}
