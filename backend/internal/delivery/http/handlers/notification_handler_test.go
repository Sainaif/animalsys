package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/notification/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNotificationHandler_GetNotificationPreferences(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mocks.NotificationUseCase)
	handler := NewNotificationHandler(mockUseCase)

	userID := primitive.NewObjectID()
	expectedPreferences := &entities.NotificationPreferences{
		UserID:       userID,
		EmailEnabled: true,
	}

	mockUseCase.On("GetNotificationPreferences", mock.Anything, userID).Return(expectedPreferences, nil)

	r := gin.Default()
	r.GET("/notifications/preferences", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, handler.GetNotificationPreferences)

	req, _ := http.NewRequest(http.MethodGet, "/notifications/preferences", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualPreferences entities.NotificationPreferences
	json.Unmarshal(w.Body.Bytes(), &actualPreferences)
	assert.Equal(t, expectedPreferences.UserID, actualPreferences.UserID)
	assert.Equal(t, expectedPreferences.EmailEnabled, actualPreferences.EmailEnabled)

	mockUseCase.AssertExpectations(t)
}

func TestNotificationHandler_UpdateNotificationPreferences(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mocks.NotificationUseCase)
	handler := NewNotificationHandler(mockUseCase)

	userID := primitive.NewObjectID()
	preferencesToUpdate := entities.NotificationPreferences{
		EmailEnabled: false,
		PushEnabled:  true,
	}

	mockUseCase.On("UpdateNotificationPreferences", mock.Anything, userID, mock.AnythingOfType("*entities.NotificationPreferences")).Return(nil)

	r := gin.Default()
	r.PUT("/notifications/preferences", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, handler.UpdateNotificationPreferences)

	body, _ := json.Marshal(preferencesToUpdate)
	req, _ := http.NewRequest(http.MethodPut, "/notifications/preferences", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestNotificationHandler_GetNotificationsByType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(mocks.NotificationUseCase)
	handler := NewNotificationHandler(mockUseCase)

	userID := primitive.NewObjectID()
	notificationType := entities.NotificationTypeInfo

	r := gin.Default()
	r.GET("/notifications/type/:type", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, handler.GetNotificationsByType)

	t.Run("should return notifications for valid type", func(t *testing.T) {
		expectedNotifications := []*entities.Notification{}
		mockUseCase.On("GetNotificationsByType", mock.Anything, userID, notificationType, int64(20), int64(0)).Return(expectedNotifications, int64(0), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/notifications/type/info", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/notifications/type/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
