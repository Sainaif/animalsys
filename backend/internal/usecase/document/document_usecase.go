package document

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentUseCase struct {
	documentRepo repositories.DocumentRepository
	auditLogRepo repositories.AuditLogRepository
}

func NewDocumentUseCase(
	documentRepo repositories.DocumentRepository,
	auditLogRepo repositories.AuditLogRepository,
) *DocumentUseCase {
	return &DocumentUseCase{
		documentRepo: documentRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateDocument creates a new document
func (uc *DocumentUseCase) CreateDocument(ctx context.Context, document *entities.Document, userID primitive.ObjectID) error {
	// Validate required fields
	if document.Title == "" {
		return errors.NewBadRequest("Document title is required")
	}

	if document.FileName == "" {
		return errors.NewBadRequest("File name is required")
	}

	if document.Type == "" {
		return errors.NewBadRequest("Document type is required")
	}

	if document.FileURL == "" {
		return errors.NewBadRequest("File URL is required")
	}

	document.UploadedBy = userID

	if err := uc.documentRepo.Create(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "document", document.Title, "").
			WithEntityID(document.ID))

	return nil
}

// ArchiveDocument archives a document
func (uc *DocumentUseCase) ArchiveDocument(ctx context.Context, documentID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can archive
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can archive this document")
	}

	document.Archive()

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "archived document").
			WithEntityID(documentID))

	return nil
}

// UnarchiveDocument unarchives a document
func (uc *DocumentUseCase) UnarchiveDocument(ctx context.Context, documentID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can unarchive
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can unarchive this document")
	}

	document.Unarchive()

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "unarchived document").
			WithEntityID(documentID))

	return nil
}

// GetDocumentByID retrieves a document by ID
func (uc *DocumentUseCase) GetDocumentByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*entities.Document, error) {
	document, err := uc.documentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check access permissions
	if !document.HasAccess(userID) {
		return nil, errors.NewForbidden("You don't have access to this document")
	}

	return document, nil
}

// UpdateDocument updates a document
func (uc *DocumentUseCase) UpdateDocument(ctx context.Context, document *entities.Document, userID primitive.ObjectID) error {
	// Validate required fields
	if document.Title == "" {
		return errors.NewBadRequest("Document title is required")
	}

	// Check if document exists
	existing, err := uc.documentRepo.FindByID(ctx, document.ID)
	if err != nil {
		return err
	}

	// Check access permissions
	if !existing.HasAccess(userID) {
		return errors.NewForbidden("You don't have permission to update this document")
	}

	// Preserve creation info
	document.UploadedBy = existing.UploadedBy
	document.UploadedAt = existing.UploadedAt

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "").
			WithEntityID(document.ID))

	return nil
}

// DeleteDocument deletes a document
func (uc *DocumentUseCase) DeleteDocument(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if document exists
	document, err := uc.documentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the uploader (only uploader can delete)
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can delete this document")
	}

	if err := uc.documentRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "document", document.Title, "").
			WithEntityID(id))

	return nil
}

// SearchDocuments searches documents with filtering and pagination, and access control
func (uc *DocumentUseCase) SearchDocuments(ctx context.Context, filter *repositories.DocumentFilter, userID primitive.ObjectID) ([]*entities.Document, int64, error) {
	filter.UserID = &userID
	return uc.documentRepo.List(ctx, filter)
}

