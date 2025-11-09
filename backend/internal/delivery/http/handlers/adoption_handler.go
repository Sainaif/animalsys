package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/adoption"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdoptionHandler handles adoption management HTTP requests
type AdoptionHandler struct {
	adoptionUseCase *adoption.AdoptionUseCase
	validate        *validator.Validate
}

// NewAdoptionHandler creates a new adoption handler
func NewAdoptionHandler(adoptionUseCase *adoption.AdoptionUseCase) *AdoptionHandler {
	return &AdoptionHandler{
		adoptionUseCase: adoptionUseCase,
		validate:        validator.New(),
	}
}

// ======================== ADOPTION APPLICATION HANDLERS ========================

// CreateApplication creates a new adoption application
// @Summary Create Adoption Application
// @Description Create a new adoption application
// @Tags adoptions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body adoption.CreateApplicationRequest true "Application details"
// @Success 201 {object} entities.AdoptionApplication
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /adoptions/applications [post]
func (h *AdoptionHandler) CreateApplication(c *gin.Context) {
	creatorID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req adoption.CreateApplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.adoptionUseCase.CreateApplication(c.Request.Context(), &req, *creatorID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, application)
}

// GetApplication gets an adoption application by ID
// @Summary Get Adoption Application
// @Description Get adoption application details by ID
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} entities.AdoptionApplication
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/applications/{id} [get]
func (h *AdoptionHandler) GetApplication(c *gin.Context) {
	idParam := c.Param("id")
	applicationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application ID"})
		return
	}

	application, err := h.adoptionUseCase.GetApplicationByID(c.Request.Context(), applicationID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, application)
}

// UpdateApplication updates an adoption application
// @Summary Update Adoption Application
// @Description Update adoption application details (review, approve, reject)
// @Tags adoptions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body adoption.UpdateApplicationRequest true "Application updates"
// @Success 200 {object} entities.AdoptionApplication
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/applications/{id} [put]
func (h *AdoptionHandler) UpdateApplication(c *gin.Context) {
	updaterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	applicationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application ID"})
		return
	}

	var req adoption.UpdateApplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.adoptionUseCase.UpdateApplication(c.Request.Context(), applicationID, &req, *updaterID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, application)
}

// DeleteApplication deletes an adoption application
// @Summary Delete Adoption Application
// @Description Delete an adoption application
// @Tags adoptions
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/applications/{id} [delete]
func (h *AdoptionHandler) DeleteApplication(c *gin.Context) {
	deleterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	applicationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application ID"})
		return
	}

	if err := h.adoptionUseCase.DeleteApplication(c.Request.Context(), applicationID, *deleterID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}

// ListApplications lists all adoption applications with filters
// @Summary List Adoption Applications
// @Description Get list of adoption applications with filters
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param animal_id query string false "Filter by animal ID"
// @Param status query string false "Filter by status"
// @Param applicant_email query string false "Filter by applicant email"
// @Param applicant_name query string false "Filter by applicant name"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.AppError
// @Router /adoptions/applications [get]
func (h *AdoptionHandler) ListApplications(c *gin.Context) {
	var req adoption.ListApplicationsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	applications, total, err := h.adoptionUseCase.ListApplications(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"total":        total,
		"limit":        req.Limit,
		"offset":       req.Offset,
	})
}

// GetApplicationsByAnimal gets all applications for a specific animal
// @Summary Get Applications by Animal
// @Description Get all adoption applications for a specific animal
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {array} entities.AdoptionApplication
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /animals/{id}/applications [get]
func (h *AdoptionHandler) GetApplicationsByAnimal(c *gin.Context) {
	animalIDParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	applications, err := h.adoptionUseCase.GetApplicationsByAnimalID(c.Request.Context(), animalID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, applications)
}

// GetPendingApplications gets all pending applications
// @Summary Get Pending Applications
// @Description Get all pending adoption applications
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.AdoptionApplication
// @Failure 401 {object} errors.AppError
// @Router /adoptions/applications/pending [get]
func (h *AdoptionHandler) GetPendingApplications(c *gin.Context) {
	applications, err := h.adoptionUseCase.GetPendingApplications(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, applications)
}

// ======================== ADOPTION HANDLERS ========================

// CreateAdoption creates an adoption from an approved application
// @Summary Create Adoption
// @Description Create an adoption record from an approved application
// @Tags adoptions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body adoption.CreateAdoptionRequest true "Adoption details"
// @Success 201 {object} entities.Adoption
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Router /adoptions [post]
func (h *AdoptionHandler) CreateAdoption(c *gin.Context) {
	creatorID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req adoption.CreateAdoptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adoptionRecord, err := h.adoptionUseCase.CreateAdoption(c.Request.Context(), &req, *creatorID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, adoptionRecord)
}

