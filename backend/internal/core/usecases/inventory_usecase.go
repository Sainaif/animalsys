package usecases

import (
	"context"
	"errors"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type InventoryUseCase struct {
	inventoryRepo      interfaces.InventoryRepository
	stockMovementRepo  interfaces.StockMovementRepository
	auditRepo          interfaces.AuditLogRepository
}

func NewInventoryUseCase(
	inventoryRepo interfaces.InventoryRepository,
	stockMovementRepo interfaces.StockMovementRepository,
	auditRepo interfaces.AuditLogRepository,
) *InventoryUseCase {
	return &InventoryUseCase{
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
		auditRepo:         auditRepo,
	}
}

func (uc *InventoryUseCase) CreateItem(ctx context.Context, req *entities.InventoryCreateRequest, createdBy string) (*entities.InventoryItem, error) {
	item := entities.NewInventoryItem(
		req.Name,
		req.Category,
		req.Unit,
	)

	item.SKU = req.SKU
	item.Description = req.Description
	item.CurrentStock = req.InitialStock
	item.MinimumStock = req.MinimumStock
	item.UnitPrice = req.UnitPrice
	item.Supplier = req.Supplier
	item.Location = req.Location
	item.ExpiryDate = req.ExpiryDate
	item.Notes = req.Notes
	item.CreatedBy = createdBy

	if err := uc.inventoryRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	// Record initial stock as a stock movement if > 0
	if req.InitialStock > 0 {
		movement := entities.NewStockMovement(
			item.ID.Hex(),
			item.Name,
			entities.MovementTypeIn,
			req.InitialStock,
			"Initial stock",
		)
		movement.RecordedBy = createdBy
		uc.stockMovementRepo.Create(ctx, movement)
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "inventory", item.ID.Hex(), "Inventory item created")
	uc.auditRepo.Create(ctx, auditLog)

	return item, nil
}

func (uc *InventoryUseCase) GetItemByID(ctx context.Context, id string) (*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetByID(ctx, id)
}

func (uc *InventoryUseCase) ListItems(ctx context.Context, filter *entities.InventoryFilter) ([]*entities.InventoryItem, int64, error) {
	return uc.inventoryRepo.List(ctx, filter)
}

func (uc *InventoryUseCase) UpdateItem(ctx context.Context, id string, req *entities.InventoryUpdateRequest, updatedBy string) (*entities.InventoryItem, error) {
	item, err := uc.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.SKU != "" {
		item.SKU = req.SKU
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.MinimumStock > 0 {
		item.MinimumStock = req.MinimumStock
	}
	if req.UnitPrice > 0 {
		item.UnitPrice = req.UnitPrice
	}
	if req.Supplier != "" {
		item.Supplier = req.Supplier
	}
	if req.Location != "" {
		item.Location = req.Location
	}
	if req.ExpiryDate != nil {
		item.ExpiryDate = req.ExpiryDate
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}
	item.UpdatedBy = updatedBy

	if err := uc.inventoryRepo.Update(ctx, id, item); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "inventory", id, "Inventory item updated")
	uc.auditRepo.Create(ctx, auditLog)

	return item, nil
}

func (uc *InventoryUseCase) DeleteItem(ctx context.Context, id string, deletedBy string) error {
	if err := uc.inventoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "inventory", id, "Inventory item deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *InventoryUseCase) RecordStockMovement(ctx context.Context, req *entities.StockMovementCreateRequest, recordedBy string) (*entities.StockMovement, error) {
	// Verify item exists and get current stock
	item, err := uc.inventoryRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}

	// Validate stock for outgoing movements
	if req.MovementType == entities.MovementTypeOut && item.CurrentStock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	movement := entities.NewStockMovement(
		req.ItemID,
		item.Name,
		req.MovementType,
		req.Quantity,
		req.Reason,
	)
	movement.Reference = req.Reference
	movement.Notes = req.Notes
	movement.RecordedBy = recordedBy

	if err := uc.stockMovementRepo.Create(ctx, movement); err != nil {
		return nil, err
	}

	// Update item stock
	if req.MovementType == entities.MovementTypeIn {
		uc.inventoryRepo.UpdateStock(ctx, req.ItemID, req.Quantity)
	} else {
		uc.inventoryRepo.UpdateStock(ctx, req.ItemID, -req.Quantity)
	}

	// Audit
	auditLog := entities.NewAuditLog(recordedBy, "", "", entities.ActionCreate, "stock_movement", movement.ID.Hex(), "Stock movement recorded")
	uc.auditRepo.Create(ctx, auditLog)

	return movement, nil
}

func (uc *InventoryUseCase) GetStockMovementByID(ctx context.Context, id string) (*entities.StockMovement, error) {
	return uc.stockMovementRepo.GetByID(ctx, id)
}

func (uc *InventoryUseCase) GetItemStockMovements(ctx context.Context, itemID string, limit, offset int) ([]*entities.StockMovement, int64, error) {
	return uc.stockMovementRepo.GetByItemID(ctx, itemID, limit, offset)
}

func (uc *InventoryUseCase) GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetLowStockItems(ctx)
}

func (uc *InventoryUseCase) GetExpiringItems(ctx context.Context, days int) ([]*entities.InventoryItem, error) {
	return uc.inventoryRepo.GetExpiringItems(ctx, days)
}

func (uc *InventoryUseCase) GetInventoryStatistics(ctx context.Context) (map[string]interface{}, error) {
	// Get all items
	items, total, err := uc.inventoryRepo.List(ctx, &entities.InventoryFilter{
		Limit:  0,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalValue := 0.0
	byCategory := make(map[string]int)
	lowStockCount := 0

	for _, item := range items {
		totalValue += float64(item.CurrentStock) * item.UnitPrice
		byCategory[string(item.Category)]++
		if item.CurrentStock <= item.MinimumStock {
			lowStockCount++
		}
	}

	stats := map[string]interface{}{
		"total_items":      total,
		"total_value":      totalValue,
		"by_category":      byCategory,
		"low_stock_count":  lowStockCount,
	}

	return stats, nil
}

func (uc *InventoryUseCase) AdjustStock(ctx context.Context, itemID string, newQuantity int, reason, adjustedBy string) error {
	item, err := uc.inventoryRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}

	oldQuantity := item.CurrentStock
	diff := newQuantity - oldQuantity

	// Record as adjustment movement
	var movementType entities.MovementType
	quantity := diff
	if diff > 0 {
		movementType = entities.MovementTypeIn
	} else {
		movementType = entities.MovementTypeOut
		quantity = -diff
	}

	movement := entities.NewStockMovement(
		itemID,
		item.Name,
		movementType,
		quantity,
		reason,
	)
	movement.RecordedBy = adjustedBy
	movement.Notes = "Stock adjustment"

	if err := uc.stockMovementRepo.Create(ctx, movement); err != nil {
		return err
	}

	// Update item stock directly to the new quantity
	if err := uc.inventoryRepo.UpdateStock(ctx, itemID, diff); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(adjustedBy, "", "", entities.ActionUpdate, "inventory", itemID, "Stock adjusted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}
