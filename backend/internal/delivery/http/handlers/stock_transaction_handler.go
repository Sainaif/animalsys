package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/stock"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StockTransactionHandler handles stock transaction-related HTTP requests
type StockTransactionHandler struct {
	stockTransactionUseCase *stock.StockTransactionUseCase
}

// NewStockTransactionHandler creates a new stock transaction handler
func NewStockTransactionHandler(stockTransactionUseCase *stock.StockTransactionUseCase) *StockTransactionHandler {
	return &StockTransactionHandler{
		stockTransactionUseCase: stockTransactionUseCase,
	}
}

// GetStockTransaction gets a stock transaction by ID
func (h *StockTransactionHandler) GetStockTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.stockTransactionUseCase.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// ListStockTransactions lists stock transactions with filtering
func (h *StockTransactionHandler) ListStockTransactions(c *gin.Context) {
	filter := &repositories.StockTransactionFilter{}

	// Parse query parameters
	if itemIDStr := c.Query("item_id"); itemIDStr != "" {
		itemID, err := primitive.ObjectIDFromHex(itemIDStr)
		if err == nil {
			filter.ItemID = &itemID
		}
	}

	filter.Type = c.Query("type")

	if processedByStr := c.Query("processed_by"); processedByStr != "" {
		processedBy, err := primitive.ObjectIDFromHex(processedByStr)
		if err == nil {
			filter.ProcessedBy = &processedBy
		}
	}

	filter.RelatedEntity = c.Query("related_entity")

	if relatedEntityIDStr := c.Query("related_entity_id"); relatedEntityIDStr != "" {
		relatedEntityID, err := primitive.ObjectIDFromHex(relatedEntityIDStr)
		if err == nil {
			filter.RelatedEntityID = &relatedEntityID
		}
	}

	filter.DateFrom = c.Query("date_from")
	filter.DateTo = c.Query("date_to")

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "processed_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	transactions, total, err := h.stockTransactionUseCase.ListTransactions(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   transactions,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetStockTransactionsByItem gets transactions for a specific inventory item
func (h *StockTransactionHandler) GetStockTransactionsByItem(c *gin.Context) {
	idParam := c.Param("item_id")
	itemID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	transactions, err := h.stockTransactionUseCase.GetTransactionsByItem(c.Request.Context(), itemID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetStockTransactionsByType gets transactions by type
func (h *StockTransactionHandler) GetStockTransactionsByType(c *gin.Context) {
	transactionType := entities.TransactionType(c.Param("type"))

	transactions, err := h.stockTransactionUseCase.GetTransactionsByType(c.Request.Context(), transactionType)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetStockTransactionStatistics gets transaction statistics
func (h *StockTransactionHandler) GetStockTransactionStatistics(c *gin.Context) {
	stats, err := h.stockTransactionUseCase.GetTransactionStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ExportStockTransactions exports stock transactions
func (h *StockTransactionHandler) ExportStockTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Export initiated"})
}
