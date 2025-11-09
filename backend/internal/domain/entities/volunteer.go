package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VolunteerStatus represents the status of a volunteer
type VolunteerStatus string

const (
	VolunteerStatusActive    VolunteerStatus = "active"
	VolunteerStatusInactive  VolunteerStatus = "inactive"
	VolunteerStatusOnHold    VolunteerStatus = "on_hold"    // Temporarily unavailable
	VolunteerStatusSuspended VolunteerStatus = "suspended"  // Disciplinary action
	VolunteerStatusAlumni    VolunteerStatus = "alumni"     // Former volunteer
)

// SkillLevel represents proficiency level in a skill
type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "beginner"
	SkillLevelIntermediate SkillLevel = "intermediate"
	SkillLevelAdvanced     SkillLevel = "advanced"
	SkillLevelExpert       SkillLevel = "expert"
)

// VolunteerSkill represents a skill with proficiency level
type VolunteerSkill struct {
	Name  string     `json:"name" bson:"name"`
	Level SkillLevel `json:"level" bson:"level"`
}

// Availability represents volunteer availability
type Availability struct {
	DayOfWeek   string `json:"day_of_week" bson:"day_of_week"`       // Monday, Tuesday, etc.
	StartTime   string `json:"start_time" bson:"start_time"`         // "09:00"
	EndTime     string `json:"end_time" bson:"end_time"`             // "17:00"
	IsAvailable bool   `json:"is_available" bson:"is_available"`
}

// EmergencyContact represents emergency contact information
type EmergencyContact struct {
	Name         string `json:"name" bson:"name"`
	Relationship string `json:"relationship" bson:"relationship"`
	Phone        string `json:"phone" bson:"phone"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
}

// Certification represents a certification or training completion
type Certification struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	IssuingOrg     string             `json:"issuing_org,omitempty" bson:"issuing_org,omitempty"`
	IssueDate      time.Time          `json:"issue_date" bson:"issue_date"`
	ExpirationDate *time.Time         `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`
	CertificateURL string             `json:"certificate_url,omitempty" bson:"certificate_url,omitempty"`
	IsExpired      bool               `json:"is_expired" bson:"is_expired"`
}

// BackgroundCheck represents background check information
type BackgroundCheck struct {
	CompletedDate  time.Time  `json:"completed_date" bson:"completed_date"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"`
	Status         string     `json:"status" bson:"status"` // passed, failed, pending
	CheckedBy      string     `json:"checked_by,omitempty" bson:"checked_by,omitempty"`
	Notes          string     `json:"notes,omitempty" bson:"notes,omitempty"`
}

// Volunteer represents a volunteer in the system
type Volunteer struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	FirstName  string          `json:"first_name" bson:"first_name"`
	LastName   string          `json:"last_name" bson:"last_name"`
	Email      string          `json:"email" bson:"email"`
	Phone      string          `json:"phone,omitempty" bson:"phone,omitempty"`
	DateOfBirth *time.Time     `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	Status     VolunteerStatus `json:"status" bson:"status"`

	// Address
	Address string `json:"address,omitempty" bson:"address,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	State   string `json:"state,omitempty" bson:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty" bson:"zip_code,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`

	// Emergency Contact
	EmergencyContact EmergencyContact `json:"emergency_contact" bson:"emergency_contact"`

	// Volunteer Details
	Skills           []VolunteerSkill `json:"skills,omitempty" bson:"skills,omitempty"`
	Interests        []string         `json:"interests,omitempty" bson:"interests,omitempty"`
	Availability     []Availability   `json:"availability,omitempty" bson:"availability,omitempty"`
	PreferredTasks   []string         `json:"preferred_tasks,omitempty" bson:"preferred_tasks,omitempty"`

	// Onboarding
	ApplicationDate  time.Time          `json:"application_date" bson:"application_date"`
	ApprovalDate     *time.Time         `json:"approval_date,omitempty" bson:"approval_date,omitempty"`
	OrientationDate  *time.Time         `json:"orientation_date,omitempty" bson:"orientation_date,omitempty"`
	BackgroundCheck  *BackgroundCheck   `json:"background_check,omitempty" bson:"background_check,omitempty"`
	Certifications   []Certification    `json:"certifications,omitempty" bson:"certifications,omitempty"`

	// Activity Tracking
	TotalHours       float64   `json:"total_hours" bson:"total_hours"`
	EventsAttended   int       `json:"events_attended" bson:"events_attended"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty" bson:"last_activity_date,omitempty"`

	// Performance
	Rating           float64 `json:"rating,omitempty" bson:"rating,omitempty"`           // Average rating from supervisors
	ReviewCount      int     `json:"review_count" bson:"review_count"`
	CommendationCount int    `json:"commendation_count" bson:"commendation_count"`
	WarningCount     int     `json:"warning_count" bson:"warning_count"`

	// User Association (if volunteer has a user account)
	UserID *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`

	// Additional Information
	Notes             string `json:"notes,omitempty" bson:"notes,omitempty"`
	TShirtSize        string `json:"tshirt_size,omitempty" bson:"tshirt_size,omitempty"`
	DietaryRestrictions string `json:"dietary_restrictions,omitempty" bson:"dietary_restrictions,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// GetFullName returns the volunteer's full name
