package inventory

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryUseCase interface {
	CreateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error
	GetInventoryItemByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error)
	GetInventoryItemBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error)
	UpdateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error
	DeleteInventoryItem(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error
	ListInventoryItems(ctx context.Context, filter *repositories.InventoryFilter) ([]*entities.InventoryItem, int64, error)
	GetInventoryItemsByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, int64, error)
	GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error)
	GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error)
	GetInventoryStatistics(ctx context.Context) (*repositories.InventoryStatistics, error)
	AddStock(ctx context.Context, itemID primitive.ObjectID, quantity, unitCost float64, userID primitive.ObjectID, reference, notes string) error
	RemoveStock(ctx context.Context, itemID primitive.ObjectID, quantity float64, userID primitive.ObjectID, reason, reference, notes string) error
	AdjustStock(ctx context.Context, itemID primitive.ObjectID, newQuantity float64, userID primitive.ObjectID, reason, notes string) error
	ActivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error
	DeactivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error
	GetInventoryHistory(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, int64, error)
}
