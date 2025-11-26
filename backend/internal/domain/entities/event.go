package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeFundraising  EventType = "fundraising"
	EventTypeAdoption     EventType = "adoption"      // Adoption event/fair
	EventTypeEducational  EventType = "educational"   // Educational program
	EventTypeVolunteer    EventType = "volunteer"     // Volunteer appreciation/training
	EventTypeCommunity    EventType = "community"     // Community outreach
	EventTypeSocial       EventType = "social"        // Social gathering
	EventTypeShopping     EventType = "shopping"
	EventTypeContribution EventType = "contribution"
	EventTypeOther        EventType = "other"
)

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusScheduled EventStatus = "scheduled"
	EventStatusActive    EventStatus = "active"      // Currently happening
	EventStatusCompleted EventStatus = "completed"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPostponed EventStatus = "postponed"
)

// EventLocation represents the location of an event
type EventLocation struct {
	Name      string  `json:"name" bson:"name"`                               // e.g., "Foundation Headquarters"
	Address   string  `json:"address" bson:"address"`
	City      string  `json:"city" bson:"city"`
	State     string  `json:"state" bson:"state"`
	ZipCode   string  `json:"zip_code" bson:"zip_code"`
	Country   string  `json:"country" bson:"country"`
	Latitude  float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

// EventRegistration represents registration requirements
type EventRegistration struct {
	Required       bool      `json:"required" bson:"required"`
	MaxAttendees   int       `json:"max_attendees,omitempty" bson:"max_attendees,omitempty"`
	CurrentCount   int       `json:"current_count" bson:"current_count"`
	RegistrationFee float64  `json:"registration_fee,omitempty" bson:"registration_fee,omitempty"`
	Deadline       *time.Time `json:"deadline,omitempty" bson:"deadline,omitempty"`
}

// Event represents a foundation event
type Event struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Name        MultilingualName `json:"name" bson:"name"`
	Description MultilingualName `json:"description" bson:"description"`
	Type        EventType        `json:"type" bson:"type"`
	Status      EventStatus      `json:"status" bson:"status"`

	// Schedule
	StartDate time.Time  `json:"start_date" bson:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Duration  int        `json:"duration" bson:"duration"` // Duration in minutes

	// Location
	Location    EventLocation `json:"location" bson:"location"`
	VirtualLink string        `json:"virtual_link,omitempty" bson:"virtual_link,omitempty"` // For online events

	// Registration
	Registration EventRegistration `json:"registration" bson:"registration"`

	// Details
	ImageURL     string   `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Tags         []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Public       bool     `json:"public" bson:"public"` // Publicly visible
	Featured     bool     `json:"featured" bson:"featured"`

	// Organization
	Organizer         primitive.ObjectID   `json:"organizer" bson:"organizer"` // Staff organizing the event
	ContactEmail      string               `json:"contact_email,omitempty" bson:"contact_email,omitempty"`
	ContactPhone      string               `json:"contact_phone,omitempty" bson:"contact_phone,omitempty"`
	RequiredVolunteers int                 `json:"required_volunteers" bson:"required_volunteers"`
	AssignedVolunteers []primitive.ObjectID `json:"assigned_volunteers,omitempty" bson:"assigned_volunteers,omitempty"`

	// Campaign Association
	CampaignID *primitive.ObjectID `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"` // Link to fundraising campaign

	// Financial
	Budget      float64 `json:"budget,omitempty" bson:"budget,omitempty"`
	FundsRaised float64 `json:"funds_raised,omitempty" bson:"funds_raised,omitempty"`

	// Statistics
	AttendeeCount  int `json:"attendee_count" bson:"attendee_count"`
	VolunteerCount int `json:"volunteer_count" bson:"volunteer_count"`
	AnimalsAdopted int `json:"animals_adopted,omitempty" bson:"animals_adopted,omitempty"`

	// Additional Information
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsUpcoming checks if the event is upcoming
func (e *Event) IsUpcoming() bool {
	return e.StartDate.After(time.Now()) && e.Status == EventStatusScheduled
}

// IsActive checks if the event is currently active
func (e *Event) IsActive() bool {
	now := time.Now()
	if e.Status != EventStatusActive {
		return false
	}
	if now.Before(e.StartDate) {
		return false
	}
	if e.EndDate != nil && now.After(*e.EndDate) {
		return false
	}
	return true
}

// IsFull checks if the event has reached max capacity
func (e *Event) IsFull() bool {
	if !e.Registration.Required || e.Registration.MaxAttendees == 0 {
		return false
	}
	return e.Registration.CurrentCount >= e.Registration.MaxAttendees
}

// NeedsVolunteers checks if more volunteers are needed
func (e *Event) NeedsVolunteers() bool {
	return len(e.AssignedVolunteers) < e.RequiredVolunteers
}

// UpdateStatistics updates event statistics
func (e *Event) UpdateStatistics(attendees, volunteers int, fundsRaised float64, animalsAdopted int) {
	e.AttendeeCount = attendees
	e.VolunteerCount = volunteers
	if fundsRaised > 0 {
		if e.Type == EventTypeShopping || e.Type == EventTypeContribution {
			e.Budget += fundsRaised
		} else {
			e.FundsRaised += fundsRaised
		}
	}
	if animalsAdopted > 0 {
		e.AnimalsAdopted += animalsAdopted
	}
}

// NewEvent creates a new event
func NewEvent(
	name MultilingualName,
	eventType EventType,
	startDate time.Time,
	organizer primitive.ObjectID,
	createdBy primitive.ObjectID,
) *Event {
	now := time.Now()
	return &Event{
		Name:               name,
		Type:               eventType,
		Status:             EventStatusDraft,
		StartDate:          startDate,
		Duration:           60, // Default 1 hour
		Public:             false,
		Featured:           false,
		Organizer:          organizer,
		RequiredVolunteers: 0,
		Registration: EventRegistration{
			Required:     false,
			CurrentCount: 0,
		},
		AttendeeCount:  0,
		VolunteerCount: 0,
		Budget:         0,
		FundsRaised:    0,
		AnimalsAdopted: 0,
		CreatedBy:      createdBy,
		UpdatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