// GetDocumentsByRelatedEntity gets documents by related entity, filtered by access
func (uc *DocumentUseCase) GetDocumentsByRelatedEntity(ctx context.Context, entityType string, entityID primitive.ObjectID, userID primitive.ObjectID) ([]*entities.Document, error) {
	// This repository method doesn't support access control, so we filter here.
	// For a production system, this should be moved to the repository layer.
	documents, err := uc.documentRepo.GetByRelatedEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}

	accessibleDocs := make([]*entities.Document, 0)
	for _, doc := range documents {
		if doc.HasAccess(userID) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

// GetDocumentsByType gets documents by type, filtered by access
func (uc *DocumentUseCase) GetDocumentsByType(ctx context.Context, docType entities.DocumentType, userID primitive.ObjectID) ([]*entities.Document, error) {
	// This repository method doesn't support access control, so we filter here.
	// For a production system, this should be moved to the repository layer.
	documents, err := uc.documentRepo.GetByType(ctx, docType)
	if err != nil {
		return nil, err
	}

	accessibleDocs := make([]*entities.Document, 0)
	for _, doc := range documents {
		if doc.HasAccess(userID) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

// GetDocumentsByCategory gets documents by category, filtered by access
func (uc *DocumentUseCase) GetDocumentsByCategory(ctx context.Context, category entities.DocumentType, userID primitive.ObjectID) ([]*entities.Document, error) {
	// Assuming category is treated the same as type for now
	// This repository method doesn't support access control, so we filter here.
	// For a production system, this should be moved to the repository layer.
	documents, err := uc.documentRepo.GetByType(ctx, category)
	if err != nil {
		return nil, err
	}

	accessibleDocs := make([]*entities.Document, 0)
	for _, doc := range documents {
		if doc.HasAccess(userID) {
			accessibleDocs = append(accessibleDocs, doc)
		}
	}

	return accessibleDocs, nil
}

// GetDocumentsByUploader gets documents uploaded by a specific user
func (uc *DocumentUseCase) GetDocumentsByUploader(ctx context.Context, userID primitive.ObjectID) ([]*entities.Document, error) {
	return uc.documentRepo.GetByUploader(ctx, userID)
}

// GetPublicDocuments gets all public documents
func (uc *DocumentUseCase) GetPublicDocuments(ctx context.Context) ([]*entities.Document, error) {
	return uc.documentRepo.GetPublicDocuments(ctx)
}

// GetExpiredDocuments gets all expired documents
func (uc *DocumentUseCase) GetExpiredDocuments(ctx context.Context) ([]*entities.Document, error) {
	return uc.documentRepo.GetExpiredDocuments(ctx)
}

// GetExpiringSoonDocuments gets documents expiring soon
func (uc *DocumentUseCase) GetExpiringSoonDocuments(ctx context.Context) ([]*entities.Document, error) {
	return uc.documentRepo.GetExpiringSoonDocuments(ctx)
}

// GetDocumentVersions gets all versions of a document
func (uc *DocumentUseCase) GetDocumentVersions(ctx context.Context, documentID primitive.ObjectID, userID primitive.ObjectID) ([]*entities.Document, error) {
	// Check if user has access to the main document
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	if !document.HasAccess(userID) {
		return nil, errors.NewForbidden("You don't have access to this document")
	}

	return uc.documentRepo.GetDocumentVersions(ctx, documentID)
}

// GetDocumentStatistics gets document statistics
func (uc *DocumentUseCase) GetDocumentStatistics(ctx context.Context) (*repositories.DocumentStatistics, error) {
	return uc.documentRepo.GetDocumentStatistics(ctx)
}

// DownloadDocument records a download and increments the counter
func (uc *DocumentUseCase) DownloadDocument(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*entities.Document, error) {
	document, err := uc.documentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check access permissions
	if !document.HasAccess(userID) {
		return nil, errors.NewForbidden("You don't have access to this document")
	}

	// Increment download count
	if err := uc.documentRepo.IncrementDownloadCount(ctx, id); err != nil {
		return nil, err
	}

	// Refresh document to get updated count
	document, err = uc.documentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionView, "document", document.Title, "downloaded document").
			WithEntityID(id))

	return document, nil
}

// CreateNewVersion creates a new version of a document
func (uc *DocumentUseCase) CreateNewVersion(ctx context.Context, documentID primitive.ObjectID, filePath string, fileSize int64, userID primitive.ObjectID) (*entities.Document, error) {
	// Get the existing document
	existingDoc, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to update
	if !existingDoc.HasAccess(userID) {
		return nil, errors.NewForbidden("You don't have permission to update this document")
	}

	if filePath == "" {
		return nil, errors.NewBadRequest("File path is required")
	}

	// Derive MIME type and filename from the path
	fileName := filePath
	mimeType := "application/octet-stream"

	// Create new version
	newVersion := existingDoc.CreateNewVersion(fileName, fileSize, mimeType, filePath, userID)

	if err := uc.documentRepo.Create(ctx, newVersion); err != nil {
		return nil, err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "document", newVersion.Title, "created new version").
			WithEntityID(newVersion.ID))

	return newVersion, nil
}

// GrantAccess grants access to a document for a user
func (uc *DocumentUseCase) GrantAccess(ctx context.Context, documentID, targetUserID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can grant access
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can grant access to this document")
	}

	document.GrantAccess(targetUserID)

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "granted access").
			WithEntityID(documentID))

	return nil
}

// RevokeAccess revokes access to a document for a user
func (uc *DocumentUseCase) RevokeAccess(ctx context.Context, documentID, targetUserID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can revoke access
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can revoke access to this document")
	}

	document.RevokeAccess(targetUserID)

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "revoked access").
			WithEntityID(documentID))

	return nil
}

// MakePublic makes a document public
func (uc *DocumentUseCase) MakePublic(ctx context.Context, documentID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can make public
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can make this document public")
	}

	document.MakePublic()

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "made public").
			WithEntityID(documentID))

	return nil
}

// MakePrivate makes a document private
func (uc *DocumentUseCase) MakePrivate(ctx context.Context, documentID, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can make private
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can make this document private")
	}

	document.MakePrivate()

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "made private").
			WithEntityID(documentID))

	return nil
}

// SetExpiration sets expiration date for a document
func (uc *DocumentUseCase) SetExpiration(ctx context.Context, documentID primitive.ObjectID, expiresAt time.Time, userID primitive.ObjectID) error {
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return err
	}

	// Only uploader can set expiration
	if document.UploadedBy != userID {
		return errors.NewForbidden("Only the uploader can set expiration for this document")
	}

	document.ExpiresAt = &expiresAt
	document.UpdatedAt = time.Now()

	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "document", document.Title, "set expiration date").
			WithEntityID(documentID))

	return nil
}
