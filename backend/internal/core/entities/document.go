package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeMedical  DocumentType = "medical"
	DocumentTypeLegal    DocumentType = "legal"
	DocumentTypeContract DocumentType = "contract"
	DocumentTypeInvoice  DocumentType = "invoice"
	DocumentTypeReport   DocumentType = "report"
	DocumentTypePhoto    DocumentType = "photo"
	DocumentTypeOther    DocumentType = "other"
)

// Document represents a document in the system
type Document struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileName          string             `bson:"file_name" json:"file_name"`
	OriginalFileName  string             `bson:"original_file_name" json:"original_file_name"`
	FileType          DocumentType       `bson:"file_type" json:"file_type"`
	ContentType       string             `bson:"content_type" json:"content_type"` // MIME type
	FileSize          int64              `bson:"file_size" json:"file_size"`       // in bytes
	FilePath          string             `bson:"file_path" json:"file_path"`
	FileURL           string             `bson:"file_url,omitempty" json:"file_url,omitempty"`
	GridFSID          string             `bson:"gridfs_id,omitempty" json:"gridfs_id,omitempty"`
	RelatedEntityType string             `bson:"related_entity_type,omitempty" json:"related_entity_type,omitempty"` // animal, adoption, user, etc.
	RelatedEntityID   string             `bson:"related_entity_id,omitempty" json:"related_entity_id,omitempty"`
	Tags              []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	ExpiryDate        *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty"`
	UploadedBy        string             `bson:"uploaded_by" json:"uploaded_by"`
	UploadedAt        time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// DocumentUploadRequest represents document upload request
type DocumentUploadRequest struct {
	FileType          DocumentType `json:"file_type" validate:"required,oneof=medical legal contract invoice report photo other"`
	RelatedEntityType string       `json:"related_entity_type,omitempty"`
	RelatedEntityID   string       `json:"related_entity_id,omitempty"`
	Tags              []string     `json:"tags,omitempty"`
	ExpiryDate        *time.Time   `json:"expiry_date,omitempty"`
	Description       string       `json:"description,omitempty"`
}

// DocumentFilter represents filters for querying documents
type DocumentFilter struct {
	FileType          DocumentType
	RelatedEntityType string
	RelatedEntityID   string
	Tags              []string
	Search            string
	Limit             int
	Offset            int
	SortBy            string
	SortOrder         string
}

// NewDocument creates a new document
func NewDocument(fileName, originalFileName string, fileType DocumentType, contentType string, fileSize int64, uploadedBy string) *Document {
	now := time.Now()
	return &Document{
		ID:               primitive.NewObjectID(),
		FileName:         fileName,
		OriginalFileName: originalFileName,
		FileType:         fileType,
		ContentType:      contentType,
		FileSize:         fileSize,
		UploadedBy:       uploadedBy,
		UploadedAt:       now,
		UpdatedAt:        now,
	}
}
