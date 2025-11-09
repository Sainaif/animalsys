package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeAnimal     ReportType = "animal"
	ReportTypeAdoption   ReportType = "adoption"
	ReportTypeDonation   ReportType = "donation"
	ReportTypeVeterinary ReportType = "veterinary"
	ReportTypeVolunteer  ReportType = "volunteer"
	ReportTypeEvent      ReportType = "event"
	ReportTypeCampaign   ReportType = "campaign"
	ReportTypeFinancial  ReportType = "financial"
	ReportTypeCustom     ReportType = "custom"
)

// ReportFormat represents the output format of a report
type ReportFormat string

const (
	ReportFormatJSON ReportFormat = "json"
	ReportFormatCSV  ReportFormat = "csv"
	ReportFormatPDF  ReportFormat = "pdf"
	ReportFormatXLSX ReportFormat = "xlsx"
)

// ReportStatus represents the status of report generation
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusRunning   ReportStatus = "running"
	ReportStatusCompleted ReportStatus = "completed"
	ReportStatusFailed    ReportStatus = "failed"
)

// ReportScheduleFrequency represents how often a report should run
type ReportScheduleFrequency string

const (
	ReportFrequencyNone    ReportScheduleFrequency = "none"
	ReportFrequencyDaily   ReportScheduleFrequency = "daily"
	ReportFrequencyWeekly  ReportScheduleFrequency = "weekly"
	ReportFrequencyMonthly ReportScheduleFrequency = "monthly"
	ReportFrequencyYearly  ReportScheduleFrequency = "yearly"
)

// ReportSchedule represents the scheduling configuration for a report
type ReportSchedule struct {
	Enabled      bool                    `json:"enabled" bson:"enabled"`
	Frequency    ReportScheduleFrequency `json:"frequency" bson:"frequency"`
	DayOfWeek    int                     `json:"day_of_week,omitempty" bson:"day_of_week,omitempty"`       // 0-6 (Sunday-Saturday)
	DayOfMonth   int                     `json:"day_of_month,omitempty" bson:"day_of_month,omitempty"`     // 1-31
	Time         string                  `json:"time,omitempty" bson:"time,omitempty"`                     // HH:MM format
	NextRunAt    *time.Time              `json:"next_run_at,omitempty" bson:"next_run_at,omitempty"`
	LastRunAt    *time.Time              `json:"last_run_at,omitempty" bson:"last_run_at,omitempty"`
	Recipients   []string                `json:"recipients,omitempty" bson:"recipients,omitempty"`         // Email addresses
}

