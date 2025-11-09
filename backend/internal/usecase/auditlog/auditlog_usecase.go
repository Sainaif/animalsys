package auditlog

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLogUseCase struct {
	auditLogRepo repositories.AuditLogRepository
	userRepo     repositories.UserRepository
}

func NewAuditLogUseCase(
	auditLogRepo repositories.AuditLogRepository,
	userRepo repositories.UserRepository,
) *AuditLogUseCase {
	return &AuditLogUseCase{
		auditLogRepo: auditLogRepo,
		userRepo:     userRepo,
	}
}

// GetAuditLogByID retrieves an audit log by ID
func (uc *AuditLogUseCase) GetAuditLogByID(ctx context.Context, id primitive.ObjectID) (*entities.AuditLog, error) {
	log, err := uc.auditLogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return log, nil
}

// ListAuditLogs lists audit logs with filtering and pagination
func (uc *AuditLogUseCase) ListAuditLogs(ctx context.Context, filter *repositories.AuditLogFilter) ([]*entities.AuditLog, int64, error) {
	return uc.auditLogRepo.List(ctx, *filter)
}

// GetUserActivity gets all audit logs for a specific user
func (uc *AuditLogUseCase) GetUserActivity(ctx context.Context, userID primitive.ObjectID, limit, offset int64) ([]*entities.AuditLog, int64, error) {
	filter := repositories.AuditLogFilter{
		UserID: &userID,
		Limit:  limit,
		Offset: offset,
	}

	return uc.auditLogRepo.List(ctx, filter)
}

// GetEntityHistory gets all audit logs for a specific entity
func (uc *AuditLogUseCase) GetEntityHistory(ctx context.Context, entityType string, entityID primitive.ObjectID, limit, offset int64) ([]*entities.AuditLog, int64, error) {
	filter := repositories.AuditLogFilter{
		EntityType: entityType,
		EntityID:   &entityID,
		Limit:      limit,
		Offset:     offset,
	}

	return uc.auditLogRepo.List(ctx, filter)
}

// GetRecentActivity gets recent audit logs across all users
func (uc *AuditLogUseCase) GetRecentActivity(ctx context.Context, limit int64) ([]*entities.AuditLog, int64, error) {
	filter := repositories.AuditLogFilter{
		Limit:  limit,
		Offset: 0,
	}

	return uc.auditLogRepo.List(ctx, filter)
}

// GetActionLogs gets all audit logs for a specific action type
func (uc *AuditLogUseCase) GetActionLogs(ctx context.Context, action string, limit, offset int64) ([]*entities.AuditLog, int64, error) {
	filter := repositories.AuditLogFilter{
		Action: action,
		Limit:  limit,
		Offset: offset,
	}

	return uc.auditLogRepo.List(ctx, filter)
}

// GetLogsForDateRange gets audit logs within a date range
func (uc *AuditLogUseCase) GetLogsForDateRange(ctx context.Context, fromDate, toDate time.Time, limit, offset int64) ([]*entities.AuditLog, int64, error) {
	from := primitive.NewDateTimeFromTime(fromDate)
	to := primitive.NewDateTimeFromTime(toDate)

	filter := repositories.AuditLogFilter{
		FromDate: &from,
		ToDate:   &to,
		Limit:    limit,
		Offset:   offset,
	}

	return uc.auditLogRepo.List(ctx, filter)
}

// GetAuditStatistics gets audit log statistics
func (uc *AuditLogUseCase) GetAuditStatistics(ctx context.Context) (*AuditStatistics, error) {
	// Get all logs (limited to recent activity for statistics)
	filter := repositories.AuditLogFilter{
		Limit:  1000, // Sample size for statistics
		Offset: 0,
	}

	logs, total, err := uc.auditLogRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	stats := &AuditStatistics{
		TotalLogs:    total,
		ByAction:     make(map[string]int64),
		ByEntityType: make(map[string]int64),
		ByUser:       make(map[string]int64),
	}

	// Count by action
	for _, log := range logs {
		stats.ByAction[string(log.Action)]++
		stats.ByEntityType[log.EntityType]++

		// Get user info for statistics
		if user, err := uc.userRepo.FindByID(ctx, log.UserID); err == nil {
			userName := user.FirstName + " " + user.LastName
			stats.ByUser[userName]++
		}
	}

	// Get today's activity
	today := time.Now().Truncate(24 * time.Hour)
	todayDateTime := primitive.NewDateTimeFromTime(today)
	todayFilter := repositories.AuditLogFilter{
		FromDate: &todayDateTime,
		Limit:    10000,
		Offset:   0,
	}
	_, todayCount, err := uc.auditLogRepo.List(ctx, todayFilter)
	if err == nil {
		stats.LogsToday = todayCount
	}

	// Get this week's activity
	weekStart := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	weekDateTime := primitive.NewDateTimeFromTime(weekStart)
	weekFilter := repositories.AuditLogFilter{
		FromDate: &weekDateTime,
		Limit:    10000,
		Offset:   0,
	}
	_, weekCount, err := uc.auditLogRepo.List(ctx, weekFilter)
	if err == nil {
		stats.LogsThisWeek = weekCount
	}

	// Get this month's activity
	monthStart := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	monthDateTime := primitive.NewDateTimeFromTime(monthStart)
	monthFilter := repositories.AuditLogFilter{
		FromDate: &monthDateTime,
		Limit:    10000,
		Offset:   0,
	}
	_, monthCount, err := uc.auditLogRepo.List(ctx, monthFilter)
	if err == nil {
		stats.LogsThisMonth = monthCount
	}

	return stats, nil
}

// DeleteOldLogs deletes audit logs older than specified days (admin only)
func (uc *AuditLogUseCase) DeleteOldLogs(ctx context.Context, days int) (int64, error) {
	return uc.auditLogRepo.DeleteOlderThan(ctx, days)
}

// AuditStatistics represents audit log statistics
type AuditStatistics struct {
	TotalLogs      int64            `json:"total_logs"`
	LogsToday      int64            `json:"logs_today"`
	LogsThisWeek   int64            `json:"logs_this_week"`
	LogsThisMonth  int64            `json:"logs_this_month"`
	ByAction       map[string]int64 `json:"by_action"`
	ByEntityType   map[string]int64 `json:"by_entity_type"`
	ByUser         map[string]int64 `json:"by_user"`
}
