package inventory

import (
	"context"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInventoryUseCase_GetOutOfStockItems(t *testing.T) {
	inventoryRepo := new(mocks.InventoryRepository)
	stockTransactionRepo := new(mocks.StockTransactionRepository)
	useCase := NewInventoryUseCase(inventoryRepo, stockTransactionRepo, nil)

	expectedItems := []*entities.InventoryItem{
		{ID: primitive.NewObjectID(), Name: "Item 1"},
	}
	expectedTotal := int64(1)

	inventoryRepo.On("GetOutOfStockItems", mock.Anything).Return(expectedItems, expectedTotal, nil)

	items, total, err := useCase.GetOutOfStockItems(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	assert.Equal(t, expectedTotal, total)
	inventoryRepo.AssertExpectations(t)
}

func TestInventoryUseCase_GetInventoryHistory(t *testing.T) {
	inventoryRepo := new(mocks.InventoryRepository)
	stockTransactionRepo := new(mocks.StockTransactionRepository)
	useCase := NewInventoryUseCase(inventoryRepo, stockTransactionRepo, nil)

	itemID := primitive.NewObjectID()
	expectedHistory := []*entities.StockTransaction{
		{ID: primitive.NewObjectID(), ItemID: itemID},
	}
	expectedTotal := int64(1)

	stockTransactionRepo.On("List", mock.Anything, mock.MatchedBy(func(filter *repositories.StockTransactionFilter) bool {
		return filter.ItemID.Hex() == itemID.Hex() && filter.SortBy == "created_at" && filter.SortOrder == "desc"
	})).Return(expectedHistory, expectedTotal, nil)

	history, total, err := useCase.GetInventoryHistory(context.Background(), itemID)

	assert.NoError(t, err)
	assert.Equal(t, expectedHistory, history)
	assert.Equal(t, expectedTotal, total)
	stockTransactionRepo.AssertExpectations(t)
}

func TestInventoryUseCase_GetInventoryItemsByCategory(t *testing.T) {
	inventoryRepo := new(mocks.InventoryRepository)
	stockTransactionRepo := new(mocks.StockTransactionRepository)
	useCase := NewInventoryUseCase(inventoryRepo, stockTransactionRepo, nil)

	category := entities.ItemCategory("Food")
	expectedItems := []*entities.InventoryItem{
		{ID: primitive.NewObjectID(), Name: "Item 1", Category: category},
	}
	expectedTotal := int64(1)

	inventoryRepo.On("GetByCategory", mock.Anything, category).Return(expectedItems, nil)

	items, total, err := useCase.GetInventoryItemsByCategory(context.Background(), category)

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	assert.Equal(t, expectedTotal, total)
	inventoryRepo.AssertExpectations(t)
}
