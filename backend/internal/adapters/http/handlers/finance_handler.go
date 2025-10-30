package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type FinanceHandler struct {
	financeUseCase *usecases.FinanceUseCase
}

func NewFinanceHandler(financeUseCase *usecases.FinanceUseCase) *FinanceHandler {
	return &FinanceHandler{
		financeUseCase: financeUseCase,
	}
}

func (h *FinanceHandler) Create(c *gin.Context) {
	var req entities.FinanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	transaction, err := h.financeUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *FinanceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.financeUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *FinanceHandler) List(c *gin.Context) {
	filter := &entities.FinanceFilter{
		Type:       entities.TransactionType(c.Query("type")),
		Category:   entities.FinanceCategory(c.Query("category")),
		StartDate:  c.Query("start_date"),
		EndDate:    c.Query("end_date"),
		FiscalYear: c.Query("fiscal_year"),
		Limit:      parseIntQuery(c.Query("limit"), 10),
		Offset:     parseIntQuery(c.Query("offset"), 0),
	}

	transactions, total, err := h.financeUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   transactions,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *FinanceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.FinanceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	transaction, err := h.financeUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *FinanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.financeUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

func (h *FinanceHandler) GetSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_date and end_date are required"})
		return
	}

	summary, err := h.financeUseCase.GetSummary(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *FinanceHandler) GetReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_date and end_date are required"})
		return
	}

	report, err := h.financeUseCase.GetFinancialReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
