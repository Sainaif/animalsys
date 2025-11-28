package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryRepository defines the interface for inventory item data access
type InventoryRepository interface {
	Create(ctx context.Context, item *entities.InventoryItem) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error)
	FindBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error)
	Update(ctx context.Context, item *entities.InventoryItem) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter *InventoryFilter) ([]*entities.InventoryItem, int64, error)
	GetByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, error)
	GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error)
	GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error)
	GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error)
	GetInventoryStatistics(ctx context.Context) (*InventoryStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// StockTransactionRepository defines the interface for stock transaction data access
type StockTransactionRepository interface {
	Create(ctx context.Context, transaction *entities.StockTransaction) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.StockTransaction, error)
	List(ctx context.Context, filter *StockTransactionFilter) ([]*entities.StockTransaction, int64, error)
	GetByItem(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, error)
	GetByType(ctx context.Context, transactionType entities.TransactionType) ([]*entities.StockTransaction, error)
	GetTransactionStatistics(ctx context.Context) (*TransactionStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// InventoryFilter defines filter criteria for listing inventory items
type InventoryFilter struct {
	Category       string
	SubCategory    string
	IsActive       *bool
	IsLowStock     *bool
	IsExpired      *bool
	IsExpiringSoon *bool
	Location       string
	Search         string
	Tags           []string
	Limit          int64
	Offset         int64
	SortBy         string
	SortOrder      string
}

// StockTransactionFilter defines filter criteria for listing transactions
type StockTransactionFilter struct {
	ItemID          *primitive.ObjectID
	Type            string
	ProcessedBy     *primitive.ObjectID
	RelatedEntity   string
	RelatedEntityID *primitive.ObjectID
	DateFrom        string
	DateTo          string
	Limit           int64
	Offset          int64
	SortBy          string
	SortOrder       string
}

// InventoryStatistics represents inventory statistics
type InventoryStatistics struct {
	TotalItems          int64            `json:"total_items"`
	ActiveItems         int64            `json:"active_items"`
	ByCategory          map[string]int64 `json:"by_category"`
	LowStockItems       int64            `json:"low_stock_items"`
	ExpiredItems        int64            `json:"expired_items"`
	ExpiringSoonItems   int64            `json:"expiring_soon_items"`
	TotalValue          float64          `json:"total_value"`
	ItemsNeedingReorder int64            `json:"items_needing_reorder"`
}

// TransactionStatistics represents transaction statistics
type TransactionStatistics struct {
	TotalTransactions  int64            `json:"total_transactions"`
	ByType             map[string]int64 `json:"by_type"`
	TotalStockIn       float64          `json:"total_stock_in"`
	TotalStockOut      float64          `json:"total_stock_out"`
	TotalWaste         float64          `json:"total_waste"`
	TransactionsToday  int64            `json:"transactions_today"`
	TransactionsThisWeek int64          `json:"transactions_this_week"`
	TransactionsThisMonth int64         `json:"transactions_this_month"`
}
