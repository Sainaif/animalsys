package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type inventoryRepository struct {
	collection *mongo.Collection
}

func NewInventoryRepository(db *mongo.Database) interfaces.InventoryRepository {
	return &inventoryRepository{
		collection: db.Collection("inventory"),
	}
}

func (r *inventoryRepository) Create(ctx context.Context, item *entities.InventoryItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *inventoryRepository) GetByID(ctx context.Context, id string) (*entities.InventoryItem, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid inventory item ID")
	}

	var item entities.InventoryItem
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("inventory item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) Update(ctx context.Context, id string, item *entities.InventoryItem) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid inventory item ID")
	}

	item.UpdatedAt = time.Now()
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": item})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("inventory item not found")
	}
	return nil
}

func (r *inventoryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid inventory item ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("inventory item not found")
	}
	return nil
}

func (r *inventoryRepository) List(ctx context.Context, itemType entities.InventoryItemType, limit, offset int) ([]*entities.InventoryItem, int64, error) {
	filter := bson.M{}
	if itemType != "" {
		filter["type"] = itemType
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"name": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
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

func (r *inventoryRepository) GetLowStock(ctx context.Context) ([]*entities.InventoryItem, error) {
	pipeline := mongo.Pipeline{
		{{"$addFields", bson.M{
			"is_low_stock": bson.M{"$lte": []interface{}{"$stock_level", "$min_stock"}},
		}}},
		{{"$match", bson.M{"is_low_stock": true}}},
		{{"$sort", bson.M{"stock_level": 1}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
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

func (r *inventoryRepository) GetOutOfStock(ctx context.Context) ([]*entities.InventoryItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"stock_level": 0})
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

func (r *inventoryRepository) GetExpiringSoon(ctx context.Context, days int) ([]*entities.InventoryItem, error) {
	expiryDate := time.Now().AddDate(0, 0, days)

	filter := bson.M{
		"expiry_date": bson.M{
			"$lte": expiryDate,
			"$gte": time.Now(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"expiry_date": 1}))
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

func (r *inventoryRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid inventory item ID")
	}

	update := bson.M{
		"$set": bson.M{
			"stock_level": quantity,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("inventory item not found")
	}
	return nil
}

type stockMovementRepository struct {
	collection *mongo.Collection
}

func NewStockMovementRepository(db *mongo.Database) interfaces.StockMovementRepository {
	return &stockMovementRepository{
		collection: db.Collection("stock_movements"),
	}
}

func (r *stockMovementRepository) Create(ctx context.Context, movement *entities.StockMovement) error {
	movement.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, movement)
	return err
}

func (r *stockMovementRepository) GetByID(ctx context.Context, id string) (*entities.StockMovement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid stock movement ID")
	}

	var movement entities.StockMovement
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("stock movement not found")
		}
		return nil, err
	}
	return &movement, nil
}

func (r *stockMovementRepository) GetByItemID(ctx context.Context, itemID string, limit, offset int) ([]*entities.StockMovement, int64, error) {
	filter := bson.M{"item_id": itemID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var movements []*entities.StockMovement
	if err = cursor.All(ctx, &movements); err != nil {
		return nil, 0, err
	}
	return movements, total, nil
}

func (r *stockMovementRepository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.StockMovement, error) {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	filter := bson.M{"created_at": bson.M{"$gte": start, "$lte": end}}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movements []*entities.StockMovement
	if err = cursor.All(ctx, &movements); err != nil {
		return nil, err
	}
	return movements, nil
}
