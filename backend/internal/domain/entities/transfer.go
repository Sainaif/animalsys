package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransferDirection represents the direction of the transfer
type TransferDirection string

const (
	TransferDirectionIncoming TransferDirection = "incoming"
	TransferDirectionOutgoing TransferDirection = "outgoing"
)

// TransferStatus represents the status of an animal transfer
type TransferStatus string

const (
	TransferStatusPending    TransferStatus = "pending"
	TransferStatusApproved   TransferStatus = "approved"
	TransferStatusInTransit  TransferStatus = "in_transit"
	TransferStatusCompleted  TransferStatus = "completed"
	TransferStatusCancelled  TransferStatus = "cancelled"
	TransferStatusRejected   TransferStatus = "rejected"
)

// TransferReason represents the reason for transfer
type TransferReason string

const (
	TransferReasonOvercapacity     TransferReason = "overcapacity"
	TransferReasonSpecialization   TransferReason = "specialization"
	TransferReasonMedical          TransferReason = "medical"
	TransferReasonAdoption         TransferReason = "adoption"
	TransferReasonFoster           TransferReason = "foster"
	TransferReasonBehavioral       TransferReason = "behavioral"
	TransferReasonReunification    TransferReason = "reunification"
	TransferReasonSanctuary        TransferReason = "sanctuary"
	TransferReasonOther            TransferReason = "other"
)

