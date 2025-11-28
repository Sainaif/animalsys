package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonorRepository defines the interface for donor data access
type DonorRepository interface {
	// Create creates a new donor
	Create(ctx context.Context, donor *entities.Donor) error

	// FindByID finds a donor by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error)

	// FindManyByIDs finds multiple donors by their IDs
	FindManyByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*entities.Donor, error)

	// Update updates an existing donor
	Update(ctx context.Context, donor *entities.Donor) error

	// Delete deletes a donor by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of donors with pagination and filters
	List(ctx context.Context, filter *DonorFilter) ([]*entities.Donor, int64, error)

	// FindByEmail finds a donor by email
	FindByEmail(ctx context.Context, email string) (*entities.Donor, error)

	// GetMajorDonors returns all major donors
	GetMajorDonors(ctx context.Context) ([]*entities.Donor, error)

	// GetLapsedDonors returns donors who haven't donated in N days
	GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error)

	// GetDonorStatistics returns donor statistics
	GetDonorStatistics(ctx context.Context) (*DonorStatistics, error)

	// EnsureIndexes creates necessary indexes for the donors collection
	EnsureIndexes(ctx context.Context) error
}

// DonorFilter defines filter criteria for listing donors
type DonorFilter struct {
	Type            string
	Status          string
	Search          string // Search in name, email, organization
	MinTotalDonated *float64
	MaxTotalDonated *float64
	Tags            []string
	FromDate        *time.Time // Created from date
	ToDate          *time.Time // Created to date
	Limit           int64
	Offset          int64
	SortBy          string // Field to sort by
	SortOrder       string // "asc" or "desc"
}

// DonorStatistics represents donor statistics
type DonorStatistics struct {
	TotalDonors          int64            `json:"total_donors" bson:"total_donors"`
	ActiveDonors         int64            `json:"active_donors" bson:"active_donors"`
	InactiveDonors       int64            `json:"inactive_donors" bson:"inactive_donors"`
	MajorDonors          int64            `json:"major_donors" bson:"major_donors"`
	LapsedDonors         int64            `json:"lapsed_donors" bson:"lapsed_donors"`
	IndividualDonors     int64            `json:"individual_donors" bson:"individual_donors"`
	OrganizationDonors   int64            `json:"organization_donors" bson:"organization_donors"`
	AverageLifetimeValue float64          `json:"average_lifetime_value" bson:"average_lifetime_value"`
	TotalLifetimeValue   float64          `json:"total_lifetime_value" bson:"total_lifetime_value"`
	ByType               map[string]int64 `json:"by_type" bson:"by_type"`
	ByStatus             map[string]int64 `json:"by_status" bson:"by_status"`
}

// DonationRepository defines the interface for donation data access
type DonationRepository interface {
	// Create creates a new donation
	Create(ctx context.Context, donation *entities.Donation) error

	// FindByID finds a donation by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Donation, error)

	// Update updates an existing donation
	Update(ctx context.Context, donation *entities.Donation) error

	// Delete deletes a donation by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of donations with pagination and filters
	List(ctx context.Context, filter *DonationFilter) ([]*entities.Donation, int64, error)

	// GetByDonorID returns all donations for a specific donor
	GetByDonorID(ctx context.Context, donorID primitive.ObjectID) ([]*entities.Donation, error)

	// GetByCampaignID returns all donations for a specific campaign
	GetByCampaignID(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Donation, error)

	// AggregateDonorsByCampaignID aggregates donor data for a campaign
	AggregateDonorsByCampaignID(ctx context.Context, campaignID primitive.ObjectID, limit, offset int64) ([]*DonorAggregationResult, int64, error)

	// GetRecurringDonations returns all active recurring donations
	GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error)

	// GetPendingThankYous returns donations without thank you sent
	GetPendingThankYous(ctx context.Context) ([]*entities.Donation, error)

	// GetDonationStatistics returns donation statistics
	GetDonationStatistics(ctx context.Context) (*DonationStatistics, error)

	// GetDonationsByDateRange returns donations within a date range
	GetDonationsByDateRange(ctx context.Context, from, to time.Time) ([]*entities.Donation, error)

	// EnsureIndexes creates necessary indexes for the donations collection
	EnsureIndexes(ctx context.Context) error
}

// DonorAggregationResult represents the result of a donor aggregation query
type DonorAggregationResult struct {
	DonorID     primitive.ObjectID `bson:"_id"`
	TotalAmount float64            `bson:"total_amount"`
}

// DonationFilter defines filter criteria for listing donations
type DonationFilter struct {
	DonorID       *primitive.ObjectID
	CampaignID    *primitive.ObjectID
	Type          string
	Status        string
	MinAmount     *float64
	MaxAmount     *float64
	FromDate      *time.Time
	ToDate        *time.Time
	PaymentMethod string
	Designation   string
	IsRecurring   *bool
	Anonymous     *bool
	Limit         int64
	Offset        int64
	SortBy        string // Field to sort by
	SortOrder     string // "asc" or "desc"
}

