package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeMedical        DocumentType = "medical"
	DocumentTypeAdoption       DocumentType = "adoption"
	DocumentTypeLegal          DocumentType = "legal"
	DocumentTypeFinancial      DocumentType = "financial"
	DocumentTypePolicy         DocumentType = "policy"
	DocumentTypeContract       DocumentType = "contract"
	DocumentTypeCertificate    DocumentType = "certificate"
	DocumentTypeReport         DocumentType = "report"
	DocumentTypePhoto          DocumentType = "photo"
	DocumentTypeVideo          DocumentType = "video"
	DocumentTypeOther          DocumentType = "other"
)

// Document represents a file/document in the system
type Document struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Type        DocumentType       `json:"type" bson:"type"`

	// File information
	FileName     string `json:"file_name" bson:"file_name"`
	FileSize     int64  `json:"file_size" bson:"file_size"`           // in bytes
	MimeType     string `json:"mime_type" bson:"mime_type"`
	FileURL      string `json:"file_url" bson:"file_url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"`

	// Related entities
	RelatedEntity   string              `json:"related_entity,omitempty" bson:"related_entity,omitempty"`         // "animal", "adoption", "donor", etc.
	RelatedEntityID *primitive.ObjectID `json:"related_entity_id,omitempty" bson:"related_entity_id,omitempty"`

	// Access control
	IsPublic     bool                  `json:"is_public" bson:"is_public"`
	IsConfidential bool                `json:"is_confidential" bson:"is_confidential"`
	AccessibleBy []primitive.ObjectID  `json:"accessible_by,omitempty" bson:"accessible_by,omitempty"` // User IDs

	// Version control
	Version       int                  `json:"version" bson:"version"`
	PreviousVersion *primitive.ObjectID `json:"previous_version,omitempty" bson:"previous_version,omitempty"`

	// Expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	IsExpired bool       `json:"is_expired" bson:"is_expired"`

	// Metadata
	Tags       []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty" bson:"metadata,omitempty"`

	// Tracking
	UploadedBy primitive.ObjectID `json:"uploaded_by" bson:"uploaded_by"`
	UploadedAt time.Time          `json:"uploaded_at" bson:"uploaded_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`

	// Download tracking
	DownloadCount int       `json:"download_count" bson:"download_count"`
	LastDownloadAt *time.Time `json:"last_download_at,omitempty" bson:"last_download_at,omitempty"`
}

// NewDocument creates a new document
func NewDocument(title, fileName string, fileSize int64, mimeType string, docType DocumentType, uploadedBy primitive.ObjectID) *Document {
	now := time.Now()
	return &Document{
		ID:             primitive.NewObjectID(),
		Title:          title,
		Type:           docType,
		FileName:       fileName,
		FileSize:       fileSize,
		MimeType:       mimeType,
		IsPublic:       false,
		IsConfidential: false,
		Version:        1,
		IsExpired:      false,
		Tags:           []string{},
		Metadata:       make(map[string]string),
		AccessibleBy:   []primitive.ObjectID{},
		UploadedBy:     uploadedBy,
		UploadedAt:     now,
		UpdatedAt:      now,
		DownloadCount:  0,
	}
}

// IncrementDownloadCount increments the download counter
func (d *Document) IncrementDownloadCount() {
	d.DownloadCount++
	now := time.Now()
	d.LastDownloadAt = &now
	d.UpdatedAt = now
}

// GrantAccessTo grants access to a user
func (d *Document) GrantAccessTo(userID primitive.ObjectID) {
	// Check if user already has access
	for _, id := range d.AccessibleBy {
		if id == userID {
			return
		}
	}
	d.AccessibleBy = append(d.AccessibleBy, userID)
	d.UpdatedAt = time.Now()
}

// RevokeAccessFrom revokes access from a user
func (d *Document) RevokeAccessFrom(userID primitive.ObjectID) {
	for i, id := range d.AccessibleBy {
		if id == userID {
			d.AccessibleBy = append(d.AccessibleBy[:i], d.AccessibleBy[i+1:]...)
			d.UpdatedAt = time.Now()
			return
		}
	}
}

// HasAccess checks if a user has access to the document
func (d *Document) HasAccess(userID primitive.ObjectID) bool {
	// Public documents are accessible to everyone
	if d.IsPublic {
		return true
	}

	// Uploader always has access
	if d.UploadedBy == userID {
		return true
	}

	// Check if user is in access list
	for _, id := range d.AccessibleBy {
		if id == userID {
			return true
		}
	}

	return false
}

// CheckExpiration checks and updates expiration status
func (d *Document) CheckExpiration() {
	if d.ExpiresAt != nil && time.Now().After(*d.ExpiresAt) {
		d.IsExpired = true
		d.UpdatedAt = time.Now()
	}
}

// CreateNewVersion creates a new version of the document
func (d *Document) CreateNewVersion(newFileName string, newFileSize int64, newMimeType, newFileURL string, uploadedBy primitive.ObjectID) *Document {
	now := time.Now()
	newDoc := &Document{
		ID:              primitive.NewObjectID(),
		Title:           d.Title,
		Description:     d.Description,
		Type:            d.Type,
		FileName:        newFileName,
		FileSize:        newFileSize,
		MimeType:        newMimeType,
		FileURL:         newFileURL,
		RelatedEntity:   d.RelatedEntity,
		RelatedEntityID: d.RelatedEntityID,
		IsPublic:        d.IsPublic,
		IsConfidential:  d.IsConfidential,
		AccessibleBy:    d.AccessibleBy,
		Version:         d.Version + 1,
		PreviousVersion: &d.ID,
		ExpiresAt:       d.ExpiresAt,
		IsExpired:       false,
		Tags:            d.Tags,
		Metadata:        d.Metadata,
		UploadedBy:      uploadedBy,
		UploadedAt:      now,
		UpdatedAt:       now,
		DownloadCount:   0,
	}
	return newDoc
}

// GrantAccess is an alias for GrantAccessTo
func (d *Document) GrantAccess(userID primitive.ObjectID) {
	d.GrantAccessTo(userID)
}

// RevokeAccess is an alias for RevokeAccessFrom
func (d *Document) RevokeAccess(userID primitive.ObjectID) {
	d.RevokeAccessFrom(userID)
}

// MakePublic makes the document public
func (d *Document) MakePublic() {
	d.IsPublic = true
	d.UpdatedAt = time.Now()
}

// MakePrivate makes the document private
func (d *Document) MakePrivate() {
	d.IsPublic = false
	d.UpdatedAt = time.Now()
}
