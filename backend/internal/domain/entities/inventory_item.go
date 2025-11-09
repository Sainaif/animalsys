package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ItemCategory represents the category of an inventory item
type ItemCategory string

const (
	ItemCategoryFood           ItemCategory = "food"
	ItemCategoryMedicine       ItemCategory = "medicine"
	ItemCategorySupplies       ItemCategory = "supplies"
	ItemCategoryToys           ItemCategory = "toys"
	ItemCategoryBedding        ItemCategory = "bedding"
	ItemCategoryCleaning       ItemCategory = "cleaning"
	ItemCategoryEquipment      ItemCategory = "equipment"
	ItemCategoryGrooming       ItemCategory = "grooming"
	ItemCategoryMedicalSupplies ItemCategory = "medical_supplies"
	ItemCategoryOffice         ItemCategory = "office"
	ItemCategoryOther          ItemCategory = "other"
)

// ItemUnit represents the unit of measurement
type ItemUnit string

const (
	ItemUnitPiece      ItemUnit = "piece"
	ItemUnitBox        ItemUnit = "box"
	ItemUnitBag        ItemUnit = "bag"
	ItemUnitBottle     ItemUnit = "bottle"
	ItemUnitCan        ItemUnit = "can"
	ItemUnitPound      ItemUnit = "lb"
	ItemUnitKilogram   ItemUnit = "kg"
	ItemUnitLiter      ItemUnit = "liter"
	ItemUnitGallon     ItemUnit = "gallon"
	ItemUnitPackage    ItemUnit = "package"
	ItemUnitRoll       ItemUnit = "roll"
	ItemUnitCase       ItemUnit = "case"
)

// InventoryItem represents an item in inventory
type InventoryItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	Barcode     string             `json:"barcode,omitempty" bson:"barcode,omitempty"`
	Category    ItemCategory       `json:"category" bson:"category"`
	SubCategory string             `json:"sub_category,omitempty" bson:"sub_category,omitempty"`

	// Unit and quantity
	Unit            ItemUnit `json:"unit" bson:"unit"`
	CurrentStock    float64  `json:"current_stock" bson:"current_stock"`
	MinimumStock    float64  `json:"minimum_stock" bson:"minimum_stock"`
	ReorderPoint    float64  `json:"reorder_point" bson:"reorder_point"`
	MaximumStock    float64  `json:"maximum_stock,omitempty" bson:"maximum_stock,omitempty"`
	ReorderQuantity float64  `json:"reorder_quantity" bson:"reorder_quantity"`

	// Cost and pricing
	UnitCost      float64 `json:"unit_cost" bson:"unit_cost"`
	TotalValue    float64 `json:"total_value" bson:"total_value"` // current_stock * unit_cost
	LastPurchasePrice float64 `json:"last_purchase_price,omitempty" bson:"last_purchase_price,omitempty"`

	// Supplier
	PreferredSupplier   string `json:"preferred_supplier,omitempty" bson:"preferred_supplier,omitempty"`
	SupplierProductCode string `json:"supplier_product_code,omitempty" bson:"supplier_product_code,omitempty"`

	// Storage
	Location       string `json:"location,omitempty" bson:"location,omitempty"` // Warehouse, Shelf A, etc.
	StorageConditions string `json:"storage_conditions,omitempty" bson:"storage_conditions,omitempty"`

	// Expiration
	HasExpiration bool       `json:"has_expiration" bson:"has_expiration"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`

	// Tracking
	IsActive      bool   `json:"is_active" bson:"is_active"`
	IsLowStock    bool   `json:"is_low_stock" bson:"is_low_stock"`
	IsExpired     bool   `json:"is_expired" bson:"is_expired"`
	IsExpiringSoon bool  `json:"is_expiring_soon" bson:"is_expiring_soon"` // Within 30 days

	// Usage tracking
	TotalUsed       float64    `json:"total_used" bson:"total_used"`
	LastUsedDate    *time.Time `json:"last_used_date,omitempty" bson:"last_used_date,omitempty"`
	LastRestockDate *time.Time `json:"last_restock_date,omitempty" bson:"last_restock_date,omitempty"`

	// Notes
	Notes string   `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags  []string `json:"tags,omitempty" bson:"tags,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewInventoryItem creates a new inventory item
func NewInventoryItem(name string, category ItemCategory, unit ItemUnit, createdBy primitive.ObjectID) *InventoryItem {
	now := time.Now()
	return &InventoryItem{
		ID:              primitive.NewObjectID(),
		Name:            name,
		Category:        category,
		Unit:            unit,
		CurrentStock:    0,
		MinimumStock:    0,
		ReorderPoint:    0,
		ReorderQuantity: 0,
		UnitCost:        0,
		TotalValue:      0,
		HasExpiration:   false,
		IsActive:        true,
		IsLowStock:      false,
		IsExpired:       false,
		IsExpiringSoon:  false,
		TotalUsed:       0,
		Tags:            []string{},
		CreatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// AddStock adds stock to the item
func (i *InventoryItem) AddStock(quantity, unitCost float64) {
	i.CurrentStock += quantity
	i.UnitCost = unitCost
	i.TotalValue = i.CurrentStock * i.UnitCost
	now := time.Now()
	i.LastRestockDate = &now
	i.UpdatedAt = now
	i.CheckLowStock()
}

// RemoveStock removes stock from the item
func (i *InventoryItem) RemoveStock(quantity float64) bool {
	if i.CurrentStock < quantity {
		return false
	}
	i.CurrentStock -= quantity
	i.TotalValue = i.CurrentStock * i.UnitCost
	i.TotalUsed += quantity
	now := time.Now()
	i.LastUsedDate = &now
	i.UpdatedAt = now
	i.CheckLowStock()
	return true
}

// CheckLowStock updates the low stock flag
func (i *InventoryItem) CheckLowStock() {
	if i.ReorderPoint > 0 {
		i.IsLowStock = i.CurrentStock <= i.ReorderPoint
	} else if i.MinimumStock > 0 {
		i.IsLowStock = i.CurrentStock <= i.MinimumStock
	}
}

// CheckExpiration checks and updates expiration status
func (i *InventoryItem) CheckExpiration() {
	if !i.HasExpiration || i.ExpirationDate == nil {
		i.IsExpired = false
		i.IsExpiringSoon = false
		return
	}

	now := time.Now()
	i.IsExpired = now.After(*i.ExpirationDate)

	thirtyDaysFromNow := now.AddDate(0, 0, 30)
	i.IsExpiringSoon = i.ExpirationDate.Before(thirtyDaysFromNow) && !i.IsExpired
	i.UpdatedAt = now
}

// Deactivate deactivates the item
func (i *InventoryItem) Deactivate() {
	i.IsActive = false
	i.UpdatedAt = time.Now()
}

// Activate activates the item
func (i *InventoryItem) Activate() {
	i.IsActive = true
	i.UpdatedAt = time.Now()
}

// NeedsReorder checks if the item needs to be reordered
func (i *InventoryItem) NeedsReorder() bool {
	if i.ReorderPoint > 0 {
		return i.CurrentStock <= i.ReorderPoint
	}
	return false
}

// GetStockPercentage returns the stock percentage relative to maximum
func (i *InventoryItem) GetStockPercentage() float64 {
	if i.MaximumStock <= 0 {
		return 0
	}
	return (i.CurrentStock / i.MaximumStock) * 100
}
