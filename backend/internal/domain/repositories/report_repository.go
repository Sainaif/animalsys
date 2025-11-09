package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportFilter represents filters for report queries
type ReportFilter struct {
	Type      string
	Format    string
	Category  string
	Active    *bool
	IsPublic  *bool
	CreatedBy *primitive.ObjectID
	Search    string
	Tags      []string
	SortBy    string
	SortOrder string
	Limit     int64
	Offset    int64
}

// ReportExecutionFilter represents filters for report execution queries
type ReportExecutionFilter struct {
	ReportID   *primitive.ObjectID
	Status     string
	ExecutedBy *primitive.ObjectID
	StartDate  *time.Time
	EndDate    *time.Time
	SortBy     string
	SortOrder  string
	Limit      int64
	Offset     int64
}

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	Create(ctx context.Context, report *entities.Report) error
	Update(ctx context.Context, report *entities.Report) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Report, error)
	List(ctx context.Context, filter *ReportFilter) ([]*entities.Report, int64, error)
	GetActiveReports(ctx context.Context) ([]*entities.Report, error)
	GetPublicReports(ctx context.Context) ([]*entities.Report, error)
	GetScheduledReports(ctx context.Context) ([]*entities.Report, error)
	GetReportsForExecution(ctx context.Context) ([]*entities.Report, error)
	IncrementExecutionCount(ctx context.Context, id primitive.ObjectID, status entities.ReportStatus) error
	EnsureIndexes(ctx context.Context) error
}

// ReportExecutionRepository defines the interface for report execution data access
type ReportExecutionRepository interface {
	Create(ctx context.Context, execution *entities.ReportExecution) error
	Update(ctx context.Context, execution *entities.ReportExecution) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.ReportExecution, error)
	List(ctx context.Context, filter *ReportExecutionFilter) ([]*entities.ReportExecution, int64, error)
	GetByReportID(ctx context.Context, reportID primitive.ObjectID) ([]*entities.ReportExecution, error)
	GetRecentExecutions(ctx context.Context, limit int) ([]*entities.ReportExecution, error)
	DeleteOlderThan(ctx context.Context, date time.Time) (int64, error)
	EnsureIndexes(ctx context.Context) error
}
