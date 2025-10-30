package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VolunteerStatus represents the status of a volunteer
type VolunteerStatus string

const (
	VolunteerStatusPending  VolunteerStatus = "pending"
	VolunteerStatusActive   VolunteerStatus = "active"
	VolunteerStatusInactive VolunteerStatus = "inactive"
	VolunteerStatusSuspended VolunteerStatus = "suspended"
)

// Volunteer represents a volunteer in the organization
type Volunteer struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID              string               `bson:"user_id" json:"user_id"`
	RegistrationDate    time.Time            `bson:"registration_date" json:"registration_date"`
	Status              VolunteerStatus      `bson:"status" json:"status"`
	Skills              []string             `bson:"skills,omitempty" json:"skills,omitempty"`
	Availability        []string             `bson:"availability,omitempty" json:"availability,omitempty"` // days of week
	EmergencyContact    string               `bson:"emergency_contact" json:"emergency_contact"`
	EmergencyPhone      string               `bson:"emergency_phone" json:"emergency_phone"`
	BackgroundCheckStatus string             `bson:"background_check_status,omitempty" json:"background_check_status,omitempty"`
	BackgroundCheckDate *time.Time           `bson:"background_check_date,omitempty" json:"background_check_date,omitempty"`
	Trainings           []Training           `bson:"trainings,omitempty" json:"trainings,omitempty"`
	TotalHours          float64              `bson:"total_hours" json:"total_hours"`
	HoursThisMonth      float64              `bson:"hours_this_month" json:"hours_this_month"`
	PerformanceRating   float64              `bson:"performance_rating,omitempty" json:"performance_rating,omitempty"`
	Notes               string               `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt           time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time            `bson:"updated_at" json:"updated_at"`
	CreatedBy           string               `bson:"created_by,omitempty" json:"created_by,omitempty"`
}

// Training represents a training or certification
type Training struct {
	Name             string     `bson:"name" json:"name"`
	CompletionDate   time.Time  `bson:"completion_date" json:"completion_date"`
	ExpiryDate       *time.Time `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CertificateURL   string     `bson:"certificate_url,omitempty" json:"certificate_url,omitempty"`
	TrainingProvider string     `bson:"training_provider,omitempty" json:"training_provider,omitempty"`
}

// VolunteerHourEntry represents hours logged by volunteer
type VolunteerHourEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VolunteerID string             `bson:"volunteer_id" json:"volunteer_id"`
	Date        time.Time          `bson:"date" json:"date"`
	CheckIn     time.Time          `bson:"check_in" json:"check_in"`
	CheckOut    *time.Time         `bson:"check_out,omitempty" json:"check_out,omitempty"`
	Hours       float64            `bson:"hours" json:"hours"`
	Activity    string             `bson:"activity" json:"activity"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty"`
	ApprovedBy  string             `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// VolunteerCreateRequest represents volunteer registration
type VolunteerCreateRequest struct {
	UserID           string   `json:"user_id" validate:"required"`
	Skills           []string `json:"skills,omitempty"`
	Availability     []string `json:"availability,omitempty"`
	EmergencyContact string   `json:"emergency_contact" validate:"required"`
	EmergencyPhone   string   `json:"emergency_phone" validate:"required"`
}

// VolunteerUpdateRequest represents volunteer update
type VolunteerUpdateRequest struct {
	Status              VolunteerStatus `json:"status,omitempty"`
	Skills              []string        `json:"skills,omitempty"`
	Availability        []string        `json:"availability,omitempty"`
	BackgroundCheckStatus string        `json:"background_check_status,omitempty"`
	PerformanceRating   float64         `json:"performance_rating,omitempty" validate:"omitempty,gte=0,lte=5"`
	Notes               string          `json:"notes,omitempty"`
}

// NewVolunteer creates a new volunteer
func NewVolunteer(userID, emergencyContact, emergencyPhone string) *Volunteer {
	now := time.Now()
	return &Volunteer{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		RegistrationDate: now,
		Status:           VolunteerStatusPending,
		EmergencyContact: emergencyContact,
		EmergencyPhone:   emergencyPhone,
		TotalHours:       0,
		HoursThisMonth:   0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// AddTraining adds a training record
func (v *Volunteer) AddTraining(training Training) {
	v.Trainings = append(v.Trainings, training)
	v.UpdatedAt = time.Now()
}

// LogHours adds hours to volunteer's record
func (v *Volunteer) LogHours(hours float64) {
	v.TotalHours += hours
	v.HoursThisMonth += hours
	v.UpdatedAt = time.Now()
}
