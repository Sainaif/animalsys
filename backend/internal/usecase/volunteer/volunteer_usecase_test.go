package volunteer

import (
	"context"
	"testing"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVolunteerUseCase_LogHours(t *testing.T) {
	mockVolunteerRepo := new(mocks.VolunteerRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewVolunteerUseCase(mockVolunteerRepo, nil, mockAuditLogRepo)

	volunteerID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	hours := 5.5
	notes := "Test notes"

	mockVolunteerRepo.On("UpdateHours", mock.Anything, volunteerID, hours, notes).Return(nil).Once()
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()
	mockVolunteerRepo.On("FindByID", mock.Anything, volunteerID).Return(&entities.Volunteer{ID: volunteerID}, nil).Once()

	volunteer, err := uc.LogHours(context.Background(), volunteerID, hours, notes, userID)

	assert.NoError(t, err)
	assert.NotNil(t, volunteer)
	mockVolunteerRepo.AssertExpectations(t)
	mockAuditLogRepo.AssertExpectations(t)
}

func TestVolunteerUseCase_GetVolunteerHours(t *testing.T) {
	mockVolunteerRepo := new(mocks.VolunteerRepository)
	uc := NewVolunteerUseCase(mockVolunteerRepo, nil, nil)

	volunteerID := primitive.NewObjectID()
	expectedHours := 10.5

	mockVolunteerRepo.On("FindByID", mock.Anything, volunteerID).Return(&entities.Volunteer{ID: volunteerID, TotalHours: expectedHours}, nil).Once()

	hours, err := uc.GetVolunteerHours(context.Background(), volunteerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedHours, hours)
	mockVolunteerRepo.AssertExpectations(t)
}

func TestVolunteerUseCase_ActivateVolunteer(t *testing.T) {
	mockVolunteerRepo := new(mocks.VolunteerRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewVolunteerUseCase(mockVolunteerRepo, nil, mockAuditLogRepo)

	volunteerID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	mockVolunteerRepo.On("FindByID", mock.Anything, volunteerID).Return(&entities.Volunteer{ID: volunteerID}, nil).Once()
	mockVolunteerRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Volunteer")).Return(nil).Once()
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

	err := uc.ActivateVolunteer(context.Background(), volunteerID, userID)

	assert.NoError(t, err)
	mockVolunteerRepo.AssertExpectations(t)
	mockAuditLogRepo.AssertExpectations(t)
}

func TestVolunteerUseCase_DeactivateVolunteer(t *testing.T) {
	mockVolunteerRepo := new(mocks.VolunteerRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	uc := NewVolunteerUseCase(mockVolunteerRepo, nil, mockAuditLogRepo)

	volunteerID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	mockVolunteerRepo.On("FindByID", mock.Anything, volunteerID).Return(&entities.Volunteer{ID: volunteerID}, nil).Once()
	mockVolunteerRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Volunteer")).Return(nil).Once()
	mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

	err := uc.DeactivateVolunteer(context.Background(), volunteerID, userID)

	assert.NoError(t, err)
	mockVolunteerRepo.AssertExpectations(t)
	mockAuditLogRepo.AssertExpectations(t)
}

func TestVolunteerUseCase_GetVolunteersByRole(t *testing.T) {
	mockVolunteerRepo := new(mocks.VolunteerRepository)
	uc := NewVolunteerUseCase(mockVolunteerRepo, nil, nil)

	role := entities.VolunteerRoleAnimalCare
	expectedVolunteers := []*entities.Volunteer{{ID: primitive.NewObjectID()}}
	expectedTotal := int64(1)

	mockVolunteerRepo.On("List", mock.Anything, mock.AnythingOfType("*repositories.VolunteerFilter")).Return(expectedVolunteers, expectedTotal, nil).Once()

	volunteers, total, err := uc.GetVolunteersByRole(context.Background(), role)

	assert.NoError(t, err)
	assert.Equal(t, expectedVolunteers, volunteers)
	assert.Equal(t, expectedTotal, total)
	mockVolunteerRepo.AssertExpectations(t)
}
