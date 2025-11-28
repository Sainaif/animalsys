package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/veterinary"
	"github.com/sainaif/animalsys/backend/internal/usecase/veterinary/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVeterinaryHandler_ListVeterinaryRecords(t *testing.T) {
	mockUseCase := new(mocks.VeterinaryUseCase)
	handler := NewVeterinaryHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/records", handler.ListVeterinaryRecords)

	records := []*entities.VeterinaryRecord{
		{RecordType: "visit", Visit: &entities.VeterinaryVisit{ID: primitive.NewObjectID()}},
	}

	mockUseCase.On("ListVeterinaryRecords", mock.Anything, mock.AnythingOfType("*veterinary.ListRecordsRequest")).Return(records, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/records?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["total"])
	assert.Equal(t, float64(10), response["limit"])
	assert.Len(t, response["records"], 1)

	mockUseCase.AssertExpectations(t)
}

func TestVeterinaryHandler_CreateVeterinaryRecord(t *testing.T) {
	mockUseCase := new(mocks.VeterinaryUseCase)
	handler := NewVeterinaryHandler(mockUseCase)
	userID := primitive.NewObjectID()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/records", handler.CreateVeterinaryRecord)

	t.Run("create visit", func(t *testing.T) {
		visitReq := &veterinary.CreateVisitRequest{AnimalID: primitive.NewObjectID().Hex(), VisitDate: time.Now(), VeterinarianName: "Dr. Dolittle", VisitType: "check-up"}
		createReq := CreateRecordRequest{RecordType: "visit", Visit: visitReq}
		body, _ := json.Marshal(createReq)

		mockUseCase.On("CreateVisit", mock.Anything, mock.AnythingOfType("*veterinary.CreateVisitRequest"), userID).Return(&entities.VeterinaryVisit{}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("create vaccination", func(t *testing.T) {
		vaccinationReq := &veterinary.CreateVaccinationRequest{AnimalID: primitive.NewObjectID().Hex(), DateAdministered: time.Now(), VaccineName: "Rabies", VaccineType: "core", DoseNumber: 1, VeterinarianName: "Dr. Dolittle"}
		createReq := CreateRecordRequest{RecordType: "vaccination", Vaccination: vaccinationReq}
		body, _ := json.Marshal(createReq)

		mockUseCase.On("CreateVaccination", mock.Anything, mock.AnythingOfType("*veterinary.CreateVaccinationRequest"), userID).Return(&entities.Vaccination{}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
