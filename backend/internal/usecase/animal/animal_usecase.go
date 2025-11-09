package animal

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/sainaif/animalsys/backend/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnimalUseCase handles animal business logic
type AnimalUseCase struct {
	animalRepo   repositories.AnimalRepository
	auditLogRepo repositories.AuditLogRepository
	storageService *storage.StorageService
}

// NewAnimalUseCase creates a new animal use case
func NewAnimalUseCase(
	animalRepo repositories.AnimalRepository,
	auditLogRepo repositories.AuditLogRepository,
	storageService *storage.StorageService,
) *AnimalUseCase {
	return &AnimalUseCase{
		animalRepo:     animalRepo,
		auditLogRepo:   auditLogRepo,
		storageService: storageService,
	}
}

// CreateAnimalRequest represents a request to create an animal
type CreateAnimalRequest struct {
	Name           entities.MultilingualName  `json:"name" validate:"required"`
	Category       entities.AnimalCategory    `json:"category" validate:"required"`
	Species        string                     `json:"species" validate:"required"`
	Breed          string                     `json:"breed,omitempty"`
	Sex            entities.AnimalSex         `json:"sex" validate:"required"`
	Status         entities.AnimalStatus      `json:"status"`
	DateOfBirth    *time.Time                 `json:"date_of_birth,omitempty"`
	AgeEstimated   bool                       `json:"age_estimated"`
	Color          string                     `json:"color,omitempty"`
	Size           entities.AnimalSize        `json:"size,omitempty"`
	Weight         float64                    `json:"weight,omitempty"`
	Description    entities.MultilingualName  `json:"description"`
	Medical        entities.MedicalInfo       `json:"medical"`
	Behavior       entities.BehaviorInfo      `json:"behavior"`
	IntakeDate     time.Time                  `json:"intake_date" validate:"required"`
	IntakeReason   string                     `json:"intake_reason,omitempty"`
	Location       string                     `json:"location" validate:"required"`
	AdoptionFee    float64                    `json:"adoption_fee"`
	Requirements   []string                   `json:"requirements,omitempty"`
}

// UpdateAnimalRequest represents a request to update an animal
type UpdateAnimalRequest struct {
	Name           *entities.MultilingualName `json:"name,omitempty"`
	Category       *entities.AnimalCategory   `json:"category,omitempty"`
	Species        *string                    `json:"species,omitempty"`
	Breed          *string                    `json:"breed,omitempty"`
	Sex            *entities.AnimalSex        `json:"sex,omitempty"`
	Status         *entities.AnimalStatus     `json:"status,omitempty"`
	DateOfBirth    *time.Time                 `json:"date_of_birth,omitempty"`
	AgeEstimated   *bool                      `json:"age_estimated,omitempty"`
	Color          *string                    `json:"color,omitempty"`
	Size           *entities.AnimalSize       `json:"size,omitempty"`
	Weight         *float64                   `json:"weight,omitempty"`
	Description    *entities.MultilingualName `json:"description,omitempty"`
	Medical        *entities.MedicalInfo      `json:"medical,omitempty"`
	Behavior       *entities.BehaviorInfo     `json:"behavior,omitempty"`
	Location       *string                    `json:"location,omitempty"`
	AdoptionFee    *float64                   `json:"adoption_fee,omitempty"`
	Requirements   *[]string                  `json:"requirements,omitempty"`
}

// ListAnimalsRequest represents a request to list animals
type ListAnimalsRequest struct {
	Category         string   `form:"category"`
	Species          string   `form:"species"`
	Status           string   `form:"status"`
	Sex              string   `form:"sex"`
	Size             string   `form:"size"`
	AvailableOnly    bool     `form:"available_only"`
	GoodWithKids     *bool    `form:"good_with_kids"`
	GoodWithDogs     *bool    `form:"good_with_dogs"`
	GoodWithCats     *bool    `form:"good_with_cats"`
	Search           string   `form:"search"`
	AssignedCaretaker string  `form:"assigned_caretaker"`
	MinAge           *float64 `form:"min_age"`
	MaxAge           *float64 `form:"max_age"`
	Limit            int64    `form:"limit" validate:"min=0,max=100"`
	Offset           int64    `form:"offset" validate:"min=0"`
	SortBy           string   `form:"sort_by"`
	SortOrder        string   `form:"sort_order"`
}

