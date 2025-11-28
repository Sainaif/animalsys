package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/sainaif/animalsys/backend/internal/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskHandler_AddTaskComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockTaskRepo := new(mocks.TaskRepository)
		mockAuditLogRepo := new(mocks.AuditLogRepository)
		taskUseCase := task.NewTaskUseCase(mockTaskRepo, mockAuditLogRepo)
		handler := NewTaskHandler(taskUseCase)

		taskID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		commentText := "This is a test comment"

		reqBody := gin.H{"comment": commentText}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/tasks/"+taskID.Hex()+"/comments", bytes.NewBuffer(jsonBody))
		c.Params = gin.Params{gin.Param{Key: "id", Value: taskID.Hex()}}
		c.Set("user_id", userID)

		task := &entities.Task{ID: taskID, Title: "Test Task"}

		mockTaskRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)
		mockTaskRepo.On("AddTaskComment", mock.Anything, taskID, mock.AnythingOfType("*entities.TaskComment")).Return(nil)
		mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil)

		handler.AddTaskComment(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var comment entities.TaskComment
		err := json.Unmarshal(w.Body.Bytes(), &comment)
		assert.NoError(t, err)
		assert.Equal(t, commentText, comment.Text)
		assert.Equal(t, userID, comment.CreatedBy)
	})

	t.Run("invalid task id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/tasks/invalid-id/comments", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-id"}}

		handler := NewTaskHandler(nil)
		handler.AddTaskComment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTaskHandler_GetTaskComments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockTaskRepo := new(mocks.TaskRepository)
		taskUseCase := task.NewTaskUseCase(mockTaskRepo, nil)
		handler := NewTaskHandler(taskUseCase)

		taskID := primitive.NewObjectID()
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/tasks/"+taskID.Hex()+"/comments", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: taskID.Hex()}}

		comments := []entities.TaskComment{
			{ID: primitive.NewObjectID(), Text: "Comment 1", CreatedBy: userID, CreatedAt: time.Now()},
			{ID: primitive.NewObjectID(), Text: "Comment 2", CreatedBy: userID, CreatedAt: time.Now()},
		}
		task := &entities.Task{ID: taskID, Title: "Test Task", Comments: comments}

		mockTaskRepo.On("FindByID", mock.Anything, taskID).Return(task, nil)

		handler.GetTaskComments(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var res gin.H
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, float64(2), res["total"])
	})
}
