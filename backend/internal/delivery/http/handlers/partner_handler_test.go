package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/delivery/http/handlers"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories/mocks"
	"github.com/sainaif/animalsys/backend/internal/usecase/partner"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() (*gin.Engine, *mocks.PartnerRepository, *mocks.AuditLogRepository, *handlers.PartnerHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockPartnerRepo := new(mocks.PartnerRepository)
	mockAuditLogRepo := new(mocks.AuditLogRepository)
	partnerUseCase := partner.NewPartnerUseCase(mockPartnerRepo, mockAuditLogRepo)
	partnerHandler := handlers.NewPartnerHandler(partnerUseCase)

	// Middleware to set user_id for testing
	router.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID())
		c.Next()
	})

	return router, mockPartnerRepo, mockAuditLogRepo, partnerHandler
}
func TestPartnerHandler_GetPartnersAcceptingIntakes(t *testing.T) {
	router, mockPartnerRepo, _, partnerHandler := setupRouter()

	router.GET("/partners/accepting-intakes", partnerHandler.GetPartnersAcceptingIntakes)

	t.Run("Success - No capacity filter", func(t *testing.T) {
		partners := []*entities.Partner{
			{ID: primitive.NewObjectID(), Name: "Partner 1", AcceptsIntakes: true},
		}
		mockPartnerRepo.On("List", mock.Anything, mock.MatchedBy(func(filter *repositories.PartnerFilter) bool {
			return *filter.AcceptsIntakes == true && filter.HasCapacity == nil
		})).Return(partners, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/partners/accepting-intakes", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, float64(1), response["total"])
		data := response["data"].([]interface{})
		assert.Len(t, data, 1)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("Success - With capacity filter", func(t *testing.T) {
		partners := []*entities.Partner{
			{ID: primitive.NewObjectID(), Name: "Partner 2", AcceptsIntakes: true, CurrentCapacity: 5, MaxCapacity: 10},
		}
		mockPartnerRepo.On("List", mock.Anything, mock.MatchedBy(func(filter *repositories.PartnerFilter) bool {
			return *filter.AcceptsIntakes == true && *filter.HasCapacity == true
		})).Return(partners, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/partners/accepting-intakes?has_capacity=true", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, float64(1), response["total"])
		data := response["data"].([]interface{})
		assert.Len(t, data, 1)
		mockPartnerRepo.AssertExpectations(t)
	})
}

func TestPartnerHandler_RatePartner(t *testing.T) {
	router, mockPartnerRepo, mockAuditLogRepo, partnerHandler := setupRouter()
	testPartnerID := primitive.NewObjectID()
	router.POST("/partners/:id/rate", partnerHandler.RatePartner)

	t.Run("Success", func(t *testing.T) {
		mockPartner := &entities.Partner{ID: testPartnerID, Name: "Test Partner"}
		mockPartnerRepo.On("FindByID", mock.Anything, testPartnerID).Return(mockPartner, nil).Once()
		mockPartnerRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Partner")).Return(nil).Once()
		mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

		rating := 4.5
		body, _ := json.Marshal(gin.H{"rating": rating})
		req, _ := http.NewRequest(http.MethodPost, "/partners/"+testPartnerID.Hex()+"/rate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "Rating added successfully", response["message"])
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("Invalid Partner ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/partners/invalid-id/rate", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Validation Error - Rating too high", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"rating": 6.0})
		req, _ := http.NewRequest(http.MethodPost, "/partners/"+testPartnerID.Hex()+"/rate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
func TestPartnerHandler_GetPartnerStatisticsDetail(t *testing.T) {
	router, mockPartnerRepo, _, partnerHandler := setupRouter()
	testPartnerID := primitive.NewObjectID()
	router.GET("/partners/:id/statistics", partnerHandler.GetPartnerStatisticsDetail)

	t.Run("Success", func(t *testing.T) {
		mockPartner := &entities.Partner{
			ID:                 testPartnerID,
			Name:               "Stats Partner",
			Status:             entities.PartnerStatusActive,
			AcceptsIntakes:     true,
			CurrentCapacity:    10,
			MaxCapacity:        20,
			Rating:             4.2,
			TotalRatings:       50,
			TotalTransfersIn:   15,
			TotalTransfersOut:  5,
			SuccessfulPlacements: 12,
		}
		mockPartnerRepo.On("FindByID", mock.Anything, testPartnerID).Return(mockPartner, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/partners/"+testPartnerID.Hex()+"/statistics", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, testPartnerID.Hex(), response["partner_id"])
		assert.Equal(t, float64(10), response["current_capacity"])
		assert.Equal(t, float64(4.2), response["average_rating"])
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("Partner Not Found", func(t *testing.T) {
		mockPartnerRepo.On("FindByID", mock.Anything, testPartnerID).Return(nil, errors.NewNotFound("partner not found")).Once()
		req, _ := http.NewRequest(http.MethodGet, "/partners/"+testPartnerID.Hex()+"/statistics", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockPartnerRepo.AssertExpectations(t)
	})
}

func TestPartnerHandler_UpdatePartnerCapacity(t *testing.T) {
	router, mockPartnerRepo, mockAuditLogRepo, partnerHandler := setupRouter()
	testPartnerID := primitive.NewObjectID()
	router.POST("/partners/:id/update-capacity", partnerHandler.UpdatePartnerCapacity)

	t.Run("Success", func(t *testing.T) {
		mockPartner := &entities.Partner{ID: testPartnerID, Name: "Capacity Partner", MaxCapacity: 20}
		mockPartnerRepo.On("FindByID", mock.Anything, testPartnerID).Return(mockPartner, nil).Once()
		mockPartnerRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *entities.Partner) bool {
			return p.CurrentCapacity == 15
		})).Return(nil).Once()
		mockAuditLogRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.AuditLog")).Return(nil).Once()

		body, _ := json.Marshal(gin.H{"current_capacity": 15})
		req, _ := http.NewRequest(http.MethodPost, "/partners/"+testPartnerID.Hex()+"/update-capacity", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "Capacity updated successfully", response["message"])
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Exceeds Max Capacity", func(t *testing.T) {
		mockPartner := &entities.Partner{ID: testPartnerID, Name: "Capacity Partner", MaxCapacity: 20}
		mockPartnerRepo.On("FindByID", mock.Anything, testPartnerID).Return(mockPartner, nil).Once()

		body, _ := json.Marshal(gin.H{"current_capacity": 25})
		req, _ := http.NewRequest(http.MethodPost, "/partners/"+testPartnerID.Hex()+"/update-capacity", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Negative Capacity", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"current_capacity": -5})
		req, _ := http.NewRequest(http.MethodPost, "/partners/"+testPartnerID.Hex()+"/update-capacity", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