// ListAnimalsResponse represents the response for listing animals
type ListAnimalsResponse struct {
	Animals []*entities.Animal `json:"animals"`
	Total   int64              `json:"total"`
	Limit   int64              `json:"limit"`
	Offset  int64              `json:"offset"`
}

// CreateAnimal creates a new animal
func (uc *AnimalUseCase) CreateAnimal(ctx context.Context, req *CreateAnimalRequest, creatorID primitive.ObjectID) (*entities.Animal, error) {
	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = entities.AnimalStatusAvailable
	}

	animal := &entities.Animal{
		Name:         req.Name,
		Category:     req.Category,
		Species:      req.Species,
		Breed:        req.Breed,
		Sex:          req.Sex,
		Status:       status,
		DateOfBirth:  req.DateOfBirth,
		AgeEstimated: req.AgeEstimated,
		Color:        req.Color,
		Size:         req.Size,
		Weight:       req.Weight,
		Description:  req.Description,
		Medical:      req.Medical,
		Behavior:     req.Behavior,
		Shelter: entities.ShelterInfo{
			IntakeDate:   req.IntakeDate,
			IntakeReason: req.IntakeReason,
			Location:     req.Location,
		},
		Adoption: entities.AdoptionInfo{
			AdoptionFee:  req.AdoptionFee,
			Requirements: req.Requirements,
		},
		CreatedBy: creatorID,
		UpdatedBy: creatorID,
	}

	if err := uc.animalRepo.Create(ctx, animal); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "animal", "", "").
		WithEntityID(animal.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return animal, nil
}

// GetAnimalByID retrieves an animal by ID
func (uc *AnimalUseCase) GetAnimalByID(ctx context.Context, id primitive.ObjectID) (*entities.Animal, error) {
	return uc.animalRepo.FindByID(ctx, id)
}

// UpdateAnimal updates an animal
func (uc *AnimalUseCase) UpdateAnimal(ctx context.Context, id primitive.ObjectID, req *UpdateAnimalRequest, updaterID primitive.ObjectID) (*entities.Animal, error) {
	// Get existing animal
	animal, err := uc.animalRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if can be modified
	if !animal.CanBeModified() {
		return nil, errors.NewBadRequest("animal cannot be modified")
	}

	// Track changes for audit log
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.Name != nil {
		changes["name"] = req.Name
		animal.Name = *req.Name
	}
	if req.Category != nil {
		changes["category"] = req.Category
		animal.Category = *req.Category
	}
	if req.Species != nil {
		changes["species"] = req.Species
		animal.Species = *req.Species
	}
	if req.Breed != nil {
		changes["breed"] = req.Breed
		animal.Breed = *req.Breed
	}
	if req.Sex != nil {
		changes["sex"] = req.Sex
		animal.Sex = *req.Sex
	}
	if req.Status != nil {
		changes["status"] = req.Status
		animal.Status = *req.Status
	}
	if req.DateOfBirth != nil {
		changes["date_of_birth"] = req.DateOfBirth
		animal.DateOfBirth = req.DateOfBirth
	}
	if req.AgeEstimated != nil {
		animal.AgeEstimated = *req.AgeEstimated
	}
	if req.Color != nil {
		animal.Color = *req.Color
	}
	if req.Size != nil {
		animal.Size = *req.Size
	}
	if req.Weight != nil {
		animal.Weight = *req.Weight
	}
	if req.Description != nil {
		animal.Description = *req.Description
	}
	if req.Medical != nil {
		animal.Medical = *req.Medical
	}
	if req.Behavior != nil {
		animal.Behavior = *req.Behavior
	}
	if req.Location != nil {
		animal.Shelter.Location = *req.Location
	}
	if req.AdoptionFee != nil {
		animal.Adoption.AdoptionFee = *req.AdoptionFee
	}
	if req.Requirements != nil {
		animal.Adoption.Requirements = *req.Requirements
	}

	animal.UpdatedBy = updaterID

	if err := uc.animalRepo.Update(ctx, animal); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(updaterID, entities.ActionUpdate, "animal", "", "").
		WithEntityID(id).
		WithChanges(changes)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return animal, nil
}

