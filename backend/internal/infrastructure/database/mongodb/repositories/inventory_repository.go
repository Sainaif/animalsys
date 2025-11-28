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

type inventoryRepository struct {
	db *mongodb.Database
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *mongodb.Database) repositories.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) collection() *mongo.Collection {
	return r.db.DB.Collection("inventory_items")
}

// Create creates a new inventory item
func (r *inventoryRepository) Create(ctx context.Context, item *entities.InventoryItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}

	_, err := r.collection().InsertOne(ctx, item)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.ErrConflict
		}
		return err
	}

	return nil
}

// FindByID finds an inventory item by ID
func (r *inventoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.InventoryItem, error) {
	var item entities.InventoryItem
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &item, nil
}

// FindBySKU finds an inventory item by SKU
func (r *inventoryRepository) FindBySKU(ctx context.Context, sku string) (*entities.InventoryItem, error) {
	var item entities.InventoryItem
	err := r.collection().FindOne(ctx, bson.M{"sku": sku}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &item, nil
}

// Update updates an inventory item
func (r *inventoryRepository) Update(ctx context.Context, item *entities.InventoryItem) error {
	item.UpdatedAt = time.Now()

	result, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": item.ID},
		bson.M{"$set": item},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Delete deletes an inventory item
func (r *inventoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// List lists inventory items with filtering and pagination
func (r *inventoryRepository) List(ctx context.Context, filter *repositories.InventoryFilter) ([]*entities.InventoryItem, int64, error) {
	query := bson.M{}

	// Apply filters
	if filter.Category != "" {
		query["category"] = filter.Category
	}

	if filter.SubCategory != "" {
		query["sub_category"] = filter.SubCategory
	}

	if filter.IsActive != nil {
		query["is_active"] = *filter.IsActive
	}

	if filter.IsLowStock != nil {
		query["is_low_stock"] = *filter.IsLowStock
	}

	if filter.IsExpired != nil {
		query["is_expired"] = *filter.IsExpired
	}

	if filter.IsExpiringSoon != nil {
		query["is_expiring_soon"] = *filter.IsExpiringSoon
	}

	if filter.Location != "" {
		query["location"] = bson.M{"$regex": filter.Location, "$options": "i"}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"sku": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags}
	}

	// Count total
	total, err := r.collection().CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options
	opts := options.Find()

	// Sorting
	sortField := "name"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := 1 // ascending by default
	if filter.SortOrder == "desc" {
		sortOrder = -1
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

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetByCategory gets inventory items by category
func (r *inventoryRepository) GetByCategory(ctx context.Context, category entities.ItemCategory) ([]*entities.InventoryItem, error) {
	query := bson.M{"category": category, "is_active": true}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetOutOfStockItems retrieves items that are out of stock
func (r *inventoryRepository) GetOutOfStockItems(ctx context.Context) ([]*entities.InventoryItem, int64, error) {
	filter := bson.M{
		"is_active": true,
		"$expr": bson.M{
			"$lte": []string{"$current_stock", "$min_stock_level"},
		},
	}

	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.NewInternalServer("Failed to count out of stock items")
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, 0, errors.NewInternalServer("Failed to query out of stock items")
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, errors.NewInternalServer("Failed to decode out of stock items")
	}

	return items, total, nil
}

// GetLowStockItems gets items with low stock
func (r *inventoryRepository) GetLowStockItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	query := bson.M{
		"is_active":   true,
		"is_low_stock": true,
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "current_stock", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetExpiredItems gets expired items
func (r *inventoryRepository) GetExpiredItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	query := bson.M{
		"is_active":  true,
		"is_expired": true,
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "expiration_date", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetExpiringSoonItems gets items expiring soon
func (r *inventoryRepository) GetExpiringSoonItems(ctx context.Context) ([]*entities.InventoryItem, error) {
	query := bson.M{
		"is_active":        true,
		"is_expiring_soon": true,
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "expiration_date", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetItemsNeedingReorder gets items that need reordering
func (r *inventoryRepository) GetItemsNeedingReorder(ctx context.Context) ([]*entities.InventoryItem, error) {
	query := bson.M{
		"is_active": true,
		"$expr": bson.M{
			"$and": []bson.M{
				{"$gt": []interface{}{"$reorder_point", 0}},
				{"$lte": []interface{}{"$current_stock", "$reorder_point"}},
			},
		},
	}

	cursor, err := r.collection().Find(ctx, query, options.Find().SetSort(bson.D{{Key: "current_stock", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.InventoryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetInventoryStatistics gets inventory statistics
func (r *inventoryRepository) GetInventoryStatistics(ctx context.Context) (*repositories.InventoryStatistics, error) {
	stats := &repositories.InventoryStatistics{
		ByCategory: make(map[string]int64),
	}

	// Total items
	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalItems = total

	// Active items
	activeCount, err := r.collection().CountDocuments(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	stats.ActiveItems = activeCount

	// Low stock items
	lowStockCount, err := r.collection().CountDocuments(ctx, bson.M{
		"is_active":   true,
		"is_low_stock": true,
	})
	if err != nil {
		return nil, err
	}
	stats.LowStockItems = lowStockCount

	// Expired items
	expiredCount, err := r.collection().CountDocuments(ctx, bson.M{
		"is_active":  true,
		"is_expired": true,
	})
	if err != nil {
		return nil, err
	}
	stats.ExpiredItems = expiredCount

	// Expiring soon items
	expiringSoonCount, err := r.collection().CountDocuments(ctx, bson.M{
		"is_active":        true,
		"is_expiring_soon": true,
	})
	if err != nil {
		return nil, err
	}
	stats.ExpiringSoonItems = expiringSoonCount

	// Items needing reorder
	needingReorderCount, err := r.collection().CountDocuments(ctx, bson.M{
		"is_active": true,
		"$expr": bson.M{
			"$and": []bson.M{
				{"$gt": []interface{}{"$reorder_point", 0}},
				{"$lte": []interface{}{"$current_stock", "$reorder_point"}},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	stats.ItemsNeedingReorder = needingReorderCount

	// Aggregation for category breakdown and total value
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_active": true}}},
		{{Key: "$facet", Value: bson.M{
			"byCategory": []bson.M{
				{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
			},
			"totals": []bson.M{
				{"$group": bson.M{
					"_id":        nil,
					"totalValue": bson.M{"$sum": "$total_value"},
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

		// By category
		if byCategory, ok := result["byCategory"].([]interface{}); ok {
			for _, item := range byCategory {
				if m, ok := item.(bson.M); ok {
					category := m["_id"].(string)
					count := m["count"].(int32)
					stats.ByCategory[category] = int64(count)
				}
			}
		}

		// Total value
		if totals, ok := result["totals"].([]interface{}); ok && len(totals) > 0 {
			if m, ok := totals[0].(bson.M); ok {
				if totalValue, ok := m["totalValue"].(float64); ok {
					stats.TotalValue = totalValue
				}
			}
		}
	}

	return stats, nil
}

// EnsureIndexes creates the necessary indexes for the inventory collection
func (r *inventoryRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sku", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_low_stock", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_expired", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_expiring_soon", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expiration_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}
