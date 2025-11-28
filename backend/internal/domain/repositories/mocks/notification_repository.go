package mocks

import (
	"context"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationRepository struct {
	mock.Mock
}

func (m *NotificationRepository) Create(ctx context.Context, notification *entities.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *NotificationRepository) Update(ctx context.Context, notification *entities.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *NotificationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *NotificationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Notification), args.Error(1)
}

func (m *NotificationRepository) List(ctx context.Context, filter *repositories.NotificationFilter) ([]*entities.Notification, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *NotificationRepository) GetByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Notification), args.Error(1)
}

func (m *NotificationRepository) GetUnreadByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Notification), args.Error(1)
}

func (m *NotificationRepository) CountUnreadByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *NotificationRepository) MarkAsRead(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *NotificationRepository) MarkAsUnread(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *NotificationRepository) MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *NotificationRepository) Dismiss(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *NotificationRepository) DeleteExpired(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *NotificationRepository) FindByGroupKey(ctx context.Context, userID primitive.ObjectID, groupKey string) (*entities.Notification, error) {
	args := m.Called(ctx, userID, groupKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Notification), args.Error(1)
}

func (m *NotificationRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *NotificationRepository) FindPreferencesByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.NotificationPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.NotificationPreferences), args.Error(1)
}

func (m *NotificationRepository) UpsertPreferences(ctx context.Context, preferences *entities.NotificationPreferences) error {
	args := m.Called(ctx, preferences)
	return args.Error(0)
}
