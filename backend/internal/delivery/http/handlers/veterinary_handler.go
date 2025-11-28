package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/veterinary"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VeterinaryHandler wires veterinary HTTP endpoints to the use case layer.
type VeterinaryHandler struct {
	useCase  veterinary.VeterinaryUseCaseInterface
	validate *validator.Validate
}

// NewVeterinaryHandler creates a fully wired veterinary handler.
func NewVeterinaryHandler(useCase veterinary.VeterinaryUseCaseInterface) *VeterinaryHandler {
	return &VeterinaryHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

func (h *VeterinaryHandler) respondWithError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// GetVisitsByAnimal returns all veterinary visits for a specific animal.
func (h *VeterinaryHandler) GetVisitsByAnimal(c *gin.Context) {
	animalID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	visits, err := h.useCase.GetVisitsByAnimalID(c.Request.Context(), animalID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"visits": visits,
		"total":  len(visits),
	})
}

// GetVaccinationsByAnimal returns all vaccinations for a specific animal.
func (h *VeterinaryHandler) GetVaccinationsByAnimal(c *gin.Context) {
	animalID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	vaccinations, err := h.useCase.GetVaccinationsByAnimalID(c.Request.Context(), animalID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vaccinations": vaccinations,
		"total":        len(vaccinations),
	})
}

// ListVisits lists veterinary visits with optional filters.
func (h *VeterinaryHandler) ListVisits(c *gin.Context) {
	var req veterinary.ListVisitsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visits, total, err := h.useCase.ListVisits(c.Request.Context(), &req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"visits": visits,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	})
}

// GetUpcomingVisits lists upcoming scheduled visits within the provided number of days.
func (h *VeterinaryHandler) GetUpcomingVisits(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days parameter"})
		return
	}

	visits, err := h.useCase.GetUpcomingVisits(c.Request.Context(), days)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"visits": visits,
		"total":  len(visits),
	})
}

// GetVisit returns a single veterinary visit.
func (h *VeterinaryHandler) GetVisit(c *gin.Context) {
	visitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visit ID"})
		return
	}

	visit, err := h.useCase.GetVisitByID(c.Request.Context(), visitID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"visit": visit})
}

// CreateVisit creates a new veterinary visit.
func (h *VeterinaryHandler) CreateVisit(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req veterinary.CreateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit, err := h.useCase.CreateVisit(c.Request.Context(), &req, *userID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, visit)
}

// UpdateVisit updates an existing veterinary visit.
func (h *VeterinaryHandler) UpdateVisit(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	visitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visit ID"})
		return
	}

	var req veterinary.UpdateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit, err := h.useCase.UpdateVisit(c.Request.Context(), visitID, &req, *userID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, visit)
}

// DeleteVisit deletes a veterinary visit.
func (h *VeterinaryHandler) DeleteVisit(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	visitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visit ID"})
		return
	}

	if err := h.useCase.DeleteVisit(c.Request.Context(), visitID, *userID); err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "visit deleted"})
}

// ListVaccinations lists vaccination records with optional filtering.
func (h *VeterinaryHandler) ListVaccinations(c *gin.Context) {
	var req veterinary.ListVaccinationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vaccinations, total, err := h.useCase.ListVaccinations(c.Request.Context(), &req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vaccinations": vaccinations,
		"total":        total,
		"limit":        req.Limit,
		"offset":       req.Offset,
	})
}

// GetDueVaccinations lists vaccinations that are due within provided days.
func (h *VeterinaryHandler) GetDueVaccinations(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days parameter"})
		return
	}

	vaccinations, err := h.useCase.GetDueVaccinations(c.Request.Context(), days)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vaccinations": vaccinations,
		"total":        len(vaccinations),
	})
}

// GetVaccination returns a single vaccination record.
func (h *VeterinaryHandler) GetVaccination(c *gin.Context) {
	vaccinationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vaccination ID"})
		return
	}

	vaccination, err := h.useCase.GetVaccinationByID(c.Request.Context(), vaccinationID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"vaccination": vaccination})
}

// CreateVaccination creates a new vaccination record.
func (h *VeterinaryHandler) CreateVaccination(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req veterinary.CreateVaccinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vaccination, err := h.useCase.CreateVaccination(c.Request.Context(), &req, *userID)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, vaccination)
}

// DeleteVaccination deletes an existing vaccination record.
func (h *VeterinaryHandler) DeleteVaccination(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	vaccinationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vaccination ID"})
		return
	}

	if err := h.useCase.DeleteVaccination(c.Request.Context(), vaccinationID, *userID); err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vaccination deleted"})
}

// A request to create a veterinary record, which can be either a visit or a vaccination.
type CreateRecordRequest struct {
	RecordType  string                           `json:"record_type" binding:"required,oneof=visit vaccination"`
	Visit       *veterinary.CreateVisitRequest       `json:"visit,omitempty"`
	Vaccination *veterinary.CreateVaccinationRequest `json:"vaccination,omitempty"`
}

// Placeholder endpoints for future veterinary records features.
func (h *VeterinaryHandler) CreateVeterinaryRecord(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.RecordType {
	case "visit":
		if req.Visit == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "visit data is required for record_type 'visit'"})
			return
		}
		if err := h.validate.Struct(req.Visit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		record, err := h.useCase.CreateVisit(c.Request.Context(), req.Visit, *userID)
		if err != nil {
			h.respondWithError(c, err)
			return
		}
		c.JSON(http.StatusCreated, record)
	case "vaccination":
		if req.Vaccination == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vaccination data is required for record_type 'vaccination'"})
			return
		}
		if err := h.validate.Struct(req.Vaccination); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		record, err := h.useCase.CreateVaccination(c.Request.Context(), req.Vaccination, *userID)
		if err != nil {
			h.respondWithError(c, err)
			return
		}
		c.JSON(http.StatusCreated, record)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_type"})
	}
}

func (h *VeterinaryHandler) ListVeterinaryRecords(c *gin.Context) {
	var req veterinary.ListRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records, total, err := h.useCase.ListVeterinaryRecords(c.Request.Context(), &req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   total,
		"limit":   req.Limit,
		"offset":  req.Offset,
	})
}
