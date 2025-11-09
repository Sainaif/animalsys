package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type BatchHandler struct{}

func NewBatchHandler() *BatchHandler {
	return &BatchHandler{}
}

// CreateBatch creates a batch operation
func (h *BatchHandler) CreateBatch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Batch created"})
}
