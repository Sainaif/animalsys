package task

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCase struct {
	taskRepo     repositories.TaskRepository
	auditLogRepo repositories.AuditLogRepository
}

func NewTaskUseCase(
	taskRepo repositories.TaskRepository,
	auditLogRepo repositories.AuditLogRepository,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:     taskRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateTask creates a new task
func (uc *TaskUseCase) CreateTask(ctx context.Context, task *entities.Task, userID primitive.ObjectID) error {
	// Validate required fields
	if task.Title == "" {
		return errors.NewBadRequest("Task title is required")
	}

	if task.Category == "" {
		return errors.NewBadRequest("Task category is required")
	}

	if task.Priority == "" {
		return errors.NewBadRequest("Task priority is required")
	}

	task.CreatedBy = userID

	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "task", task.Title, "").
			WithEntityID(task.ID))

	return nil
}

// AddTaskComment adds a comment to a task
func (uc *TaskUseCase) AddTaskComment(ctx context.Context, taskID primitive.ObjectID, text string, userID primitive.ObjectID) (*entities.TaskComment, error) {
	if text == "" {
		return nil, errors.NewBadRequest("Comment text cannot be empty")
	}

	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	comment := task.AddComment(text, userID)

	if err := uc.taskRepo.AddTaskComment(ctx, taskID, comment); err != nil {
		return nil, err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "added comment").
			WithEntityID(taskID))

	return comment, nil
}

// GetTaskComments gets all comments for a task
func (uc *TaskUseCase) GetTaskComments(ctx context.Context, taskID primitive.ObjectID) ([]entities.TaskComment, error) {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return task.Comments, nil
}

// GetTaskByID retrieves a task by ID
func (uc *TaskUseCase) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*entities.Task, error) {
	return uc.taskRepo.FindByID(ctx, id)
}

// UpdateTask updates a task
func (uc *TaskUseCase) UpdateTask(ctx context.Context, task *entities.Task, userID primitive.ObjectID) error {
	// Validate required fields
	if task.Title == "" {
		return errors.NewBadRequest("Task title is required")
	}

	// Check if task exists
	existing, err := uc.taskRepo.FindByID(ctx, task.ID)
	if err != nil {
		return err
	}

	// Update timestamps
	task.CreatedBy = existing.CreatedBy
	task.CreatedAt = existing.CreatedAt

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "").
			WithEntityID(task.ID))

	return nil
}

// DeleteTask deletes a task
func (uc *TaskUseCase) DeleteTask(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if task exists
	task, err := uc.taskRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.taskRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "task", task.Title, "").
			WithEntityID(id))

	return nil
}

// ListTasks lists tasks with filtering and pagination
func (uc *TaskUseCase) ListTasks(ctx context.Context, filter *repositories.TaskFilter) ([]*entities.Task, int64, error) {
	return uc.taskRepo.List(ctx, filter)
}

// GetTasksByAssignee gets tasks assigned to a user
func (uc *TaskUseCase) GetTasksByAssignee(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	return uc.taskRepo.GetByAssignee(ctx, userID)
}

// GetTasksByCategory gets tasks by category
func (uc *TaskUseCase) GetTasksByCategory(ctx context.Context, category entities.TaskCategory) ([]*entities.Task, error) {
	return uc.taskRepo.GetByCategory(ctx, category)
}

// GetTasksByStatus gets tasks by status
func (uc *TaskUseCase) GetTasksByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error) {
	return uc.taskRepo.GetByStatus(ctx, status)
}

// GetTasksByRelatedEntity gets tasks by related entity
func (uc *TaskUseCase) GetTasksByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Task, error) {
	return uc.taskRepo.GetByRelatedEntity(ctx, entityType, entityID)
}

// GetOverdueTasks gets all overdue tasks
func (uc *TaskUseCase) GetOverdueTasks(ctx context.Context) ([]*entities.Task, error) {
	return uc.taskRepo.GetOverdueTasks(ctx)
}

// GetUpcomingTasks gets upcoming tasks (due within specified days)
func (uc *TaskUseCase) GetUpcomingTasks(ctx context.Context, days int) ([]*entities.Task, error) {
	return uc.taskRepo.GetUpcomingTasks(ctx, days)
}

// GetMyTasks gets tasks for a specific user (assigned or created)
func (uc *TaskUseCase) GetMyTasks(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	return uc.taskRepo.GetMyTasks(ctx, userID)
}

// GetTaskStatistics gets task statistics
func (uc *TaskUseCase) GetTaskStatistics(ctx context.Context) (*repositories.TaskStatistics, error) {
	return uc.taskRepo.GetTaskStatistics(ctx)
}

// AssignTask assigns a task to a user
func (uc *TaskUseCase) AssignTask(ctx context.Context, taskID, assigneeID, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	task.AssignedTo = &assigneeID
	task.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "assigned task").
			WithEntityID(taskID))

	return nil
}

// UnassignTask unassigns a task
func (uc *TaskUseCase) UnassignTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	task.AssignedTo = nil
	task.UpdatedAt = time.Now()

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "unassigned task").
			WithEntityID(taskID))

	return nil
}

// StartTask starts a task (changes status to in_progress)
func (uc *TaskUseCase) StartTask(ctx context.Context, taskID, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.Status != entities.TaskStatusPending {
		return errors.NewBadRequest("Task is not in pending status")
	}

	task.Start(userID)

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "started task").
			WithEntityID(taskID))

	return nil
}

// CompleteTask completes a task
func (uc *TaskUseCase) CompleteTask(ctx context.Context, taskID, userID primitive.ObjectID, completionNotes string) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.Status == entities.TaskStatusCompleted {
		return errors.NewBadRequest("Task is already completed")
	}

	task.Complete(userID)
	if completionNotes != "" {
		task.CompletionNotes = completionNotes
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "completed task").
			WithEntityID(taskID))

	return nil
}

// CancelTask cancels a task
func (uc *TaskUseCase) CancelTask(ctx context.Context, taskID, userID primitive.ObjectID, reason string) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.Status == entities.TaskStatusCompleted {
		return errors.NewBadRequest("Cannot cancel a completed task")
	}

	task.Cancel()
	if reason != "" {
		task.CompletionNotes = "Cancelled: " + reason
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "cancelled task").
			WithEntityID(taskID))

	return nil
}

// AddChecklistItem adds a checklist item to a task
func (uc *TaskUseCase) AddChecklistItem(ctx context.Context, taskID primitive.ObjectID, description string, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if description == "" {
		return errors.NewBadRequest("Checklist item description is required")
	}

	task.AddChecklistItem(description)

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "added checklist item").
			WithEntityID(taskID))

	return nil
}

// CompleteChecklistItem marks a checklist item as complete
func (uc *TaskUseCase) CompleteChecklistItem(ctx context.Context, taskID primitive.ObjectID, itemID string, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !task.CompleteChecklistItem(itemID, userID) {
		return errors.NewBadRequest("Checklist item not found")
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "completed checklist item").
			WithEntityID(taskID))

	return nil
}

// RemoveChecklistItem removes a checklist item from a task
func (uc *TaskUseCase) RemoveChecklistItem(ctx context.Context, taskID primitive.ObjectID, itemID string, userID primitive.ObjectID) error {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !task.RemoveChecklistItem(itemID) {
		return errors.NewBadRequest("Checklist item not found")
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "task", task.Title, "removed checklist item").
			WithEntityID(taskID))

	return nil
}
