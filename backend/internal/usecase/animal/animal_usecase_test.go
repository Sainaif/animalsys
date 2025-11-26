package animal

import (
	"context"
	"testing"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateAnimal_AuditLog(t *testing.T) {
	// Setup
	animalRepo := new(mocks.AnimalRepository)
	auditLogRepo := new(mocks.AuditLogRepository)
	uc := NewAnimalUseCase(animalRepo, auditLogRepo, nil)

	animalID := primitive.NewObjectID()
	updaterID := primitive.NewObjectID()
	now := time.Now()

	// Mock existing animal
	existingAnimal := &entities.Animal{
		ID:     animalID,
		Status: entities.AnimalStatusAvailable,
	}

	// Mock request with all fields
	req := &UpdateAnimalRequest{
		Name:         &entities.MultilingualName{English: "New Name"},
		Category:     &[]entities.AnimalCategory{entities.CategoryMammal}[0],
		Species:      &[]string{"New Species"}[0],
		Breed:        &[]string{"New Breed"}[0],
		Sex:          &[]entities.AnimalSex{entities.SexFemale}[0],
		Status:       &[]entities.AnimalStatus{entities.AnimalStatusAdopted}[0],
		DateOfBirth:  &now,
		AgeEstimated: &[]bool{true}[0],
		Color:        &[]string{"New Color"}[0],
		Size:         &[]entities.AnimalSize{entities.SizeLarge}[0],
		Weight:       &[]float64{25.5}[0],
		Description:  &entities.MultilingualName{English: "New Description"},
		Medical:      &entities.MedicalInfo{Vaccinated: true},
		Behavior:     &entities.BehaviorInfo{GoodWithKids: true},
		Location:     &[]string{"New Location"}[0],
		AdoptionFee:  &[]float64{200.0}[0],
		Requirements: &[]string{"New Requirements"},
	}

	animalRepo.On("FindByID", mock.Anything, animalID).Return(existingAnimal, nil)
	animalRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	auditLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Execute
	_, err := uc.UpdateAnimal(context.Background(), animalID, req, updaterID)

	// Assert
	assert.NoError(t, err)

	// Verify that the audit log was created with the correct changes
	auditLogRepo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(log *entities.AuditLog) bool {
		assert.Equal(t, updaterID, log.UserID)
		assert.Equal(t, entities.ActionUpdate, log.Action)
		assert.Equal(t, "animal", log.EntityType)
		assert.Equal(t, animalID, *log.EntityID)

		// Check that all updated fields are in the changes map
		changes := log.Changes
		assert.Contains(t, changes, "name")
		assert.Contains(t, changes, "category")
		assert.Contains(t, changes, "species")
		assert.Contains(t, changes, "breed")
		assert.Contains(t, changes, "sex")
		assert.Contains(t, changes, "status")
		assert.Contains(t, changes, "date_of_birth")
		assert.Contains(t, changes, "age_estimated")
		assert.Contains(t, changes, "color")
		assert.Contains(t, changes, "size")
		assert.Contains(t, changes, "weight")
		assert.Contains(t, changes, "description")
		assert.Contains(t, changes, "medical")
		assert.Contains(t, changes, "behavior")
		assert.Contains(t, changes, "location")
		assert.Contains(t, changes, "adoption_fee")
		assert.Contains(t, changes, "requirements")

		return true
	}))
}
