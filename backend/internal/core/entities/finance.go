package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FinanceType represents the type of financial transaction
type FinanceType string

const (
	FinanceTypeIncome  FinanceType = "income"
	FinanceTypeExpense FinanceType = "expense"
)

// FinanceCategory represents the category of financial transaction
type FinanceCategory string

// Income categories
const (
	IncomeCategoryDonations    FinanceCategory = "donations"
	IncomeCategoryGrants       FinanceCategory = "grants"
	IncomeCategoryFundraising  FinanceCategory = "fundraising"
	IncomeCategoryAdoptionFees FinanceCategory = "adoption_fees"
	IncomeCategoryOther        FinanceCategory = "other_income"
)

// Expense categories
const (
	ExpenseCategoryVeterinary ExpenseCategory = "veterinary"
	ExpenseCategoryFood       ExpenseCategory = "food"
	ExpenseCategoryUtilities  ExpenseCategory = "utilities"
	ExpenseCategorySalaries   ExpenseCategory = "salaries"
	ExpenseCategorySupplies   ExpenseCategory = "supplies"
	ExpenseCategoryMaintenance ExpenseCategory = "maintenance"
	ExpenseCategoryOther      ExpenseCategory = "other_expense"
)

type ExpenseCategory string

// Finance represents a financial transaction
type Finance struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type            FinanceType        `bson:"type" json:"type"`
	Category        string             `bson:"category" json:"category"`
	Subcategory     string             `bson:"subcategory,omitempty" json:"subcategory,omitempty"`
	Amount          float64            `bson:"amount" json:"amount"`
	Currency        string             `bson:"currency" json:"currency"`
	Date            time.Time          `bson:"date" json:"date"`
	Description     string             `bson:"description" json:"description"`
	PaymentMethod   string             `bson:"payment_method,omitempty" json:"payment_method,omitempty"` // cash, card, bank_transfer, etc.
	ReferenceNumber string             `bson:"reference_number,omitempty" json:"reference_number,omitempty"`
	ReceiptURL      string             `bson:"receipt_url,omitempty" json:"receipt_url,omitempty"`
	RelatedEntityType string           `bson:"related_entity_type,omitempty" json:"related_entity_type,omitempty"`
	RelatedEntityID string             `bson:"related_entity_id,omitempty" json:"related_entity_id,omitempty"`
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	FiscalYear      int                `bson:"fiscal_year" json:"fiscal_year"`
	Quarter         int                `bson:"quarter" json:"quarter"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy       string             `bson:"created_by" json:"created_by"`
}

// FinanceCreateRequest represents finance record creation
type FinanceCreateRequest struct {
	Type            FinanceType `json:"type" validate:"required,oneof=income expense"`
	Category        string      `json:"category" validate:"required"`
	Subcategory     string      `json:"subcategory,omitempty"`
	Amount          float64     `json:"amount" validate:"required,gt=0"`
	Currency        string      `json:"currency" validate:"required,len=3"`
	Date            time.Time   `json:"date" validate:"required"`
	Description     string      `json:"description" validate:"required"`
	PaymentMethod   string      `json:"payment_method,omitempty"`
	ReferenceNumber string      `json:"reference_number,omitempty"`
	RelatedEntityType string    `json:"related_entity_type,omitempty"`
	RelatedEntityID string      `json:"related_entity_id,omitempty"`
	Notes           string      `json:"notes,omitempty"`
}

// FinanceUpdateRequest represents finance record update
type FinanceUpdateRequest struct {
	Category        string      `json:"category,omitempty"`
	Subcategory     string      `json:"subcategory,omitempty"`
	Amount          float64     `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Date            time.Time   `json:"date,omitempty"`
	Description     string      `json:"description,omitempty"`
	PaymentMethod   string      `json:"payment_method,omitempty"`
	ReferenceNumber string      `json:"reference_number,omitempty"`
	Notes           string      `json:"notes,omitempty"`
}

// FinanceFilter represents filters for querying finances
type FinanceFilter struct {
	Type       FinanceType
	Category   string
	StartDate  time.Time
	EndDate    time.Time
	MinAmount  float64
	MaxAmount  float64
	FiscalYear int
	Quarter    int
	Search     string
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}

// NewFinance creates a new finance record
func NewFinance(financeType FinanceType, category string, amount float64, currency string, date time.Time, description, createdBy string) *Finance {
	now := time.Now()
	fiscalYear := date.Year()
	quarter := (int(date.Month()) + 2) / 3

	return &Finance{
		ID:          primitive.NewObjectID(),
		Type:        financeType,
		Category:    category,
		Amount:      amount,
		Currency:    currency,
		Date:        date,
		Description: description,
		FiscalYear:  fiscalYear,
		Quarter:     quarter,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   createdBy,
	}
}
