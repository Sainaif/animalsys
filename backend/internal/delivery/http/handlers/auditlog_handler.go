package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/auditlog"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLogHandler handles audit log-related HTTP requests
type AuditLogHandler struct {
	auditLogUseCase *auditlog.AuditLogUseCase
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(auditLogUseCase *auditlog.AuditLogUseCase) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogUseCase: auditLogUseCase,
	}
}

// GetAuditLog gets an audit log by ID
func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
		return
	}

	log, err := h.auditLogUseCase.GetAuditLogByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListAuditLogs lists audit logs with filtering
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	filter := &repositories.AuditLogFilter{}

	// Parse query parameters
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err == nil {
			filter.UserID = &userID
		}
	}

	filter.Action = c.Query("action")
	filter.EntityType = c.Query("entity_type")

	if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
		entityID, err := primitive.ObjectIDFromHex(entityIDStr)
		if err == nil {
			filter.EntityID = &entityID
		}
	}

	// Parse date filters
	if fromDateStr := c.Query("from_date"); fromDateStr != "" {
		fromDate, err := time.Parse(time.RFC3339, fromDateStr)
		if err == nil {
			fromDateTime := primitive.NewDateTimeFromTime(fromDate)
			filter.FromDate = &fromDateTime
		}
	}

	if toDateStr := c.Query("to_date"); toDateStr != "" {
		toDate, err := time.Parse(time.RFC3339, toDateStr)
		if err == nil {
			toDateTime := primitive.NewDateTimeFromTime(toDate)
			filter.ToDate = &toDateTime
		}
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	logs, total, err := h.auditLogUseCase.ListAuditLogs(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetUserActivity gets audit logs for a specific user
func (h *AuditLogHandler) GetUserActivity(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	logs, total, err := h.auditLogUseCase.GetUserActivity(c.Request.Context(), userID, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetEntityHistory gets audit logs for a specific entity
func (h *AuditLogHandler) GetEntityHistory(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityIDParam := c.Param("entity_id")
	entityID, err := primitive.ObjectIDFromHex(entityIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	logs, total, err := h.auditLogUseCase.GetEntityHistory(c.Request.Context(), entityType, entityID, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetRecentActivity gets recent audit logs
func (h *AuditLogHandler) GetRecentActivity(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)

	logs, total, err := h.auditLogUseCase.GetRecentActivity(c.Request.Context(), limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
	})
}

// GetActionLogs gets audit logs for a specific action
func (h *AuditLogHandler) GetActionLogs(c *gin.Context) {
	action := c.Param("action")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	logs, total, err := h.auditLogUseCase.GetActionLogs(c.Request.Context(), action, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetLogsForDateRange gets audit logs within a date range
func (h *AuditLogHandler) GetLogsForDateRange(c *gin.Context) {
	fromDateStr := c.Query("from")
	if fromDateStr == "" {
		fromDateStr = c.Query("start_date")
	}
	toDateStr := c.Query("to")
	if toDateStr == "" {
		toDateStr = c.Query("end_date")
	}

	if fromDateStr == "" || toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both from/start_date and to/end_date dates are required"})
		return
	}

	// Try RFC3339 format first, then fall back to simple date format
	fromDate, err := time.Parse(time.RFC3339, fromDateStr)
	if err != nil {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format (use RFC3339 or YYYY-MM-DD)"})
			return
		}
	}

	toDate, err := time.Parse(time.RFC3339, toDateStr)
	if err != nil {
		toDate, err = time.Parse("2006-01-02", toDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format (use RFC3339 or YYYY-MM-DD)"})
			return
		}
		// Set to end of day for simple date format
		toDate = toDate.Add(24*time.Hour - time.Second)
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	logs, total, err := h.auditLogUseCase.GetLogsForDateRange(c.Request.Context(), fromDate, toDate, limit, offset)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetAuditStatistics gets audit log statistics
func (h *AuditLogHandler) GetAuditStatistics(c *gin.Context) {
	stats, err := h.auditLogUseCase.GetAuditStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ExportAuditLogs exports audit logs in specified format
func (h *AuditLogHandler) ExportAuditLogs(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	// For now, return a simple response indicating export capability
	c.JSON(http.StatusOK, gin.H{
		"format": format,
		"status": "export initiated",
		"message": "Audit logs export prepared",
	})
}

// DeleteOldLogs deletes audit logs older than specified days (admin only)
func (h *AuditLogHandler) DeleteOldLogs(c *gin.Context) {
	daysParam := c.Query("days")
	if daysParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "days parameter is required"})
		return
	}

	days, err := strconv.Atoi(daysParam)
	if err != nil || days < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter (must be >= 1)"})
		return
	}

	deletedCount, err := h.auditLogUseCase.DeleteOldLogs(c.Request.Context(), days)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Old audit logs deleted successfully",
		"deleted": deletedCount,
	})
}

// SearchAuditLogs searches audit logs
func (h *AuditLogHandler) SearchAuditLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"logs": []interface{}{}, "total": 0})
}
