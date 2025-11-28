package veterinary

import (
	"context"
	"testing"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVeterinaryUseCase_ListVeterinaryRecords(t *testing.T) {
	mockVisitRepo := new(mocks.VeterinaryVisitRepository)
	mockVaccinationRepo := new(mocks.VaccinationRepository)

	useCase := NewVeterinaryUseCase(mockVisitRepo, mockVaccinationRepo, nil, nil)

	ctx := context.Background()
	animalID := primitive.NewObjectID()

	records := []*entities.VeterinaryRecord{
		{RecordType: "vaccination", Vaccination: &entities.Vaccination{ID: primitive.NewObjectID(), AnimalID: animalID, DateAdministered: time.Now()}},
		{RecordType: "visit", Visit: &entities.VeterinaryVisit{ID: primitive.NewObjectID(), AnimalID: animalID, VisitDate: time.Now().Add(-24 * time.Hour)}},
	}

	t.Run("list all records", func(t *testing.T) {
		filter := repositories.CombinedFilter{
			AnimalID:  &animalID,
			Limit:     10,
			Offset:    0,
			SortBy:    "date",
			SortOrder: "desc",
		}
		mockVisitRepo.On("ListCombined", ctx, filter).Return(records, int64(2), nil).Once()

		req := &ListRecordsRequest{
			AnimalID: animalID.Hex(),
			Limit:    10,
			Offset:   0,
		}

		result, total, err := useCase.ListVeterinaryRecords(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
		assert.Equal(t, "vaccination", result[0].RecordType)
		assert.Equal(t, "visit", result[1].RecordType)

		mockVisitRepo.AssertExpectations(t)
	})

	t.Run("list only visits with pagination", func(t *testing.T) {
		filter := repositories.VeterinaryVisitFilter{
			AnimalID:  &animalID,
			Limit:     5,
			Offset:    10,
			SortBy:    "visit_date",
			SortOrder: "desc",
		}
		visits := []*entities.VeterinaryVisit{
			{ID: primitive.NewObjectID(), AnimalID: animalID, VisitDate: time.Now().Add(-24 * time.Hour)},
		}
		mockVisitRepo.On("List", ctx, filter).Return(visits, int64(1), nil).Once()

		req := &ListRecordsRequest{
			AnimalID:   animalID.Hex(),
			RecordType: "visit",
			Limit:      5,
			Offset:     10,
		}

		result, total, err := useCase.ListVeterinaryRecords(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, "visit", result[0].RecordType)

		mockVisitRepo.AssertExpectations(t)
	})
}
