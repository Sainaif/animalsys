package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLogRepository struct {
	mock.Mock
}

func (m *AuditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *AuditLogRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.AuditLog), args.Error(1)
}

func (m *AuditLogRepository) List(ctx context.Context, filter repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *AuditLogRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AuditLogRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
