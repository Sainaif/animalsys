package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/internal/infrastructure/middleware"
	"github.com/sainaif/animalsys/backend/internal/usecase/medical"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedicalHandler struct {
	medicalUseCase *medical.MedicalUseCase
}

func NewMedicalHandler(medicalUseCase *medical.MedicalUseCase) *MedicalHandler {
	return &MedicalHandler{
		medicalUseCase: medicalUseCase,
	}
}

// ===== Medical Condition Handlers =====

// CreateCondition creates a new medical condition
func (h *MedicalHandler) CreateCondition(c *gin.Context) {
	var condition entities.MedicalCondition
	if err := c.ShouldBindJSON(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err := h.medicalUseCase.CreateCondition(c.Request.Context(), &condition, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, condition)
}

// GetCondition gets a medical condition by ID
func (h *MedicalHandler) GetCondition(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid condition ID"})
		return
	}

	condition, err := h.medicalUseCase.GetCondition(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, condition)
}

// UpdateCondition updates a medical condition
func (h *MedicalHandler) UpdateCondition(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid condition ID"})
		return
	}

	var condition entities.MedicalCondition
	if err := c.ShouldBindJSON(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	condition.ID = id
	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.UpdateCondition(c.Request.Context(), &condition, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, condition)
}

// DeleteCondition deletes a medical condition
func (h *MedicalHandler) DeleteCondition(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid condition ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.DeleteCondition(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condition deleted successfully"})
}

// ListConditions lists medical conditions with filtering
func (h *MedicalHandler) ListConditions(c *gin.Context) {
	filter := repositories.MedicalConditionFilter{}

	if animalID := c.Query("animal_id"); animalID != "" {
		id, err := primitive.ObjectIDFromHex(animalID)
		if err == nil {
			filter.AnimalID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = entities.ConditionStatus(status)
	}

	if severity := c.Query("severity"); severity != "" {
		filter.Severity = entities.ConditionSeverity(severity)
	}

	if chronic := c.Query("is_chronic"); chronic != "" {
		isChronic := chronic == "true"
		filter.IsChronic = &isChronic
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.ParseInt(limit, 10, 64); err == nil {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.ParseInt(offset, 10, 64); err == nil {
			filter.Offset = val
		}
	}

	conditions, total, err := h.medicalUseCase.ListConditions(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  conditions,
		"total": total,
	})
}

// GetAnimalConditions gets all conditions for an animal
func (h *MedicalHandler) GetAnimalConditions(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	conditions, err := h.medicalUseCase.GetAnimalConditions(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, conditions)
}

// GetActiveConditions gets active conditions for an animal
func (h *MedicalHandler) GetActiveConditions(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	conditions, err := h.medicalUseCase.GetActiveConditions(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, conditions)
}

// GetChronicConditions gets all chronic conditions
func (h *MedicalHandler) GetChronicConditions(c *gin.Context) {
	conditions, err := h.medicalUseCase.GetChronicConditions(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, conditions)
}

// ResolveCondition marks a condition as resolved
func (h *MedicalHandler) ResolveCondition(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid condition ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.ResolveCondition(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condition resolved successfully"})
}

// ===== Medication Handlers =====

// CreateMedication creates a new medication
func (h *MedicalHandler) CreateMedication(c *gin.Context) {
	var medication entities.Medication
	if err := c.ShouldBindJSON(&medication); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err := h.medicalUseCase.CreateMedication(c.Request.Context(), &medication, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, medication)
}

// GetMedication gets a medication by ID
func (h *MedicalHandler) GetMedication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		return
	}

	medication, err := h.medicalUseCase.GetMedication(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medication)
}

// UpdateMedication updates a medication
func (h *MedicalHandler) UpdateMedication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		return
	}

	var medication entities.Medication
	if err := c.ShouldBindJSON(&medication); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medication.ID = id
	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.UpdateMedication(c.Request.Context(), &medication, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medication)
}

// DeleteMedication deletes a medication
func (h *MedicalHandler) DeleteMedication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.DeleteMedication(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medication deleted successfully"})
}

// ListMedications lists medications with filtering
func (h *MedicalHandler) ListMedications(c *gin.Context) {
	filter := repositories.MedicationFilter{}

	if animalID := c.Query("animal_id"); animalID != "" {
		id, err := primitive.ObjectIDFromHex(animalID)
		if err == nil {
			filter.AnimalID = &id
		}
	}

	if conditionID := c.Query("condition_id"); conditionID != "" {
		id, err := primitive.ObjectIDFromHex(conditionID)
		if err == nil {
			filter.ConditionID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = entities.MedicationStatus(status)
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.ParseInt(limit, 10, 64); err == nil {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.ParseInt(offset, 10, 64); err == nil {
			filter.Offset = val
		}
	}

	medications, total, err := h.medicalUseCase.ListMedications(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  medications,
		"total": total,
	})
}

// GetAnimalMedications gets all medications for an animal
func (h *MedicalHandler) GetAnimalMedications(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	medications, err := h.medicalUseCase.GetAnimalMedications(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medications)
}

// GetActiveMedications gets active medications for an animal
func (h *MedicalHandler) GetActiveMedications(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	medications, err := h.medicalUseCase.GetActiveMedications(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medications)
}

// GetMedicationsDueForRefill gets medications due for refill
func (h *MedicalHandler) GetMedicationsDueForRefill(c *gin.Context) {
	medications, err := h.medicalUseCase.GetMedicationsDueForRefill(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medications)
}

// GetExpiringSoonMedications gets medications expiring within specified days
func (h *MedicalHandler) GetExpiringSoonMedications(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	medications, err := h.medicalUseCase.GetExpiringSoonMedications(c.Request.Context(), days)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, medications)
}

