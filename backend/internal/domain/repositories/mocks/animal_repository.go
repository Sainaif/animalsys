package mocks

import (
	"context"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnimalRepository struct {
	mock.Mock
}

func (m *AnimalRepository) Create(ctx context.Context, animal *entities.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *AnimalRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Animal, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Animal), args.Error(1)
}

func (m *AnimalRepository) Update(ctx context.Context, animal *entities.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *AnimalRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AnimalRepository) List(ctx context.Context, filter repositories.AnimalFilter) ([]*entities.Animal, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entities.Animal), args.Get(1).(int64), args.Error(2)
}

func (m *AnimalRepository) AddDailyNote(ctx context.Context, animalID primitive.ObjectID, note entities.DailyNote) error {
	args := m.Called(ctx, animalID, note)
	return args.Error(0)
}

func (m *AnimalRepository) UpdateImages(ctx context.Context, animalID primitive.ObjectID, images entities.AnimalImages) error {
	args := m.Called(ctx, animalID, images)
	return args.Error(0)
}

func (m *AnimalRepository) UpdateStatus(ctx context.Context, animalID primitive.ObjectID, status entities.AnimalStatus) error {
	args := m.Called(ctx, animalID, status)
	return args.Error(0)
}

func (m *AnimalRepository) GetStatistics(ctx context.Context) (*repositories.AnimalStatistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.AnimalStatistics), args.Error(1)
}

func (m *AnimalRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
