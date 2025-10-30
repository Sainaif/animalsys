package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type VeterinaryHandler struct {
	veterinaryUseCase *usecases.VeterinaryUseCase
}

func NewVeterinaryHandler(veterinaryUseCase *usecases.VeterinaryUseCase) *VeterinaryHandler {
	return &VeterinaryHandler{
		veterinaryUseCase: veterinaryUseCase,
	}
}

// Veterinary Visit endpoints

func (h *VeterinaryHandler) RecordVisit(c *gin.Context) {
	var req entities.VeterinaryVisitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	recordedBy, _ := c.Get("user_id")
	visit, err := h.veterinaryUseCase.RecordVisit(c.Request.Context(), &req, recordedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, visit)
}

func (h *VeterinaryHandler) GetVisitByID(c *gin.Context) {
	id := c.Param("id")

	visit, err := h.veterinaryUseCase.GetVisitByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, visit)
}

func (h *VeterinaryHandler) ListVisits(c *gin.Context) {
	filter := &entities.VeterinaryFilter{
		AnimalID:  c.Query("animal_id"),
		VisitType: c.Query("visit_type"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Limit:     parseIntQuery(c.Query("limit"), 10),
		Offset:    parseIntQuery(c.Query("offset"), 0),
	}

	visits, total, err := h.veterinaryUseCase.ListVisits(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   visits,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *VeterinaryHandler) UpdateVisit(c *gin.Context) {
	id := c.Param("id")

	var req entities.VeterinaryVisitUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	visit, err := h.veterinaryUseCase.UpdateVisit(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, visit)
}

func (h *VeterinaryHandler) DeleteVisit(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.veterinaryUseCase.DeleteVisit(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "visit deleted successfully"})
}

func (h *VeterinaryHandler) GetUpcomingFollowUps(c *gin.Context) {
	days := parseIntQuery(c.Query("days"), 7)

	visits, err := h.veterinaryUseCase.GetUpcomingFollowUps(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, visits)
}

func (h *VeterinaryHandler) GetStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_date and end_date are required"})
		return
	}

	stats, err := h.veterinaryUseCase.GetVeterinaryStatistics(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Vaccination endpoints

func (h *VeterinaryHandler) RecordVaccination(c *gin.Context) {
	var req entities.VaccinationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	recordedBy, _ := c.Get("user_id")
	vaccination, err := h.veterinaryUseCase.RecordVaccination(c.Request.Context(), &req, recordedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vaccination)
}

func (h *VeterinaryHandler) GetVaccinationByID(c *gin.Context) {
	id := c.Param("id")

	vaccination, err := h.veterinaryUseCase.GetVaccinationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vaccination)
}

func (h *VeterinaryHandler) ListVaccinations(c *gin.Context) {
	filter := &entities.VaccinationFilter{
		AnimalID: c.Query("animal_id"),
		Limit:    parseIntQuery(c.Query("limit"), 10),
		Offset:   parseIntQuery(c.Query("offset"), 0),
	}

	vaccinations, total, err := h.veterinaryUseCase.ListVaccinations(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   vaccinations,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *VeterinaryHandler) UpdateVaccination(c *gin.Context) {
	id := c.Param("id")

	var req entities.VaccinationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	vaccination, err := h.veterinaryUseCase.UpdateVaccination(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vaccination)
}

func (h *VeterinaryHandler) DeleteVaccination(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.veterinaryUseCase.DeleteVaccination(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vaccination deleted successfully"})
}

func (h *VeterinaryHandler) GetUpcomingVaccinations(c *gin.Context) {
	days := parseIntQuery(c.Query("days"), 30)

	vaccinations, err := h.veterinaryUseCase.GetUpcomingVaccinations(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vaccinations)
}
