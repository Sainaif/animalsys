package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeFinancial   ReportType = "financial"
	ReportTypeStatutory   ReportType = "statutory"
	ReportTypeAdoption    ReportType = "adoption"
	ReportTypeVolunteer   ReportType = "volunteer"
	ReportTypeInventory   ReportType = "inventory"
	ReportTypeCustom      ReportType = "custom"
)

// ReportStatus represents the status of report generation
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusProcessing ReportStatus = "processing"
	ReportStatusCompleted ReportStatus = "completed"
	ReportStatusFailed    ReportStatus = "failed"
)

// ReportFormat represents the output format of report
type ReportFormat string

const (
	ReportFormatPDF   ReportFormat = "pdf"
	ReportFormatExcel ReportFormat = "excel"
	ReportFormatCSV   ReportFormat = "csv"
	ReportFormatJSON  ReportFormat = "json"
)

// Report represents a generated report
type Report struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name        string                 `bson:"name" json:"name"`
	Type        ReportType             `bson:"type" json:"type"`
	Status      ReportStatus           `bson:"status" json:"status"`
	Format      ReportFormat           `bson:"format" json:"format"`
	StartDate   time.Time              `bson:"start_date" json:"start_date"`
	EndDate     time.Time              `bson:"end_date" json:"end_date"`
	Filters     map[string]interface{} `bson:"filters,omitempty" json:"filters,omitempty"`
	FileURL     string                 `bson:"file_url,omitempty" json:"file_url,omitempty"`
	FileSize    int64                  `bson:"file_size,omitempty" json:"file_size,omitempty"`
	Error       string                 `bson:"error,omitempty" json:"error,omitempty"`
	GeneratedAt *time.Time             `bson:"generated_at,omitempty" json:"generated_at,omitempty"`
	ExpiresAt   *time.Time             `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	CreatedBy   string                 `bson:"created_by" json:"created_by"`
}

// ReportGenerateRequest represents report generation request
type ReportGenerateRequest struct {
	Name      string                 `json:"name" validate:"required"`
	Type      ReportType             `json:"type" validate:"required,oneof=financial statutory adoption volunteer inventory custom"`
	Format    ReportFormat           `json:"format" validate:"required,oneof=pdf excel csv json"`
	StartDate time.Time              `json:"start_date" validate:"required"`
	EndDate   time.Time              `json:"end_date" validate:"required,gtfield=StartDate"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
}

// FinancialReportData represents financial report data
type FinancialReportData struct {
	Period        string  `json:"period"`
	TotalIncome   float64 `json:"total_income"`
	TotalExpense  float64 `json:"total_expense"`
	NetIncome     float64 `json:"net_income"`
	IncomeByCategory map[string]float64 `json:"income_by_category"`
	ExpenseByCategory map[string]float64 `json:"expense_by_category"`
}

// AdoptionReportData represents adoption report data
type AdoptionReportData struct {
	Period            string `json:"period"`
	TotalAdoptions    int    `json:"total_adoptions"`
	PendingApplications int  `json:"pending_applications"`
	SuccessRate       float64 `json:"success_rate"`
	AdoptionsBySpecies map[string]int `json:"adoptions_by_species"`
	AverageProcessTime float64 `json:"average_process_time_days"`
}

// VolunteerReportData represents volunteer report data
type VolunteerReportData struct {
	Period           string  `json:"period"`
	ActiveVolunteers int     `json:"active_volunteers"`
	TotalHours       float64 `json:"total_hours"`
	AverageHoursPerVolunteer float64 `json:"average_hours_per_volunteer"`
	NewVolunteers    int     `json:"new_volunteers"`
}

// InventoryReportData represents inventory report data
type InventoryReportData struct {
	TotalItems     int      `json:"total_items"`
	LowStockItems  int      `json:"low_stock_items"`
	OutOfStockItems int     `json:"out_of_stock_items"`
	TotalValue     float64  `json:"total_value"`
	ItemsByCategory map[string]int `json:"items_by_category"`
}

// NewReport creates a new report
func NewReport(name string, reportType ReportType, format ReportFormat, startDate, endDate time.Time, createdBy string) *Report {
	now := time.Now()
	expiresAt := now.AddDate(0, 0, 30) // Expire after 30 days

	return &Report{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Type:      reportType,
		Status:    ReportStatusPending,
		Format:    format,
		StartDate: startDate,
		EndDate:   endDate,
		ExpiresAt: &expiresAt,
		CreatedAt: now,
		CreatedBy: createdBy,
	}
}

// MarkAsCompleted marks report as completed
func (r *Report) MarkAsCompleted(fileURL string, fileSize int64) {
	now := time.Now()
	r.Status = ReportStatusCompleted
	r.FileURL = fileURL
	r.FileSize = fileSize
	r.GeneratedAt = &now
}

// MarkAsFailed marks report as failed
func (r *Report) MarkAsFailed(errorMsg string) {
	r.Status = ReportStatusFailed
	r.Error = errorMsg
}
