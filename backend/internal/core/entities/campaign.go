package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CampaignType represents the type of campaign
type CampaignType string

const (
	CampaignTypeFundraising  CampaignType = "fundraising"
	CampaignTypeAdoptionDrive CampaignType = "adoption_drive"
	CampaignTypeEvent        CampaignType = "event"
	CampaignTypeAwareness    CampaignType = "awareness"
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusPlanning   CampaignStatus = "planning"
	CampaignStatusActive     CampaignStatus = "active"
	CampaignStatusCompleted  CampaignStatus = "completed"
	CampaignStatusCancelled  CampaignStatus = "cancelled"
)

// Campaign represents a campaign
type Campaign struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Type            CampaignType       `bson:"type" json:"type"`
	Status          CampaignStatus     `bson:"status" json:"status"`
	Description     string             `bson:"description" json:"description"`
	StartDate       time.Time          `bson:"start_date" json:"start_date"`
	EndDate         time.Time          `bson:"end_date" json:"end_date"`
	Goal            float64            `bson:"goal,omitempty" json:"goal,omitempty"` // financial goal or # of adoptions
	GoalType        string             `bson:"goal_type,omitempty" json:"goal_type,omitempty"` // financial, adoptions, awareness
	CurrentProgress float64            `bson:"current_progress" json:"current_progress"`
	Budget          float64            `bson:"budget,omitempty" json:"budget,omitempty"`
	TotalExpenses   float64            `bson:"total_expenses" json:"total_expenses"`
	Milestones      []Milestone        `bson:"milestones,omitempty" json:"milestones,omitempty"`
	ImageURL        string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy       string             `bson:"created_by" json:"created_by"`
}

// Milestone represents a campaign milestone
type Milestone struct {
	Title       string     `bson:"title" json:"title"`
	Description string     `bson:"description,omitempty" json:"description,omitempty"`
	TargetDate  time.Time  `bson:"target_date" json:"target_date"`
	Completed   bool       `bson:"completed" json:"completed"`
	CompletedDate *time.Time `bson:"completed_date,omitempty" json:"completed_date,omitempty"`
}

// CampaignCreateRequest represents campaign creation
type CampaignCreateRequest struct {
	Name        string         `json:"name" validate:"required"`
	Type        CampaignType   `json:"type" validate:"required,oneof=fundraising adoption_drive event awareness"`
	Description string         `json:"description" validate:"required"`
	StartDate   time.Time      `json:"start_date" validate:"required"`
	EndDate     time.Time      `json:"end_date" validate:"required,gtfield=StartDate"`
	Goal        float64        `json:"goal,omitempty" validate:"omitempty,gt=0"`
	GoalType    string         `json:"goal_type,omitempty"`
	Budget      float64        `json:"budget,omitempty" validate:"omitempty,gte=0"`
	Milestones  []Milestone    `json:"milestones,omitempty"`
}

// CampaignUpdateRequest represents campaign update
type CampaignUpdateRequest struct {
	Name            string          `json:"name,omitempty"`
	Status          CampaignStatus  `json:"status,omitempty" validate:"omitempty,oneof=planning active completed cancelled"`
	Description     string          `json:"description,omitempty"`
	StartDate       time.Time       `json:"start_date,omitempty"`
	EndDate         time.Time       `json:"end_date,omitempty"`
	Goal            float64         `json:"goal,omitempty" validate:"omitempty,gt=0"`
	CurrentProgress float64         `json:"current_progress,omitempty" validate:"omitempty,gte=0"`
	Budget          float64         `json:"budget,omitempty" validate:"omitempty,gte=0"`
}

// NewCampaign creates a new campaign
func NewCampaign(name string, campaignType CampaignType, description string, startDate, endDate time.Time, createdBy string) *Campaign {
	now := time.Now()
	return &Campaign{
		ID:              primitive.NewObjectID(),
		Name:            name,
		Type:            campaignType,
		Status:          CampaignStatusPlanning,
		Description:     description,
		StartDate:       startDate,
		EndDate:         endDate,
		CurrentProgress: 0,
		TotalExpenses:   0,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       createdBy,
	}
}

// UpdateProgress updates campaign progress
func (c *Campaign) UpdateProgress(progress float64) {
	c.CurrentProgress = progress
	c.UpdatedAt = time.Now()
}

// AddExpense adds an expense to campaign
func (c *Campaign) AddExpense(amount float64) {
	c.TotalExpenses += amount
	c.UpdatedAt = time.Now()
}

// GetProgressPercentage calculates progress percentage
func (c *Campaign) GetProgressPercentage() float64 {
	if c.Goal == 0 {
		return 0
	}
	return (c.CurrentProgress / c.Goal) * 100
}

// IsActive checks if campaign is currently active
func (c *Campaign) IsActive() bool {
	now := time.Now()
	return c.Status == CampaignStatusActive &&
		now.After(c.StartDate) &&
		now.Before(c.EndDate)
}
