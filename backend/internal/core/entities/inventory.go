package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryItemType represents the type of inventory item
type InventoryItemType string

const (
	InventoryTypeFood      InventoryItemType = "food"
	InventoryTypeMedicine  InventoryItemType = "medicine"
	InventoryTypeSupplies  InventoryItemType = "supplies"
	InventoryTypeEquipment InventoryItemType = "equipment"
)

// InventoryItem represents an inventory item
type InventoryItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Type         InventoryItemType  `bson:"type" json:"type"`
	Category     string             `bson:"category,omitempty" json:"category,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	SKU          string             `bson:"sku,omitempty" json:"sku,omitempty"`
	StockLevel   int                `bson:"stock_level" json:"stock_level"`
	MinStock     int                `bson:"min_stock" json:"min_stock"`
	MaxStock     int                `bson:"max_stock,omitempty" json:"max_stock,omitempty"`
	Unit         string             `bson:"unit" json:"unit"` // kg, liters, pieces, etc.
	UnitPrice    float64            `bson:"unit_price,omitempty" json:"unit_price,omitempty"`
	Supplier     string             `bson:"supplier,omitempty" json:"supplier,omitempty"`
	ExpiryDate   *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Location     string             `bson:"location,omitempty" json:"location,omitempty"`
	LastRestocked *time.Time        `bson:"last_restocked,omitempty" json:"last_restocked,omitempty"`
	Notes        string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy    string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
}

// StockMovement represents a stock movement transaction
type StockMovement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ItemID      string             `bson:"item_id" json:"item_id"`
	ItemName    string             `bson:"item_name" json:"item_name"`
	Type        string             `bson:"type" json:"type"` // in, out, transfer, adjustment
	Quantity    int                `bson:"quantity" json:"quantity"`
	FromLocation string            `bson:"from_location,omitempty" json:"from_location,omitempty"`
	ToLocation  string             `bson:"to_location,omitempty" json:"to_location,omitempty"`
	Reason      string             `bson:"reason,omitempty" json:"reason,omitempty"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
}

// InventoryItemCreateRequest represents inventory item creation
type InventoryItemCreateRequest struct {
	Name        string            `json:"name" validate:"required"`
	Type        InventoryItemType `json:"type" validate:"required,oneof=food medicine supplies equipment"`
	Category    string            `json:"category,omitempty"`
	Description string            `json:"description,omitempty"`
	SKU         string            `json:"sku,omitempty"`
	StockLevel  int               `json:"stock_level" validate:"required,gte=0"`
	MinStock    int               `json:"min_stock" validate:"required,gte=0"`
	MaxStock    int               `json:"max_stock,omitempty" validate:"omitempty,gtfield=MinStock"`
	Unit        string            `json:"unit" validate:"required"`
	UnitPrice   float64           `json:"unit_price,omitempty" validate:"omitempty,gte=0"`
	Supplier    string            `json:"supplier,omitempty"`
	ExpiryDate  *time.Time        `json:"expiry_date,omitempty"`
	Location    string            `json:"location,omitempty"`
	Notes       string            `json:"notes,omitempty"`
}

// StockMovementRequest represents stock movement request
type StockMovementRequest struct {
	ItemID       string `json:"item_id" validate:"required"`
	Type         string `json:"type" validate:"required,oneof=in out transfer adjustment"`
	Quantity     int    `json:"quantity" validate:"required,gt=0"`
	FromLocation string `json:"from_location,omitempty"`
	ToLocation   string `json:"to_location,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// NewInventoryItem creates a new inventory item
func NewInventoryItem(name string, itemType InventoryItemType, stockLevel, minStock int, unit, createdBy string) *InventoryItem {
	now := time.Now()
	return &InventoryItem{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Type:       itemType,
		StockLevel: stockLevel,
		MinStock:   minStock,
		Unit:       unit,
		CreatedAt:  now,
		UpdatedAt:  now,
		CreatedBy:  createdBy,
	}
}

// IsLowStock checks if item is low on stock
func (i *InventoryItem) IsLowStock() bool {
	return i.StockLevel <= i.MinStock
}

// IsOutOfStock checks if item is out of stock
func (i *InventoryItem) IsOutOfStock() bool {
	return i.StockLevel == 0
}

// AddStock adds stock to item
func (i *InventoryItem) AddStock(quantity int) {
	i.StockLevel += quantity
	now := time.Now()
	i.LastRestocked = &now
	i.UpdatedAt = now
}

// RemoveStock removes stock from item
func (i *InventoryItem) RemoveStock(quantity int) error {
	if i.StockLevel < quantity {
		return primitive.ErrNoDocuments // or custom error
	}
	i.StockLevel -= quantity
	i.UpdatedAt = time.Now()
	return nil
}