// DeleteAnimal deletes an animal
func (uc *AnimalUseCase) DeleteAnimal(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	// Check if animal exists
	animal, err := uc.animalRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete images if any
	if animal.Images.Primary != "" {
		_ = uc.storageService.DeleteFile(ctx, animal.Images.Primary)
	}
	if len(animal.Images.Gallery) > 0 {
		_ = uc.storageService.DeleteMultipleFiles(ctx, animal.Images.Gallery)
	}

	if err := uc.animalRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "animal", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListAnimals lists animals with filters and pagination
func (uc *AnimalUseCase) ListAnimals(ctx context.Context, req *ListAnimalsRequest) (*ListAnimalsResponse, error) {
	// Set default pagination
	if req.Limit == 0 {
		req.Limit = 20
	}

	// Build filter
	filter := repositories.AnimalFilter{
		Category:         req.Category,
		Species:          req.Species,
		Status:           req.Status,
		Sex:              req.Sex,
		Size:             req.Size,
		AvailableOnly:    req.AvailableOnly,
		GoodWithKids:     req.GoodWithKids,
		GoodWithDogs:     req.GoodWithDogs,
		GoodWithCats:     req.GoodWithCats,
		Search:           req.Search,
		MinAge:           req.MinAge,
		MaxAge:           req.MaxAge,
		Limit:            req.Limit,
		Offset:           req.Offset,
		SortBy:           req.SortBy,
		SortOrder:        req.SortOrder,
	}

	// Parse assigned caretaker if provided
	if req.AssignedCaretaker != "" {
		caretakerID, err := primitive.ObjectIDFromHex(req.AssignedCaretaker)
		if err == nil {
			filter.AssignedCaretaker = &caretakerID
		}
	}

	animals, total, err := uc.animalRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListAnimalsResponse{
		Animals: animals,
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
	}, nil
}

// UploadAnimalImages uploads images for an animal
func (uc *AnimalUseCase) UploadAnimalImages(ctx context.Context, animalID primitive.ObjectID, primary *multipart.FileHeader, gallery []*multipart.FileHeader, userID primitive.ObjectID) error {
	animal, err := uc.animalRepo.FindByID(ctx, animalID)
	if err != nil {
		return err
	}

	images := animal.Images

	// Upload primary image
	if primary != nil {
		// Delete old primary image if exists
		if images.Primary != "" {
			_ = uc.storageService.DeleteFile(ctx, images.Primary)
		}

		url, err := uc.storageService.UploadImage(ctx, primary, "animals")
		if err != nil {
			return err
		}
		images.Primary = url
	}

	// Upload gallery images
	if len(gallery) > 0 {
		urls, err := uc.storageService.UploadMultipleImages(ctx, gallery, "animals")
		if err != nil {
			return err
		}
		images.Gallery = append(images.Gallery, urls...)
	}

	if err := uc.animalRepo.UpdateImages(ctx, animalID, images); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "animal", "", "").
		WithEntityID(animalID).
		WithChanges(map[string]interface{}{"images": "updated"})
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// AddDailyNote adds a daily note to an animal
func (uc *AnimalUseCase) AddDailyNote(ctx context.Context, animalID primitive.ObjectID, noteText string, userID primitive.ObjectID) error {
	note := entities.DailyNote{
		Date:      time.Now(),
		Note:      noteText,
		CreatedBy: userID,
	}

	if err := uc.animalRepo.AddDailyNote(ctx, animalID, note); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "animal", "", "").
		WithEntityID(animalID).
		WithChanges(map[string]interface{}{"daily_note": "added"})
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetStatistics retrieves animal statistics
func (uc *AnimalUseCase) GetStatistics(ctx context.Context) (*repositories.AnimalStatistics, error) {
	return uc.animalRepo.GetStatistics(ctx)
}

// GetSpeciesByCategory retrieves species for a category
func (uc *AnimalUseCase) GetSpeciesByCategory(category entities.AnimalCategory) []entities.SpeciesInfo {
	return entities.GetSpeciesByCategory(category)
}
