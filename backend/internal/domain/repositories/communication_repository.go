package repositories

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommunicationTemplateFilter represents filters for template queries
type CommunicationTemplateFilter struct {
	Type      string
	Category  string
	Active    *bool
	IsDefault *bool
	Language  string
	Search    string
	SortBy    string
	SortOrder string
	Limit     int64
	Offset    int64
}

// CommunicationTemplateRepository defines the interface for template data access
type CommunicationTemplateRepository interface {
	Create(ctx context.Context, template *entities.CommunicationTemplate) error
	Update(ctx context.Context, template *entities.CommunicationTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.CommunicationTemplate, error)
	List(ctx context.Context, filter *CommunicationTemplateFilter) ([]*entities.CommunicationTemplate, int64, error)
	GetByCategory(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) ([]*entities.CommunicationTemplate, error)
	GetDefault(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) (*entities.CommunicationTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*entities.CommunicationTemplate, error)
	IncrementUsage(ctx context.Context, id primitive.ObjectID) error
	EnsureIndexes(ctx context.Context) error
}

// CommunicationFilter represents filters for communication queries
type CommunicationFilter struct {
	Type           string
	Category       string
	Status         string
	RecipientType  string
	RecipientID    *primitive.ObjectID
	SenderID       *primitive.ObjectID
	TemplateID     *primitive.ObjectID
	CampaignID     *primitive.ObjectID
	BatchID        string
	RelatedType    string
	RelatedID      *primitive.ObjectID
	StartDate      *time.Time
	EndDate        *time.Time
	SortBy         string
	SortOrder      string
	Limit          int64
	Offset         int64
}

// CommunicationStatistics represents communication statistics
type CommunicationStatistics struct {
	TotalSent       int64
	TotalDelivered  int64
	TotalFailed     int64
	TotalBounced    int64
	TotalOpened     int64
	TotalClicked    int64
	OpenRate        float64
	ClickRate       float64
	DeliveryRate    float64
	ByType          map[string]int64
	ByCategory      map[string]int64
	ByStatus        map[string]int64
}

// CommunicationRepository defines the interface for communication data access
type CommunicationRepository interface {
	Create(ctx context.Context, communication *entities.Communication) error
	Update(ctx context.Context, communication *entities.Communication) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Communication, error)
	List(ctx context.Context, filter *CommunicationFilter) ([]*entities.Communication, int64, error)
	GetPending(ctx context.Context) ([]*entities.Communication, error)
	GetForRetry(ctx context.Context) ([]*entities.Communication, error)
	GetByRecipient(ctx context.Context, recipientType entities.RecipientType, recipientID primitive.ObjectID) ([]*entities.Communication, error)
	GetByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Communication, error)
	GetByBatch(ctx context.Context, batchID string) ([]*entities.Communication, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status entities.CommunicationStatus) error
	MarkAsOpened(ctx context.Context, id primitive.ObjectID) error
	MarkAsClicked(ctx context.Context, id primitive.ObjectID) error
	GetStatistics(ctx context.Context, startDate, endDate time.Time) (*CommunicationStatistics, error)
	EnsureIndexes(ctx context.Context) error
}

// NotificationFilter represents filters for notification queries
type NotificationFilter struct {
	UserID     *primitive.ObjectID
	Type       string
	Priority   string
	Read       *bool
	Dismissed  *bool
	Category   string
	RelatedType string
	RelatedID  *primitive.ObjectID
	SortBy     string
	SortOrder  string
	Limit      int64
	Offset     int64
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	Update(ctx context.Context, notification *entities.Notification) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Notification, error)
	List(ctx context.Context, filter *NotificationFilter) ([]*entities.Notification, int64, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error)
	GetUnreadByUser(ctx context.Context, userID primitive.ObjectID) ([]*entities.Notification, error)
	CountUnreadByUser(ctx context.Context, userID primitive.ObjectID) (int64, error)
	MarkAsRead(ctx context.Context, id primitive.ObjectID) error
	MarkAsUnread(ctx context.Context, id primitive.ObjectID) error
	MarkAllAsRead(ctx context.Context, userID primitive.ObjectID) error
	Dismiss(ctx context.Context, id primitive.ObjectID) error
	DeleteExpired(ctx context.Context) (int64, error)
	FindByGroupKey(ctx context.Context, userID primitive.ObjectID, groupKey string) (*entities.Notification, error)
	EnsureIndexes(ctx context.Context) error
}
