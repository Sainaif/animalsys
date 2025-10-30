package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShiftStatus represents the status of a shift
type ShiftStatus string

const (
	ShiftStatusAssigned        ShiftStatus = "assigned"
	ShiftStatusSwapRequested   ShiftStatus = "swap_requested"
	ShiftStatusAbsentRequested ShiftStatus = "absent_requested"
	ShiftStatusSwapped         ShiftStatus = "swapped"
	ShiftStatusAbsentApproved  ShiftStatus = "absent_approved"
	ShiftStatusCompleted       ShiftStatus = "completed"
	ShiftStatusCancelled       ShiftStatus = "cancelled"
)

// Schedule represents a shift schedule
type Schedule struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	UserName        string             `bson:"user_name" json:"user_name"`
	ShiftDate       time.Time          `bson:"shift_date" json:"shift_date"`
	StartTime       time.Time          `bson:"start_time" json:"start_time"`
	EndTime         time.Time          `bson:"end_time" json:"end_time"`
	Task            string             `bson:"task" json:"task"`
	Status          ShiftStatus        `bson:"status" json:"status"`
	SwapRequestedBy string             `bson:"swap_requested_by,omitempty" json:"swap_requested_by,omitempty"`
	SwapRequestedWith string           `bson:"swap_requested_with,omitempty" json:"swap_requested_with,omitempty"`
	AbsenceReason   string             `bson:"absence_reason,omitempty" json:"absence_reason,omitempty"`
	Notes           string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy       string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	ApprovedBy      string             `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
}

// ScheduleCreateRequest represents schedule creation
type ScheduleCreateRequest struct {
	UserID    string    `json:"user_id" validate:"required"`
	ShiftDate time.Time `json:"shift_date" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Task      string    `json:"task" validate:"required"`
	Notes     string    `json:"notes,omitempty"`
}

// ScheduleUpdateRequest represents schedule update
type ScheduleUpdateRequest struct {
	ShiftDate time.Time   `json:"shift_date,omitempty"`
	StartTime time.Time   `json:"start_time,omitempty"`
	EndTime   time.Time   `json:"end_time,omitempty"`
	Task      string      `json:"task,omitempty"`
	Status    ShiftStatus `json:"status,omitempty"`
	Notes     string      `json:"notes,omitempty"`
}

// SwapRequest represents a shift swap request
type SwapRequest struct {
	RequestedWith string `json:"requested_with" validate:"required"`
	Reason        string `json:"reason,omitempty"`
}

// AbsenceRequest represents an absence request
type AbsenceRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// NewSchedule creates a new schedule
func NewSchedule(userID, userName string, shiftDate, startTime, endTime time.Time, task string) *Schedule {
	now := time.Now()
	return &Schedule{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		UserName:  userName,
		ShiftDate: shiftDate,
		StartTime: startTime,
		EndTime:   endTime,
		Task:      task,
		Status:    ShiftStatusAssigned,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// RequestSwap requests a shift swap
func (s *Schedule) RequestSwap(requestedWith, requestedBy string) {
	s.Status = ShiftStatusSwapRequested
	s.SwapRequestedWith = requestedWith
	s.SwapRequestedBy = requestedBy
	s.UpdatedAt = time.Now()
}

// ApproveSwap approves a shift swap
func (s *Schedule) ApproveSwap(approvedBy string) {
	s.Status = ShiftStatusSwapped
	s.ApprovedBy = approvedBy
	s.UpdatedAt = time.Now()
}

// RequestAbsence requests an absence
func (s *Schedule) RequestAbsence(reason string) {
	s.Status = ShiftStatusAbsentRequested
	s.AbsenceReason = reason
	s.UpdatedAt = time.Now()
}

// ApproveAbsence approves an absence
func (s *Schedule) ApproveAbsence(approvedBy string) {
	s.Status = ShiftStatusAbsentApproved
	s.ApprovedBy = approvedBy
	s.UpdatedAt = time.Now()
}
