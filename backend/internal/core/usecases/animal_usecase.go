package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type AnimalUseCase struct {
	animalRepo interfaces.AnimalRepository
	auditRepo  interfaces.AuditLogRepository
}

func NewAnimalUseCase(animalRepo interfaces.AnimalRepository, auditRepo interfaces.AuditLogRepository) *AnimalUseCase {
	return &AnimalUseCase{
		animalRepo: animalRepo,
		auditRepo:  auditRepo,
	}
}

func (uc *AnimalUseCase) Create(ctx context.Context, req *entities.AnimalCreateRequest, createdBy string) (*entities.Animal, error) {
	animal := entities.NewAnimal(req.Name, req.Species, req.Gender)

	animal.Breed = req.Breed
	animal.DateOfBirth = req.DateOfBirth
	animal.AgeYears = req.AgeYears
	animal.AgeMonths = req.AgeMonths
	animal.Color = req.Color
	animal.Size = req.Size
	animal.Weight = req.Weight
	animal.IntakeReason = req.IntakeReason
	animal.Description = req.Description
	animal.Characteristics = req.Characteristics
	animal.SpecialNeeds = req.SpecialNeeds
	animal.Microchipped = req.Microchipped
	animal.MicrochipNumber = req.MicrochipNumber
	animal.Neutered = req.Neutered
	animal.Vaccinated = req.Vaccinated
	animal.AdoptionFee = req.AdoptionFee
	animal.CreatedBy = createdBy

	if err := uc.animalRepo.Create(ctx, animal); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "animal", animal.ID.Hex(), "Animal created")
	uc.auditRepo.Create(ctx, auditLog)

	return animal, nil
}

func (uc *AnimalUseCase) GetByID(ctx context.Context, id string) (*entities.Animal, error) {
	return uc.animalRepo.GetByID(ctx, id)
}

func (uc *AnimalUseCase) List(ctx context.Context, filter *entities.AnimalFilter) ([]*entities.Animal, int64, error) {
	return uc.animalRepo.List(ctx, filter)
}

func (uc *AnimalUseCase) Update(ctx context.Context, id string, req *entities.AnimalUpdateRequest, updatedBy string) (*entities.Animal, error) {
	animal, err := uc.animalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		animal.Name = req.Name
	}
	if req.Breed != "" {
		animal.Breed = req.Breed
	}
	if req.Color != "" {
		animal.Color = req.Color
	}
	if req.Size != "" {
		animal.Size = req.Size
	}
	if req.Weight > 0 {
		animal.Weight = req.Weight
	}
	if req.Status != "" {
		animal.Status = req.Status
	}
	if req.Description != "" {
		animal.Description = req.Description
	}
	if req.Characteristics != nil {
		animal.Characteristics = req.Characteristics
	}
	if req.SpecialNeeds != "" {
		animal.SpecialNeeds = req.SpecialNeeds
	}
	if req.Microchipped != nil {
		animal.Microchipped = *req.Microchipped
	}
	if req.MicrochipNumber != "" {
		animal.MicrochipNumber = req.MicrochipNumber
	}
	if req.Neutered != nil {
		animal.Neutered = *req.Neutered
	}
	if req.Vaccinated != nil {
		animal.Vaccinated = *req.Vaccinated
	}
	if req.AdoptionFee > 0 {
		animal.AdoptionFee = req.AdoptionFee
	}
	animal.UpdatedBy = updatedBy

	if err := uc.animalRepo.Update(ctx, id, animal); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "animal", id, "Animal updated")
	uc.auditRepo.Create(ctx, auditLog)

	return animal, nil
}

func (uc *AnimalUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.animalRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "animal", id, "Animal deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *AnimalUseCase) AddMedicalRecord(ctx context.Context, id string, record entities.MedicalRecord, addedBy string) error {
	if err := uc.animalRepo.AddMedicalRecord(ctx, id, record); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(addedBy, "", "", entities.ActionUpdate, "animal", id, "Medical record added")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *AnimalUseCase) AddPhoto(ctx context.Context, id string, photoURL string, addedBy string) error {
	if err := uc.animalRepo.AddPhoto(ctx, id, photoURL); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(addedBy, "", "", entities.ActionUpdate, "animal", id, "Photo added")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *AnimalUseCase) GetAvailableForAdoption(ctx context.Context, limit, offset int) ([]*entities.Animal, int64, error) {
	return uc.animalRepo.GetAvailableForAdoption(ctx, limit, offset)
}
