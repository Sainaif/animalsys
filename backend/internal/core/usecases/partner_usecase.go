package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type PartnerUseCase struct {
	partnerRepo interfaces.PartnerRepository
	auditRepo   interfaces.AuditLogRepository
}

func NewPartnerUseCase(
	partnerRepo interfaces.PartnerRepository,
	auditRepo interfaces.AuditLogRepository,
) *PartnerUseCase {
	return &PartnerUseCase{
		partnerRepo: partnerRepo,
		auditRepo:   auditRepo,
	}
}

func (uc *PartnerUseCase) Create(ctx context.Context, req *entities.PartnerCreateRequest, createdBy string) (*entities.Partner, error) {
	partner := entities.NewPartner(
		req.Name,
		req.Type,
		req.ContactPerson,
		req.Email,
		req.Phone,
	)

	partner.Address = req.Address
	partner.Website = req.Website
	partner.Description = req.Description
	partner.Services = req.Services
	partner.AgreementDetails = req.AgreementDetails
	partner.Notes = req.Notes
	partner.CreatedBy = createdBy

	if err := uc.partnerRepo.Create(ctx, partner); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "partner", partner.ID.Hex(), "Partner created: "+req.Name)
	uc.auditRepo.Create(ctx, auditLog)

	return partner, nil
}

func (uc *PartnerUseCase) GetByID(ctx context.Context, id string) (*entities.Partner, error) {
	return uc.partnerRepo.GetByID(ctx, id)
}

func (uc *PartnerUseCase) List(ctx context.Context, partnerType entities.PartnerType, limit, offset int) ([]*entities.Partner, int64, error) {
	return uc.partnerRepo.List(ctx, partnerType, limit, offset)
}

func (uc *PartnerUseCase) Update(ctx context.Context, id string, req *entities.PartnerUpdateRequest, updatedBy string) (*entities.Partner, error) {
	partner, err := uc.partnerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		partner.Name = req.Name
	}
	if req.Type != "" {
		partner.Type = req.Type
	}
	if req.ContactPerson != "" {
		partner.ContactPerson = req.ContactPerson
	}
	if req.Email != "" {
		partner.Email = req.Email
	}
	if req.Phone != "" {
		partner.Phone = req.Phone
	}
	if req.Address != "" {
		partner.Address = req.Address
	}
	if req.Website != "" {
		partner.Website = req.Website
	}
	if req.Description != "" {
		partner.Description = req.Description
	}
	if req.Services != nil {
		partner.Services = req.Services
	}
	if req.AgreementDetails != "" {
		partner.AgreementDetails = req.AgreementDetails
	}
	if req.Active != nil {
		partner.Active = *req.Active
	}
	if req.Notes != "" {
		partner.Notes = req.Notes
	}
	partner.UpdatedBy = updatedBy

	if err := uc.partnerRepo.Update(ctx, id, partner); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "partner", id, "Partner updated")
	uc.auditRepo.Create(ctx, auditLog)

	return partner, nil
}

func (uc *PartnerUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.partnerRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "partner", id, "Partner deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *PartnerUseCase) GetActivePartners(ctx context.Context) ([]*entities.Partner, error) {
	return uc.partnerRepo.GetActivePartners(ctx)
}

func (uc *PartnerUseCase) AddCollaboration(ctx context.Context, id string, req *entities.CollaborationCreateRequest, addedBy string) error {
	collaboration := entities.Collaboration{
		Date:        req.Date,
		Type:        req.Type,
		Description: req.Description,
		Outcome:     req.Outcome,
		Value:       req.Value,
	}

	if err := uc.partnerRepo.AddCollaboration(ctx, id, collaboration); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(addedBy, "", "", entities.ActionUpdate, "partner", id, "Collaboration added")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *PartnerUseCase) GetPartnerStatistics(ctx context.Context) (map[string]interface{}, error) {
	partners, total, err := uc.partnerRepo.List(ctx, "", 0, 0)
	if err != nil {
		return nil, err
	}

	activeCount := 0
	byType := make(map[string]int)
	totalCollaborations := 0

	for _, partner := range partners {
		if partner.Active {
			activeCount++
		}
		byType[string(partner.Type)]++
		totalCollaborations += len(partner.CollaborationHistory)
	}

	stats := map[string]interface{}{
		"total_partners":        total,
		"active_partners":       activeCount,
		"by_type":               byType,
		"total_collaborations":  totalCollaborations,
	}

	return stats, nil
}
