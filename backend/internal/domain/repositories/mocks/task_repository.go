package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository struct {
	mock.Mock
}

func (m *TaskRepository) Create(ctx context.Context, task *entities.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Task), args.Error(1)
}

func (m *TaskRepository) Update(ctx context.Context, task *entities.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *TaskRepository) List(ctx context.Context, filter *repositories.TaskFilter) ([]*entities.Task, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.Task), args.Get(1).(int64), args.Error(2)
}

func (m *TaskRepository) GetByAssignee(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetByCategory(ctx context.Context, category entities.TaskCategory) ([]*entities.Task, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetByStatus(ctx context.Context, status entities.TaskStatus) ([]*entities.Task, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetOverdueTasks(ctx context.Context) ([]*entities.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetUpcomingTasks(ctx context.Context, days int) ([]*entities.Task, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Task, error) {
	args := m.Called(ctx, entityType, entityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetMyTasks(ctx context.Context, userID primitive.ObjectID) ([]*entities.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Task), args.Error(1)
}

func (m *TaskRepository) GetTaskStatistics(ctx context.Context) (*repositories.TaskStatistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TaskStatistics), args.Error(1)
}

func (m *TaskRepository) AddTaskComment(ctx context.Context, taskID primitive.ObjectID, comment *entities.TaskComment) error {
	args := m.Called(ctx, taskID, comment)
	return args.Error(0)
}

func (m *TaskRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
