package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonationType represents the type of donation
type DonationType string

const (
	DonationTypeMonetary  DonationType = "monetary"
	DonationTypeInKind    DonationType = "in_kind"    // Goods/services
	DonationTypeRecurring DonationType = "recurring"
	DonationTypeMemorial  DonationType = "memorial"   // In memory of
	DonationTypeHonor     DonationType = "honor"      // In honor of
	DonationTypeMatching  DonationType = "matching"   // Employer matching
)

// DonationStatus represents the status of a donation
type DonationStatus string

const (
	DonationStatusPending   DonationStatus = "pending"
	DonationStatusCompleted DonationStatus = "completed"
	DonationStatusFailed    DonationStatus = "failed"
	DonationStatusRefunded  DonationStatus = "refunded"
	DonationStatusCancelled DonationStatus = "cancelled"
)

// PaymentMethodType represents the payment method
type PaymentMethodType string

const (
	PaymentMethodCash         PaymentMethodType = "cash"
	PaymentMethodCheck        PaymentMethodType = "check"
	PaymentMethodCreditCard   PaymentMethodType = "credit_card"
	PaymentMethodDebitCard    PaymentMethodType = "debit_card"
	PaymentMethodBankTransfer PaymentMethodType = "bank_transfer"
	PaymentMethodPayPal       PaymentMethodType = "paypal"
	PaymentMethodVenmo        PaymentMethodType = "venmo"
	PaymentMethodOther        PaymentMethodType = "other"
)

// RecurrenceFrequency represents how often recurring donations occur
type RecurrenceFrequency string

const (
	RecurrenceWeekly    RecurrenceFrequency = "weekly"
	RecurrenceBiWeekly  RecurrenceFrequency = "bi_weekly"
	RecurrenceMonthly   RecurrenceFrequency = "monthly"
	RecurrenceQuarterly RecurrenceFrequency = "quarterly"
	RecurrenceYearly    RecurrenceFrequency = "yearly"
)

