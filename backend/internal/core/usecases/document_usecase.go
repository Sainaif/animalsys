package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type DocumentUseCase struct {
	documentRepo interfaces.DocumentRepository
	auditRepo    interfaces.AuditLogRepository
}

func NewDocumentUseCase(
	documentRepo interfaces.DocumentRepository,
	auditRepo interfaces.AuditLogRepository,
) *DocumentUseCase {
	return &DocumentUseCase{
		documentRepo: documentRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *DocumentUseCase) Create(ctx context.Context, req *entities.DocumentCreateRequest, uploadedBy string) (*entities.Document, error) {
	document := entities.NewDocument(
		req.Name,
		req.Type,
		req.FileURL,
		req.FileSize,
		req.MimeType,
	)

	document.EntityType = req.EntityType
	document.EntityID = req.EntityID
	document.Description = req.Description
	document.Tags = req.Tags
	document.ExpiryDate = req.ExpiryDate
	document.UploadedBy = uploadedBy

	if err := uc.documentRepo.Create(ctx, document); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(uploadedBy, "", "", entities.ActionCreate, "document", document.ID.Hex(), "Document uploaded: "+req.Name)
	uc.auditRepo.Create(ctx, auditLog)

	return document, nil
}

func (uc *DocumentUseCase) GetByID(ctx context.Context, id string) (*entities.Document, error) {
	return uc.documentRepo.GetByID(ctx, id)
}

func (uc *DocumentUseCase) List(ctx context.Context, filter *entities.DocumentFilter) ([]*entities.Document, int64, error) {
	return uc.documentRepo.List(ctx, filter)
}

func (uc *DocumentUseCase) Update(ctx context.Context, id string, req *entities.DocumentUpdateRequest, updatedBy string) (*entities.Document, error) {
	document, err := uc.documentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		document.Name = req.Name
	}
	if req.Description != "" {
		document.Description = req.Description
	}
	if req.Tags != nil {
		document.Tags = req.Tags
	}
	if req.ExpiryDate != nil {
		document.ExpiryDate = req.ExpiryDate
	}
	document.UpdatedBy = updatedBy

	if err := uc.documentRepo.Update(ctx, id, document); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "document", id, "Document updated")
	uc.auditRepo.Create(ctx, auditLog)

	return document, nil
}

func (uc *DocumentUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.documentRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "document", id, "Document deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *DocumentUseCase) GetByEntity(ctx context.Context, entityType, entityID string, limit, offset int) ([]*entities.Document, int64, error) {
	return uc.documentRepo.GetByEntity(ctx, entityType, entityID, limit, offset)
}

func (uc *DocumentUseCase) GetByType(ctx context.Context, docType entities.DocumentType, limit, offset int) ([]*entities.Document, int64, error) {
	return uc.documentRepo.GetByType(ctx, docType, limit, offset)
}

func (uc *DocumentUseCase) GetExpiringSoon(ctx context.Context, days int) ([]*entities.Document, error) {
	return uc.documentRepo.GetExpiringSoon(ctx, days)
}

func (uc *DocumentUseCase) SearchByTags(ctx context.Context, tags []string, limit, offset int) ([]*entities.Document, int64, error) {
	return uc.documentRepo.SearchByTags(ctx, tags, limit, offset)
}