// Transfer represents an animal transfer between organizations
type Transfer struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Direction TransferDirection  `json:"direction" bson:"direction"`
	Status    TransferStatus     `json:"status" bson:"status"`
	Reason    TransferReason     `json:"reason" bson:"reason"`

	// Animal being transferred
	AnimalID primitive.ObjectID `json:"animal_id" bson:"animal_id"`

	// Partner organization
	PartnerID primitive.ObjectID `json:"partner_id" bson:"partner_id"`

	// From/To details
	FromOrganization string `json:"from_organization" bson:"from_organization"`
	ToOrganization   string `json:"to_organization" bson:"to_organization"`

	// Contact persons
	FromContact ContactPerson `json:"from_contact,omitempty" bson:"from_contact,omitempty"`
	ToContact   ContactPerson `json:"to_contact,omitempty" bson:"to_contact,omitempty"`

	// Scheduling
	RequestedDate  time.Time  `json:"requested_date" bson:"requested_date"`
	ApprovedDate   *time.Time `json:"approved_date,omitempty" bson:"approved_date,omitempty"`
	ScheduledDate  *time.Time `json:"scheduled_date,omitempty" bson:"scheduled_date,omitempty"`
	CompletedDate  *time.Time `json:"completed_date,omitempty" bson:"completed_date,omitempty"`
	CancelledDate  *time.Time `json:"cancelled_date,omitempty" bson:"cancelled_date,omitempty"`

	// Transport
	TransportMethod  string  `json:"transport_method,omitempty" bson:"transport_method,omitempty"` // "ground", "air", "volunteer"
	TransportBy      string  `json:"transport_by,omitempty" bson:"transport_by,omitempty"`
	TransportCost    float64 `json:"transport_cost,omitempty" bson:"transport_cost,omitempty"`
	EstimatedDuration int    `json:"estimated_duration,omitempty" bson:"estimated_duration,omitempty"` // in minutes

	// Medical/Health
	HealthCertificate     bool       `json:"health_certificate" bson:"health_certificate"`
	HealthCertificateURL  string     `json:"health_certificate_url,omitempty" bson:"health_certificate_url,omitempty"`
	VaccinationRecords    bool       `json:"vaccination_records" bson:"vaccination_records"`
	MedicalRecordsURL     string     `json:"medical_records_url,omitempty" bson:"medical_records_url,omitempty"`
	RequiresQuarantine    bool       `json:"requires_quarantine" bson:"requires_quarantine"`
	QuarantineDays        int        `json:"quarantine_days,omitempty" bson:"quarantine_days,omitempty"`

	// Conditions and Terms
	ReturnPolicy      string `json:"return_policy,omitempty" bson:"return_policy,omitempty"`
	FollowUpRequired  bool   `json:"follow_up_required" bson:"follow_up_required"`
	FollowUpDate      *time.Time `json:"follow_up_date,omitempty" bson:"follow_up_date,omitempty"`
	SpecialConditions string `json:"special_conditions,omitempty" bson:"special_conditions,omitempty"`

	// Financial
	TransferFee      float64 `json:"transfer_fee,omitempty" bson:"transfer_fee,omitempty"`
	FeeWaived        bool    `json:"fee_waived" bson:"fee_waived"`
	PaymentStatus    string  `json:"payment_status,omitempty" bson:"payment_status,omitempty"` // "pending", "paid", "waived"

	// Documents
	Documents []string `json:"documents,omitempty" bson:"documents,omitempty"` // Document IDs

	// Notes
	Notes            string `json:"notes,omitempty" bson:"notes,omitempty"`
	ReasonDetails    string `json:"reason_details,omitempty" bson:"reason_details,omitempty"`
	CancellationReason string `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`

	// Tracking
	RequestedBy  primitive.ObjectID  `json:"requested_by" bson:"requested_by"`
	ApprovedBy   *primitive.ObjectID `json:"approved_by,omitempty" bson:"approved_by,omitempty"`
	CompletedBy  *primitive.ObjectID `json:"completed_by,omitempty" bson:"completed_by,omitempty"`
	CancelledBy  *primitive.ObjectID `json:"cancelled_by,omitempty" bson:"cancelled_by,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewTransfer creates a new transfer
func NewTransfer(animalID, partnerID primitive.ObjectID, direction TransferDirection, reason TransferReason, requestedBy primitive.ObjectID) *Transfer {
	now := time.Now()
	return &Transfer{
		ID:                  primitive.NewObjectID(),
		AnimalID:            animalID,
		PartnerID:           partnerID,
		Direction:           direction,
		Status:              TransferStatusPending,
		Reason:              reason,
		RequestedDate:       now,
		HealthCertificate:   false,
		VaccinationRecords:  false,
		RequiresQuarantine:  false,
		FollowUpRequired:    false,
		FeeWaived:           false,
		Documents:           []string{},
		RequestedBy:         requestedBy,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// Approve approves the transfer
func (t *Transfer) Approve(approvedBy primitive.ObjectID) {
	t.Status = TransferStatusApproved
	now := time.Now()
	t.ApprovedDate = &now
	t.ApprovedBy = &approvedBy
	t.UpdatedAt = now
}

// StartTransit marks the transfer as in transit
func (t *Transfer) StartTransit() {
	t.Status = TransferStatusInTransit
	t.UpdatedAt = time.Now()
}

// Complete completes the transfer
func (t *Transfer) Complete(completedBy primitive.ObjectID) {
	t.Status = TransferStatusCompleted
	now := time.Now()
	t.CompletedDate = &now
	t.CompletedBy = &completedBy
	t.UpdatedAt = now
}

// Cancel cancels the transfer
func (t *Transfer) Cancel(cancelledBy primitive.ObjectID, reason string) {
	t.Status = TransferStatusCancelled
	now := time.Now()
	t.CancelledDate = &now
	t.CancelledBy = &cancelledBy
	t.CancellationReason = reason
	t.UpdatedAt = now
}

// Reject rejects the transfer
func (t *Transfer) Reject(reason string) {
	t.Status = TransferStatusRejected
	t.CancellationReason = reason
	t.UpdatedAt = time.Now()
}

// Schedule schedules the transfer
func (t *Transfer) Schedule(scheduledDate time.Time) {
	t.ScheduledDate = &scheduledDate
	t.UpdatedAt = time.Now()
}

// IsOverdue checks if the transfer is overdue
func (t *Transfer) IsOverdue() bool {
	if t.ScheduledDate == nil || t.Status == TransferStatusCompleted || t.Status == TransferStatusCancelled {
		return false
	}
	return time.Now().After(*t.ScheduledDate)
}

// RequiresFollowUp checks if follow-up is needed
func (t *Transfer) RequiresFollowUp() bool {
	if !t.FollowUpRequired || t.FollowUpDate == nil {
		return false
	}
	return time.Now().After(*t.FollowUpDate) && t.Status == TransferStatusCompleted
}

// GetDuration calculates the duration from request to completion
func (t *Transfer) GetDuration() *time.Duration {
	if t.CompletedDate == nil {
		return nil
	}
	duration := t.CompletedDate.Sub(t.RequestedDate)
	return &duration
}