// PaymentInfo represents payment details
type PaymentInfo struct {
	Method            PaymentMethodType `json:"method" bson:"method"`
	TransactionID     string            `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	CheckNumber       string            `json:"check_number,omitempty" bson:"check_number,omitempty"`
	LastFourDigits    string            `json:"last_four_digits,omitempty" bson:"last_four_digits,omitempty"` // For cards
	ProcessorResponse string            `json:"processor_response,omitempty" bson:"processor_response,omitempty"`
}

// RecurringInfo represents recurring donation details
type RecurringInfo struct {
	Frequency     RecurrenceFrequency `json:"frequency" bson:"frequency"`
	StartDate     time.Time           `json:"start_date" bson:"start_date"`
	EndDate       *time.Time          `json:"end_date,omitempty" bson:"end_date,omitempty"`
	NextBillingDate *time.Time        `json:"next_billing_date,omitempty" bson:"next_billing_date,omitempty"`
	Active        bool                `json:"active" bson:"active"`
	FailureCount  int                 `json:"failure_count" bson:"failure_count"`
}

// InKindItem represents an in-kind donation item
type InKindItem struct {
	Description    string  `json:"description" bson:"description"`
	Quantity       int     `json:"quantity" bson:"quantity"`
	EstimatedValue float64 `json:"estimated_value" bson:"estimated_value"`
	Category       string  `json:"category,omitempty" bson:"category,omitempty"` // e.g., "food", "supplies", "equipment"
}

// TaxReceipt represents tax receipt information
type TaxReceipt struct {
	ReceiptNumber string     `json:"receipt_number,omitempty" bson:"receipt_number,omitempty"`
	SentDate      *time.Time `json:"sent_date,omitempty" bson:"sent_date,omitempty"`
	SentMethod    string     `json:"sent_method,omitempty" bson:"sent_method,omitempty"` // email, mail
	ReceiptURL    string     `json:"receipt_url,omitempty" bson:"receipt_url,omitempty"`
}

// Donation represents a donation to the foundation
type Donation struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Donor Information
	DonorID    primitive.ObjectID `json:"donor_id" bson:"donor_id"`
	DonorName  string             `json:"donor_name" bson:"donor_name"`  // Cached for performance
	DonorEmail string             `json:"donor_email" bson:"donor_email"` // Cached for performance
	Anonymous  bool               `json:"anonymous" bson:"anonymous"`

	// Donation Details
	Type              DonationType   `json:"type" bson:"type"`
	Status            DonationStatus `json:"status" bson:"status"`
	Amount            float64        `json:"amount" bson:"amount"`
	Currency          string         `json:"currency" bson:"currency"` // Default: USD
	DonationDate      time.Time      `json:"donation_date" bson:"donation_date"`

	// Campaign/Designation
	CampaignID   *primitive.ObjectID `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"`
	CampaignName string              `json:"campaign_name,omitempty" bson:"campaign_name,omitempty"`
	Designation  string              `json:"designation,omitempty" bson:"designation,omitempty"` // e.g., "general", "medical", "building"
	Restricted   bool                `json:"restricted" bson:"restricted"` // Restricted to specific use

	// Payment Information
	Payment      PaymentInfo `json:"payment" bson:"payment"`
	PaymentDate  *time.Time  `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	Fee          float64     `json:"fee" bson:"fee"`           // Processing fee
	NetAmount    float64     `json:"net_amount" bson:"net_amount"` // Amount - Fee

	// Recurring Donation
	IsRecurring   bool           `json:"is_recurring" bson:"is_recurring"`
	RecurringInfo *RecurringInfo `json:"recurring_info,omitempty" bson:"recurring_info,omitempty"`

	// In-Kind Donation
	InKindItems []InKindItem `json:"in_kind_items,omitempty" bson:"in_kind_items,omitempty"`

	// Memorial/Honor Information
	InMemoryOf string `json:"in_memory_of,omitempty" bson:"in_memory_of,omitempty"`
	InHonorOf  string `json:"in_honor_of,omitempty" bson:"in_honor_of,omitempty"`
	NotifyName string `json:"notify_name,omitempty" bson:"notify_name,omitempty"`     // Who to notify
	NotifyEmail string `json:"notify_email,omitempty" bson:"notify_email,omitempty"` // Email to notify

	// Tax Receipt
	TaxDeductible bool       `json:"tax_deductible" bson:"tax_deductible"`
	TaxReceipt    TaxReceipt `json:"tax_receipt" bson:"tax_receipt"`

	// Acknowledgment
	ThankYouSent     bool       `json:"thank_you_sent" bson:"thank_you_sent"`
	ThankYouSentDate *time.Time `json:"thank_you_sent_date,omitempty" bson:"thank_you_sent_date,omitempty"`

	// Matching Gift
	MatchingGift         bool    `json:"matching_gift" bson:"matching_gift"`
	MatchingCompany      string  `json:"matching_company,omitempty" bson:"matching_company,omitempty"`
	MatchingAmount       float64 `json:"matching_amount,omitempty" bson:"matching_amount,omitempty"`
	MatchingSubmitted    bool    `json:"matching_submitted" bson:"matching_submitted"`

	// Additional Information
	Notes         string   `json:"notes,omitempty" bson:"notes,omitempty"`
	Source        string   `json:"source,omitempty" bson:"source,omitempty"` // online, event, mail, etc.
	Tags          []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Attachments   []string `json:"attachments,omitempty" bson:"attachments,omitempty"` // URLs to documents

	// Metadata
	ProcessedBy primitive.ObjectID `json:"processed_by" bson:"processed_by"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy   primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsCompleted checks if the donation is completed
func (d *Donation) IsCompleted() bool {
	return d.Status == DonationStatusCompleted
}

// CalculateNetAmount calculates the net amount after fees
func (d *Donation) CalculateNetAmount() {
	d.NetAmount = d.Amount - d.Fee
	if d.NetAmount < 0 {
		d.NetAmount = 0
	}
}

// GenerateReceiptNumber generates a unique receipt number
func (d *Donation) GenerateReceiptNumber() string {
	year := d.DonationDate.Year()
	return primitive.NewObjectID().Hex()[:8] + "-" + string(rune(year))
}

// NewDonation creates a new donation
func NewDonation(
	donorID primitive.ObjectID,
	amount float64,
	donationType DonationType,
	createdBy primitive.ObjectID,
) *Donation {
	now := time.Now()
	donation := &Donation{
		DonorID:       donorID,
		Type:          donationType,
		Status:        DonationStatusPending,
		Amount:        amount,
		Currency:      "USD",
		DonationDate:  now,
		TaxDeductible: true,
		ThankYouSent:  false,
		MatchingGift:  false,
		Restricted:    false,
		ProcessedBy:   createdBy,
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	donation.CalculateNetAmount()
	return donation
}