func (v *Volunteer) GetFullName() string {
	return v.FirstName + " " + v.LastName
}

// IsActive checks if the volunteer is currently active
func (v *Volunteer) IsActive() bool {
	return v.Status == VolunteerStatusActive
}

// HasSkill checks if volunteer has a specific skill
func (v *Volunteer) HasSkill(skillName string) bool {
	for _, skill := range v.Skills {
		if skill.Name == skillName {
			return true
		}
	}
	return false
}

// GetSkillLevel returns the level of a specific skill
func (v *Volunteer) GetSkillLevel(skillName string) (SkillLevel, bool) {
	for _, skill := range v.Skills {
		if skill.Name == skillName {
			return skill.Level, true
		}
	}
	return "", false
}

// IsAvailable checks if volunteer is available on a specific day
func (v *Volunteer) IsAvailable(dayOfWeek string) bool {
	for _, avail := range v.Availability {
		if avail.DayOfWeek == dayOfWeek && avail.IsAvailable {
			return true
		}
	}
	return false
}

// NeedsBackgroundCheck checks if background check is expired or missing
func (v *Volunteer) NeedsBackgroundCheck() bool {
	if v.BackgroundCheck == nil {
		return true
	}
	if v.BackgroundCheck.ExpirationDate != nil && v.BackgroundCheck.ExpirationDate.Before(time.Now()) {
		return true
	}
	return v.BackgroundCheck.Status != "passed"
}

// HasExpiredCertifications checks if any certifications are expired
func (v *Volunteer) HasExpiredCertifications() bool {
	now := time.Now()
	for _, cert := range v.Certifications {
		if cert.ExpirationDate != nil && cert.ExpirationDate.Before(now) {
			return true
		}
	}
	return false
}

// UpdateExpiredCertifications updates the IsExpired flag for all certifications
func (v *Volunteer) UpdateExpiredCertifications() {
	now := time.Now()
	for i := range v.Certifications {
		if v.Certifications[i].ExpirationDate != nil {
			v.Certifications[i].IsExpired = v.Certifications[i].ExpirationDate.Before(now)
		}
	}
}

// AddHours adds volunteer hours and updates activity date
func (v *Volunteer) AddHours(hours float64) {
	v.TotalHours += hours
	now := time.Now()
	v.LastActivityDate = &now
}

// IncrementEventsAttended increments the events attended counter
func (v *Volunteer) IncrementEventsAttended() {
	v.EventsAttended++
	now := time.Now()
	v.LastActivityDate = &now
}

// UpdateRating updates the volunteer's average rating
func (v *Volunteer) UpdateRating(newRating float64) {
	totalRating := v.Rating * float64(v.ReviewCount)
	v.ReviewCount++
	v.Rating = (totalRating + newRating) / float64(v.ReviewCount)
}

// AddCommendation adds a commendation to the volunteer's record
func (v *Volunteer) AddCommendation() {
	v.CommendationCount++
}

// AddWarning adds a warning to the volunteer's record
func (v *Volunteer) AddWarning() {
	v.WarningCount++
	// Automatically suspend if warnings exceed threshold
	if v.WarningCount >= 3 {
		v.Status = VolunteerStatusSuspended
	}
}

// Approve approves the volunteer application
func (v *Volunteer) Approve() {
	now := time.Now()
	v.ApprovalDate = &now
	v.Status = VolunteerStatusActive
}

// NewVolunteer creates a new volunteer
func NewVolunteer(
	firstName string,
	lastName string,
	email string,
	createdBy primitive.ObjectID,
) *Volunteer {
	now := time.Now()
	return &Volunteer{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		Status:            VolunteerStatusInactive, // Starts inactive until approved
		ApplicationDate:   now,
		TotalHours:        0,
		EventsAttended:    0,
		Rating:            0,
		ReviewCount:       0,
		CommendationCount: 0,
		WarningCount:      0,
		Skills:            []VolunteerSkill{},
		Interests:         []string{},
		Availability:      []Availability{},
		PreferredTasks:    []string{},
		Certifications:    []Certification{},
		CreatedBy:         createdBy,
		UpdatedBy:         createdBy,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
