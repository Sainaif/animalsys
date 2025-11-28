package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/volunteer/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVolunteerHandler_LogVolunteerHours(t *testing.T) {
	mockUseCase := new(mocks.VolunteerUseCase)
	handler := NewVolunteerHandler(mockUseCase)
	userID := primitive.NewObjectID()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/volunteers/:id/log-hours", handler.LogVolunteerHours)

	volunteerID := primitive.NewObjectID()
	logReq := gin.H{
		"hours": 5.5,
		"notes": "Helped with the adoption event",
	}
	body, _ := json.Marshal(logReq)

	mockUseCase.On("LogHours", mock.Anything, volunteerID, 5.5, "Helped with the adoption event", userID).Return(&entities.Volunteer{ID: volunteerID, TotalHours: 5.5}, nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/volunteers/"+volunteerID.Hex()+"/log-hours", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestVolunteerHandler_GetVolunteerHours(t *testing.T) {
	mockUseCase := new(mocks.VolunteerUseCase)
	handler := NewVolunteerHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/volunteers/:id/hours", handler.GetVolunteerHours)

	volunteerID := primitive.NewObjectID()

	mockUseCase.On("GetVolunteerHours", mock.Anything, volunteerID).Return(10.5, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/volunteers/"+volunteerID.Hex()+"/hours", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10.5, response["total_hours"])

	mockUseCase.AssertExpectations(t)
}

func TestVolunteerHandler_ActivateVolunteer(t *testing.T) {
	mockUseCase := new(mocks.VolunteerUseCase)
	handler := NewVolunteerHandler(mockUseCase)
	userID := primitive.NewObjectID()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/volunteers/:id/activate", handler.ActivateVolunteer)

	volunteerID := primitive.NewObjectID()

	mockUseCase.On("ActivateVolunteer", mock.Anything, volunteerID, userID).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/volunteers/"+volunteerID.Hex()+"/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestVolunteerHandler_DeactivateVolunteer(t *testing.T) {
	mockUseCase := new(mocks.VolunteerUseCase)
	handler := NewVolunteerHandler(mockUseCase)
	userID := primitive.NewObjectID()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/volunteers/:id/deactivate", handler.DeactivateVolunteer)

	volunteerID := primitive.NewObjectID()

	mockUseCase.On("DeactivateVolunteer", mock.Anything, volunteerID, userID).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodPost, "/volunteers/"+volunteerID.Hex()+"/deactivate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestVolunteerHandler_GetVolunteersByRole(t *testing.T) {
	mockUseCase := new(mocks.VolunteerUseCase)
	handler := NewVolunteerHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/volunteers/role/:role", handler.GetVolunteersByRole)

	role := entities.VolunteerRoleAnimalCare
	volunteers := []*entities.Volunteer{
		{ID: primitive.NewObjectID(), Roles: []entities.VolunteerRole{role}},
	}

	mockUseCase.On("GetVolunteersByRole", mock.Anything, role).Return(volunteers, int64(1), nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/volunteers/role/"+string(role), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["total"])
	assert.Len(t, response["volunteers"], 1)

	mockUseCase.AssertExpectations(t)
}
