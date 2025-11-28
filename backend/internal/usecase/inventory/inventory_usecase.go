package inventory

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryUseCase struct {
	inventoryRepo        repositories.InventoryRepository
	stockTransactionRepo repositories.StockTransactionRepository
	auditLogRepo         repositories.AuditLogRepository
}

func NewInventoryUseCase(
	inventoryRepo repositories.InventoryRepository,
	stockTransactionRepo repositories.StockTransactionRepository,
	auditLogRepo repositories.AuditLogRepository,
) *InventoryUseCase {
	return &InventoryUseCase{
		inventoryRepo:        inventoryRepo,
		stockTransactionRepo: stockTransactionRepo,
		auditLogRepo:         auditLogRepo,
	}
}

// CreateInventoryItem creates a new inventory item
func (uc *InventoryUseCase) CreateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error {
	// Validate required fields
	if item.Name == "" {
		return errors.NewBadRequest("Item name is required")
	}

	if item.Category == "" {
		return errors.NewBadRequest("Item category is required")
	}

	if item.Unit == "" {
		return errors.NewBadRequest("Item unit is required")
	}

	// Set metadata
	item.CreatedBy = userID
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Set default values
	if !item.IsActive {
		item.IsActive = true
	}

	// Initialize arrays
	if item.Tags == nil {
		item.Tags = []string{}
	}

	// Calculate total value
	item.TotalValue = item.CurrentStock * item.UnitCost

	// Check low stock and expiration
	item.CheckLowStock()
	item.CheckExpiration()

	if err := uc.inventoryRepo.Create(ctx, item); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "inventory", item.Name, "").
			WithEntityID(item.ID))

	return nil
}

// GetInventoryItemByID retrieves an inventory item by ID
func (uc *InventoryUseCase) GetInventoryItemByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error) {
	item, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetInventoryItemBySKU retrieves an inventory item by SKU
func (uc *InventoryUseCase) GetInventoryItemBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error) {
	item, err := uc.inventoryRepo.FindBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateInventoryItem updates an inventory item
func (uc *InventoryUseCase) UpdateInventoryItem(ctx context.Context, item *entities.InventoryItem, userID primitive.ObjectID) error {
	// Validate required fields
	if item.Name == "" {
		return errors.NewBadRequest("Item name is required")
	}

	// Check if item exists
	existing, err := uc.inventoryRepo.FindByID(ctx, item.ID)
	if err != nil {
		return err
	}

	// Preserve creation info and statistics
	item.CreatedBy = existing.CreatedBy
	item.CreatedAt = existing.CreatedAt
	item.TotalUsed = existing.TotalUsed
	item.LastUsedDate = existing.LastUsedDate
	item.LastRestockDate = existing.LastRestockDate

	// Recalculate values
	item.TotalValue = item.CurrentStock * item.UnitCost
	item.CheckLowStock()
	item.CheckExpiration()

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "").
			WithEntityID(item.ID))

	return nil
}

// DeleteInventoryItem deletes an inventory item
func (uc *InventoryUseCase) DeleteInventoryItem(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if item exists
	item, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.inventoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "inventory", item.Name, "").
			WithEntityID(id))

	return nil
}

// ListInventoryItems lists inventory items with filtering and pagination
func (uc *InventoryUseCase) ListInventoryItems(ctx context.Context, filter *repositories.InventoryFilter) ([]*entities.InventoryItem, int64, error) {
	return uc.inventoryRepo.List(ctx, filter)
}

// GetOutOfStockItems gets out-of-stock items
func (uc *InventoryUseCase) GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error) {
	return uc.inventoryRepo.GetOutOfStockItems(ctx)
}

// GetInventoryHistory gets the history of an inventory item
func (uc *InventoryUseCase) GetInventoryHistory(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, int64, error) {
	filter := &repositories.StockTransactionFilter{
		ItemID:    &itemID,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return uc.stockTransactionRepo.List(ctx, filter)
}

// GetInventoryItemsByCategory gets inventory items by category
func (uc *InventoryUseCase) GetInventoryItemsByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, int64, error) {
	items, err := uc.inventoryRepo.GetByCategory(ctx, category)
	if err != nil {
		return nil, 0, err
	}
	return items, int64(len(items)), nil
}

// GetLowStockItems gets items with low stock
func (uc *InventoryUseCase) GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetLowStockItems(ctx)
}

