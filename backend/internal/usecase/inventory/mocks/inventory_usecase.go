package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryUseCase struct {
	mock.Mock
}

func (m *InventoryUseCase) CreateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error {
	args := m.Called(ctx, item, userID)
	return args.Error(0)
}

func (m *InventoryUseCase) GetInventoryItemByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) GetInventoryItemBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) UpdateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error {
	args := m.Called(ctx, item, userID)
	return args.Error(0)
}

func (m *InventoryUseCase) DeleteInventoryItem(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *InventoryUseCase) ListInventoryItems(ctx context.Context, filter *repositories.InventoryFilter) ([]*entities.InventoryItem, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.InventoryItem), args.Get(1).(int64), args.Error(2)
}

func (m *InventoryUseCase) GetInventoryItemsByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*entities.InventoryItem), args.Get(1).(int64), args.Error(2)
}

func (m *InventoryUseCase) GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Get(1).(int64), args.Error(2)
}

func (m *InventoryUseCase) GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.InventoryItem), args.Error(1)
}

func (m *InventoryUseCase) GetInventoryStatistics(ctx context.Context) (*repositories.InventoryStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(*repositories.InventoryStatistics), args.Error(1)
}

func (m *InventoryUseCase) AddStock(ctx context.Context, itemID primitive.ObjectID, quantity, unitCost float64, userID primitive.ObjectID, reference, notes string) error {
	args := m.Called(ctx, itemID, quantity, unitCost, userID, reference, notes)
	return args.Error(0)
}

func (m *InventoryUseCase) RemoveStock(ctx context.Context, itemID primitive.ObjectID, quantity float64, userID primitive.ObjectID, reason, reference, notes string) error {
	args := m.Called(ctx, itemID, quantity, userID, reason, reference, notes)
	return args.Error(0)
}

func (m *InventoryUseCase) AdjustStock(ctx context.Context, itemID primitive.ObjectID, newQuantity float64, userID primitive.ObjectID, reason, notes string) error {
	args := m.Called(ctx, itemID, newQuantity, userID, reason, notes)
	return args.Error(0)
}

func (m *InventoryUseCase) ActivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error {
	args := m.Called(ctx, itemID, userID)
	return args.Error(0)
}

func (m *InventoryUseCase) DeactivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error {
	args := m.Called(ctx, itemID, userID)
	return args.Error(0)
}

func (m *InventoryUseCase) GetInventoryHistory(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, int64, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]*entities.StockTransaction), args.Get(1).(int64), args.Error(2)
}
