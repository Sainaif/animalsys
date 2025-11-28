package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepository struct {
	db *mongodb.Database
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *mongodb.Database) repositories.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) collection() *mongo.Collection {
	return r.db.DB.Collection(mongodb.Collections.Tasks)
}

// Create creates a new task
func (r *taskRepository) Create(ctx context.Context, task *entities.Task) error {
	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	_, err := r.collection().InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// FindByID finds a task by ID
func (r *taskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Task, error) {
	var task entities.Task

	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFound("Task not found")
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	return &task, nil
}

// Update updates an existing task
func (r *taskRepository) Update(ctx context.Context, task *entities.Task) error {
	task.UpdatedAt = time.Now()

	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": task}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Task not found")
	}

	return nil
}

// Delete deletes a task by ID
func (r *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewNotFound("Task not found")
	}

	return nil
}

// List returns a list of tasks with pagination and filters
func (r *taskRepository) List(ctx context.Context, filter *repositories.TaskFilter) ([]*entities.Task, int64, error) {
	query := bson.M{}

	if filter.AssignedTo != nil {
		query["assigned_to"] = filter.AssignedTo
	}

	if filter.AssignedBy != nil {
		query["assigned_by"] = filter.AssignedBy
	}

	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Priority != "" {
		query["priority"] = filter.Priority
	}

	if filter.RelatedEntity != "" {
		query["related_entity"] = filter.RelatedEntity
	}

	if filter.RelatedEntityID != nil {
		query["related_entity_id"] = filter.RelatedEntityID
	}

	if filter.IsOverdue != nil && *filter.IsOverdue {
		query["due_date"] = bson.M{"$lt": time.Now()}
		query["status"] = bson.M{"$nin": []string{
			string(entities.TaskStatusCompleted),
			string(entities.TaskStatusCancelled),
		}}
	}

	if filter.DueAfter != nil {
		query["due_date"] = bson.M{"$gte": filter.DueAfter}
	}

	if filter.DueBefore != nil {
		if existing, ok := query["due_date"].(bson.M); ok {
			existing["$lte"] = filter.DueBefore
		} else {
			query["due_date"] = bson.M{"$lte": filter.DueBefore}
		}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"title": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Find with pagination
	opts := options.Find()

	if filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
	}

	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	// Sorting
	sortOrder := 1
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
		sortOrder = -1
	}

	opts.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, total, nil
}

