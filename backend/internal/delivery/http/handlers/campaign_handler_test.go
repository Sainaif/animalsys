package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/campaign"
	"github.com/sainaif/animalsys/backend/internal/usecase/campaign/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCampaignHandler_GetCampaignDonors(t *testing.T) {
	mockUseCase := new(mocks.CampaignUseCase)
	handler := NewCampaignHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/campaigns/:id/donors", handler.GetCampaignDonors)

	campaignID := primitive.NewObjectID()
	donors := []*campaign.CampaignDonor{
		{Donor: &entities.Donor{ID: primitive.NewObjectID()}, TotalAmount: 100},
	}

	mockUseCase.On("GetCampaignDonors", mock.Anything, campaignID, int64(20), int64(0)).Return(donors, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/campaigns/"+campaignID.Hex()+"/donors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["total"])
	assert.Len(t, response["donors"], 1)

	mockUseCase.AssertExpectations(t)
}

func TestCampaignHandler_UpdateCampaignAmount(t *testing.T) {
	mockUseCase := new(mocks.CampaignUseCase)
	handler := NewCampaignHandler(mockUseCase)
	userID := primitive.NewObjectID()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(func(c *gin.Context) { // Mock middleware
		c.Set("user_id", userID)
		c.Next()
	})
	router.POST("/campaigns/:id/update-amount", handler.UpdateCampaignAmount)

	campaignID := primitive.NewObjectID()
	payload := gin.H{"amount": 500.0}
	body, _ := json.Marshal(payload)

	mockUseCase.On("UpdateCampaignAmount", mock.Anything, campaignID, 500.0, userID).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/campaigns/"+campaignID.Hex()+"/update-amount", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Campaign amount updated successfully", response["message"])

	mockUseCase.AssertExpectations(t)
}

func TestCampaignHandler_GetCampaignProgress(t *testing.T) {
	mockUseCase := new(mocks.CampaignUseCase)
	handler := NewCampaignHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/campaigns/:id/progress", handler.GetCampaignProgress)

	campaignID := primitive.NewObjectID()
	progress := &campaign.CampaignProgress{
		GoalAmount:    1000,
		CurrentAmount: 250,
		Percentage:    25,
		DonationCount: 10,
	}

	mockUseCase.On("GetCampaignProgress", mock.Anything, campaignID).Return(progress, nil)

	req, _ := http.NewRequest(http.MethodGet, "/campaigns/"+campaignID.Hex()+"/progress", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response campaign.CampaignProgress
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, progress, &response)

	mockUseCase.AssertExpectations(t)
}

func TestCampaignHandler_ShareCampaign(t *testing.T) {
	mockUseCase := new(mocks.CampaignUseCase)
	handler := NewCampaignHandler(mockUseCase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/campaigns/:id/share", handler.ShareCampaign)

	campaignID := primitive.NewObjectID()
	shareable := &campaign.CampaignShareable{
		Title:       "Test Campaign",
		Description: "Test Description",
		URL:         "/campaigns/" + campaignID.Hex(),
	}

	mockUseCase.On("ShareCampaign", mock.Anything, campaignID).Return(shareable, nil)

	req, _ := http.NewRequest(http.MethodPost, "/campaigns/"+campaignID.Hex()+"/share", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response campaign.CampaignShareable
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, shareable, &response)

	mockUseCase.AssertExpectations(t)
}
