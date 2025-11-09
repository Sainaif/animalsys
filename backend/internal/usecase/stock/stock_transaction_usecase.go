package stock

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockTransactionUseCase struct {
	stockTransactionRepo repositories.StockTransactionRepository
	inventoryRepo        repositories.InventoryRepository
}

func NewStockTransactionUseCase(
	stockTransactionRepo repositories.StockTransactionRepository,
	inventoryRepo repositories.InventoryRepository,
) *StockTransactionUseCase {
	return &StockTransactionUseCase{
		stockTransactionRepo: stockTransactionRepo,
		inventoryRepo:        inventoryRepo,
	}
}

// GetTransactionByID retrieves a stock transaction by ID
func (uc *StockTransactionUseCase) GetTransactionByID(ctx context.Context, id primitive.ObjectID) (*entities.StockTransaction, error) {
	transaction, err := uc.stockTransactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// ListTransactions lists stock transactions with filtering and pagination
func (uc *StockTransactionUseCase) ListTransactions(ctx context.Context, filter *repositories.StockTransactionFilter) ([]*entities.StockTransaction, int64, error) {
	return uc.stockTransactionRepo.List(ctx, filter)
}

// GetTransactionsByItem gets transactions for a specific inventory item
func (uc *StockTransactionUseCase) GetTransactionsByItem(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, error) {
	return uc.stockTransactionRepo.GetByItem(ctx, itemID)
}

// GetTransactionsByType gets transactions by type
func (uc *StockTransactionUseCase) GetTransactionsByType(ctx context.Context, transactionType entities.TransactionType) ([]*entities.StockTransaction, error) {
	return uc.stockTransactionRepo.GetByType(ctx, transactionType)
}

// GetTransactionStatistics gets transaction statistics
func (uc *StockTransactionUseCase) GetTransactionStatistics(ctx context.Context) (*repositories.TransactionStatistics, error) {
	return uc.stockTransactionRepo.GetTransactionStatistics(ctx)
}
