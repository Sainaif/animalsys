package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/transfer"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransferHandler handles transfer-related HTTP requests
type TransferHandler struct {
	transferUseCase *transfer.TransferUseCase
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(transferUseCase *transfer.TransferUseCase) *TransferHandler {
	return &TransferHandler{
		transferUseCase: transferUseCase,
	}
}

// CreateTransfer creates a new transfer
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var transferReq entities.Transfer
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.CreateTransfer(c.Request.Context(), &transferReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, transferReq)
}

// GetTransfer gets a transfer by ID
func (h *TransferHandler) GetTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	transfer, err := h.transferUseCase.GetTransferByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfer)
}

// UpdateTransfer updates a transfer
func (h *TransferHandler) UpdateTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	var transferReq entities.Transfer
	if err := c.ShouldBindJSON(&transferReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transferReq.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.UpdateTransfer(c.Request.Context(), &transferReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transferReq)
}

// DeleteTransfer deletes a transfer
func (h *TransferHandler) DeleteTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.DeleteTransfer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer deleted successfully"})
}

// ListTransfers lists transfers with filtering
func (h *TransferHandler) ListTransfers(c *gin.Context) {
	filter := &repositories.TransferFilter{}

	// Parse query parameters
	filter.Direction = c.Query("direction")
	filter.Status = c.Query("status")
	filter.Reason = c.Query("reason")
	filter.Search = c.Query("search")

	if animalIDStr := c.Query("animal_id"); animalIDStr != "" {
		animalID, err := primitive.ObjectIDFromHex(animalIDStr)
		if err == nil {
			filter.AnimalID = &animalID
		}
	}

	if partnerIDStr := c.Query("partner_id"); partnerIDStr != "" {
		partnerID, err := primitive.ObjectIDFromHex(partnerIDStr)
		if err == nil {
			filter.PartnerID = &partnerID
		}
	}

	if requestedByStr := c.Query("requested_by"); requestedByStr != "" {
		requestedBy, err := primitive.ObjectIDFromHex(requestedByStr)
		if err == nil {
			filter.RequestedBy = &requestedBy
		}
	}

	if approvedByStr := c.Query("approved_by"); approvedByStr != "" {
		approvedBy, err := primitive.ObjectIDFromHex(approvedByStr)
		if err == nil {
			filter.ApprovedBy = &approvedBy
		}
	}

	if scheduledAfterStr := c.Query("scheduled_after"); scheduledAfterStr != "" {
		scheduledAfter, err := time.Parse(time.RFC3339, scheduledAfterStr)
		if err == nil {
			filter.ScheduledAfter = &scheduledAfter
		}
	}

	if scheduledBeforeStr := c.Query("scheduled_before"); scheduledBeforeStr != "" {
		scheduledBefore, err := time.Parse(time.RFC3339, scheduledBeforeStr)
		if err == nil {
			filter.ScheduledBefore = &scheduledBefore
		}
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "scheduled_date")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	transfers, total, err := h.transferUseCase.ListTransfers(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   transfers,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetTransfersByAnimal gets transfers for a specific animal
func (h *TransferHandler) GetTransfersByAnimal(c *gin.Context) {
	idParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	transfers, err := h.transferUseCase.GetTransfersByAnimal(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetTransfersByPartner gets transfers for a specific partner
func (h *TransferHandler) GetTransfersByPartner(c *gin.Context) {
	idParam := c.Param("partner_id")
	partnerID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	transfers, err := h.transferUseCase.GetTransfersByPartner(c.Request.Context(), partnerID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetTransfersByStatus gets transfers by status
func (h *TransferHandler) GetTransfersByStatus(c *gin.Context) {
	status := entities.TransferStatus(c.Param("status"))

	transfers, err := h.transferUseCase.GetTransfersByStatus(c.Request.Context(), status)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetPendingTransfers gets all pending transfers
func (h *TransferHandler) GetPendingTransfers(c *gin.Context) {
	transfers, err := h.transferUseCase.GetPendingTransfers(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetUpcomingTransfers gets transfers scheduled within the next N days
func (h *TransferHandler) GetUpcomingTransfers(c *gin.Context) {
	days := 7 // default
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil {
			days = d
		}
	}

	transfers, err := h.transferUseCase.GetUpcomingTransfers(c.Request.Context(), days)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetOverdueTransfers gets transfers that are overdue
func (h *TransferHandler) GetOverdueTransfers(c *gin.Context) {
	transfers, err := h.transferUseCase.GetOverdueTransfers(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetRequiringFollowUp gets transfers requiring follow-up
func (h *TransferHandler) GetRequiringFollowUp(c *gin.Context) {
	transfers, err := h.transferUseCase.GetRequiringFollowUp(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// GetTransferStatistics gets transfer statistics
func (h *TransferHandler) GetTransferStatistics(c *gin.Context) {
	stats, err := h.transferUseCase.GetTransferStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ApproveTransfer approves a transfer
func (h *TransferHandler) ApproveTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.ApproveTransfer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer approved successfully"})
}

// RejectTransfer rejects a transfer
func (h *TransferHandler) RejectTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.RejectTransfer(c.Request.Context(), id, userID, req.Reason); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer rejected successfully"})
}

// StartTransit marks the transfer as in transit
func (h *TransferHandler) StartTransit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.StartTransit(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer started successfully"})
}

// CompleteTransfer completes a transfer
func (h *TransferHandler) CompleteTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.CompleteTransfer(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer completed successfully"})
}

// CancelTransfer cancels a transfer
func (h *TransferHandler) CancelTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.CancelTransfer(c.Request.Context(), id, userID, req.Reason); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer cancelled successfully"})
}

// ScheduleTransfer schedules a transfer for a specific date
func (h *TransferHandler) ScheduleTransfer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transfer ID"})
		return
	}

	var req struct {
		ScheduledDate string `json:"scheduled_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduledDate, err := time.Parse(time.RFC3339, req.ScheduledDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.transferUseCase.ScheduleTransfer(c.Request.Context(), id, userID, scheduledDate); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer scheduled successfully"})
}

// InitiateTransfer initiates a transfer
func (h *TransferHandler) InitiateTransfer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transfer initiated"})
}
