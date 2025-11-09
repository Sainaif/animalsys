package report

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportUseCase struct {
	reportRepo          repositories.ReportRepository
	reportExecutionRepo repositories.ReportExecutionRepository
	auditLogRepo        repositories.AuditLogRepository
}

func NewReportUseCase(
	reportRepo repositories.ReportRepository,
	reportExecutionRepo repositories.ReportExecutionRepository,
	auditLogRepo repositories.AuditLogRepository,
) *ReportUseCase {
	return &ReportUseCase{
		reportRepo:          reportRepo,
		reportExecutionRepo: reportExecutionRepo,
		auditLogRepo:        auditLogRepo,
	}
}

// CreateReport creates a new report
func (uc *ReportUseCase) CreateReport(ctx context.Context, report *entities.Report, userID primitive.ObjectID) error {
	if report.Name == "" {
		return errors.NewBadRequest("Report name is required")
	}

	report.CreatedBy = userID

	if err := uc.reportRepo.Create(ctx, report); err != nil {
		return err
	}

	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "report", "", "").
			WithEntityID(report.ID))

	return nil
}

// UpdateReport updates an existing report
func (uc *ReportUseCase) UpdateReport(ctx context.Context, report *entities.Report, userID primitive.ObjectID) error {
	if _, err := uc.reportRepo.FindByID(ctx, report.ID); err != nil {
		return err
	}

	if err := uc.reportRepo.Update(ctx, report); err != nil {
		return err
	}

	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "report", "", "").
			WithEntityID(report.ID))

	return nil
}

// GetReportByID retrieves a report by ID
func (uc *ReportUseCase) GetReportByID(ctx context.Context, id primitive.ObjectID) (*entities.Report, error) {
	return uc.reportRepo.FindByID(ctx, id)
}

// ListReports lists reports with filtering
func (uc *ReportUseCase) ListReports(ctx context.Context, filter *repositories.ReportFilter) ([]*entities.Report, int64, error) {
	return uc.reportRepo.List(ctx, filter)
}

// DeleteReport deletes a report
func (uc *ReportUseCase) DeleteReport(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	report, err := uc.reportRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.reportRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "report", "", "").
			WithEntityID(report.ID))

	return nil
}

// ExecuteReport executes a report and creates an execution record
func (uc *ReportUseCase) ExecuteReport(ctx context.Context, reportID primitive.ObjectID, parameters map[string]interface{}, userID primitive.ObjectID) (*entities.ReportExecution, error) {
	report, err := uc.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	if !report.Active {
		return nil, errors.NewBadRequest("Report is not active")
	}

	execution := entities.NewReportExecution(reportID, userID, parameters)
	execution.MarkAsRunning()

	if err := uc.reportExecutionRepo.Create(ctx, execution); err != nil {
		return nil, err
	}

	// Note: Actual report generation would happen here
	// For now, we'll mark it as completed immediately
	// In production, this would be handled by a background job
	execution.MarkAsCompleted(0, "", 0)
	_ = uc.reportExecutionRepo.Update(ctx, execution)

	_ = uc.reportRepo.IncrementExecutionCount(ctx, reportID, entities.ReportStatusCompleted)

	return execution, nil
}

// GetReportExecutions retrieves executions for a report
func (uc *ReportUseCase) GetReportExecutions(ctx context.Context, reportID primitive.ObjectID) ([]*entities.ReportExecution, error) {
	return uc.reportExecutionRepo.GetByReportID(ctx, reportID)
}

// GetRecentExecutions retrieves recent report executions
func (uc *ReportUseCase) GetRecentExecutions(ctx context.Context, limit int) ([]*entities.ReportExecution, error) {
	return uc.reportExecutionRepo.GetRecentExecutions(ctx, limit)
}

// GetActiveReports retrieves all active reports
func (uc *ReportUseCase) GetActiveReports(ctx context.Context) ([]*entities.Report, error) {
	return uc.reportRepo.GetActiveReports(ctx)
}

// GetPublicReports retrieves all public reports
func (uc *ReportUseCase) GetPublicReports(ctx context.Context) ([]*entities.Report, error) {
	return uc.reportRepo.GetPublicReports(ctx)
}

// CleanupOldExecutions deletes report executions older than the specified number of days
func (uc *ReportUseCase) CleanupOldExecutions(ctx context.Context, days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return uc.reportExecutionRepo.DeleteOlderThan(ctx, cutoffDate)
}
