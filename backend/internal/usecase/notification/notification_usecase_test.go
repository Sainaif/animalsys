package notification

import (
	"context"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	repoMocks "github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNotificationUseCase_GetNotificationPreferences(t *testing.T) {
	mockNotificationRepo := new(repoMocks.NotificationRepository)
	mockAuditLogRepo := new(repoMocks.AuditLogRepository)
	useCase := NewNotificationUseCase(mockNotificationRepo, mockAuditLogRepo)

	userID := primitive.NewObjectID()

	t.Run("should return existing preferences", func(t *testing.T) {
		expectedPreferences := &entities.NotificationPreferences{
			UserID:       userID,
			EmailEnabled: true,
			PushEnabled:  true,
		}

		mockNotificationRepo.On("FindPreferencesByUserID", mock.Anything, userID).Return(expectedPreferences, nil).Once()

		preferences, err := useCase.GetNotificationPreferences(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPreferences, preferences)
		mockNotificationRepo.AssertExpectations(t)
	})

	t.Run("should create and return default preferences if none exist", func(t *testing.T) {
		mockNotificationRepo.On("FindPreferencesByUserID", mock.Anything, userID).Return(nil, nil).Once()
		mockNotificationRepo.On("UpsertPreferences", mock.Anything, mock.AnythingOfType("*entities.NotificationPreferences")).Return(nil).Once()

		preferences, err := useCase.GetNotificationPreferences(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, preferences)
		assert.Equal(t, userID, preferences.UserID)
		assert.True(t, preferences.EmailEnabled)
		assert.False(t, preferences.PushEnabled)
		mockNotificationRepo.AssertExpectations(t)
	})
}

func TestNotificationUseCase_UpdateNotificationPreferences(t *testing.T) {
	mockNotificationRepo := new(repoMocks.NotificationRepository)
	mockAuditLogRepo := new(repoMocks.AuditLogRepository)
	useCase := NewNotificationUseCase(mockNotificationRepo, mockAuditLogRepo)

	userID := primitive.NewObjectID()
	preferencesToUpdate := &entities.NotificationPreferences{
		UserID:      userID,
		EmailEnabled: false,
		PushEnabled:  true,
	}

	mockNotificationRepo.On("UpsertPreferences", mock.Anything, preferencesToUpdate).Return(nil).Once()

	err := useCase.UpdateNotificationPreferences(context.Background(), userID, preferencesToUpdate)

	assert.NoError(t, err)
	mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationUseCase_GetNotificationsByType(t *testing.T) {
	mockNotificationRepo := new(repoMocks.NotificationRepository)
	mockAuditLogRepo := new(repoMocks.AuditLogRepository)
	useCase := NewNotificationUseCase(mockNotificationRepo, mockAuditLogRepo)

	userID := primitive.NewObjectID()
	notificationType := entities.NotificationTypeInfo
	limit := int64(10)
	offset := int64(0)

	expectedNotifications := []*entities.Notification{
		{
			ID:      primitive.NewObjectID(),
			UserID:  userID,
			Type:    notificationType,
			Title:   "Test",
			Message: "Test message",
		},
	}
	expectedTotal := int64(1)

	mockNotificationRepo.On("List", mock.Anything, mock.Anything).Return(expectedNotifications, expectedTotal, nil).Once()

	notifications, total, err := useCase.GetNotificationsByType(context.Background(), userID, notificationType, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedNotifications, notifications)
	assert.Equal(t, expectedTotal, total)
	mockNotificationRepo.AssertExpectations(t)
}
