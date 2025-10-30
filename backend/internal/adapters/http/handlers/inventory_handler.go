package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type InventoryHandler struct {
	inventoryUseCase *usecases.InventoryUseCase
}

func NewInventoryHandler(inventoryUseCase *usecases.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
	}
}

// Inventory Item endpoints

func (h *InventoryHandler) CreateItem(c *gin.Context) {
	var req entities.InventoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	item, err := h.inventoryUseCase.CreateItem(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *InventoryHandler) GetItemByID(c *gin.Context) {
	id := c.Param("id")

	item, err := h.inventoryUseCase.GetItemByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *InventoryHandler) ListItems(c *gin.Context) {
	filter := &entities.InventoryFilter{
		Category: entities.InventoryCategory(c.Query("category")),
		Search:   c.Query("search"),
		Limit:    parseIntQuery(c.Query("limit"), 10),
		Offset:   parseIntQuery(c.Query("offset"), 0),
	}

	items, total, err := h.inventoryUseCase.ListItems(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")

	var req entities.InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	item, err := h.inventoryUseCase.UpdateItem(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.inventoryUseCase.DeleteItem(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}

func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	items, err := h.inventoryUseCase.GetLowStockItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) GetExpiringItems(c *gin.Context) {
	days := parseIntQuery(c.Query("days"), 30)

	items, err := h.inventoryUseCase.GetExpiringItems(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) GetStatistics(c *gin.Context) {
	stats, err := h.inventoryUseCase.GetInventoryStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Stock Movement endpoints

func (h *InventoryHandler) RecordStockMovement(c *gin.Context) {
	var req entities.StockMovementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	recordedBy, _ := c.Get("user_id")
	movement, err := h.inventoryUseCase.RecordStockMovement(c.Request.Context(), &req, recordedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movement)
}

func (h *InventoryHandler) GetItemStockMovements(c *gin.Context) {
	itemID := c.Param("id")
	limit := parseIntQuery(c.Query("limit"), 10)
	offset := parseIntQuery(c.Query("offset"), 0)

	movements, total, err := h.inventoryUseCase.GetItemStockMovements(c.Request.Context(), itemID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   movements,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	itemID := c.Param("id")

	var req struct {
		NewQuantity int    `json:"new_quantity" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	adjustedBy, _ := c.Get("user_id")
	err := h.inventoryUseCase.AdjustStock(c.Request.Context(), itemID, req.NewQuantity, req.Reason, adjustedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock adjusted successfully"})
}
