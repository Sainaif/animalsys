package task

import (
	"context"
	"net/http"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	apperrors "github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskUseCase_AddTaskComment(t *testing.T) {
	mockTaskRepo := new(mocks.TaskRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	useCase := NewTaskUseCase(mockTaskRepo, mockAuditLogRepo)

	ctx := context.Background()
	taskID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	commentText := "A valid comment"

	t.Run("success", func(t *testing.T) {
		task := &entities.Task{ID: taskID, Title: "Test Task"}

		mockTaskRepo.On("FindByID", ctx, taskID).Return(task, nil).Once()
		mockTaskRepo.On("AddTaskComment", ctx, taskID, mock.AnythingOfType("*entities.TaskComment")).Return(nil).Once()
		mockAuditLogRepo.On("Create", ctx, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

		comment, err := useCase.AddTaskComment(ctx, taskID, commentText, userID)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, commentText, comment.Text)
		assert.Equal(t, userID, comment.CreatedBy)
		mockTaskRepo.AssertExpectations(t)
		mockAuditLogRepo.AssertExpectations(t)
	})

	t.Run("error - empty comment text", func(t *testing.T) {
		comment, err := useCase.AddTaskComment(ctx, taskID, "", userID)

		assert.Error(t, err)
		assert.Nil(t, comment)
		appErr, ok := err.(*apperrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.Code)
	})

	t.Run("error - task not found", func(t *testing.T) {
		mockTaskRepo.On("FindByID", ctx, taskID).Return(nil, apperrors.NewNotFound("Task not found")).Once()

		comment, err := useCase.AddTaskComment(ctx, taskID, commentText, userID)

		assert.Error(t, err)
		assert.Nil(t, comment)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUseCase_GetTaskComments(t *testing.T) {
	mockTaskRepo := new(mocks.TaskRepository)
	useCase := NewTaskUseCase(mockTaskRepo, nil)

	ctx := context.Background()
	taskID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {
		taskComments := []entities.TaskComment{{Text: "First comment"}, {Text: "Second comment"}}
		task := &entities.Task{ID: taskID, Comments: taskComments}

		mockTaskRepo.On("FindByID", ctx, taskID).Return(task, nil).Once()

		comments, err := useCase.GetTaskComments(ctx, taskID)

		assert.NoError(t, err)
		assert.Len(t, comments, 2)
		assert.Equal(t, "First comment", comments[0].Text)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("error - task not found", func(t *testing.T) {
		mockTaskRepo.On("FindByID", ctx, taskID).Return(nil, apperrors.NewNotFound("Task not found")).Once()

		comments, err := useCase.GetTaskComments(ctx, taskID)

		assert.Error(t, err)
		assert.Nil(t, comments)
		mockTaskRepo.AssertExpectations(t)
	})
}
