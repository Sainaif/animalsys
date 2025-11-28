package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLogRepository is a mock type for the AuditLogRepository type
type AuditLogRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, log
func (m *AuditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

// FindByID provides a mock function with given fields: ctx, id
func (m *AuditLogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuditLog), args.Error(1)
}

// List provides a mock function with given fields: ctx, filter
func (m *AuditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entities.AuditLog), args.Get(1).(int64), args.Error(2)
}

// DeleteOlderThan provides a mock function with given fields: ctx, days
func (m *AuditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(int64), args.Error(1)
}

// EnsureIndexes provides a mock function with given fields: ctx
func (m *AuditLogRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
