package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/usecases"
)

type AnimalHandler struct {
	animalUseCase *usecases.AnimalUseCase
}

func NewAnimalHandler(animalUseCase *usecases.AnimalUseCase) *AnimalHandler {
	return &AnimalHandler{
		animalUseCase: animalUseCase,
	}
}

// Create godoc
// @Summary Create animal
// @Description Create a new animal record
// @Tags animals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entities.AnimalCreateRequest true "Animal data"
// @Success 201 {object} entities.Animal
// @Failure 400 {object} ErrorResponse
// @Router /animals [post]
func (h *AnimalHandler) Create(c *gin.Context) {
	var req entities.AnimalCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdBy, _ := c.Get("user_id")
	animal, err := h.animalUseCase.Create(c.Request.Context(), &req, createdBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, animal)
}

// GetByID godoc
// @Summary Get animal by ID
// @Description Get animal details by ID
// @Tags animals
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {object} entities.Animal
// @Failure 404 {object} ErrorResponse
// @Router /animals/{id} [get]
func (h *AnimalHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	animal, err := h.animalUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}

// List godoc
// @Summary List animals
// @Description Get list of animals with filtering and pagination
// @Tags animals
// @Produce json
// @Param species query string false "Filter by species"
// @Param status query string false "Filter by status"
// @Param gender query string false "Filter by gender"
// @Param search query string false "Search in name/chip_number"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /animals [get]
func (h *AnimalHandler) List(c *gin.Context) {
	filter := &entities.AnimalFilter{
		Species: c.Query("species"),
		Status:  entities.AnimalStatus(c.Query("status")),
		Gender:  c.Query("gender"),
		Search:  c.Query("search"),
		Limit:   parseIntQuery(c.Query("limit"), 10),
		Offset:  parseIntQuery(c.Query("offset"), 0),
	}

	animals, total, err := h.animalUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   animals,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Update godoc
// @Summary Update animal
// @Description Update animal details
// @Tags animals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Animal ID"
// @Param request body entities.AnimalUpdateRequest true "Update data"
// @Success 200 {object} entities.Animal
// @Failure 400 {object} ErrorResponse
// @Router /animals/{id} [put]
func (h *AnimalHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entities.AnimalUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	updatedBy, _ := c.Get("user_id")
	animal, err := h.animalUseCase.Update(c.Request.Context(), id, &req, updatedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}

// Delete godoc
// @Summary Delete animal
// @Description Delete animal by ID
// @Tags animals
// @Security BearerAuth
// @Param id path string true "Animal ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /animals/{id} [delete]
func (h *AnimalHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedBy, _ := c.Get("user_id")
	err := h.animalUseCase.Delete(c.Request.Context(), id, deletedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "animal deleted successfully",
	})
}

// AddMedicalRecord godoc
// @Summary Add medical record
// @Description Add a medical record to an animal
// @Tags animals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Animal ID"
// @Param request body entities.MedicalRecord true "Medical record data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /animals/{id}/medical-records [post]
func (h *AnimalHandler) AddMedicalRecord(c *gin.Context) {
	id := c.Param("id")

	var record entities.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	addedBy, _ := c.Get("user_id")
	err := h.animalUseCase.AddMedicalRecord(c.Request.Context(), id, record, addedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "medical record added successfully",
	})
}

// AddPhoto godoc
// @Summary Add photo
// @Description Add a photo URL to an animal
// @Tags animals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Animal ID"
// @Param request body map[string]string true "Photo data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /animals/{id}/photos [post]
func (h *AnimalHandler) AddPhoto(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PhotoURL string `json:"photo_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	addedBy, _ := c.Get("user_id")
	err := h.animalUseCase.AddPhoto(c.Request.Context(), id, req.PhotoURL, addedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "photo added successfully",
	})
}

// GetAvailableForAdoption godoc
// @Summary Get available animals
// @Description Get animals available for adoption
// @Tags animals
// @Produce json
// @Success 200 {object} []entities.Animal
// @Failure 500 {object} ErrorResponse
// @Router /animals/available [get]
func (h *AnimalHandler) GetAvailableForAdoption(c *gin.Context) {
	animals, err := h.animalUseCase.GetAvailableForAdoption(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, animals)
}
