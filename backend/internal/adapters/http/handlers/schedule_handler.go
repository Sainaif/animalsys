package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type ScheduleHandler struct {
	scheduleUseCase *usecases.ScheduleUseCase
}

func NewScheduleHandler(scheduleUseCase *usecases.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleUseCase: scheduleUseCase,
	}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req entities.ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	schedule, err := h.scheduleUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.scheduleUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) List(c *gin.Context) {
	filter := &entities.ScheduleFilter{
		UserID:    c.Query("user_id"),
		ShiftType: c.Query("shift_type"),
		Status:    entities.ShiftStatus(c.Query("status")),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Limit:     parseIntQuery(c.Query("limit"), 10),
		Offset:    parseIntQuery(c.Query("offset"), 0),
	}

	schedules, total, err := h.scheduleUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   schedules,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.ScheduleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	schedule, err := h.scheduleUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.scheduleUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted successfully"})
}

func (h *ScheduleHandler) RequestSwap(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		TargetUserID string `json:"target_user_id" binding:"required"`
		Reason       string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	requestingUserID, _ := c.Get("user_id")
	err := h.scheduleUseCase.RequestSwap(c.Request.Context(), id, requestingUserID.(string), req.TargetUserID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "swap requested successfully"})
}

func (h *ScheduleHandler) ApproveSwap(c *gin.Context) {
	id := c.Param("id")

	approvedBy, _ := c.Get("user_id")
	err := h.scheduleUseCase.ApproveSwap(c.Request.Context(), id, approvedBy.(string), approvedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "swap approved successfully"})
}

func (h *ScheduleHandler) RejectSwap(c *gin.Context) {
	id := c.Param("id")

	rejectedBy, _ := c.Get("user_id")
	err := h.scheduleUseCase.RejectSwap(c.Request.Context(), id, rejectedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "swap rejected successfully"})
}

func (h *ScheduleHandler) MarkComplete(c *gin.Context) {
	id := c.Param("id")

	completedBy, _ := c.Get("user_id")
	err := h.scheduleUseCase.MarkComplete(c.Request.Context(), id, completedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shift marked as complete"})
}

func (h *ScheduleHandler) MarkAbsent(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	markedBy, _ := c.Get("user_id")
	err := h.scheduleUseCase.MarkAbsent(c.Request.Context(), id, req.Reason, markedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shift marked as absent"})
}
