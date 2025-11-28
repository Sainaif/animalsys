package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VeterinaryVisitRepository struct {
	mock.Mock
}

func (m *VeterinaryVisitRepository) Create(ctx context.Context, visit *entities.VeterinaryVisit) error {
	args := m.Called(ctx, visit)
	return args.Error(0)
}

func (m *VeterinaryVisitRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryVisitRepository) Update(ctx context.Context, visit *entities.VeterinaryVisit) error {
	args := m.Called(ctx, visit)
	return args.Error(0)
}

func (m *VeterinaryVisitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *VeterinaryVisitRepository) List(ctx context.Context, filter repositories.VeterinaryVisitFilter) ([]*entities.VeterinaryVisit, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Get(1).(int64), args.Error(2)
}

func (m *VeterinaryVisitRepository) GetByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, animalID)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryVisitRepository) GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*entities.VeterinaryVisit), args.Error(1)
}

func (m *VeterinaryVisitRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *VeterinaryVisitRepository) ListCombined(ctx context.Context, filter repositories.CombinedFilter) ([]*entities.VeterinaryRecord, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.VeterinaryRecord), args.Get(1).(int64), args.Error(2)
}
