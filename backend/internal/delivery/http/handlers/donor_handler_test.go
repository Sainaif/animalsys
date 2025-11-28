package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/delivery/http/dto"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockDonorUseCase struct {
	mock.Mock
}

func (m *mockDonorUseCase) GetTopDonors(ctx context.Context, limit int) ([]*dto.TopDonorResponse, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*dto.TopDonorResponse), args.Error(1)
}

func (m *mockDonorUseCase) GetRecurringDonors(ctx context.Context) ([]*dto.RecurringDonorResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*dto.RecurringDonorResponse), args.Error(1)
}

func (m *mockDonorUseCase) UpdateCommunicationPreferences(ctx context.Context, donorID primitive.ObjectID, prefs entities.DonorPreferences, userID primitive.ObjectID) error {
	args := m.Called(ctx, donorID, prefs, userID)
	return args.Error(0)
}
func (m *mockDonorUseCase) CreateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error {
	panic("implement me")
}
func (m *mockDonorUseCase) GetDonor(ctx context.Context, id primitive.ObjectID) (*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorUseCase) UpdateDonor(ctx context.Context, donor *entities.Donor, userID primitive.ObjectID) error {
	panic("implement me")
}
func (m *mockDonorUseCase) DeleteDonor(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	panic("implement me")
}
func (m *mockDonorUseCase) ListDonors(ctx context.Context, filter *repositories.DonorFilter) ([]*entities.Donor, int64, error) {
	panic("implement me")
}
func (m *mockDonorUseCase) GetMajorDonors(ctx context.Context) ([]*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorUseCase) GetLapsedDonors(ctx context.Context, days int) ([]*entities.Donor, error) {
	panic("implement me")
}
func (m *mockDonorUseCase) GetDonorStatistics(ctx context.Context) (*repositories.DonorStatistics, error) {
	panic("implement me")
}
func (m *mockDonorUseCase) UpdateDonorEngagement(ctx context.Context, id primitive.ObjectID, volunteerHours int, eventsAttended int, userID primitive.ObjectID) error {
	panic("implement me")
}

func TestDonorHandler_GetTopDonors(t *testing.T) {
	mockUseCase := new(mockDonorUseCase)
	handler := NewDonorHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/donors/top", handler.GetTopDonors)

	req, _ := http.NewRequest(http.MethodGet, "/donors/top?limit=5", nil)
	w := httptest.NewRecorder()

	mockUseCase.On("GetTopDonors", mock.Anything, 5).Return([]*dto.TopDonorResponse{}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDonorHandler_GetRecurringDonors(t *testing.T) {
	mockUseCase := new(mockDonorUseCase)
	handler := NewDonorHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/donors/recurring", handler.GetRecurringDonors)

	req, _ := http.NewRequest(http.MethodGet, "/donors/recurring", nil)
	w := httptest.NewRecorder()

	mockUseCase.On("GetRecurringDonors", mock.Anything).Return([]*dto.RecurringDonorResponse{}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDonorHandler_UpdateCommunicationPreferences(t *testing.T) {
	mockUseCase := new(mockDonorUseCase)
	handler := NewDonorHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/donors/:id/communication-preferences", func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID())
		handler.UpdateCommunicationPreferences(c)
	})

	donorID := primitive.NewObjectID()
	prefs := entities.DonorPreferences{EmailOptIn: true, SmsOptIn: false}
	body, _ := json.Marshal(prefs)

	req, _ := http.NewRequest(http.MethodPost, "/donors/"+donorID.Hex()+"/communication-preferences", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mockUseCase.On("UpdateCommunicationPreferences", mock.Anything, donorID, prefs, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
