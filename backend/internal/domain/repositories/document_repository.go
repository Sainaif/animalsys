package repositories

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentRepository defines the interface for document data access
type DocumentRepository interface {
	// Create creates a new document
	Create(ctx context.Context, document *entities.Document) error

	// FindByID finds a document by ID
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Document, error)

	// Update updates an existing document
	Update(ctx context.Context, document *entities.Document) error

	// Delete deletes a document by ID
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List returns a list of documents with pagination and filters
	List(ctx context.Context, filter *DocumentFilter) ([]*entities.Document, int64, error)

	// GetByRelatedEntity returns documents related to a specific entity
	GetByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID) ([]*entities.Document, error)

	// GetByType returns documents by type
	GetByType(ctx context.Context, docType entities.DocumentType) ([]*entities.Document, error)

	// GetByUploader returns documents uploaded by a specific user
	GetByUploader(ctx context.Context, userID primitive.ObjectID) ([]*entities.Document, error)

	// GetPublicDocuments returns all public documents
	GetPublicDocuments(ctx context.Context) ([]*entities.Document, error)

	// GetExpiredDocuments returns all expired documents
	GetExpiredDocuments(ctx context.Context) ([]*entities.Document, error)

	// GetExpiringSoonDocuments returns documents expiring soon
	GetExpiringSoonDocuments(ctx context.Context) ([]*entities.Document, error)

	// GetDocumentVersions returns all versions of a document
	GetDocumentVersions(ctx context.Context, documentID primitive.ObjectID) ([]*entities.Document, error)

	// IncrementDownloadCount increments the download counter
	IncrementDownloadCount(ctx context.Context, id primitive.ObjectID) error

	// GetDocumentStatistics returns document statistics
	GetDocumentStatistics(ctx context.Context) (*DocumentStatistics, error)

	// EnsureIndexes creates necessary indexes for the documents collection
	EnsureIndexes(ctx context.Context) error
}

// DocumentFilter defines filter criteria for listing documents
type DocumentFilter struct {
	Type            string
	RelatedEntity   string
	RelatedEntityID *primitive.ObjectID
	UploadedBy      *primitive.ObjectID
	IsPublic        *bool
	IsConfidential  *bool
	IsExpired       *bool
	Search          string
	Tags            []string
	IncludeArchived bool
	UserID          *primitive.ObjectID // For access control
	Limit           int64
	Offset          int64
	SortBy          string // Field to sort by
	SortOrder       string // "asc" or "desc"
}

// DocumentStatistics represents document statistics
type DocumentStatistics struct {
	TotalDocuments   int64            `json:"total_documents"`
	ByType           map[string]int64 `json:"by_type"`
	TotalSize        int64            `json:"total_size"` // in bytes
	PublicDocuments  int64            `json:"public_documents"`
	PrivateDocuments int64            `json:"private_documents"`
	ExpiredDocuments int64            `json:"expired_documents"`
	TotalDownloads   int64            `json:"total_downloads"`
}
