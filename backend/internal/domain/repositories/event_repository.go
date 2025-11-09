package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventFilter represents filters for event queries
type EventFilter struct {
	Type      string
	Status    string
	Search    string
	Public    *bool
	Featured  *bool
	StartDate *time.Time
	EndDate   *time.Time
	SortBy    string
	SortOrder string
	Limit     int64
	Offset    int64
}

// EventStatistics represents event statistics
type EventStatistics struct {
	TotalEvents         int64
	UpcomingEvents      int64
	ActiveEvents        int64
	CompletedEvents     int64
	TotalAttendees      int64
	TotalVolunteers     int64
	TotalFundsRaised    float64
	TotalAnimalsAdopted int64
	AverageAttendees    float64
	ByType              map[string]int64
	ByStatus            map[string]int64
}

// EventRepository defines the interface for event data access
type EventRepository interface {
	Create(ctx context.Context, event *entities.Event) error
	Update(ctx context.Context, event *entities.Event) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Event, error)
	List(ctx context.Context, filter *EventFilter) ([]*entities.Event, int64, error)
	GetUpcomingEvents(ctx context.Context) ([]*entities.Event, error)
	GetActiveEvents(ctx context.Context) ([]*entities.Event, error)
	GetCompletedEvents(ctx context.Context, limit int) ([]*entities.Event, error)
	GetPublicEvents(ctx context.Context) ([]*entities.Event, error)
	GetFeaturedEvents(ctx context.Context) ([]*entities.Event, error)
	GetEventsByOrganizer(ctx context.Context, organizerID primitive.ObjectID) ([]*entities.Event, error)
	GetEventsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Event, error)
	GetEventsNeedingVolunteers(ctx context.Context) ([]*entities.Event, error)
	UpdateEventStatistics(ctx context.Context, eventID primitive.ObjectID, attendees, volunteers int, fundsRaised float64, animalsAdopted int) error
	GetEventStatistics(ctx context.Context) (*EventStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// VolunteerFilter represents filters for volunteer queries
type VolunteerFilter struct {
	Status     string
	Skills     []string
	Search     string
	HasUserID  *bool
	SortBy     string
	SortOrder  string
	Limit      int64
	Offset     int64
}

// VolunteerStatistics represents volunteer statistics
type VolunteerStatistics struct {
	TotalVolunteers     int64
	ActiveVolunteers    int64
	InactiveVolunteers  int64
	TotalHours          float64
	AverageHours        float64
	TotalEvents         int64
	AverageRating       float64
	ByStatus            map[string]int64
	TopVolunteers       []*entities.Volunteer // Top 10 by hours
}

// VolunteerRepository defines the interface for volunteer data access
type VolunteerRepository interface {
	Create(ctx context.Context, volunteer *entities.Volunteer) error
	Update(ctx context.Context, volunteer *entities.Volunteer) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Volunteer, error)
	FindByEmail(ctx context.Context, email string) (*entities.Volunteer, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*entities.Volunteer, error)
	List(ctx context.Context, filter *VolunteerFilter) ([]*entities.Volunteer, int64, error)
	GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error)
	GetVolunteersBySkill(ctx context.Context, skill string) ([]*entities.Volunteer, error)
	GetVolunteersNeedingBackgroundCheck(ctx context.Context) ([]*entities.Volunteer, error)
	GetVolunteersWithExpiredCertifications(ctx context.Context) ([]*entities.Volunteer, error)
	GetTopVolunteers(ctx context.Context, limit int) ([]*entities.Volunteer, error)
	UpdateHours(ctx context.Context, volunteerID primitive.ObjectID, hours float64) error
	IncrementEventsAttended(ctx context.Context, volunteerID primitive.ObjectID) error
	GetVolunteerStatistics(ctx context.Context) (*VolunteerStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// EventAttendanceFilter represents filters for event attendance queries
type EventAttendanceFilter struct {
	EventID      *primitive.ObjectID
	VolunteerID  *primitive.ObjectID
	UserID       *primitive.ObjectID
	DonorID      *primitive.ObjectID
	AttendeeType string
	Status       string
	SortBy       string
	SortOrder    string
	Limit        int64
	Offset       int64
}

// EventAttendanceRepository defines the interface for event attendance data access
type EventAttendanceRepository interface {
	Create(ctx context.Context, attendance *entities.EventAttendance) error
	Update(ctx context.Context, attendance *entities.EventAttendance) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.EventAttendance, error)
	List(ctx context.Context, filter *EventAttendanceFilter) ([]*entities.EventAttendance, int64, error)
	GetAttendanceByEvent(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error)
	GetAttendanceByVolunteer(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.EventAttendance, error)
	GetAttendanceByDonor(ctx context.Context, donorID primitive.ObjectID) ([]*entities.EventAttendance, error)
	GetConfirmedAttendees(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error)
	GetNoShows(ctx context.Context, eventID primitive.ObjectID) ([]*entities.EventAttendance, error)
	GetPendingPayments(ctx context.Context) ([]*entities.EventAttendance, error)
	CountAttendeesByEvent(ctx context.Context, eventID primitive.ObjectID) (int64, error)
	EnsureIndexes(ctx context.Context) error
}

// VolunteerAssignmentFilter represents filters for volunteer assignment queries
type VolunteerAssignmentFilter struct {
	VolunteerID *primitive.ObjectID
	EventID     *primitive.ObjectID
	AnimalID    *primitive.ObjectID
	CampaignID  *primitive.ObjectID
	Type        string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
	SortBy      string
	SortOrder   string
	Limit       int64
	Offset      int64
}

// VolunteerAssignmentRepository defines the interface for volunteer assignment data access
type VolunteerAssignmentRepository interface {
	Create(ctx context.Context, assignment *entities.VolunteerAssignment) error
	Update(ctx context.Context, assignment *entities.VolunteerAssignment) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.VolunteerAssignment, error)
	List(ctx context.Context, filter *VolunteerAssignmentFilter) ([]*entities.VolunteerAssignment, int64, error)
	GetAssignmentsByVolunteer(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	GetAssignmentsByEvent(ctx context.Context, eventID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	GetUpcomingAssignments(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	GetActiveAssignments(ctx context.Context, volunteerID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	GetAssignmentsNeedingReminder(ctx context.Context) ([]*entities.VolunteerAssignment, error)
	GetCompletedAssignmentsByVolunteer(ctx context.Context, volunteerID primitive.ObjectID, startDate, endDate time.Time) ([]*entities.VolunteerAssignment, error)
	GetAssignmentsByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	GetAssignmentsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.VolunteerAssignment, error)
	EnsureIndexes(ctx context.Context) error
}