// DonationStatistics represents donation statistics
type DonationStatistics struct {
	TotalDonations     int64            `json:"total_donations" bson:"total_donations"`
	TotalAmount        float64          `json:"total_amount" bson:"total_amount"`
	AverageDonation    float64          `json:"average_donation" bson:"average_donation"`
	LargestDonation    float64          `json:"largest_donation" bson:"largest_donation"`
	SmallestDonation   float64          `json:"smallest_donation" bson:"smallest_donation"`
	DonationsToday     int64            `json:"donations_today" bson:"donations_today"`
	DonationsThisWeek  int64            `json:"donations_this_week" bson:"donations_this_week"`
	DonationsThisMonth int64            `json:"donations_this_month" bson:"donations_this_month"`
	DonationsThisYear  int64            `json:"donations_this_year" bson:"donations_this_year"`
	AmountToday        float64          `json:"amount_today" bson:"amount_today"`
	AmountThisWeek     float64          `json:"amount_this_week" bson:"amount_this_week"`
	AmountThisMonth    float64          `json:"amount_this_month" bson:"amount_this_month"`
	AmountThisYear     float64          `json:"amount_this_year" bson:"amount_this_year"`
	RecurringDonations int64            `json:"recurring_donations" bson:"recurring_donations"`
	RecurringAmount    float64          `json:"recurring_amount" bson:"recurring_amount"`
	ByType             map[string]int64 `json:"by_type" bson:"by_type"`
	ByStatus           map[string]int64 `json:"by_status" bson:"by_status"`
	ByPaymentMethod    map[string]int64 `json:"by_payment_method" bson:"by_payment_method"`
}

// CampaignRepository defines the interface for campaign data access
type CampaignRepository interface {
	// Create creates a new campaign
	Create(ctx context.Context, campaign *entities.Campaign) error

	// FindByID finds a campaign by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Campaign, error)

	// Update updates an existing campaign
	Update(ctx context.Context, campaign *entities.Campaign) error

	// Delete deletes a campaign by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of campaigns with pagination and filters
	List(ctx context.Context, filter *CampaignFilter) ([]*entities.Campaign, int64, error)

	// GetActiveCampaigns returns all active campaigns
	GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error)

	// GetFeaturedCampaigns returns featured campaigns
	GetFeaturedCampaigns(ctx context.Context) ([]*entities.Campaign, error)

	// GetPublicCampaigns returns public campaigns
	GetPublicCampaigns(ctx context.Context) ([]*entities.Campaign, error)

	// GetCampaignsByManager returns campaigns managed by a specific user
	GetCampaignsByManager(ctx context.Context, managerID primitive.ObjectID) ([]*entities.Campaign, error)

	// UpdateCampaignStats updates campaign statistics after a donation
	UpdateCampaignStats(ctx context.Context, campaignID primitive.ObjectID, donationAmount float64, isNewDonor bool) error

	// GetCampaignStatistics returns campaign statistics
	GetCampaignStatistics(ctx context.Context) (*CampaignStatistics, error)

	// EnsureIndexes creates necessary indexes for the campaigns collection
	EnsureIndexes(ctx context.Context) error
}

// CampaignFilter defines filter criteria for listing campaigns
type CampaignFilter struct {
	Type          string
	Status        string
	Public        *bool
	Featured      *bool
	ManagerID     *primitive.ObjectID
	StartDateFrom *time.Time
	StartDateTo   *time.Time
	EndDateFrom   *time.Time
	EndDateTo     *time.Time
	GoalAmountMin float64
	GoalAmountMax float64
	Tags          []string
	Search        string
	Limit         int64
	Offset        int64
	SortBy        string // Field to sort by
	SortOrder     string // "asc" or "desc"
}

// CampaignStatistics represents campaign statistics
type CampaignStatistics struct {
	TotalCampaigns    int64            `json:"total_campaigns" bson:"total_campaigns"`
	ActiveCampaigns   int64            `json:"active_campaigns" bson:"active_campaigns"`
	TotalGoalAmount   float64          `json:"total_goal_amount" bson:"total_goal_amount"`
	TotalRaisedAmount float64          `json:"total_raised_amount" bson:"total_raised_amount"`
	AverageProgress   float64          `json:"average_progress" bson:"average_progress"`
	TotalDonors       int64            `json:"total_donors" bson:"total_donors"`
	ByType            map[string]int64 `json:"by_type" bson:"by_type"`
	ByStatus          map[string]int64 `json:"by_status" bson:"by_status"`
}
