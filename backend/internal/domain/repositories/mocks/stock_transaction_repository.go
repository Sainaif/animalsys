package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockTransactionRepository struct {
	mock.Mock
}

func (m *StockTransactionRepository) Create(ctx context.Context, transaction *entities.StockTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *StockTransactionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.StockTransaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.StockTransaction), args.Error(1)
}

func (m *StockTransactionRepository) List(ctx context.Context, filter *repositories.StockTransactionFilter) ([]*entities.StockTransaction, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.StockTransaction), args.Get(1).(int64), args.Error(2)
}

func (m *StockTransactionRepository) GetByItem(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]*entities.StockTransaction), args.Error(1)
}

func (m *StockTransactionRepository) GetByType(ctx context.Context, transactionType entities.TransactionType) ([]*entities.StockTransaction, error) {
	args := m.Called(ctx, transactionType)
	return args.Get(0).([]*entities.StockTransaction), args.Error(1)
}

func (m *StockTransactionRepository) GetTransactionStatistics(ctx context.Context) (*repositories.TransactionStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(*repositories.TransactionStatistics), args.Error(1)
}

func (m *StockTransactionRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
