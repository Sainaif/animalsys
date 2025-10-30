package usecases

import (
	"context"
	"errors"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type AdoptionUseCase struct {
	adoptionRepo interfaces.AdoptionRepository
	animalRepo   interfaces.AnimalRepository
	auditRepo    interfaces.AuditLogRepository
}

func NewAdoptionUseCase(
	adoptionRepo interfaces.AdoptionRepository,
	animalRepo interfaces.AnimalRepository,
	auditRepo interfaces.AuditLogRepository,
) *AdoptionUseCase {
	return &AdoptionUseCase{
		adoptionRepo: adoptionRepo,
		animalRepo:   animalRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *AdoptionUseCase) Create(ctx context.Context, req *entities.AdoptionCreateRequest, applicantID string) (*entities.Adoption, error) {
	// Get animal
	animal, err := uc.animalRepo.GetByID(ctx, req.AnimalID)
	if err != nil {
		return nil, err
	}

	// Check if animal is available
	if animal.Status != entities.AnimalStatusAvailable {
		return nil, errors.New("animal is not available for adoption")
	}

	// Create adoption application
	adoption := entities.NewAdoption(
		req.AnimalID,
		animal.Name,
		req.ApplicantName,
		req.ApplicantEmail,
		req.ApplicantPhone,
		req.ApplicantAddress,
		animal.AdoptionFee,
	)
	adoption.ApplicantID = applicantID
	adoption.ApplicationData = req.ApplicationData

	if err := uc.adoptionRepo.Create(ctx, adoption); err != nil {
		return nil, err
	}

	// Update animal status to reserved
	uc.animalRepo.UpdateStatus(ctx, req.AnimalID, entities.AnimalStatusReserved)

	// Audit
	auditLog := entities.NewAuditLog(applicantID, req.ApplicantEmail, "user", entities.ActionCreate, "adoption", adoption.ID.Hex(), "Adoption application submitted")
	uc.auditRepo.Create(ctx, auditLog)

	return adoption, nil
}

func (uc *AdoptionUseCase) GetByID(ctx context.Context, id string) (*entities.Adoption, error) {
	return uc.adoptionRepo.GetByID(ctx, id)
}

func (uc *AdoptionUseCase) List(ctx context.Context, filter *entities.AdoptionFilter) ([]*entities.Adoption, int64, error) {
	return uc.adoptionRepo.List(ctx, filter)
}

func (uc *AdoptionUseCase) Update(ctx context.Context, id string, req *entities.AdoptionUpdateRequest, processedBy string) (*entities.Adoption, error) {
	adoption, err := uc.adoptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != "" {
		adoption.Status = req.Status
	}
	if req.InterviewDate != nil {
		adoption.InterviewDate = req.InterviewDate
	}
	if req.RejectionReason != "" {
		adoption.RejectionReason = req.RejectionReason
	}
	if req.Notes != "" {
		adoption.Notes = req.Notes
	}
	if req.FeePaid != nil {
		adoption.FeePaid = *req.FeePaid
	}
	adoption.ProcessedBy = processedBy

	if err := uc.adoptionRepo.Update(ctx, id, adoption); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(processedBy, "", "", entities.ActionUpdate, "adoption", id, "Adoption application updated")
	uc.auditRepo.Create(ctx, auditLog)

	return adoption, nil
}

func (uc *AdoptionUseCase) Approve(ctx context.Context, id string, processedBy string) error {
	adoption, err := uc.adoptionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	adoption.Approve(processedBy)

	if err := uc.adoptionRepo.Update(ctx, id, adoption); err != nil {
		return err
	}

	// Update animal status
	uc.animalRepo.UpdateStatus(ctx, adoption.AnimalID, entities.AnimalStatusAdopted)

	// Audit
	auditLog := entities.NewAuditLog(processedBy, "", "", entities.ActionUpdate, "adoption", id, "Adoption approved")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *AdoptionUseCase) Reject(ctx context.Context, id string, reason string, processedBy string) error {
	adoption, err := uc.adoptionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	adoption.Reject(reason, processedBy)

	if err := uc.adoptionRepo.Update(ctx, id, adoption); err != nil {
		return err
	}

	// Update animal status back to available
	uc.animalRepo.UpdateStatus(ctx, adoption.AnimalID, entities.AnimalStatusAvailable)

	// Audit
	auditLog := entities.NewAuditLog(processedBy, "", "", entities.ActionUpdate, "adoption", id, "Adoption rejected")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *AdoptionUseCase) GetPendingApplications(ctx context.Context, limit, offset int) ([]*entities.Adoption, int64, error) {
	return uc.adoptionRepo.GetPendingApplications(ctx, limit, offset)
}
