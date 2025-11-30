package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationUseCase struct {
	mock.Mock
}

func (m *NotificationUseCase) GetNotificationPreferences(ctx context.Context, userID primitive.ObjectID) (*entities.NotificationPreferences, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entities.NotificationPreferences), args.Error(1)
}

func (m *NotificationUseCase) UpdateNotificationPreferences(ctx context.Context, userID primitive.ObjectID, preferences *entities.NotificationPreferences) error {
	args := m.Called(ctx, userID, preferences)
	return args.Error(0)
}

func (m *NotificationUseCase) GetNotificationsByType(ctx context.Context, userID primitive.ObjectID, notificationType entities.NotificationType, limit, offset int64) ([]*entities.Notification, int64, error) {
	args := m.Called(ctx, userID, notificationType, limit, offset)
	return args.Get(0).([]*entities.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *NotificationUseCase) CreateNotification(ctx context.Context, notification *entities.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *NotificationUseCase) GetNotificationByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*entities.Notification, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(*entities.Notification), args.Error(1)
}

func (m *NotificationUseCase) GetUserNotifications(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Notification), args.Error(1)
}

func (m *NotificationUseCase) MarkAsRead(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *NotificationUseCase) ListNotifications(ctx context.Context, filter *repositories.NotificationFilter) ([]*entities.Notification, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *NotificationUseCase) GetUnreadNotifications(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Notification), args.Error(1)
}

func (m *NotificationUseCase) GetUnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *NotificationUseCase) MarkAsUnread(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *NotificationUseCase) DismissNotification(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *NotificationUseCase) DeleteNotification(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
