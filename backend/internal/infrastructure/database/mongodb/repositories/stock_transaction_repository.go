package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/database/mongodb"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type stockTransactionRepository struct {
	db *mongodb.Database
}

// NewStockTransactionRepository creates a new stock transaction repository
func NewStockTransactionRepository(db *mongodb.Database) repositories.StockTransactionRepository {
	return &stockTransactionRepository{db: db}
}

func (r *stockTransactionRepository) collection() *mongo.Collection {
	return r.db.DB.Collection("stock_transactions")
}

// Create creates a new stock transaction
func (r *stockTransactionRepository) Create(ctx context.Context, transaction *entities.StockTransaction) error {
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}

	_, err := r.collection().InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a stock transaction by ID
func (r *stockTransactionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.StockTransaction, error) {
	var transaction entities.StockTransaction
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &transaction, nil
}

// List lists stock transactions with filtering and pagination
func (r *stockTransactionRepository) List(ctx context.Context, filter *repositories.StockTransactionFilter) ([]*entities.StockTransaction, int64, error) {
	query := bson.M{}

	// Apply filters
	if filter.ItemID != nil {
		query["item_id"] = filter.ItemID
	}

	if filter.Type != "" {
		query["type"] = filter.Type
	}

	if filter.ProcessedBy != nil {
		query["processed_by"] = filter.ProcessedBy
	}

	if filter.RelatedEntity != "" {
		query["related_entity"] = filter.RelatedEntity
	}

	if filter.RelatedEntityID != nil {
		query["related_entity_id"] = filter.RelatedEntityID
	}

	// Date filters
	if filter.DateFrom != "" || filter.DateTo != "" {
		dateQuery := bson.M{}
		if filter.DateFrom != "" {
			dateFrom, err := time.Parse(time.RFC3339, filter.DateFrom)
			if err == nil {
				dateQuery["$gte"] = dateFrom
			}
		}
		if filter.DateTo != "" {
			dateTo, err := time.Parse(time.RFC3339, filter.DateTo)
			if err == nil {
				dateQuery["$lte"] = dateTo
			}
		}
		if len(dateQuery) > 0 {
			query["processed_at"] = dateQuery
		}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options
	opts := options.Find()

	// Sorting
	sortField := "processed_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1 // descending by default
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Pagination
	if filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	// Execute query
	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var transactions []*entities.StockTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetByItem gets transactions for a specific inventory item
func (r *stockTransactionRepository) GetByItem(ctx context.Context, itemID primitive.ObjectID) ([]*entities.StockTransaction, error) {
	query := bson.M{"item_id": itemID}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "processed_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*entities.StockTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetByType gets transactions by type
func (r *stockTransactionRepository) GetByType(ctx context.Context, transactionType entities.TransactionType) ([]*entities.StockTransaction, error) {
	query := bson.M{"type": transactionType}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "processed_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*entities.StockTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionStatistics gets transaction statistics
func (r *stockTransactionRepository) GetTransactionStatistics(ctx context.Context) (*repositories.TransactionStatistics, error) {
	stats := &repositories.TransactionStatistics{
		ByType: make(map[string]int64),
	}

	// Total transactions
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalTransactions = total

	// Transactions today
	today := time.Now().Truncate(24 * time.Hour)
	todayCount, err := r.collection().CountDocuments(ctx, bson.M{
		"processed_at": bson.M{"$gte": today},
	})
	if err != nil {
		return nil, err
	}
	stats.TransactionsToday = todayCount

	// Transactions this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	weekCount, err := r.collection().CountDocuments(ctx, bson.M{
		"processed_at": bson.M{"$gte": weekAgo},
	})
	if err != nil {
		return nil, err
	}
	stats.TransactionsThisWeek = weekCount

	// Transactions this month
	monthAgo := time.Now().AddDate(0, -1, 0)
	monthCount, err := r.collection().CountDocuments(ctx, bson.M{
		"processed_at": bson.M{"$gte": monthAgo},
	})
	if err != nil {
		return nil, err
	}
	stats.TransactionsThisMonth = monthCount

	// Aggregation for type breakdown and totals
	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.M{
			"byType": []bson.M{
				{"$group": bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}},
			},
			"stockIn": []bson.M{
				{"$match": bson.M{"type": bson.M{"$in": []entities.TransactionType{
					entities.TransactionTypeIn,
					entities.TransactionTypeDonation,
				}}}},
				{"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": "$quantity"},
				}},
			},
			"stockOut": []bson.M{
				{"$match": bson.M{"type": bson.M{"$in": []entities.TransactionType{
					entities.TransactionTypeOut,
					entities.TransactionTypeReturn,
				}}}},
				{"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": "$quantity"},
				}},
			},
			"waste": []bson.M{
				{"$match": bson.M{"type": entities.TransactionTypeWaste}},
				{"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": "$quantity"},
				}},
			},
		}}},
	}

	cursor, err := r.collection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		result := results[0]

		// By type
		if byType, ok := result["byType"].([]interface{}); ok {
			for _, item := range byType {
				if m, ok := item.(bson.M); ok {
					transactionType := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByType[transactionType] = int64(count)
				}
			}
		}

		// Stock in
		if stockIn, ok := result["stockIn"].([]interface{}); ok && len(stockIn) > 0 {
			if m, ok := stockIn[0].(bson.M); ok {
				if total, ok := m["total"].(float64); ok {
					stats.TotalStockIn = total
				}
			}
		}

		// Stock out
		if stockOut, ok := result["stockOut"].([]interface{}); ok && len(stockOut) > 0 {
			if m, ok := stockOut[0].(bson.M); ok {
				if total, ok := m["total"].(float64); ok {
					stats.TotalStockOut = total
				}
			}
		}

		// Waste
		if waste, ok := result["waste"].([]interface{}); ok && len(waste) > 0 {
			if m, ok := waste[0].(bson.M); ok {
				if total, ok := m["total"].(float64); ok {
					stats.TotalWaste = total
				}
			}
		}
	}

	return stats, nil
}

// EnsureIndexes creates the necessary indexes for the stock transactions collection
func (r *stockTransactionRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "item_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "processed_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "processed_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "related_entity", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
