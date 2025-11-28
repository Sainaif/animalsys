package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusOnHold     TaskStatus = "on_hold"
)

// TaskPriority represents the priority of a task
type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityUrgent   TaskPriority = "urgent"
)

// TaskCategory represents the category of a task
type TaskCategory string

const (
	TaskCategoryAnimalCare    TaskCategory = "animal_care"
	TaskCategoryMedical       TaskCategory = "medical"
	TaskCategoryAdoption      TaskCategory = "adoption"
	TaskCategoryMaintenance   TaskCategory = "maintenance"
	TaskCategoryAdministrative TaskCategory = "administrative"
	TaskCategoryEvent         TaskCategory = "event"
	TaskCategoryFundraising   TaskCategory = "fundraising"
	TaskCategoryVolunteer     TaskCategory = "volunteer"
	TaskCategoryOther         TaskCategory = "other"
)

// Task represents a task or to-do item
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Category    TaskCategory       `json:"category" bson:"category"`
	Status      TaskStatus         `json:"status" bson:"status"`
	Priority    TaskPriority       `json:"priority" bson:"priority"`

	// Assignment
	AssignedTo   *primitive.ObjectID `json:"assigned_to,omitempty" bson:"assigned_to,omitempty"`
	AssignedBy   primitive.ObjectID  `json:"assigned_by" bson:"assigned_by"`
	AssignedAt   time.Time           `json:"assigned_at" bson:"assigned_at"`

	// Related entities
	RelatedEntity     string              `json:"related_entity,omitempty" bson:"related_entity,omitempty"`         // "animal", "adoption", "event", etc.
	RelatedEntityID   *primitive.ObjectID `json:"related_entity_id,omitempty" bson:"related_entity_id,omitempty"`

	// Scheduling
	DueDate       *time.Time              `json:"due_date,omitempty" bson:"due_date,omitempty"`
	StartDate     *time.Time              `json:"start_date,omitempty" bson:"start_date,omitempty"`
	StartedAt     *time.Time              `json:"started_at,omitempty" bson:"started_at,omitempty"`
	StartedBy     *primitive.ObjectID     `json:"started_by,omitempty" bson:"started_by,omitempty"`
	CompletedAt   *time.Time              `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CompletedBy   *primitive.ObjectID     `json:"completed_by,omitempty" bson:"completed_by,omitempty"`
	CompletionNotes string                `json:"completion_notes,omitempty" bson:"completion_notes,omitempty"`

	// Recurrence
	IsRecurring   bool                `json:"is_recurring" bson:"is_recurring"`
	RecurringRule *RecurringRule      `json:"recurring_rule,omitempty" bson:"recurring_rule,omitempty"`

	// Checklist
	Checklist []ChecklistItem `json:"checklist,omitempty" bson:"checklist,omitempty"`

	// Notes and attachments
	Notes       string        `json:"notes,omitempty" bson:"notes,omitempty"`
	Attachments []string      `json:"attachments,omitempty" bson:"attachments,omitempty"` // Document IDs
	Comments    []TaskComment `json:"comments,omitempty" bson:"comments,omitempty"`

	// Metadata
	Tags      []string  `json:"tags,omitempty" bson:"tags,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// RecurringRule defines how a task recurs