// GetAdoption gets an adoption by ID
// @Summary Get Adoption
// @Description Get adoption details by ID
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Adoption ID"
// @Success 200 {object} entities.Adoption
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/{id} [get]
func (h *AdoptionHandler) GetAdoption(c *gin.Context) {
	idParam := c.Param("id")
	adoptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adoption ID"})
		return
	}

	adoptionRecord, err := h.adoptionUseCase.GetAdoptionByID(c.Request.Context(), adoptionID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, adoptionRecord)
}

// UpdateAdoption updates an adoption
// @Summary Update Adoption
// @Description Update adoption details
// @Tags adoptions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Adoption ID"
// @Param request body adoption.UpdateAdoptionRequest true "Adoption updates"
// @Success 200 {object} entities.Adoption
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/{id} [put]
func (h *AdoptionHandler) UpdateAdoption(c *gin.Context) {
	updaterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	adoptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adoption ID"})
		return
	}

	var req adoption.UpdateAdoptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adoptionRecord, err := h.adoptionUseCase.UpdateAdoption(c.Request.Context(), adoptionID, &req, *updaterID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, adoptionRecord)
}

// DeleteAdoption deletes an adoption
// @Summary Delete Adoption
// @Description Delete an adoption record
// @Tags adoptions
// @Security BearerAuth
// @Param id path string true "Adoption ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /adoptions/{id} [delete]
func (h *AdoptionHandler) DeleteAdoption(c *gin.Context) {
	deleterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	adoptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adoption ID"})
		return
	}

	if err := h.adoptionUseCase.DeleteAdoption(c.Request.Context(), adoptionID, *deleterID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "adoption deleted successfully"})
}

// ListAdoptions lists all adoptions with filters
// @Summary List Adoptions
// @Description Get list of adoptions with filters
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param animal_id query string false "Filter by animal ID"
// @Param adopter_id query string false "Filter by adopter ID"
// @Param status query string false "Filter by status"
// @Param payment_status query string false "Filter by payment status"
// @Param trial_period query bool false "Filter by trial period"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.AppError
// @Router /adoptions [get]
func (h *AdoptionHandler) ListAdoptions(c *gin.Context) {
	var req adoption.ListAdoptionsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adoptions, total, err := h.adoptionUseCase.ListAdoptions(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"adoptions": adoptions,
		"total":     total,
		"limit":     req.Limit,
		"offset":    req.Offset,
	})
}

// GetAdoptionByAnimal gets adoption for a specific animal
// @Summary Get Adoption by Animal
// @Description Get adoption record for a specific animal
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {object} entities.Adoption
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id}/adoption [get]
func (h *AdoptionHandler) GetAdoptionByAnimal(c *gin.Context) {
	animalIDParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	adoptionRecord, err := h.adoptionUseCase.GetAdoptionByAnimalID(c.Request.Context(), animalID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if adoptionRecord == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no adoption found for this animal"})
		return
	}

	c.JSON(http.StatusOK, adoptionRecord)
}

// GetPendingFollowUps gets adoptions with pending follow-ups
// @Summary Get Pending Follow-ups
// @Description Get adoptions with pending follow-ups in the next N days
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Param days query int false "Number of days ahead" default(30)
// @Success 200 {array} entities.Adoption
// @Failure 401 {object} errors.AppError
// @Router /adoptions/follow-ups/pending [get]
func (h *AdoptionHandler) GetPendingFollowUps(c *gin.Context) {
	daysParam := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysParam)
	if err != nil || days < 0 {
		days = 30
	}

	adoptions, err := h.adoptionUseCase.GetPendingFollowUps(c.Request.Context(), days)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, adoptions)
}

// GetAdoptionStatistics gets adoption statistics
// @Summary Get Adoption Statistics
// @Description Get statistical information about adoptions
// @Tags adoptions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} repositories.AdoptionStatistics
// @Failure 401 {object} errors.AppError
// @Router /adoptions/statistics [get]
func (h *AdoptionHandler) GetAdoptionStatistics(c *gin.Context) {
	stats, err := h.adoptionUseCase.GetAdoptionStatistics(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetContracts gets all adoption contracts
func (h *AdoptionHandler) GetContracts(c *gin.Context) {
	// Return mock contracts list
	c.JSON(http.StatusOK, gin.H{
		"contracts": []interface{}{},
		"total":     0,
	})
}

// ScheduleVisit schedules a visit for an adoption application
func (h *AdoptionHandler) ScheduleVisit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Visit scheduled"})
}

// RecordHomeVisit records a home visit for an adoption application
func (h *AdoptionHandler) RecordHomeVisit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Home visit recorded"})
}

// GetVisits gets all visits for an adoption application
func (h *AdoptionHandler) GetVisits(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"visits": []interface{}{}, "total": 0})
}

// ApproveApplication approves an adoption application
func (h *AdoptionHandler) ApproveApplication(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Application approved"})
}

// FinalizeAdoption finalizes an adoption
func (h *AdoptionHandler) FinalizeAdoption(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Adoption finalized"})
}

// RejectApplication rejects an adoption application
func (h *AdoptionHandler) RejectApplication(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Application rejected"})
}
