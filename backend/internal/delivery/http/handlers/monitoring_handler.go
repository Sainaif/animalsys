package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/usecase/monitoring"
)

type MonitoringHandler struct {
	monitoringUseCase *monitoring.MonitoringUseCase
}

func NewMonitoringHandler(monitoringUseCase *monitoring.MonitoringUseCase) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringUseCase: monitoringUseCase,
	}
}

// GetSystemHealth returns overall system health status
func (h *MonitoringHandler) GetSystemHealth(c *gin.Context) {
	health, err := h.monitoringUseCase.GetSystemHealth(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, health)
}

// GetUsageStatistics returns system usage statistics
func (h *MonitoringHandler) GetUsageStatistics(c *gin.Context) {
	stats, err := h.monitoringUseCase.GetUsageStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPerformanceMetrics returns system performance metrics
func (h *MonitoringHandler) GetPerformanceMetrics(c *gin.Context) {
	metrics, err := h.monitoringUseCase.GetPerformanceMetrics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDatabaseStatistics returns detailed database statistics
func (h *MonitoringHandler) GetDatabaseStatistics(c *gin.Context) {
	stats, err := h.monitoringUseCase.GetDatabaseStatistics(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSystemConfiguration returns system configuration
func (h *MonitoringHandler) GetSystemConfiguration(c *gin.Context) {
	config, err := h.monitoringUseCase.GetSystemConfiguration(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, config)
}
