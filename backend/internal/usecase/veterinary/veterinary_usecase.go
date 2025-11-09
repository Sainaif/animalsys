package veterinary

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VeterinaryUseCase handles veterinary business logic
type VeterinaryUseCase struct {
	visitRepo        repositories.VeterinaryVisitRepository
	vaccinationRepo  repositories.VaccinationRepository
	animalRepo       repositories.AnimalRepository
	auditLogRepo     repositories.AuditLogRepository
}

// NewVeterinaryUseCase creates a new veterinary use case
func NewVeterinaryUseCase(
	visitRepo repositories.VeterinaryVisitRepository,
	vaccinationRepo repositories.VaccinationRepository,
	animalRepo repositories.AnimalRepository,
	auditLogRepo repositories.AuditLogRepository,
) *VeterinaryUseCase {
	return &VeterinaryUseCase{
		visitRepo:       visitRepo,
		vaccinationRepo: vaccinationRepo,
		animalRepo:      animalRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// CreateVisitRequest represents a request to create a veterinary visit
type CreateVisitRequest struct {
	AnimalID          string                     `json:"animal_id" validate:"required"`
	VisitType         entities.VisitType         `json:"visit_type" validate:"required"`
	VisitDate         time.Time                  `json:"visit_date" validate:"required"`
	VeterinarianName  string                     `json:"veterinarian_name" validate:"required"`
	ClinicName        string                     `json:"clinic_name,omitempty"`
	ClinicAddress     string                     `json:"clinic_address,omitempty"`
	ClinicPhone       string                     `json:"clinic_phone,omitempty"`
	ChiefComplaint    string                     `json:"chief_complaint,omitempty"`
	VitalSigns        entities.VitalSigns        `json:"vital_signs"`
	PhysicalExamNotes string                     `json:"physical_exam_notes,omitempty"`
	Diagnosis         string                     `json:"diagnosis,omitempty"`
	Treatment         string                     `json:"treatment,omitempty"`
	Prescriptions     []entities.Prescription    `json:"prescriptions,omitempty"`
	TestsOrdered      []string                   `json:"tests_ordered,omitempty"`
	VaccinationsGiven []string                   `json:"vaccinations_given,omitempty"`
	FollowUpRequired  bool                       `json:"follow_up_required"`
	FollowUpDate      *time.Time                 `json:"follow_up_date,omitempty"`
	FollowUpNotes     string                     `json:"follow_up_notes,omitempty"`
	Cost              float64                    `json:"cost"`
	PaymentStatus     string                     `json:"payment_status,omitempty"`
}

// UpdateVisitRequest represents a request to update a veterinary visit
type UpdateVisitRequest struct {
	Status            *entities.VisitStatus      `json:"status,omitempty"`
	VitalSigns        *entities.VitalSigns       `json:"vital_signs,omitempty"`
	PhysicalExamNotes *string                    `json:"physical_exam_notes,omitempty"`
	Diagnosis         *string                    `json:"diagnosis,omitempty"`
	Treatment         *string                    `json:"treatment,omitempty"`
	Prescriptions     *[]entities.Prescription   `json:"prescriptions,omitempty"`
	TestResults       *[]entities.TestResult     `json:"test_results,omitempty"`
	FollowUpRequired  *bool                      `json:"follow_up_required,omitempty"`
	FollowUpDate      *time.Time                 `json:"follow_up_date,omitempty"`
	Cost              *float64                   `json:"cost,omitempty"`
	PaymentStatus     *string                    `json:"payment_status,omitempty"`
}

// CreateVaccinationRequest represents a request to create a vaccination record
type CreateVaccinationRequest struct {
	AnimalID         string                     `json:"animal_id" validate:"required"`
	VaccineType      entities.VaccinationType   `json:"vaccine_type" validate:"required"`
	VaccineName      string                     `json:"vaccine_name" validate:"required"`
	Manufacturer     string                     `json:"manufacturer,omitempty"`
	LotNumber        string                     `json:"lot_number,omitempty"`
	DoseNumber       int                        `json:"dose_number" validate:"required,min=1"`
	TotalDoses       int                        `json:"total_doses,omitempty"`
	DateAdministered time.Time                  `json:"date_administered" validate:"required"`
	NextDueDate      *time.Time                 `json:"next_due_date,omitempty"`
	ExpirationDate   *time.Time                 `json:"expiration_date,omitempty"`
	VeterinarianName string                     `json:"veterinarian_name" validate:"required"`
	ClinicName       string                     `json:"clinic_name,omitempty"`
	Cost             float64                    `json:"cost,omitempty"`
	Notes            string                     `json:"notes,omitempty"`
}

// ListVisitsRequest represents a request to list veterinary visits
type ListVisitsRequest struct {
	AnimalID         string     `form:"animal_id"`
	VisitType        string     `form:"visit_type"`
	Status           string     `form:"status"`
	FromDate         *time.Time `form:"from_date"`
	ToDate           *time.Time `form:"to_date"`
	VeterinarianName string     `form:"veterinarian_name"`
	Limit            int64      `form:"limit"`
	Offset           int64      `form:"offset"`
	SortBy           string     `form:"sort_by"`
	SortOrder        string     `form:"sort_order"`
}

// ListVaccinationsRequest represents a request to list vaccinations
type ListVaccinationsRequest struct {
	AnimalID    string     `form:"animal_id"`
	VaccineType string     `form:"vaccine_type"`
	FromDate    *time.Time `form:"from_date"`
	ToDate      *time.Time `form:"to_date"`
	Limit       int64      `form:"limit"`
	Offset      int64      `form:"offset"`
	SortBy      string     `form:"sort_by"`
	SortOrder   string     `form:"sort_order"`
}

// CreateVisit creates a new veterinary visit
func (uc *VeterinaryUseCase) CreateVisit(ctx context.Context, req *CreateVisitRequest, creatorID primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	// Parse animal ID
	animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
	if err != nil {
		return nil, errors.NewBadRequest("invalid animal ID")
	}

	// Verify animal exists
	_, err = uc.animalRepo.FindByID(ctx, animalID)
	if err != nil {
		return nil, err
	}

	visit := &entities.VeterinaryVisit{
		AnimalID:          animalID,
		VisitType:         req.VisitType,
		Status:            entities.VisitStatusCompleted,
		VisitDate:         req.VisitDate,
		VeterinarianName:  req.VeterinarianName,
		ClinicName:        req.ClinicName,
		ClinicAddress:     req.ClinicAddress,
		ClinicPhone:       req.ClinicPhone,
		ChiefComplaint:    req.ChiefComplaint,
		VitalSigns:        req.VitalSigns,
		PhysicalExamNotes: req.PhysicalExamNotes,
		Diagnosis:         req.Diagnosis,
		Treatment:         req.Treatment,
		Prescriptions:     req.Prescriptions,
		TestsOrdered:      req.TestsOrdered,
		VaccinationsGiven: req.VaccinationsGiven,
		FollowUpRequired:  req.FollowUpRequired,
		FollowUpDate:      req.FollowUpDate,
		FollowUpNotes:     req.FollowUpNotes,
		Cost:              req.Cost,
		PaymentStatus:     req.PaymentStatus,
		CreatedBy:         creatorID,
		UpdatedBy:         creatorID,
	}

	if err := uc.visitRepo.Create(ctx, visit); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "veterinary_visit", "", "").
		WithEntityID(visit.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return visit, nil
}

// GetVisitByID retrieves a veterinary visit by ID
func (uc *VeterinaryUseCase) GetVisitByID(ctx context.Context, id primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	return uc.visitRepo.FindByID(ctx, id)
}

// UpdateVisit updates a veterinary visit
func (uc *VeterinaryUseCase) UpdateVisit(ctx context.Context, id primitive.ObjectID, req *UpdateVisitRequest, updaterID primitive.ObjectID) (*entities.VeterinaryVisit, error) {
	visit, err := uc.visitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track changes
	changes := make(map[string]interface{})

	if req.Status != nil {
		changes["status"] = req.Status
		visit.Status = *req.Status
	}
	if req.VitalSigns != nil {
		visit.VitalSigns = *req.VitalSigns
	}
	if req.PhysicalExamNotes != nil {
		visit.PhysicalExamNotes = *req.PhysicalExamNotes
	}
	if req.Diagnosis != nil {
		visit.Diagnosis = *req.Diagnosis
	}
	if req.Treatment != nil {
		visit.Treatment = *req.Treatment
	}
	if req.Prescriptions != nil {
		visit.Prescriptions = *req.Prescriptions
	}
	if req.TestResults != nil {
		visit.TestResults = *req.TestResults
	}
	if req.FollowUpRequired != nil {
		visit.FollowUpRequired = *req.FollowUpRequired
	}
	if req.FollowUpDate != nil {
		visit.FollowUpDate = req.FollowUpDate
	}
	if req.Cost != nil {
		visit.Cost = *req.Cost
	}
	if req.PaymentStatus != nil {
		visit.PaymentStatus = *req.PaymentStatus
	}

	visit.UpdatedBy = updaterID

	if err := uc.visitRepo.Update(ctx, visit); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(updaterID, entities.ActionUpdate, "veterinary_visit", "", "").
		WithEntityID(id).
		WithChanges(changes)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return visit, nil
}

// DeleteVisit deletes a veterinary visit
func (uc *VeterinaryUseCase) DeleteVisit(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	if _, err := uc.visitRepo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.visitRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "veterinary_visit", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListVisits lists veterinary visits with filters
func (uc *VeterinaryUseCase) ListVisits(ctx context.Context, req *ListVisitsRequest) ([]*entities.VeterinaryVisit, int64, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}

	filter := repositories.VeterinaryVisitFilter{
		VisitType:        req.VisitType,
		Status:           req.Status,
		FromDate:         req.FromDate,
		ToDate:           req.ToDate,
		VeterinarianName: req.VeterinarianName,
		Limit:            req.Limit,
		Offset:           req.Offset,
		SortBy:           req.SortBy,
		SortOrder:        req.SortOrder,
	}

	if req.AnimalID != "" {
		animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
		if err == nil {
			filter.AnimalID = &animalID
		}
	}

	return uc.visitRepo.List(ctx, filter)
}

// GetVisitsByAnimalID retrieves all visits for an animal
func (uc *VeterinaryUseCase) GetVisitsByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.VeterinaryVisit, error) {
	return uc.visitRepo.GetByAnimalID(ctx, animalID)
}

// CreateVaccination creates a new vaccination record
func (uc *VeterinaryUseCase) CreateVaccination(ctx context.Context, req *CreateVaccinationRequest, creatorID primitive.ObjectID) (*entities.Vaccination, error) {
	// Parse animal ID
	animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
	if err != nil {
		return nil, errors.NewBadRequest("invalid animal ID")
	}

	// Verify animal exists
	_, err = uc.animalRepo.FindByID(ctx, animalID)
	if err != nil {
		return nil, err
	}

	vaccination := &entities.Vaccination{
		AnimalID:         animalID,
		VaccineType:      req.VaccineType,
		VaccineName:      req.VaccineName,
		Manufacturer:     req.Manufacturer,
		LotNumber:        req.LotNumber,
		DoseNumber:       req.DoseNumber,
		TotalDoses:       req.TotalDoses,
		DateAdministered: req.DateAdministered,
		NextDueDate:      req.NextDueDate,
		ExpirationDate:   req.ExpirationDate,
		VeterinarianName: req.VeterinarianName,
		ClinicName:       req.ClinicName,
		Cost:             req.Cost,
		Notes:            req.Notes,
		CreatedBy:        creatorID,
		UpdatedBy:        creatorID,
	}

	if err := uc.vaccinationRepo.Create(ctx, vaccination); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "vaccination", "", "").
		WithEntityID(vaccination.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return vaccination, nil
}

// GetVaccinationByID retrieves a vaccination by ID
func (uc *VeterinaryUseCase) GetVaccinationByID(ctx context.Context, id primitive.ObjectID) (*entities.Vaccination, error) {
	return uc.vaccinationRepo.FindByID(ctx, id)
}

// DeleteVaccination deletes a vaccination
func (uc *VeterinaryUseCase) DeleteVaccination(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	if _, err := uc.vaccinationRepo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.vaccinationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "vaccination", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListVaccinations lists vaccinations with filters
func (uc *VeterinaryUseCase) ListVaccinations(ctx context.Context, req *ListVaccinationsRequest) ([]*entities.Vaccination, int64, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}

	filter := repositories.VaccinationFilter{
		VaccineType: req.VaccineType,
		FromDate:    req.FromDate,
		ToDate:      req.ToDate,
		Limit:       req.Limit,
		Offset:      req.Offset,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}

	if req.AnimalID != "" {
		animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
		if err == nil {
			filter.AnimalID = &animalID
		}
	}

	return uc.vaccinationRepo.List(ctx, filter)
}

// GetVaccinationsByAnimalID retrieves all vaccinations for an animal
func (uc *VeterinaryUseCase) GetVaccinationsByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Vaccination, error) {
	return uc.vaccinationRepo.GetByAnimalID(ctx, animalID)
}

// GetDueVaccinations retrieves vaccinations due within specified days
func (uc *VeterinaryUseCase) GetDueVaccinations(ctx context.Context, days int) ([]*entities.Vaccination, error) {
	return uc.vaccinationRepo.GetDueVaccinations(ctx, days)
}

// GetUpcomingVisits retrieves upcoming scheduled visits
func (uc *VeterinaryUseCase) GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error) {
	return uc.visitRepo.GetUpcomingVisits(ctx, days)
}