// ReportExecution represents a single execution of a report
type ReportExecution struct {
	ID            primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	ReportID      primitive.ObjectID     `json:"report_id" bson:"report_id"`
	Status        ReportStatus           `json:"status" bson:"status"`
	StartedAt     time.Time              `json:"started_at" bson:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Duration      int64                  `json:"duration,omitempty" bson:"duration,omitempty"` // milliseconds
	RecordCount   int64                  `json:"record_count,omitempty" bson:"record_count,omitempty"`
	FileURL       string                 `json:"file_url,omitempty" bson:"file_url,omitempty"`
	FileSize      int64                  `json:"file_size,omitempty" bson:"file_size,omitempty"` // bytes
	ErrorMessage  string                 `json:"error_message,omitempty" bson:"error_message,omitempty"`
	ExecutedBy    primitive.ObjectID     `json:"executed_by" bson:"executed_by"`
	Parameters    map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
}

// Report represents a saved or scheduled report configuration
type Report struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Basic Information
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
	Type        ReportType `json:"type" bson:"type"`
	Format      ReportFormat `json:"format" bson:"format"`

	// Configuration
	Filters    map[string]interface{} `json:"filters,omitempty" bson:"filters,omitempty"`
	Columns    []string               `json:"columns,omitempty" bson:"columns,omitempty"`
	SortBy     string                 `json:"sort_by,omitempty" bson:"sort_by,omitempty"`
	SortOrder  string                 `json:"sort_order,omitempty" bson:"sort_order,omitempty"`
	Limit      int64                  `json:"limit,omitempty" bson:"limit,omitempty"`

	// Scheduling
	Schedule ReportSchedule `json:"schedule,omitempty" bson:"schedule,omitempty"`

	// Status
	Active    bool   `json:"active" bson:"active"`
	IsPublic  bool   `json:"is_public" bson:"is_public"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`

	// Statistics
	ExecutionCount    int64      `json:"execution_count" bson:"execution_count"`
	LastExecutionAt   *time.Time `json:"last_execution_at,omitempty" bson:"last_execution_at,omitempty"`
	LastExecutionStatus ReportStatus `json:"last_execution_status,omitempty" bson:"last_execution_status,omitempty"`

	// Metadata
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Category string   `json:"category,omitempty" bson:"category,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewReport creates a new report
func NewReport(name string, reportType ReportType, format ReportFormat, createdBy primitive.ObjectID) *Report {
	now := time.Now()
	return &Report{
		Name:      name,
		Type:      reportType,
		Format:    format,
		Active:    true,
		IsPublic:  false,
		CreatedBy: createdBy,
		Filters:   make(map[string]interface{}),
		Tags:      []string{},
		Schedule: ReportSchedule{
			Enabled:   false,
			Frequency: ReportFrequencyNone,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CalculateNextRun calculates the next scheduled run time
func (r *Report) CalculateNextRun() *time.Time {
	if !r.Schedule.Enabled || r.Schedule.Frequency == ReportFrequencyNone {
		return nil
	}

	now := time.Now()
	var nextRun time.Time

	switch r.Schedule.Frequency {
	case ReportFrequencyDaily:
		// Run at the specified time tomorrow
		nextRun = now.AddDate(0, 0, 1)
	case ReportFrequencyWeekly:
		// Run on the specified day of week
		daysUntil := (r.Schedule.DayOfWeek - int(now.Weekday()) + 7) % 7
		if daysUntil == 0 {
			daysUntil = 7
		}
		nextRun = now.AddDate(0, 0, daysUntil)
	case ReportFrequencyMonthly:
		// Run on the specified day of month
		nextRun = now.AddDate(0, 1, 0)
		nextRun = time.Date(nextRun.Year(), nextRun.Month(), r.Schedule.DayOfMonth, 0, 0, 0, 0, nextRun.Location())
	case ReportFrequencyYearly:
		nextRun = now.AddDate(1, 0, 0)
	default:
		return nil
	}

	return &nextRun
}

// IncrementExecutionCount increments the execution count and updates last execution info
func (r *Report) IncrementExecutionCount(status ReportStatus) {
	r.ExecutionCount++
	now := time.Now()
	r.LastExecutionAt = &now
	r.LastExecutionStatus = status
	r.UpdatedAt = now
}

// NewReportExecution creates a new report execution record
func NewReportExecution(reportID primitive.ObjectID, executedBy primitive.ObjectID, parameters map[string]interface{}) *ReportExecution {
	now := time.Now()
	return &ReportExecution{
		ReportID:   reportID,
		Status:     ReportStatusPending,
		StartedAt:  now,
		ExecutedBy: executedBy,
		Parameters: parameters,
		CreatedAt:  now,
	}
}

// MarkAsRunning marks the execution as running
func (e *ReportExecution) MarkAsRunning() {
	e.Status = ReportStatusRunning
}

// MarkAsCompleted marks the execution as completed
func (e *ReportExecution) MarkAsCompleted(recordCount int64, fileURL string, fileSize int64) {
	now := time.Now()
	e.Status = ReportStatusCompleted
	e.CompletedAt = &now
	e.Duration = now.Sub(e.StartedAt).Milliseconds()
	e.RecordCount = recordCount
	e.FileURL = fileURL
	e.FileSize = fileSize
}

// MarkAsFailed marks the execution as failed
func (e *ReportExecution) MarkAsFailed(errorMessage string) {
	now := time.Now()
	e.Status = ReportStatusFailed
	e.CompletedAt = &now
	e.Duration = now.Sub(e.StartedAt).Milliseconds()
	e.ErrorMessage = errorMessage
}
