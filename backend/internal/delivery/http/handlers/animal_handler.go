package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/animal"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnimalHandler handles animal management HTTP requests
type AnimalHandler struct {
	animalUseCase *animal.AnimalUseCase
	validate      *validator.Validate
}

// NewAnimalHandler creates a new animal handler
func NewAnimalHandler(animalUseCase *animal.AnimalUseCase) *AnimalHandler {
	return &AnimalHandler{
		animalUseCase: animalUseCase,
		validate:      validator.New(),
	}
}

// ListAnimals lists all animals with filters and pagination
// @Summary List Animals
// @Description Get list of animals with filters and pagination
// @Tags animals
// @Security BearerAuth
// @Produce json
// @Param category query string false "Filter by category"
// @Param species query string false "Filter by species"
// @Param status query string false "Filter by status"
// @Param sex query string false "Filter by sex"
// @Param size query string false "Filter by size"
// @Param available_only query bool false "Show only available animals"
// @Param search query string false "Search in name and description"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} animal.ListAnimalsResponse
// @Failure 401 {object} errors.AppError
// @Router /animals [get]
func (h *AnimalHandler) ListAnimals(c *gin.Context) {
	var req animal.ListAnimalsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.animalUseCase.ListAnimals(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAnimal gets an animal by ID
// @Summary Get Animal
// @Description Get animal details by ID
// @Tags animals
// @Security BearerAuth
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {object} entities.Animal
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id} [get]
func (h *AnimalHandler) GetAnimal(c *gin.Context) {
	idParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	animal, err := h.animalUseCase.GetAnimalByID(c.Request.Context(), animalID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, animal)
}

// CreateAnimal creates a new animal
// @Summary Create Animal
// @Description Create a new animal
// @Tags animals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body animal.CreateAnimalRequest true "Animal details"
// @Success 201 {object} entities.Animal
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Router /animals [post]
func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	creatorID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req animal.CreateAnimalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAnimal, err := h.animalUseCase.CreateAnimal(c.Request.Context(), &req, *creatorID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, newAnimal)
}

// UpdateAnimal updates an animal
// @Summary Update Animal
// @Description Update animal details
// @Tags animals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Animal ID"
// @Param request body animal.UpdateAnimalRequest true "Animal details to update"
// @Success 200 {object} entities.Animal
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id} [put]
func (h *AnimalHandler) UpdateAnimal(c *gin.Context) {
	updaterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	var req animal.UpdateAnimalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAnimal, err := h.animalUseCase.UpdateAnimal(c.Request.Context(), animalID, &req, *updaterID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedAnimal)
}

// DeleteAnimal deletes an animal
// @Summary Delete Animal
// @Description Delete an animal
// @Tags animals
// @Security BearerAuth
// @Param id path string true "Animal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id} [delete]
func (h *AnimalHandler) DeleteAnimal(c *gin.Context) {
	deleterID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	if err := h.animalUseCase.DeleteAnimal(c.Request.Context(), animalID, *deleterID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "animal deleted successfully"})
}

// UploadAnimalImages uploads images for an animal
// @Summary Upload Animal Images
// @Description Upload primary and gallery images for an animal
// @Tags animals
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Animal ID"
// @Param primary formData file false "Primary image"
// @Param gallery formData file false "Gallery images" collectionFormat(multi)
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id}/images [post]
func (h *AnimalHandler) UploadAnimalImages(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	// Get primary image
	primary, _ := c.FormFile("primary")

	// Get gallery images
	form, err := c.MultipartForm()
	var gallery []*multipart.FileHeader
	if err == nil {
		gallery = form.File["gallery"]
	}

	if primary == nil && len(gallery) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no images provided"})
		return
	}

	if err := h.animalUseCase.UploadAnimalImages(c.Request.Context(), animalID, primary, gallery, *userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "images uploaded successfully"})
}

// AddDailyNote adds a daily note to an animal
// @Summary Add Daily Note
// @Description Add a daily observation note to an animal
// @Tags animals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Animal ID"
// @Param request body map[string]string true "Note content"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /animals/{id}/notes [post]
func (h *AnimalHandler) AddDailyNote(c *gin.Context) {
	userID, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	animalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	var req struct {
		Note string `json:"note" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.animalUseCase.AddDailyNote(c.Request.Context(), animalID, req.Note, *userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note added successfully"})
}

// GetStatistics gets animal statistics
// @Summary Get Animal Statistics
// @Description Get statistical information about animals
// @Tags animals
// @Security BearerAuth
// @Produce json
// @Success 200 {object} repositories.AnimalStatistics
// @Failure 401 {object} errors.AppError
// @Router /animals/statistics [get]
func (h *AnimalHandler) GetStatistics(c *gin.Context) {
	stats, err := h.animalUseCase.GetStatistics(c.Request.Context())
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

// GetSpecies gets available species for a category
// @Summary Get Species
// @Description Get available species for a given category
// @Tags animals
// @Produce json
// @Param category query string true "Animal category"
// @Success 200 {array} entities.SpeciesInfo
// @Failure 400 {object} errors.AppError
// @Router /animals/species [get]
func (h *AnimalHandler) GetSpecies(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	species := h.animalUseCase.GetSpeciesByCategory(entities.AnimalCategory(category))
	c.JSON(http.StatusOK, species)
}
