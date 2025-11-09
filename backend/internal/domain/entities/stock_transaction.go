package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionType represents the type of stock transaction
type TransactionType string

const (
	TransactionTypeIn         TransactionType = "in"          // Stock added
	TransactionTypeOut        TransactionType = "out"         // Stock removed
	TransactionTypeAdjustment TransactionType = "adjustment"  // Inventory count adjustment
	TransactionTypeTransfer   TransactionType = "transfer"    // Transfer between locations
	TransactionTypeWaste      TransactionType = "waste"       // Expired or damaged items
	TransactionTypeReturn     TransactionType = "return"      // Return to supplier
	TransactionTypeDonation   TransactionType = "donation"    // Donated items received
)

// StockTransaction represents a stock movement transaction
type StockTransaction struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ItemID primitive.ObjectID `json:"item_id" bson:"item_id"`
	Type   TransactionType    `json:"type" bson:"type"`

	// Quantity
	Quantity      float64 `json:"quantity" bson:"quantity"`
	UnitCost      float64 `json:"unit_cost,omitempty" bson:"unit_cost,omitempty"`
	TotalCost     float64 `json:"total_cost,omitempty" bson:"total_cost,omitempty"`

	// Stock levels at time of transaction
	StockBefore float64 `json:"stock_before" bson:"stock_before"`
	StockAfter  float64 `json:"stock_after" bson:"stock_after"`

	// Transaction details
	Reason        string `json:"reason,omitempty" bson:"reason,omitempty"`
	Reference     string `json:"reference,omitempty" bson:"reference,omitempty"` // PO number, invoice, etc.

	// Related entities
	RelatedEntity   string              `json:"related_entity,omitempty" bson:"related_entity,omitempty"`         // "animal", "purchase_order", "supplier"
	RelatedEntityID *primitive.ObjectID `json:"related_entity_id,omitempty" bson:"related_entity_id,omitempty"`

	// Location
	FromLocation string `json:"from_location,omitempty" bson:"from_location,omitempty"`
	ToLocation   string `json:"to_location,omitempty" bson:"to_location,omitempty"`

	// Supplier info (for stock in)
	SupplierName string     `json:"supplier_name,omitempty" bson:"supplier_name,omitempty"`
	SupplierID   *primitive.ObjectID `json:"supplier_id,omitempty" bson:"supplier_id,omitempty"`

	// Expiration (for items with expiration)
	ExpirationDate *time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`

	// Tracking
	ProcessedBy primitive.ObjectID `json:"processed_by" bson:"processed_by"`
	ProcessedAt time.Time          `json:"processed_at" bson:"processed_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// NewStockTransaction creates a new stock transaction
func NewStockTransaction(itemID primitive.ObjectID, transactionType TransactionType, quantity, stockBefore float64, processedBy primitive.ObjectID) *StockTransaction {
	now := time.Now()
	stockAfter := stockBefore

	// Calculate stock after based on transaction type
	switch transactionType {
	case TransactionTypeIn, TransactionTypeDonation:
		stockAfter += quantity
	case TransactionTypeOut, TransactionTypeWaste, TransactionTypeReturn:
		stockAfter -= quantity
	}

	return &StockTransaction{
		ID:          primitive.NewObjectID(),
		ItemID:      itemID,
		Type:        transactionType,
		Quantity:    quantity,
		StockBefore: stockBefore,
		StockAfter:  stockAfter,
		ProcessedBy: processedBy,
		ProcessedAt: now,
		CreatedAt:   now,
	}
}

// SetCost sets the cost information
func (st *StockTransaction) SetCost(unitCost float64) {
	st.UnitCost = unitCost
	st.TotalCost = unitCost * st.Quantity
}

// SetSupplier sets the supplier information
func (st *StockTransaction) SetSupplier(supplierID primitive.ObjectID, supplierName string) {
	st.SupplierID = &supplierID
	st.SupplierName = supplierName
}

// SetLocation sets location information for transfers
func (st *StockTransaction) SetLocation(from, to string) {
	st.FromLocation = from
	st.ToLocation = to
}

// SetRelatedEntity sets the related entity
func (st *StockTransaction) SetRelatedEntity(entityType string, entityID primitive.ObjectID) {
	st.RelatedEntity = entityType
	st.RelatedEntityID = &entityID
}
