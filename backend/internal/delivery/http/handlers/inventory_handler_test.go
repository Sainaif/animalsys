package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/usecase/inventory/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInventoryHandler_GetOutOfStockItems(t *testing.T) {
	useCase := new(mocks.InventoryUseCase)
	handler := NewInventoryHandler(useCase)

	expectedItems := []*entities.InventoryItem{
		{ID: primitive.NewObjectID(), Name: "Item 1"},
	}
	expectedTotal := int64(1)

	useCase.On("GetOutOfStockItems", mock.Anything).Return(expectedItems, expectedTotal, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler.GetOutOfStockItems(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(expectedTotal), response["total"])
}

func TestInventoryHandler_GetInventoryHistory(t *testing.T) {
	useCase := new(mocks.InventoryUseCase)
	handler := NewInventoryHandler(useCase)

	itemID := primitive.NewObjectID()
	expectedHistory := []*entities.StockTransaction{
		{ID: primitive.NewObjectID(), ItemID: itemID},
	}
	expectedTotal := int64(1)

	useCase.On("GetInventoryHistory", mock.Anything, itemID).Return(expectedHistory, expectedTotal, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: itemID.Hex()}}

	handler.GetInventoryHistory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(expectedTotal), response["total"])
}
