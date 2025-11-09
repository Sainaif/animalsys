package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssignmentType represents the type of volunteer assignment
type AssignmentType string

const (
	AssignmentTypeEvent      AssignmentType = "event"       // Event volunteering
	AssignmentTypeAnimalCare AssignmentType = "animal_care" // Daily animal care
	AssignmentTypeFostering  AssignmentType = "fostering"   // Animal fostering
	AssignmentTypeAdminTask  AssignmentType = "admin_task"  // Administrative tasks
	AssignmentTypeTransport  AssignmentType = "transport"   // Animal transport
	AssignmentTypeFundraising AssignmentType = "fundraising" // Fundraising activities
	AssignmentTypeEducation  AssignmentType = "education"   // Educational programs
	AssignmentTypeOther      AssignmentType = "other"
)

// AssignmentStatus represents the status of a volunteer assignment
type AssignmentStatus string

const (
	AssignmentStatusAssigned   AssignmentStatus = "assigned"   // Assigned but not confirmed
	AssignmentStatusConfirmed  AssignmentStatus = "confirmed"  // Volunteer confirmed
	AssignmentStatusInProgress AssignmentStatus = "in_progress" // Currently working
	AssignmentStatusCompleted  AssignmentStatus = "completed"  // Task completed
	AssignmentStatusCancelled  AssignmentStatus = "cancelled"  // Cancelled by volunteer or staff
	AssignmentStatusNoShow     AssignmentStatus = "no_show"    // Volunteer didn't show up
)

