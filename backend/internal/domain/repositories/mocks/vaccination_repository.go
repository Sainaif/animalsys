package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VaccinationRepository struct {
	mock.Mock
}

func (m *VaccinationRepository) Create(ctx context.Context, vaccination *entities.Vaccination) error {
	args := m.Called(ctx, vaccination)
	return args.Error(0)
}

func (m *VaccinationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Vaccination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Vaccination), args.Error(1)
}

func (m *VaccinationRepository) Update(ctx context.Context, vaccination *entities.Vaccination) error {
	args := m.Called(ctx, vaccination)
	return args.Error(0)
}

func (m *VaccinationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *VaccinationRepository) List(ctx context.Context, filter repositories.VaccinationFilter) ([]*entities.Vaccination, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Vaccination), args.Get(1).(int64), args.Error(2)
}

func (m *VaccinationRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Vaccination, error) {
	args := m.Called(ctx, animalID)
	return args.Get(0).([]*entities.Vaccination), args.Error(1)
}

func (m *VaccinationRepository) GetDueVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*entities.Vaccination), args.Error(1)
}

func (m *VaccinationRepository) GetByVaccineType(ctx context.Context, animalID primitive.ObjectID, vaccineType entities.VaccinationType) ([]*entities.Vaccination, error) {
	args := m.Called(ctx, animalID, vaccineType)
	return args.Get(0).([]*entities.Vaccination), args.Error(1)
}

func (m *VaccinationRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
