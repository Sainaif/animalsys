package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, task *entities.Task) error

	// FindByID finds a task by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Task, error)

	// Update updates an existing task
	Update(ctx context.Context, task *entities.Task) error

	// Delete deletes a task by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of tasks with pagination and filters
	List(ctx context.Context, filter *TaskFilter) ([]*entities.Task, int64, error)

	// GetByAssignee returns tasks assigned to a specific user
	GetByAssignee(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error)

	// GetByCategory returns tasks by category
	GetByCategory(ctx context.Context, category entities.TaskCategory) ([]*entities.Task, error)

	// GetByStatus returns tasks by status
	GetByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error)

	// GetOverdueTasks returns all overdue tasks
	GetOverdueTasks(ctx context.Context) ([]*entities.Task, error)

	// GetUpcomingTasks returns tasks due within the specified days
	GetUpcomingTasks(ctx context.Context, days int) ([]*entities.Task, error)

	// GetByRelatedEntity returns tasks related to a specific entity
	GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Task, error)

	// GetMyTasks returns tasks assigned to or created by a user
	GetMyTasks(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error)

	// GetTaskStatistics returns task statistics
	GetTaskStatistics(ctx context.Context) (*TaskStatistics, error)

	// AddTaskComment adds a comment to a task
	AddTaskComment(ctx context.Context, taskID primitive.ObjectID, comment *entities.TaskComment) error

	// EnsureIndexes creates necessary indexes for the tasks collection
	EnsureIndexes(ctx context.Context) error
}

// TaskFilter defines filter criteria for listing tasks
type TaskFilter struct {
	AssignedTo      *primitive.ObjectID
	AssignedBy      *primitive.ObjectID
	CreatedBy       *primitive.ObjectID
	Category        string
	Status          string
	Priority        string
	RelatedEntity   string
	RelatedEntityID *primitive.ObjectID
	IsOverdue       *bool
	IsRecurring     *bool
	DueAfter        *time.Time
	DueBefore       *time.Time
	Search          string
	Tags            []string
	Limit           int64
	Offset          int64
	SortBy          string // Field to sort by
	SortOrder       string // "asc" or "desc"
}

// TaskStatistics represents task statistics
type TaskStatistics struct {
	TotalTasks       int64            `json:"total_tasks"`
	ByStatus         map[string]int64 `json:"by_status"`
	ByPriority       map[string]int64 `json:"by_priority"`
	ByCategory       map[string]int64 `json:"by_category"`
	OverdueTasks     int64            `json:"overdue_tasks"`
	DueTodayTasks    int64            `json:"due_today_tasks"`
	CompletedThisWeek int64           `json:"completed_this_week"`
	CompletedThisMonth int64          `json:"completed_this_month"`
}
