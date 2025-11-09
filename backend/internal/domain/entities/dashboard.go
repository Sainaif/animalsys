package entities

import (
	"time"
)

// DashboardMetrics represents the complete dashboard metrics
type DashboardMetrics struct {
	// Overview
	Overview OverviewMetrics `json:"overview"`

	// Module-specific metrics
	Animals     AnimalMetrics     `json:"animals"`
	Adoptions   AdoptionMetrics   `json:"adoptions"`
	Donations   DonationMetrics   `json:"donations"`
	Volunteers  VolunteerMetrics  `json:"volunteers"`
	Events      EventMetrics      `json:"events"`
	Veterinary  VeterinaryMetrics `json:"veterinary"`

	// Recent activity
	RecentActivity []Activity `json:"recent_activity"`

	// Trends
	Trends TrendMetrics `json:"trends"`

	// Generated at
	GeneratedAt time.Time `json:"generated_at"`
}

// OverviewMetrics represents high-level overview statistics
type OverviewMetrics struct {
	TotalAnimals        int64   `json:"total_animals"`
	AnimalsInShelter    int64   `json:"animals_in_shelter"`
	AnimalsAdopted      int64   `json:"animals_adopted"`
	AnimalsInFoster     int64   `json:"animals_in_foster"`

	TotalAdoptions      int64   `json:"total_adoptions"`
	AdoptionsThisMonth  int64   `json:"adoptions_this_month"`

	TotalDonations      float64 `json:"total_donations"`
	DonationsThisMonth  float64 `json:"donations_this_month"`

	ActiveVolunteers    int64   `json:"active_volunteers"`
	TotalVolunteerHours float64 `json:"total_volunteer_hours"`

	UpcomingEvents      int64   `json:"upcoming_events"`
	ActiveCampaigns     int64   `json:"active_campaigns"`
}

// AnimalMetrics represents animal-related statistics
type AnimalMetrics struct {
	Total              int64            `json:"total"`
	InShelter          int64            `json:"in_shelter"`
	Adopted            int64            `json:"adopted"`
	InFoster           int64            `json:"in_foster"`
	Quarantine         int64            `json:"quarantine"`
	Medical            int64            `json:"medical"`

	BySpecies          map[string]int64 `json:"by_species"`
	ByStatus           map[string]int64 `json:"by_status"`
	ByGender           map[string]int64 `json:"by_gender"`
	ByAgeGroup         map[string]int64 `json:"by_age_group"`

	AverageDaysInShelter float64        `json:"average_days_in_shelter"`
	NewArrivalsThisWeek  int64          `json:"new_arrivals_this_week"`
	NewArrivalsThisMonth int64          `json:"new_arrivals_this_month"`
}

// AdoptionMetrics represents adoption-related statistics
type AdoptionMetrics struct {
	TotalAdoptions      int64            `json:"total_adoptions"`
	ThisWeek            int64            `json:"this_week"`
	ThisMonth           int64            `json:"this_month"`
	ThisYear            int64            `json:"this_year"`

	PendingApplications int64            `json:"pending_applications"`
	ApprovedApplications int64           `json:"approved_applications"`

	BySpecies           map[string]int64 `json:"by_species"`
	SuccessRate         float64          `json:"success_rate"`
	AverageProcessingDays float64        `json:"average_processing_days"`

	PendingFollowUps    int64            `json:"pending_follow_ups"`
}

// DonationMetrics represents donation-related statistics
type DonationMetrics struct {
	TotalAmount         float64          `json:"total_amount"`
	ThisWeek            float64          `json:"this_week"`
	ThisMonth           float64          `json:"this_month"`
	ThisYear            float64          `json:"this_year"`

	TotalDonors         int64            `json:"total_donors"`
	NewDonorsThisMonth  int64            `json:"new_donors_this_month"`
	RecurringDonors     int64            `json:"recurring_donors"`
	MajorDonors         int64            `json:"major_donors"`

	AverageDonation     float64          `json:"average_donation"`
	LargestDonation     float64          `json:"largest_donation"`

	ByType              map[string]float64 `json:"by_type"`
	ByMethod            map[string]float64 `json:"by_method"`

	ActiveCampaigns     int64            `json:"active_campaigns"`
	CampaignGoalsTotal  float64          `json:"campaign_goals_total"`
	CampaignRaisedTotal float64          `json:"campaign_raised_total"`
}

// VolunteerMetrics represents volunteer-related statistics
type VolunteerMetrics struct {
	TotalVolunteers     int64            `json:"total_volunteers"`
	ActiveVolunteers    int64            `json:"active_volunteers"`
	NewThisMonth        int64            `json:"new_this_month"`

	TotalHours          float64          `json:"total_hours"`
	HoursThisWeek       float64          `json:"hours_this_week"`
	HoursThisMonth      float64          `json:"hours_this_month"`
	AverageHoursPerVolunteer float64     `json:"average_hours_per_volunteer"`

	TotalEvents         int64            `json:"total_events"`
	AverageRating       float64          `json:"average_rating"`

	ByStatus            map[string]int64 `json:"by_status"`
	TopVolunteers       []TopVolunteer   `json:"top_volunteers"`
}

// TopVolunteer represents a top volunteer for metrics
type TopVolunteer struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	TotalHours float64 `json:"total_hours"`
	Rating     float64 `json:"rating"`
}

// EventMetrics represents event-related statistics
type EventMetrics struct {
	TotalEvents         int64            `json:"total_events"`
	UpcomingEvents      int64            `json:"upcoming_events"`
	CompletedEvents     int64            `json:"completed_events"`

	TotalAttendees      int64            `json:"total_attendees"`
	TotalVolunteers     int64            `json:"total_volunteers"`

	ByType              map[string]int64 `json:"by_type"`
	ByStatus            map[string]int64 `json:"by_status"`

	AverageAttendance   float64          `json:"average_attendance"`
	TotalFundsRaised    float64          `json:"total_funds_raised"`
}

// VeterinaryMetrics represents veterinary-related statistics
type VeterinaryMetrics struct {
	TotalVisits         int64            `json:"total_visits"`
	ThisWeek            int64            `json:"this_week"`
	ThisMonth           int64            `json:"this_month"`
	UpcomingVisits      int64            `json:"upcoming_visits"`

	TotalVaccinations   int64            `json:"total_vaccinations"`
	VaccinationsDue     int64            `json:"vaccinations_due"`

	ByType              map[string]int64 `json:"by_type"`
	TotalCost           float64          `json:"total_cost"`
	AverageCostPerVisit float64          `json:"average_cost_per_visit"`
}

// Activity represents a recent activity item
type Activity struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id,omitempty"`
	UserName    string    `json:"user_name,omitempty"`
	EntityID    string    `json:"entity_id,omitempty"`
	EntityType  string    `json:"entity_type,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Icon        string    `json:"icon,omitempty"`
}

// TrendMetrics represents trend data over time
type TrendMetrics struct {
	Adoptions  []TrendDataPoint `json:"adoptions"`
	Donations  []TrendDataPoint `json:"donations"`
	Animals    []TrendDataPoint `json:"animals"`
	Volunteers []TrendDataPoint `json:"volunteers"`
}

// TrendDataPoint represents a single data point in a trend
type TrendDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Label string  `json:"label,omitempty"`
}

// ChartData represents generic chart data
type ChartData struct {
	Labels  []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

// Dataset represents a dataset for charts
type Dataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor string    `json:"background_color,omitempty"`
	BorderColor     string    `json:"border_color,omitempty"`
}
