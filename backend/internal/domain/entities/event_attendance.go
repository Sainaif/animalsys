package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AttendeeType represents the type of attendee
type AttendeeType string

const (
	AttendeeTypeVolunteer AttendeeType = "volunteer"
	AttendeeTypeStaff     AttendeeType = "staff"
	AttendeeTypePublic    AttendeeType = "public"
	AttendeeTypeDonor     AttendeeType = "donor"
	AttendeeTypeSponsor   AttendeeType = "sponsor"
)

// AttendanceStatus represents the status of attendance
type AttendanceStatus string

const (
	AttendanceStatusRegistered AttendanceStatus = "registered"
	AttendanceStatusConfirmed  AttendanceStatus = "confirmed"
	AttendanceStatusAttended   AttendanceStatus = "attended"
	AttendanceStatusNoShow     AttendanceStatus = "no_show"
	AttendanceStatusCancelled  AttendanceStatus = "cancelled"
)

// EventAttendance represents an attendee's registration and attendance for an event
type EventAttendance struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Event and Attendee Information
	EventID      primitive.ObjectID `json:"event_id" bson:"event_id"`
	AttendeeType AttendeeType       `json:"attendee_type" bson:"attendee_type"`

	// Reference to the attendee (could be volunteer, user, donor, etc.)
	VolunteerID *primitive.ObjectID `json:"volunteer_id,omitempty" bson:"volunteer_id,omitempty"`
	UserID      *primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	DonorID     *primitive.ObjectID `json:"donor_id,omitempty" bson:"donor_id,omitempty"`

	// Guest information (for public attendees without accounts)
	GuestName  string `json:"guest_name,omitempty" bson:"guest_name,omitempty"`
	GuestEmail string `json:"guest_email,omitempty" bson:"guest_email,omitempty"`
	GuestPhone string `json:"guest_phone,omitempty" bson:"guest_phone,omitempty"`

	// Registration Details
	Status           AttendanceStatus `json:"status" bson:"status"`
	RegistrationDate time.Time        `json:"registration_date" bson:"registration_date"`
	ConfirmationDate *time.Time       `json:"confirmation_date,omitempty" bson:"confirmation_date,omitempty"`
	CancellationDate *time.Time       `json:"cancellation_date,omitempty" bson:"cancellation_date,omitempty"`

	// Payment (if registration has a fee)
	RegistrationFee     float64 `json:"registration_fee" bson:"registration_fee"`
	PaymentStatus       string  `json:"payment_status,omitempty" bson:"payment_status,omitempty"` // pending, paid, refunded
	PaymentDate         *time.Time `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	PaymentMethod       string  `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	TransactionID       string  `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`

	// Attendance Tracking
	CheckInTime  *time.Time `json:"check_in_time,omitempty" bson:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty" bson:"check_out_time,omitempty"`

	// Additional Information
	NumberOfGuests      int      `json:"number_of_guests" bson:"number_of_guests"`           // Additional guests brought
	SpecialRequirements string   `json:"special_requirements,omitempty" bson:"special_requirements,omitempty"`
	DietaryRestrictions string   `json:"dietary_restrictions,omitempty" bson:"dietary_restrictions,omitempty"`
	Notes               string   `json:"notes,omitempty" bson:"notes,omitempty"`

	// Feedback (collected after event)
	Rating   int    `json:"rating,omitempty" bson:"rating,omitempty"`           // 1-5 stars
	Feedback string `json:"feedback,omitempty" bson:"feedback,omitempty"`
	FeedbackDate *time.Time `json:"feedback_date,omitempty" bson:"feedback_date,omitempty"`

	// Metadata
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// CheckIn marks the attendee as checked in
func (ea *EventAttendance) CheckIn() {
	now := time.Now()
	ea.CheckInTime = &now
	ea.Status = AttendanceStatusAttended
}

// CheckOut marks the attendee as checked out
func (ea *EventAttendance) CheckOut() {
	now := time.Now()
	ea.CheckOutTime = &now
}

// Confirm confirms the attendance
func (ea *EventAttendance) Confirm() {
	now := time.Now()
	ea.ConfirmationDate = &now
	ea.Status = AttendanceStatusConfirmed
}

// Cancel cancels the attendance
func (ea *EventAttendance) Cancel() {
	now := time.Now()
	ea.CancellationDate = &now
	ea.Status = AttendanceStatusCancelled
}

// MarkNoShow marks the attendee as no-show
func (ea *EventAttendance) MarkNoShow() {
	ea.Status = AttendanceStatusNoShow
}

// IsAttended checks if the attendee actually attended
func (ea *EventAttendance) IsAttended() bool {
	return ea.Status == AttendanceStatusAttended && ea.CheckInTime != nil
}

// GetDuration calculates the duration of attendance in minutes
func (ea *EventAttendance) GetDuration() float64 {
	if ea.CheckInTime == nil {
		return 0
	}

	endTime := time.Now()
	if ea.CheckOutTime != nil {
		endTime = *ea.CheckOutTime
	}

	return endTime.Sub(*ea.CheckInTime).Minutes()
}

// SubmitFeedback submits feedback for the event
func (ea *EventAttendance) SubmitFeedback(rating int, feedback string) {
	ea.Rating = rating
	ea.Feedback = feedback
	now := time.Now()
	ea.FeedbackDate = &now
}

// IsPaid checks if the registration fee has been paid
func (ea *EventAttendance) IsPaid() bool {
	if ea.RegistrationFee == 0 {
		return true // No fee required
	}
	return ea.PaymentStatus == "paid"
}

// MarkPaid marks the registration fee as paid
func (ea *EventAttendance) MarkPaid(paymentMethod string, transactionID string) {
	now := time.Now()
	ea.PaymentStatus = "paid"
	ea.PaymentDate = &now
	ea.PaymentMethod = paymentMethod
	ea.TransactionID = transactionID
}

// NewEventAttendance creates a new event attendance record
func NewEventAttendance(
	eventID primitive.ObjectID,
	attendeeType AttendeeType,
	registrationFee float64,
	createdBy primitive.ObjectID,
) *EventAttendance {
	now := time.Now()

	paymentStatus := ""
	if registrationFee > 0 {
		paymentStatus = "pending"
	} else {
		paymentStatus = "paid" // Free event
	}

	return &EventAttendance{
		EventID:          eventID,
		AttendeeType:     attendeeType,
		Status:           AttendanceStatusRegistered,
		RegistrationDate: now,
		RegistrationFee:  registrationFee,
		PaymentStatus:    paymentStatus,
		NumberOfGuests:   0,
		Rating:           0,
		CreatedBy:        createdBy,
		UpdatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
