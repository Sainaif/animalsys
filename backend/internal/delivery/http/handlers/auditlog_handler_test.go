package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/auditlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockAuditLogRepository struct {
	mock.Mock
}

func (m *mockAuditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *mockAuditLogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.AuditLog), args.Error(1)
}

func (m *mockAuditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *mockAuditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	args := m.Called(ctx, days)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockAuditLogRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestSearchAuditLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mockAuditLogRepository)
	auditLogUseCase := auditlog.NewAuditLogUseCase(mockRepo, nil)
	handler := NewAuditLogHandler(auditLogUseCase)

	router := gin.Default()
	router.GET("/search", handler.SearchAuditLogs)

	t.Run("successful search", func(t *testing.T) {
		userID := primitive.NewObjectID()
		expectedLogs := []*entities.AuditLog{
			{ID: primitive.NewObjectID(), UserID: userID, Action: "login", EntityType: "user"},
		}

		mockRepo.On("List", mock.Anything, mock.AnythingOfType("repositories.AuditLogFilter")).Return(expectedLogs, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/search?user_id="+userID.Hex(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, float64(1), response["total"])
		logs := response["logs"].([]interface{})
		assert.Len(t, logs, 1)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/search?user_id=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid from_date", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/search?from_date=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
