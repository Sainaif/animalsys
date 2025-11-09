package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/inventory"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryHandler handles inventory-related HTTP requests
type InventoryHandler struct {
	inventoryUseCase *inventory.InventoryUseCase
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryUseCase *inventory.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
	}
}

// CreateInventoryItem creates a new inventory item
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	var itemReq entities.InventoryItem
	if err := c.ShouldBindJSON(&itemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.CreateInventoryItem(c.Request.Context(), &itemReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, itemReq)
}

// GetInventoryItem gets an inventory item by ID
func (h *InventoryHandler) GetInventoryItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.inventoryUseCase.GetInventoryItemByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetInventoryItemBySKU gets an inventory item by SKU
func (h *InventoryHandler) GetInventoryItemBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	item, err := h.inventoryUseCase.GetInventoryItemBySKU(c.Request.Context(), sku)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateInventoryItem updates an inventory item
func (h *InventoryHandler) UpdateInventoryItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var itemReq entities.InventoryItem
	if err := c.ShouldBindJSON(&itemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemReq.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.UpdateInventoryItem(c.Request.Context(), &itemReq, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, itemReq)
}

// DeleteInventoryItem deletes an inventory item
func (h *InventoryHandler) DeleteInventoryItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.DeleteInventoryItem(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item deleted successfully"})
}

// ListInventoryItems lists inventory items with filtering
func (h *InventoryHandler) ListInventoryItems(c *gin.Context) {
	filter := &repositories.InventoryFilter{}

	// Parse query parameters
	filter.Category = c.Query("category")
	filter.SubCategory = c.Query("sub_category")
	filter.Location = c.Query("location")
	filter.Search = c.Query("search")

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	if isLowStockStr := c.Query("is_low_stock"); isLowStockStr != "" {
		isLowStock := isLowStockStr == "true"
		filter.IsLowStock = &isLowStock
	}

	if isExpiredStr := c.Query("is_expired"); isExpiredStr != "" {
		isExpired := isExpiredStr == "true"
		filter.IsExpired = &isExpired
	}

	if isExpiringSoonStr := c.Query("is_expiring_soon"); isExpiringSoonStr != "" {
		isExpiringSoon := isExpiringSoonStr == "true"
		filter.IsExpiringSoon = &isExpiringSoon
	}

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "name")
	filter.SortOrder = c.DefaultQuery("sort_order", "asc")

	items, total, err := h.inventoryUseCase.ListInventoryItems(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetInventoryItemsByCategory gets inventory items by category
func (h *InventoryHandler) GetInventoryItemsByCategory(c *gin.Context) {
	category := entities.ItemCategory(c.Param("category"))

	items, err := h.inventoryUseCase.GetInventoryItemsByCategory(c.Request.Context(), category)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetLowStockItems gets items with low stock
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	items, err := h.inventoryUseCase.GetLowStockItems(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetExpiredItems gets expired items
func (h *InventoryHandler) GetExpiredItems(c *gin.Context) {
	items, err := h.inventoryUseCase.GetExpiredItems(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetExpiringSoonItems gets items expiring soon
func (h *InventoryHandler) GetExpiringSoonItems(c *gin.Context) {
	items, err := h.inventoryUseCase.GetExpiringSoonItems(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemsNeedingReorder gets items that need reordering
func (h *InventoryHandler) GetItemsNeedingReorder(c *gin.Context) {
	items, err := h.inventoryUseCase.GetItemsNeedingReorder(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetInventoryStatistics gets inventory statistics
func (h *InventoryHandler) GetInventoryStatistics(c *gin.Context) {
	stats, err := h.inventoryUseCase.GetInventoryStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// AddStock adds stock to an inventory item
func (h *InventoryHandler) AddStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Quantity  float64 `json:"quantity" binding:"required,gt=0"`
		UnitCost  float64 `json:"unit_cost" binding:"required,gte=0"`
		Reference string  `json:"reference"`
		Notes     string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.AddStock(c.Request.Context(), id, req.Quantity, req.UnitCost, userID, req.Reference, req.Notes); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully"})
}

// RemoveStock removes stock from an inventory item
func (h *InventoryHandler) RemoveStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Quantity  float64 `json:"quantity" binding:"required,gt=0"`
		Reason    string  `json:"reason"`
		Reference string  `json:"reference"`
		Notes     string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.RemoveStock(c.Request.Context(), id, req.Quantity, userID, req.Reason, req.Reference, req.Notes); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed successfully"})
}

// AdjustStock adjusts stock (inventory count correction)
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		NewQuantity float64 `json:"new_quantity" binding:"required,gte=0"`
		Reason      string  `json:"reason" binding:"required"`
		Notes       string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.AdjustStock(c.Request.Context(), id, req.NewQuantity, userID, req.Reason, req.Notes); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock adjusted successfully"})
}

// ActivateItem activates an inventory item
func (h *InventoryHandler) ActivateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.ActivateItem(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item activated successfully"})
}

// DeactivateItem deactivates an inventory item
func (h *InventoryHandler) DeactivateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.inventoryUseCase.DeactivateItem(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deactivated successfully"})
}


// GetOutOfStockItems gets all out-of-stock items
func (h *InventoryHandler) GetOutOfStockItems(c *gin.Context) {
	// Return mock out-of-stock items list
	c.JSON(http.StatusOK, gin.H{
		"items": []interface{}{},
		"total": 0,
	})
}


// GetInventoryHistory gets history for an inventory item
func (h *InventoryHandler) GetInventoryHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"history": []interface{}{}, "total": 0})
}

// GetInventoryByCategory gets inventory items by category
func (h *InventoryHandler) GetInventoryByCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": 0})
}