// GetExpiredItems gets expired items
func (uc *InventoryUseCase) GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetExpiredItems(ctx)
}

// GetExpiringSoonItems gets items expiring soon
func (uc *InventoryUseCase) GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetExpiringSoonItems(ctx)
}

// GetItemsNeedingReorder gets items that need reordering
func (uc *InventoryUseCase) GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetItemsNeedingReorder(ctx)
}

// GetInventoryStatistics gets inventory statistics
func (uc *InventoryUseCase) GetInventoryStatistics(ctx context.Context) (*repositories.InventoryStatistics, error) {
	return uc.inventoryRepo.GetInventoryStatistics(ctx)
}

// AddStock adds stock to an inventory item
func (uc *InventoryUseCase) AddStock(ctx context.Context, itemID primitive.ObjectID, quantity, unitCost float64, userID primitive.ObjectID, reference, notes string) error {
	if quantity <= 0 {
		return errors.NewBadRequest("Quantity must be greater than 0")
	}

	item, err := uc.inventoryRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	if !item.IsActive {
		return errors.NewBadRequest("Cannot add stock to inactive item")
	}

	stockBefore := item.CurrentStock

	// Add stock to item
	item.AddStock(quantity, unitCost)

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create stock transaction
	transaction := entities.NewStockTransaction(itemID, entities.TransactionTypeIn, quantity, stockBefore, userID)
	transaction.SetCost(unitCost)
	transaction.Reference = reference
	transaction.Notes = notes

	if err := uc.stockTransactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "added stock").
			WithEntityID(itemID))

	return nil
}

// RemoveStock removes stock from an inventory item
func (uc *InventoryUseCase) RemoveStock(ctx context.Context, itemID primitive.ObjectID, quantity float64, userID primitive.ObjectID, reason, reference, notes string) error {
	if quantity <= 0 {
		return errors.NewBadRequest("Quantity must be greater than 0")
	}

	item, err := uc.inventoryRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	if !item.IsActive {
		return errors.NewBadRequest("Cannot remove stock from inactive item")
	}

	if item.CurrentStock < quantity {
		return errors.NewBadRequest("Insufficient stock")
	}

	stockBefore := item.CurrentStock

	// Remove stock from item
	if !item.RemoveStock(quantity) {
		return errors.NewBadRequest("Failed to remove stock")
	}

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create stock transaction
	transaction := entities.NewStockTransaction(itemID, entities.TransactionTypeOut, quantity, stockBefore, userID)
	transaction.Reason = reason
	transaction.Reference = reference
	transaction.Notes = notes

	if err := uc.stockTransactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "removed stock").
			WithEntityID(itemID))

	return nil
}

// AdjustStock adjusts stock (inventory count correction)
func (uc *InventoryUseCase) AdjustStock(ctx context.Context, itemID primitive.ObjectID, newQuantity float64, userID primitive.ObjectID, reason, notes string) error {
	if newQuantity < 0 {
		return errors.NewBadRequest("Quantity cannot be negative")
	}

	item, err := uc.inventoryRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	stockBefore := item.CurrentStock
	difference := newQuantity - stockBefore

	// Update stock
	item.CurrentStock = newQuantity
	item.TotalValue = item.CurrentStock * item.UnitCost
	item.UpdatedAt = time.Now()
	item.CheckLowStock()

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create stock transaction
	transaction := entities.NewStockTransaction(itemID, entities.TransactionTypeAdjustment, difference, stockBefore, userID)
	transaction.Reason = reason
	transaction.Notes = notes

	if err := uc.stockTransactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "adjusted stock").
			WithEntityID(itemID))

	return nil
}

// ActivateItem activates an inventory item
func (uc *InventoryUseCase) ActivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error {
	item, err := uc.inventoryRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	if item.IsActive {
		return errors.NewBadRequest("Item is already active")
	}

	item.Activate()

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "activated item").
			WithEntityID(itemID))

	return nil
}

// DeactivateItem deactivates an inventory item
func (uc *InventoryUseCase) DeactivateItem(ctx context.Context, itemID, userID primitive.ObjectID) error {
	item, err := uc.inventoryRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	if !item.IsActive {
		return errors.NewBadRequest("Item is already inactive")
	}

	item.Deactivate()

	if err := uc.inventoryRepo.Update(ctx, item); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "inventory", item.Name, "deactivated item").
			WithEntityID(itemID))

	return nil
}
