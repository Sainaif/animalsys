package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryRepository struct {
	mock.Mock
}

func (m *InventoryRepository) Create(ctx context.Context, item *entities.InventoryItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *InventoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) FindBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) Update(ctx context.Context, item *entities.InventoryItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *InventoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *InventoryRepository) List(ctx context.Context, filter *repositories.InventoryFilter) ([]*entities.InventoryItem, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.InventoryItem), args.Get(1).(int64), args.Error(2)
}

func (m *InventoryRepository) GetByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Get(1).(int64), args.Error(2)
}

func (m *InventoryRepository) GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryRepository) GetInventoryStatistics(ctx context.Context) (*repositories.InventoryStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(*repositories.InventoryStatistics), args.Error(1)
}

func (m *InventoryRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
