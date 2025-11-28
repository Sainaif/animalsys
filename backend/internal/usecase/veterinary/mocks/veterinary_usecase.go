package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/veterinary"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VeterinaryUseCase struct {
	mock.Mock
}

func (m *VeterinaryUseCase) CreateVisit(ctx context.Context, req *veterinary.CreateVisitRequest, creatorID primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, req, creatorID)
	return args.Get(0).(*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryUseCase) GetVisitByID(ctx context.Context, id primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryUseCase) UpdateVisit(ctx context.Context, id primitive.ObjectID, req *veterinary.UpdateVisitRequest, updaterID primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, id, req, updaterID)
	return args.Get(0).(*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryUseCase) DeleteVisit(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	args := m.Called(ctx, id, deleterID)
	return args.Error(0)
}

func (m *VeterinaryUseCase) ListVisits(ctx context.Context, req *veterinary.ListVisitsRequest) ([]*entities.VeterinaryVisit, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Get(1).(int64), args.Error(2)
}

func (m *VeterinaryUseCase) GetVisitsByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, animalID)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryUseCase) CreateVaccination(ctx context.Context, req *veterinary.CreateVaccinationRequest, creatorID primitive.ObjectID) (*entities.Vaccination, error) {
	args := m.Called(ctx, req, creatorID)
	return args.Get(0).(*entities.Vaccination), args.Error(1)
}

func (m *VeterinaryUseCase) GetVaccinationByID(ctx context.Context, id primitive.ObjectID) (*entities.Vaccination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Vaccination), args.Error(1)
}

func (m *VeterinaryUseCase) DeleteVaccination(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	args := m.Called(ctx, id, deleterID)
	return args.Error(0)
}

func (m *VeterinaryUseCase) ListVaccinations(ctx context.Context, req *veterinary.ListVaccinationsRequest) ([]*entities.Vaccination, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*entities.Vaccination), args.Get(1).(int64), args.Error(2)
}

func (m *VeterinaryUseCase) GetVaccinationsByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Vaccination, error) {
	args := m.Called(ctx, animalID)
	return args.Get(0).([]*entities.Vaccination), args.Error(1)
}

func (m *VeterinaryUseCase) GetDueVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*entities.Vaccination), args.Error(1)
}

func (m *VeterinaryUseCase) GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryUseCase) ListVeterinaryRecords(ctx context.Context, req *veterinary.ListRecordsRequest) ([]*entities.VeterinaryRecord, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*entities.VeterinaryRecord), args.Get(1).(int64), args.Error(2)
}