// GetByAssignee returns tasks assigned to a specific user
func (r *taskRepository) GetByAssignee(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	filter := bson.M{"assigned_to": userID}
	opts := options.Find().SetSort(bson.D{{Key: "due_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetByCategory returns tasks by category
func (r *taskRepository) GetByCategory(ctx context.Context, category entities.TaskCategory) ([]*entities.Task, error) {
	filter := bson.M{"category": category}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetByStatus returns tasks by status
func (r *taskRepository) GetByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error) {
	filter := bson.M{"status": status}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetOverdueTasks returns all overdue tasks
func (r *taskRepository) GetOverdueTasks(ctx context.Context) ([]*entities.Task, error) {
	filter := bson.M{
		"due_date": bson.M{"$lt": time.Now()},
		"status": bson.M{"$nin": []string{
			string(entities.TaskStatusCompleted),
			string(entities.TaskStatusCancelled),
		}},
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find overdue tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetUpcomingTasks returns tasks due within the specified days
func (r *taskRepository) GetUpcomingTasks(ctx context.Context, days int) ([]*entities.Task, error) {
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	filter := bson.M{
		"due_date": bson.M{
			"$gte": now,
			"$lte": futureDate,
		},
		"status": bson.M{"$nin": []string{
			string(entities.TaskStatusCompleted),
			string(entities.TaskStatusCancelled),
		}},
	}

	opts := options.Find().SetSort(bson.D{{Key: "due_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find upcoming tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetByRelatedEntity returns tasks related to a specific entity
func (r *taskRepository) GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Task, error) {
	filter := bson.M{
		"related_entity":    entityType,
		"related_entity_id": entityID,
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetMyTasks returns tasks assigned to or created by a user
func (r *taskRepository) GetMyTasks(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"assigned_to": userID},
			{"created_by": userID},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "due_date", Value: 1}})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find my tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entities.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// GetTaskStatistics returns task statistics
func (r *taskRepository) GetTaskStatistics(ctx context.Context) (*repositories.TaskStatistics, error) {
	stats := &repositories.TaskStatistics{
		ByStatus:   make(map[string]int64),
		ByPriority: make(map[string]int64),
		ByCategory: make(map[string]int64),
	}

	// Total tasks
	total, _ := r.collection().CountDocuments(ctx, bson.M{})
	stats.TotalTasks = total

	// By status, priority, category using aggregation
	pipeline := []bson.M{
		{
			"$facet": bson.M{
				"byStatus": []bson.M{
					{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
				},
				"byPriority": []bson.M{
					{"$group": bson.M{"_id": "$priority", "count": bson.M{"$sum": 1}}},
				},
				"byCategory": []bson.M{
					{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
				},
			},
		},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)

		var results []bson.M
		if err := cursor.All(ctx, &results); err == nil && len(results) > 0 {
			result := results[0]

			if byStatus, ok := result["byStatus"].([]interface{}); ok {
				for _, item := range byStatus {
					if m, ok := item.(bson.M); ok {
						status := m["_id"].(string)
						count := m["count"].(int32)
						stats.ByStatus[status] = int64(count)
					}
				}
			}

			if byPriority, ok := result["byPriority"].([]interface{}); ok {
				for _, item := range byPriority {
					if m, ok := item.(bson.M); ok {
						priority := m["_id"].(string)
						count := m["count"].(int32)
						stats.ByPriority[priority] = int64(count)
					}
				}
			}

			if byCategory, ok := result["byCategory"].([]interface{}); ok {
				for _, item := range byCategory {
					if m, ok := item.(bson.M); ok {
						category := m["_id"].(string)
						count := m["count"].(int32)
						stats.ByCategory[category] = int64(count)
					}
				}
			}
		}
	}

	// Overdue tasks
	overdueCount, _ := r.collection().CountDocuments(ctx, bson.M{
		"due_date": bson.M{"$lt": time.Now()},
		"status": bson.M{"$nin": []string{
			string(entities.TaskStatusCompleted),
			string(entities.TaskStatusCancelled),
		}},
	})
	stats.OverdueTasks = overdueCount

	// Due today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)
	dueTodayCount, _ := r.collection().CountDocuments(ctx, bson.M{
		"due_date": bson.M{"$gte": today, "$lt": tomorrow},
	})
	stats.DueTodayTasks = dueTodayCount

	// Completed this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	completedThisWeek, _ := r.collection().CountDocuments(ctx, bson.M{
		"status":       string(entities.TaskStatusCompleted),
		"completed_at": bson.M{"$gte": weekAgo},
	})
	stats.CompletedThisWeek = completedThisWeek

	// Completed this month
	monthAgo := time.Now().AddDate(0, -1, 0)
	completedThisMonth, _ := r.collection().CountDocuments(ctx, bson.M{
		"status":       string(entities.TaskStatusCompleted),
		"completed_at": bson.M{"$gte": monthAgo},
	})
	stats.CompletedThisMonth = completedThisMonth

	return stats, nil
}

// AddTaskComment adds a comment to a task
func (r *taskRepository) AddTaskComment(ctx context.Context, taskID primitive.ObjectID, comment *entities.TaskComment) error {
	if comment.ID.IsZero() {
		comment.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$push": bson.M{"comments": comment},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add comment to task: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewNotFound("Task not found")
	}

	return nil
}

// EnsureIndexes creates necessary indexes for the tasks collection
func (r *taskRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "assigned_to", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "priority", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "due_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "related_entity", Value: 1}, {Key: "related_entity_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create task indexes: %w", err)
	}

	return nil
}
