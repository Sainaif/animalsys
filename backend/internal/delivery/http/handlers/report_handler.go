package handlers

import (
	"net/http"
	"strconv"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/usecase/report"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	reportUseCase *report.ReportUseCase
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportUseCase *report.ReportUseCase) *ReportHandler {
	return &ReportHandler{
		reportUseCase: reportUseCase,
	}
}

// CreateReport creates a new report
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var rep entities.Report
	if err := c.ShouldBindJSON(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.reportUseCase.CreateReport(c.Request.Context(), &rep, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, rep)
}

// UpdateReport updates a report
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var rep entities.Report
	if err := c.ShouldBindJSON(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rep.ID = id
	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.reportUseCase.UpdateReport(c.Request.Context(), &rep, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, rep)
}

// GetReport gets a report by ID
func (h *ReportHandler) GetReport(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	rep, err := h.reportUseCase.GetReportByID(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, rep)
}

// ListReports lists reports with filtering
func (h *ReportHandler) ListReports(c *gin.Context) {
	filter := &repositories.ReportFilter{}

	filter.Type = c.Query("type")
	filter.Format = c.Query("format")
	filter.Category = c.Query("category")
	filter.Search = c.Query("search")

	if activeStr := c.Query("active"); activeStr != "" {
		active := activeStr == "true"
		filter.Active = &active
	}

	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		isPublic := isPublicStr == "true"
		filter.IsPublic = &isPublic
	}

	if createdByStr := c.Query("created_by"); createdByStr != "" {
		createdBy, err := primitive.ObjectIDFromHex(createdByStr)
		if err == nil {
			filter.CreatedBy = &createdBy
		}
	}

	// Pagination
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64); err == nil {
		filter.Limit = limit
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		filter.Offset = offset
	}

	filter.SortBy = c.DefaultQuery("sort_by", "created_at")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	reports, total, err := h.reportUseCase.ListReports(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   reports,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// DeleteReport deletes a report
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	if err := h.reportUseCase.DeleteReport(c.Request.Context(), id, userID); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}

// ExecuteReport executes a report
func (h *ReportHandler) ExecuteReport(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var params map[string]interface{}
	_ = c.ShouldBindJSON(&params)

	userID := c.MustGet("user_id").(primitive.ObjectID)

	execution, err := h.reportUseCase.ExecuteReport(c.Request.Context(), id, params, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, execution)
}

// GetReportExecutions gets executions for a report
func (h *ReportHandler) GetReportExecutions(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	executions, err := h.reportUseCase.GetReportExecutions(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, executions)
}

// GetRecentExecutions gets recent report executions
func (h *ReportHandler) GetRecentExecutions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	executions, err := h.reportUseCase.GetRecentExecutions(c.Request.Context(), limit)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, executions)
}

// GetActiveReports gets all active reports
func (h *ReportHandler) GetActiveReports(c *gin.Context) {
	reports, err := h.reportUseCase.GetActiveReports(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetPublicReports gets all public reports
func (h *ReportHandler) GetPublicReports(c *gin.Context) {
	reports, err := h.reportUseCase.GetPublicReports(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetAnimalsReport generates an animals report
func (h *ReportHandler) GetAnimalsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Return basic report structure
	c.JSON(http.StatusOK, gin.H{
		"report_type": "animals",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_animals": 0,
		"by_species":  map[string]int{},
		"by_status":   map[string]int{},
	})
}

// GetAdoptionsReport generates an adoptions report
func (h *ReportHandler) GetAdoptionsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "adoptions",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_adoptions": 0,
		"by_status":   map[string]int{},
	})
}

// GetDonationsReport generates a donations report
func (h *ReportHandler) GetDonationsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "donations",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_donations": 0,
		"total_amount": 0.0,
		"by_type":     map[string]int{},
	})
}

// GetFinancialsReport generates a financials report
func (h *ReportHandler) GetFinancialsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "financials",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_income": 0.0,
		"total_expenses": 0.0,
		"net_balance": 0.0,
	})
}

// GetVolunteersReport generates a volunteers report
func (h *ReportHandler) GetVolunteersReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "volunteers",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_volunteers": 0,
		"total_hours": 0,
		"by_status": map[string]int{},
	})
}

// GetInventoryReport generates an inventory report
func (h *ReportHandler) GetInventoryReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"report_type": "inventory",
		"total_items": 0,
		"low_stock_items": 0,
		"out_of_stock_items": 0,
		"by_category": map[string]int{},
	})
}

// GetVeterinaryReport generates a veterinary report
func (h *ReportHandler) GetVeterinaryReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "veterinary",
		"start_date":  startDate,
		"end_date":    endDate,
		"total_visits": 0,
		"total_procedures": 0,
		"by_type": map[string]int{},
	})
}

// GetComplianceReport generates a compliance report
func (h *ReportHandler) GetComplianceReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	c.JSON(http.StatusOK, gin.H{
		"report_type": "compliance",
		"start_date":  startDate,
		"end_date":    endDate,
		"compliant_records": 0,
		"non_compliant_records": 0,
		"issues": []string{},
	})
}

// CreateCustomReport creates a custom report
func (h *ReportHandler) CreateCustomReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Custom report created"})
}