// VolunteerAssignment represents a task or event assigned to a volunteer
type VolunteerAssignment struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Volunteer and Assignment Details
	VolunteerID    primitive.ObjectID `json:"volunteer_id" bson:"volunteer_id"`
	Type           AssignmentType     `json:"type" bson:"type"`
	Status         AssignmentStatus   `json:"status" bson:"status"`

	// Related Resources
	EventID   *primitive.ObjectID `json:"event_id,omitempty" bson:"event_id,omitempty"`     // If event assignment
	AnimalID  *primitive.ObjectID `json:"animal_id,omitempty" bson:"animal_id,omitempty"`   // If animal care
	CampaignID *primitive.ObjectID `json:"campaign_id,omitempty" bson:"campaign_id,omitempty"` // If campaign work

	// Assignment Details
	Title       string `json:"title" bson:"title"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`

	// Schedule
	StartDate time.Time  `json:"start_date" bson:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Duration  int        `json:"duration,omitempty" bson:"duration,omitempty"` // Expected duration in minutes

	// Location (if applicable)
	Location string `json:"location,omitempty" bson:"location,omitempty"`

	// Assignment Management
	AssignedBy    primitive.ObjectID `json:"assigned_by" bson:"assigned_by"`
	AssignedDate  time.Time          `json:"assigned_date" bson:"assigned_date"`
	ConfirmedDate *time.Time         `json:"confirmed_date,omitempty" bson:"confirmed_date,omitempty"`

	// Execution Tracking
	StartTime    *time.Time `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty" bson:"end_time,omitempty"`
	ActualHours  float64    `json:"actual_hours" bson:"actual_hours"`

	// Requirements and Skills
	RequiredSkills []string `json:"required_skills,omitempty" bson:"required_skills,omitempty"`
	RequiredCerts  []string `json:"required_certs,omitempty" bson:"required_certs,omitempty"`

	// Performance
	CompletionNotes string  `json:"completion_notes,omitempty" bson:"completion_notes,omitempty"`
	Rating          int     `json:"rating,omitempty" bson:"rating,omitempty"` // 1-5 rating from supervisor
	ReviewedBy      *primitive.ObjectID `json:"reviewed_by,omitempty" bson:"reviewed_by,omitempty"`
	ReviewDate      *time.Time `json:"review_date,omitempty" bson:"review_date,omitempty"`
	ReviewNotes     string  `json:"review_notes,omitempty" bson:"review_notes,omitempty"`

	// Cancellation
	CancellationReason string     `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`
	CancelledBy        *primitive.ObjectID `json:"cancelled_by,omitempty" bson:"cancelled_by,omitempty"`
	CancelledDate      *time.Time `json:"cancelled_date,omitempty" bson:"cancelled_date,omitempty"`

	// Notifications
	ReminderSent bool       `json:"reminder_sent" bson:"reminder_sent"`
	ReminderDate *time.Time `json:"reminder_date,omitempty" bson:"reminder_date,omitempty"`

	// Additional Information
	Notes string `json:"notes,omitempty" bson:"notes,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Confirm confirms the volunteer assignment
func (va *VolunteerAssignment) Confirm() {
	now := time.Now()
	va.ConfirmedDate = &now
	va.Status = AssignmentStatusConfirmed
}

// Start starts the assignment
func (va *VolunteerAssignment) Start() {
	now := time.Now()
	va.StartTime = &now
	va.Status = AssignmentStatusInProgress
}

// Complete completes the assignment
func (va *VolunteerAssignment) Complete(completionNotes string) {
	now := time.Now()
	va.EndTime = &now
	va.Status = AssignmentStatusCompleted
	va.CompletionNotes = completionNotes

	// Calculate actual hours if start time is set
	if va.StartTime != nil {
		va.ActualHours = now.Sub(*va.StartTime).Hours()
	}
}

// Cancel cancels the assignment
func (va *VolunteerAssignment) Cancel(reason string, cancelledBy primitive.ObjectID) {
	now := time.Now()
	va.Status = AssignmentStatusCancelled
	va.CancellationReason = reason
	va.CancelledBy = &cancelledBy
	va.CancelledDate = &now
}

// MarkNoShow marks the volunteer as no-show
func (va *VolunteerAssignment) MarkNoShow() {
	va.Status = AssignmentStatusNoShow
}

// Review adds a review to the assignment
func (va *VolunteerAssignment) Review(rating int, reviewNotes string, reviewedBy primitive.ObjectID) {
	now := time.Now()
	va.Rating = rating
	va.ReviewNotes = reviewNotes
	va.ReviewedBy = &reviewedBy
	va.ReviewDate = &now
}

// SendReminder marks that a reminder has been sent
func (va *VolunteerAssignment) SendReminder() {
	now := time.Now()
	va.ReminderSent = true
	va.ReminderDate = &now
}

// IsActive checks if the assignment is currently active
func (va *VolunteerAssignment) IsActive() bool {
	return va.Status == AssignmentStatusConfirmed || va.Status == AssignmentStatusInProgress
}

// IsUpcoming checks if the assignment is upcoming
func (va *VolunteerAssignment) IsUpcoming() bool {
	return va.StartDate.After(time.Now()) && (va.Status == AssignmentStatusAssigned || va.Status == AssignmentStatusConfirmed)
}

// IsPast checks if the assignment date has passed
func (va *VolunteerAssignment) IsPast() bool {
	if va.EndDate != nil {
		return va.EndDate.Before(time.Now())
	}
	return va.StartDate.Before(time.Now())
}

// GetDuration returns the actual or expected duration in hours
func (va *VolunteerAssignment) GetDuration() float64 {
	if va.ActualHours > 0 {
		return va.ActualHours
	}
	if va.Duration > 0 {
		return float64(va.Duration) / 60.0
	}
	return 0
}

// NeedsReminder checks if a reminder should be sent (24 hours before start)
func (va *VolunteerAssignment) NeedsReminder() bool {
	if va.ReminderSent {
		return false
	}
	if va.Status != AssignmentStatusConfirmed && va.Status != AssignmentStatusAssigned {
		return false
	}

	reminderTime := va.StartDate.Add(-24 * time.Hour)
	return time.Now().After(reminderTime) && time.Now().Before(va.StartDate)
}

// NewVolunteerAssignment creates a new volunteer assignment
func NewVolunteerAssignment(
	volunteerID primitive.ObjectID,
	assignmentType AssignmentType,
	title string,
	startDate time.Time,
	assignedBy primitive.ObjectID,
) *VolunteerAssignment {
	now := time.Now()
	return &VolunteerAssignment{
		VolunteerID:    volunteerID,
		Type:           assignmentType,
		Status:         AssignmentStatusAssigned,
		Title:          title,
		StartDate:      startDate,
		AssignedBy:     assignedBy,
		AssignedDate:   now,
		ActualHours:    0,
		Rating:         0,
		ReminderSent:   false,
		RequiredSkills: []string{},
		RequiredCerts:  []string{},
		CreatedBy:      assignedBy,
		UpdatedBy:      assignedBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
