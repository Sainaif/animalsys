package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/task"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	taskUseCase *task.TaskUseCase
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskUseCase *task.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq entities.Task
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.CreateTask(c.Request.Context(), &taskReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, taskReq)
}

// GetTask gets a task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskUseCase.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates a task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var taskReq entities.Task
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskReq.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.UpdateTask(c.Request.Context(), &taskReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, taskReq)
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.DeleteTask(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ListTasks lists tasks with filtering
func (h *TaskHandler) ListTasks(c *gin.Context) {
	filter := &repositories.TaskFilter{}

	// Parse query parameters
	filter.Category = c.Query("category")
	filter.Status = c.Query("status")
	filter.Priority = c.Query("priority")
	filter.RelatedEntity = c.Query("related_entity")
	filter.Search = c.Query("search")

	if assignedToStr := c.Query("assigned_to"); assignedToStr != "" {
		assignedTo, err := primitive.ObjectIDFromHex(assignedToStr)
		if err == nil {
			filter.AssignedTo = &assignedTo
		}
	}

	if createdByStr := c.Query("created_by"); createdByStr != "" {
		createdBy, err := primitive.ObjectIDFromHex(createdByStr)
		if err == nil {
			filter.CreatedBy = &createdBy
		}
	}

	if relatedIDStr := c.Query("related_entity_id"); relatedIDStr != "" {
		relatedID, err := primitive.ObjectIDFromHex(relatedIDStr)
		if err == nil {
			filter.RelatedEntityID = &relatedID
		}
	}

	if overdueStr := c.Query("overdue"); overdueStr != "" {
		overdue := overdueStr == "true"
		filter.IsOverdue = &overdue
	}

	if recurringStr := c.Query("recurring"); recurringStr != "" {
		recurring := recurringStr == "true"
		filter.IsRecurring = &recurring
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "due_date")
	filter.SortOrder = c.DefaultQuery("sort_order", "asc")

	tasks, total, err := h.taskUseCase.ListTasks(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tasks,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetMyTasks gets tasks for the current user
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	tasks, err := h.taskUseCase.GetMyTasks(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTasksByAssignee gets tasks assigned to a user
func (h *TaskHandler) GetTasksByAssignee(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tasks, err := h.taskUseCase.GetTasksByAssignee(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetOverdueTasks gets all overdue tasks
func (h *TaskHandler) GetOverdueTasks(c *gin.Context) {
	tasks, err := h.taskUseCase.GetOverdueTasks(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetUpcomingTasks gets upcoming tasks
func (h *TaskHandler) GetUpcomingTasks(c *gin.Context) {
	days := 7 // default to 7 days
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	tasks, err := h.taskUseCase.GetUpcomingTasks(c.Request.Context(), days)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTaskStatistics gets task statistics
func (h *TaskHandler) GetTaskStatistics(c *gin.Context) {
	stats, err := h.taskUseCase.GetTaskStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// AssignTask assigns a task to a user
func (h *TaskHandler) AssignTask(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		AssigneeID string `json:"assignee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assigneeID, err := primitive.ObjectIDFromHex(req.AssigneeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.AssignTask(c.Request.Context(), taskID, assigneeID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task assigned successfully"})
}

// UnassignTask unassigns a task
func (h *TaskHandler) UnassignTask(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.UnassignTask(c.Request.Context(), taskID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task unassigned successfully"})
}

// StartTask starts a task
func (h *TaskHandler) StartTask(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.StartTask(c.Request.Context(), taskID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task started successfully"})
}

// CompleteTask completes a task
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		CompletionNotes string `json:"completion_notes"`
	}
	_ = c.ShouldBindJSON(&req) // Optional notes

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.CompleteTask(c.Request.Context(), taskID, userID, req.CompletionNotes); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task completed successfully"})
}

// CancelTask cancels a task
func (h *TaskHandler) CancelTask(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req) // Optional reason

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.CancelTask(c.Request.Context(), taskID, userID, req.Reason); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task cancelled successfully"})
}

// AddChecklistItem adds a checklist item to a task
func (h *TaskHandler) AddChecklistItem(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.AddChecklistItem(c.Request.Context(), taskID, req.Description, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checklist item added successfully"})
}

// CompleteChecklistItem completes a checklist item
func (h *TaskHandler) CompleteChecklistItem(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.CompleteChecklistItem(c.Request.Context(), taskID, itemID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checklist item completed successfully"})
}

// RemoveChecklistItem removes a checklist item
func (h *TaskHandler) RemoveChecklistItem(c *gin.Context) {
	taskIDParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(taskIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.taskUseCase.RemoveChecklistItem(c.Request.Context(), taskID, itemID, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checklist item removed successfully"})
}

// AddTaskComment adds a comment to a task
func (h *TaskHandler) AddTaskComment(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	comment, err := h.taskUseCase.AddTaskComment(c.Request.Context(), taskID, req.Comment, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetTaskComments gets comments for a task
func (h *TaskHandler) GetTaskComments(c *gin.Context) {
	idParam := c.Param("id")
	taskID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	comments, err := h.taskUseCase.GetTaskComments(c.Request.Context(), taskID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": len(comments),
	})
}
