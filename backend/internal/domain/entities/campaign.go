package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

// CampaignType represents the type of campaign
type CampaignType string

const (
	CampaignTypeGeneral      CampaignType = "general"
	CampaignTypeCapital      CampaignType = "capital"       // Building/facility
	CampaignTypeEmergency    CampaignType = "emergency"
	CampaignTypeAnnual       CampaignType = "annual"
	CampaignTypeEvent        CampaignType = "event"
	CampaignTypeMembership   CampaignType = "membership"
	CampaignTypeEndOfYear    CampaignType = "end_of_year"
)

// Campaign represents a fundraising campaign
type Campaign struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Name        MultilingualName `json:"name" bson:"name"`
	Description MultilingualName `json:"description" bson:"description"`
	Type        CampaignType       `json:"type" bson:"type"`
	Status      CampaignStatus     `json:"status" bson:"status"`

	// Goals and Tracking
	GoalAmount      float64 `json:"goal_amount" bson:"goal_amount"`
	CurrentAmount   float64 `json:"current_amount" bson:"current_amount"`
	DonorCount      int     `json:"donor_count" bson:"donor_count"`
	DonationCount   int     `json:"donation_count" bson:"donation_count"`
	AverageDonation float64 `json:"average_donation" bson:"average_donation"`

	// Timeline
	StartDate time.Time  `json:"start_date" bson:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`

	// Media and Content
	ImageURL    string   `json:"image_url,omitempty" bson:"image_url,omitempty"`
	VideoURL    string   `json:"video_url,omitempty" bson:"video_url,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`

	// Visibility
	Public   bool `json:"public" bson:"public"` // Publicly visible on website
	Featured bool `json:"featured" bson:"featured"` // Featured campaign

	// Contact and Management
	Manager      primitive.ObjectID  `json:"manager" bson:"manager"` // Staff managing campaign
	ContactEmail string              `json:"contact_email,omitempty" bson:"contact_email,omitempty"`
	ContactPhone string              `json:"contact_phone,omitempty" bson:"contact_phone,omitempty"`

	// Additional Information
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetProgressPercentage returns the progress towards goal as a percentage
func (c *Campaign) GetProgressPercentage() float64 {
	if c.GoalAmount == 0 {
		return 0
	}
	return (c.CurrentAmount / c.GoalAmount) * 100
}

// IsActive checks if the campaign is currently active
func (c *Campaign) IsActive() bool {
	if c.Status != CampaignStatusActive {
		return false
	}

	now := time.Now()
	if now.Before(c.StartDate) {
		return false
	}

	if c.EndDate != nil && now.After(*c.EndDate) {
		return false
	}

	return true
}

// IsCompleted checks if the campaign has reached its goal
func (c *Campaign) IsCompleted() bool {
	return c.Status == CampaignStatusCompleted || c.CurrentAmount >= c.GoalAmount
}

// UpdateStats updates campaign statistics
func (c *Campaign) UpdateStats(donationAmount float64, isNewDonor bool) {
	c.CurrentAmount += donationAmount
	c.DonationCount++

	if isNewDonor {
		c.DonorCount++
	}

	c.AverageDonation = c.CurrentAmount / float64(c.DonationCount)

	// Auto-complete if goal reached
	if c.CurrentAmount >= c.GoalAmount && c.Status == CampaignStatusActive {
		c.Status = CampaignStatusCompleted
	}
}

// NewCampaign creates a new campaign
func NewCampaign(
	name MultilingualName,
	goalAmount float64,
	campaignType CampaignType,
	manager primitive.ObjectID,
	createdBy primitive.ObjectID,
) *Campaign {
	now := time.Now()
	return &Campaign{
		Name:            name,
		Type:            campaignType,
		Status:          CampaignStatusDraft,
		GoalAmount:      goalAmount,
		CurrentAmount:   0,
		DonorCount:      0,
		DonationCount:   0,
		AverageDonation: 0,
		StartDate:       now,
		Public:          false,
		Featured:        false,
		Manager:         manager,
		CreatedBy:       createdBy,
		UpdatedBy:       createdBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
