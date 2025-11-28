package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/event/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", primitive.NewObjectID())
	})
	return r
}

func TestEventHandler_GetPastEvents(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	r.GET("/events/past", handler.GetPastEvents)

	expectedEvents := []*entities.Event{{ID: primitive.NewObjectID()}}
	mockUseCase.On("GetPastEvents", mock.Anything, 20).Return(expectedEvents, nil)

	req, _ := http.NewRequest(http.MethodGet, "/events/past", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(len(expectedEvents)), res["total"])

	mockUseCase.AssertExpectations(t)
}

func TestEventHandler_RegisterForEvent(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	eventID := primitive.NewObjectID()
	r.POST("/events/:id/register", handler.RegisterForEvent)

	reqBody := entities.EventAttendance{AttendeeType: entities.AttendeeTypePublic}
	resBody := &entities.EventAttendance{ID: primitive.NewObjectID()}

	mockUseCase.On("RegisterForEvent", mock.Anything, eventID, mock.AnythingOfType("entities.EventAttendance"), mock.AnythingOfType("primitive.ObjectID")).Return(resBody, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/events/"+eventID.Hex()+"/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res entities.EventAttendance
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, resBody.ID, res.ID)

	mockUseCase.AssertExpectations(t)
}

func TestEventHandler_GetEventRegistrations(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	eventID := primitive.NewObjectID()
	r.GET("/events/:id/registrations", handler.GetEventRegistrations)

	expectedRegs := []*entities.EventAttendance{{ID: primitive.NewObjectID()}}
	mockUseCase.On("GetEventRegistrations", mock.Anything, eventID, int64(20), int64(0)).Return(expectedRegs, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.Hex()+"/registrations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(1), res["total"])

	mockUseCase.AssertExpectations(t)
}

func TestEventHandler_GetEventStatisticsDetail(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	eventID := primitive.NewObjectID()
	r.GET("/events/:id/statistics", handler.GetEventStatisticsDetail)

	event := &entities.Event{ID: eventID, AttendeeCount: 10}
	mockUseCase.On("GetEventStatisticsDetail", mock.Anything, eventID).Return(event, nil)

	req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.Hex()+"/statistics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(10), res["attendee_count"])

	mockUseCase.AssertExpectations(t)
}

func TestEventHandler_PublishEvent(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	eventID := primitive.NewObjectID()
	r.POST("/events/:id/publish", handler.PublishEvent)

	mockUseCase.On("PublishEvent", mock.Anything, eventID, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/events/"+eventID.Hex()+"/publish", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "Event published successfully", res["message"])

	mockUseCase.AssertExpectations(t)
}

func TestEventHandler_SendEventReminder(t *testing.T) {
	mockUseCase := new(mocks.EventUseCase)
	handler := NewEventHandler(mockUseCase)

	r := setupRouter()
	eventID := primitive.NewObjectID()
	r.POST("/events/:id/send-reminder", handler.SendEventReminder)

	mockUseCase.On("SendEventReminder", mock.Anything, eventID).Return(5, nil)

	req, _ := http.NewRequest(http.MethodPost, "/events/"+eventID.Hex()+"/send-reminder", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(5), res["recipients_count"])

	mockUseCase.AssertExpectations(t)
}