type RecurringRule struct {
	Frequency    string `json:"frequency" bson:"frequency"`         // "daily", "weekly", "monthly", "yearly"
	Interval     int    `json:"interval" bson:"interval"`           // Every X days/weeks/months
	DayOfWeek    *int   `json:"day_of_week,omitempty" bson:"day_of_week,omitempty"`     // 0-6 for weekly
	DayOfMonth   *int   `json:"day_of_month,omitempty" bson:"day_of_month,omitempty"`   // 1-31 for monthly
	EndDate      *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Occurrences  *int   `json:"occurrences,omitempty" bson:"occurrences,omitempty"`     // Stop after N occurrences
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// ChecklistItem represents an item in a task checklist
type ChecklistItem struct {
	ID          string    `json:"id" bson:"id"`
	Text        string    `json:"text" bson:"text"`
	IsCompleted bool      `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CompletedBy *primitive.ObjectID `json:"completed_by,omitempty" bson:"completed_by,omitempty"`
}

// NewTask creates a new task
func NewTask(title string, category TaskCategory, priority TaskPriority, assignedBy primitive.ObjectID) *Task {
	now := time.Now()
	return &Task{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Category:    category,
		Status:      TaskStatusPending,
		Priority:    priority,
		AssignedBy:  assignedBy,
		AssignedAt:  now,
		IsRecurring: false,
		Checklist:   []ChecklistItem{},
		Attachments: []string{},
		Tags:        []string{},
		CreatedBy:   assignedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AssignTo assigns the task to a user
func (t *Task) AssignTo(userID primitive.ObjectID) {
	t.AssignedTo = &userID
	t.UpdatedAt = time.Now()
}

// UpdateStatus updates the task status
func (t *Task) UpdateStatus(status TaskStatus) {
	t.Status = status
	t.UpdatedAt = time.Now()

	if status == TaskStatusCompleted {
		now := time.Now()
		t.CompletedAt = &now
	}
}

// AddChecklistItem adds an item to the checklist
func (t *Task) AddChecklistItem(text string) {
	item := ChecklistItem{
		ID:          primitive.NewObjectID().Hex(),
		Text:        text,
		IsCompleted: false,
	}
	t.Checklist = append(t.Checklist, item)
	t.UpdatedAt = time.Now()
}

// AddComment adds a comment to the task
func (t *Task) AddComment(text string, userID primitive.ObjectID) *TaskComment {
	comment := &TaskComment{
		ID:        primitive.NewObjectID(),
		Text:      text,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	t.Comments = append(t.Comments, *comment)
	t.UpdatedAt = time.Now()
	return comment
}

// CompleteChecklistItem marks a checklist item as completed
func (t *Task) CompleteChecklistItem(itemID string, userID primitive.ObjectID) bool {
	for i := range t.Checklist {
		if t.Checklist[i].ID == itemID {
			t.Checklist[i].IsCompleted = true
			now := time.Now()
			t.Checklist[i].CompletedAt = &now
			t.Checklist[i].CompletedBy = &userID
			t.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// IsOverdue checks if the task is overdue
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TaskStatusCompleted || t.Status == TaskStatusCancelled {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// GetCompletionPercentage calculates the completion percentage based on checklist
func (t *Task) GetCompletionPercentage() float64 {
	if len(t.Checklist) == 0 {
		if t.Status == TaskStatusCompleted {
			return 100.0
		}
		return 0.0
	}

	completed := 0
	for _, item := range t.Checklist {
		if item.IsCompleted {
			completed++
		}
	}

	return float64(completed) / float64(len(t.Checklist)) * 100.0
}

// Start starts a task (changes status to in_progress)
func (t *Task) Start(userID primitive.ObjectID) {
	t.Status = TaskStatusInProgress
	now := time.Now()
	t.StartedAt = &now
	t.StartedBy = &userID
	t.UpdatedAt = now
}

// Complete completes a task
func (t *Task) Complete(userID primitive.ObjectID) {
	t.Status = TaskStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
	t.CompletedBy = &userID
	t.UpdatedAt = now
}

// Cancel cancels a task
func (t *Task) Cancel() {
	t.Status = TaskStatusCancelled
	t.UpdatedAt = time.Now()
}

// RemoveChecklistItem removes a checklist item
func (t *Task) RemoveChecklistItem(itemID string) bool {
	for i := range t.Checklist {
		if t.Checklist[i].ID == itemID {
			t.Checklist = append(t.Checklist[:i], t.Checklist[i+1:]...)
			t.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}