// RecordMedicationAdministration records medication administration
func (h *MedicalHandler) RecordMedicationAdministration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		return
	}

	var req struct {
		DosageGiven string `json:"dosage_given" binding:"required"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	log := entities.AdministrationLog{
		AdministeredAt: time.Now(),
		AdministeredBy: *userID,
		DosageGiven:    req.DosageGiven,
		Notes:          req.Notes,
	}

	err = h.medicalUseCase.RecordMedicationAdministration(c.Request.Context(), id, log)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Administration recorded successfully"})
}

// RefillMedication processes a medication refill
func (h *MedicalHandler) RefillMedication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.RefillMedication(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medication refilled successfully"})
}

// ===== Treatment Plan Handlers =====

// CreateTreatmentPlan creates a new treatment plan
func (h *MedicalHandler) CreateTreatmentPlan(c *gin.Context) {
	var plan entities.TreatmentPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err := h.medicalUseCase.CreateTreatmentPlan(c.Request.Context(), &plan, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// GetTreatmentPlan gets a treatment plan by ID
func (h *MedicalHandler) GetTreatmentPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	plan, err := h.medicalUseCase.GetTreatmentPlan(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, plan)
}

// UpdateTreatmentPlan updates a treatment plan
func (h *MedicalHandler) UpdateTreatmentPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	var plan entities.TreatmentPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.ID = id
	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.UpdateTreatmentPlan(c.Request.Context(), &plan, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, plan)
}

// DeleteTreatmentPlan deletes a treatment plan
func (h *MedicalHandler) DeleteTreatmentPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.DeleteTreatmentPlan(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treatment plan deleted successfully"})
}

// ListTreatmentPlans lists treatment plans with filtering
func (h *MedicalHandler) ListTreatmentPlans(c *gin.Context) {
	filter := repositories.TreatmentPlanFilter{}

	if animalID := c.Query("animal_id"); animalID != "" {
		id, err := primitive.ObjectIDFromHex(animalID)
		if err == nil {
			filter.AnimalID = &id
		}
	}

	if conditionID := c.Query("condition_id"); conditionID != "" {
		id, err := primitive.ObjectIDFromHex(conditionID)
		if err == nil {
			filter.ConditionID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = entities.TreatmentPlanStatus(status)
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.ParseInt(limit, 10, 64); err == nil {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.ParseInt(offset, 10, 64); err == nil {
			filter.Offset = val
		}
	}

	plans, total, err := h.medicalUseCase.ListTreatmentPlans(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  plans,
		"total": total,
	})
}

// GetAnimalTreatmentPlans gets all treatment plans for an animal
func (h *MedicalHandler) GetAnimalTreatmentPlans(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	plans, err := h.medicalUseCase.GetAnimalTreatmentPlans(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetActiveTreatmentPlans gets active treatment plans for an animal
func (h *MedicalHandler) GetActiveTreatmentPlans(c *gin.Context) {
	animalIDParam := c.Param("animal_id")
	animalID, err := primitive.ObjectIDFromHex(animalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	plans, err := h.medicalUseCase.GetActiveTreatmentPlans(c.Request.Context(), animalID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, plans)
}

// AddProgressNote adds a progress note to a treatment plan
func (h *MedicalHandler) AddProgressNote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	var req struct {
		Note        string `json:"note" binding:"required"`
		Improvement string `json:"improvement" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	note := entities.ProgressNote{
		Date:        time.Now(),
		RecordedBy:  *userID,
		Note:        req.Note,
		Improvement: req.Improvement,
	}

	err = h.medicalUseCase.AddProgressNote(c.Request.Context(), id, note)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress note added successfully"})
}

// CompleteTreatmentPlan marks a treatment plan as completed
func (h *MedicalHandler) CompleteTreatmentPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.CompleteTreatmentPlan(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treatment plan completed successfully"})
}

// ActivateTreatmentPlan activates a draft treatment plan
func (h *MedicalHandler) ActivateTreatmentPlan(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid treatment plan ID"})
		return
	}

	userID, _ := middleware.GetUserFromContext(c)

	err = h.medicalUseCase.ActivateTreatmentPlan(c.Request.Context(), id, *userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treatment plan activated successfully"})
}
