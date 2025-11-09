package adoption

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdoptionUseCase handles adoption business logic
type AdoptionUseCase struct {
	applicationRepo repositories.AdoptionApplicationRepository
	adoptionRepo    repositories.AdoptionRepository
	animalRepo      repositories.AnimalRepository
	auditLogRepo    repositories.AuditLogRepository
}

// NewAdoptionUseCase creates a new adoption use case
func NewAdoptionUseCase(
	applicationRepo repositories.AdoptionApplicationRepository,
	adoptionRepo repositories.AdoptionRepository,
	animalRepo repositories.AnimalRepository,
	auditLogRepo repositories.AuditLogRepository,
) *AdoptionUseCase {
	return &AdoptionUseCase{
		applicationRepo: applicationRepo,
		adoptionRepo:    adoptionRepo,
		animalRepo:      animalRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// CreateApplicationRequest represents a request to create an adoption application
type CreateApplicationRequest struct {
	AnimalID              string                           `json:"animal_id" validate:"required"`
	Applicant             entities.ApplicantInfo           `json:"applicant" validate:"required"`
	Address               entities.AddressInfo             `json:"address" validate:"required"`
	Housing               entities.HousingInfo             `json:"housing" validate:"required"`
	HouseholdSize         int                              `json:"household_size" validate:"required,min=1"`
	HouseholdMembers      []entities.HouseholdMember       `json:"household_members,omitempty"`
	HasChildren           bool                             `json:"has_children"`
	ChildrenAges          []int                            `json:"children_ages,omitempty"`
	CurrentPets           []entities.CurrentPet            `json:"current_pets,omitempty"`
	PreviousPets          string                           `json:"previous_pets,omitempty"`
	PetExperience         string                           `json:"pet_experience,omitempty"`
	SurrenderedPets       bool                             `json:"surrendered_pets"`
	SurrenderReason       string                           `json:"surrender_reason,omitempty"`
	ReasonForAdoption     string                           `json:"reason_for_adoption" validate:"required"`
	PetLocation           string                           `json:"pet_location" validate:"required"`
	AloneTime             int                              `json:"alone_time"`
	ActivityLevel         string                           `json:"activity_level,omitempty"`
	PreparedFor           []string                         `json:"prepared_for,omitempty"`
	References            []entities.Reference             `json:"references,omitempty"`
	HasVeterinarian       bool                             `json:"has_veterinarian"`
	VetName               string                           `json:"vet_name,omitempty"`
	VetPhone              string                           `json:"vet_phone,omitempty"`
	VetAddress            string                           `json:"vet_address,omitempty"`
	AgreesToHomeVisit     bool                             `json:"agrees_to_home_visit"`
	AgreesToFollowUp      bool                             `json:"agrees_to_follow_up"`
	AgreesToReturnPolicy  bool                             `json:"agrees_to_return_policy"`
	UnderstandsCommitment bool                             `json:"understands_commitment"`
	AdditionalInfo        string                           `json:"additional_info,omitempty"`
}

// UpdateApplicationRequest represents a request to update an adoption application
type UpdateApplicationRequest struct {
	Status          *entities.ApplicationStatus `json:"status,omitempty"`
	ReviewNotes     *string                     `json:"review_notes,omitempty"`
	RejectionReason *string                     `json:"rejection_reason,omitempty"`
	HomeVisitDate   *time.Time                  `json:"home_visit_date,omitempty"`
	HomeVisitNotes  *string                     `json:"home_visit_notes,omitempty"`
	InterviewDate   *time.Time                  `json:"interview_date,omitempty"`
	InterviewNotes  *string                     `json:"interview_notes,omitempty"`
}

// CreateAdoptionRequest represents a request to create an adoption record
type CreateAdoptionRequest struct {
	ApplicationID      string    `json:"application_id" validate:"required"`
	AdoptionFee        float64   `json:"adoption_fee" validate:"required,min=0"`
	TrialPeriod        bool      `json:"trial_period"`
	TrialPeriodDays    int       `json:"trial_period_days,omitempty"`
	PaymentStatus      string    `json:"payment_status"`
	AmountPaid         float64   `json:"amount_paid" validate:"min=0"`
	PaymentMethod      string    `json:"payment_method,omitempty"`
	ContractURL        string    `json:"contract_url,omitempty"`
	ScheduleFollowUps  bool      `json:"schedule_follow_ups"`
	FollowUpIntervals  []int     `json:"follow_up_intervals,omitempty"` // Days: [7, 30, 90]
}

// UpdateAdoptionRequest represents a request to update an adoption
type UpdateAdoptionRequest struct {
	Status           *entities.AdoptionStatus  `json:"status,omitempty"`
	PaymentStatus    *entities.PaymentStatus   `json:"payment_status,omitempty"`
	AmountPaid       *float64                  `json:"amount_paid,omitempty"`
	PaymentDate      *time.Time                `json:"payment_date,omitempty"`
	PaymentMethod    *string                   `json:"payment_method,omitempty"`
	ContractSigned   *time.Time                `json:"contract_signed,omitempty"`
	ReturnDate       *time.Time                `json:"return_date,omitempty"`
	ReturnReason     *string                   `json:"return_reason,omitempty"`
	ReturnNotes      *string                   `json:"return_notes,omitempty"`
	Notes            *string                   `json:"notes,omitempty"`
}

// ListApplicationsRequest represents a request to list adoption applications
type ListApplicationsRequest struct {
	AnimalID       string     `form:"animal_id"`
	Status         string     `form:"status"`
	ApplicantEmail string     `form:"applicant_email"`
	ApplicantName  string     `form:"applicant_name"`
	FromDate       *time.Time `form:"from_date"`
	ToDate         *time.Time `form:"to_date"`
	Limit          int64      `form:"limit"`
	Offset         int64      `form:"offset"`
	SortBy         string     `form:"sort_by"`
	SortOrder      string     `form:"sort_order"`
}

// ListAdoptionsRequest represents a request to list adoptions
type ListAdoptionsRequest struct {
	AnimalID      string     `form:"animal_id"`
	AdopterID     string     `form:"adopter_id"`
	ApplicationID string     `form:"application_id"`
	Status        string     `form:"status"`
	PaymentStatus string     `form:"payment_status"`
	FromDate      *time.Time `form:"from_date"`
	ToDate        *time.Time `form:"to_date"`
	TrialPeriod   *bool      `form:"trial_period"`
	Limit         int64      `form:"limit"`
	Offset        int64      `form:"offset"`
	SortBy        string     `form:"sort_by"`
	SortOrder     string     `form:"sort_order"`
}

// CreateApplication creates a new adoption application
func (uc *AdoptionUseCase) CreateApplication(ctx context.Context, req *CreateApplicationRequest, creatorID primitive.ObjectID) (*entities.AdoptionApplication, error) {
	// Parse animal ID
	animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
	if err != nil {
		return nil, errors.NewBadRequest("invalid animal ID")
	}

	// Verify animal exists and is available
	animal, err := uc.animalRepo.FindByID(ctx, animalID)
	if err != nil {
		return nil, err
	}

	if animal.Status != entities.AnimalStatusAvailable {
		return nil, errors.NewBadRequest("animal is not available for adoption")
	}

	application := entities.NewAdoptionApplication(animalID, req.Applicant, creatorID)
	application.Address = req.Address
	application.Housing = req.Housing
	application.HouseholdSize = req.HouseholdSize
	application.HouseholdMembers = req.HouseholdMembers
	application.HasChildren = req.HasChildren
	application.ChildrenAges = req.ChildrenAges
	application.CurrentPets = req.CurrentPets
	application.PreviousPets = req.PreviousPets
	application.PetExperience = req.PetExperience
	application.SurrenderedPets = req.SurrenderedPets
	application.SurrenderReason = req.SurrenderReason
	application.ReasonForAdoption = req.ReasonForAdoption
	application.PetLocation = req.PetLocation
	application.AloneTime = req.AloneTime
	application.ActivityLevel = req.ActivityLevel
	application.PreparedFor = req.PreparedFor
	application.References = req.References
	application.HasVeterinarian = req.HasVeterinarian
	application.VetName = req.VetName
	application.VetPhone = req.VetPhone
	application.VetAddress = req.VetAddress
	application.AgreesToHomeVisit = req.AgreesToHomeVisit
	application.AgreesToFollowUp = req.AgreesToFollowUp
	application.AgreesToReturnPolicy = req.AgreesToReturnPolicy
	application.UnderstandsCommitment = req.UnderstandsCommitment
	application.AdditionalInfo = req.AdditionalInfo

	if err := uc.applicationRepo.Create(ctx, application); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "adoption_application", "", "").
		WithEntityID(application.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return application, nil
}

// GetApplicationByID retrieves an adoption application by ID
func (uc *AdoptionUseCase) GetApplicationByID(ctx context.Context, id primitive.ObjectID) (*entities.AdoptionApplication, error) {
	return uc.applicationRepo.FindByID(ctx, id)
}

// UpdateApplication updates an adoption application
func (uc *AdoptionUseCase) UpdateApplication(ctx context.Context, id primitive.ObjectID, req *UpdateApplicationRequest, updaterID primitive.ObjectID) (*entities.AdoptionApplication, error) {
	application, err := uc.applicationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track changes
	changes := make(map[string]interface{})
	now := time.Now()

	if req.Status != nil {
		changes["status"] = *req.Status
		application.Status = *req.Status

		// Update status-specific dates
		switch *req.Status {
		case entities.ApplicationStatusUnderReview:
			application.ReviewDate = &now
		case entities.ApplicationStatusApproved:
			application.ApprovalDate = &now
			application.ReviewDate = &now
		case entities.ApplicationStatusRejected:
			application.RejectionDate = &now
			application.ReviewDate = &now
		}
	}

	if req.ReviewNotes != nil {
		application.ReviewNotes = *req.ReviewNotes
		application.ReviewedBy = &updaterID
	}

	if req.RejectionReason != nil {
		application.RejectionReason = *req.RejectionReason
	}

	if req.HomeVisitDate != nil {
		application.HomeVisitDate = req.HomeVisitDate
	}

	if req.HomeVisitNotes != nil {
		application.HomeVisitNotes = *req.HomeVisitNotes
	}

	if req.InterviewDate != nil {
		application.InterviewDate = req.InterviewDate
	}

	if req.InterviewNotes != nil {
		application.InterviewNotes = *req.InterviewNotes
	}

	application.UpdatedBy = updaterID

	if err := uc.applicationRepo.Update(ctx, application); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(updaterID, entities.ActionUpdate, "adoption_application", "", "").
		WithEntityID(id).
		WithChanges(changes)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return application, nil
}

// DeleteApplication deletes an adoption application
func (uc *AdoptionUseCase) DeleteApplication(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	if _, err := uc.applicationRepo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.applicationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "adoption_application", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListApplications lists adoption applications with filters
func (uc *AdoptionUseCase) ListApplications(ctx context.Context, req *ListApplicationsRequest) ([]*entities.AdoptionApplication, int64, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}

	filter := repositories.AdoptionApplicationFilter{
		Status:         req.Status,
		ApplicantEmail: req.ApplicantEmail,
		ApplicantName:  req.ApplicantName,
		FromDate:       req.FromDate,
		ToDate:         req.ToDate,
		Limit:          req.Limit,
		Offset:         req.Offset,
		SortBy:         req.SortBy,
		SortOrder:      req.SortOrder,
	}

	if req.AnimalID != "" {
		animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
		if err == nil {
			filter.AnimalID = &animalID
		}
	}

	return uc.applicationRepo.List(ctx, filter)
}

// GetApplicationsByAnimalID retrieves all applications for an animal
func (uc *AdoptionUseCase) GetApplicationsByAnimalID(ctx context.Context, animalID primitive.ObjectID) ([]*entities.AdoptionApplication, error) {
	return uc.applicationRepo.GetByAnimalID(ctx, animalID)
}

// GetPendingApplications retrieves all pending applications
func (uc *AdoptionUseCase) GetPendingApplications(ctx context.Context) ([]*entities.AdoptionApplication, error) {
	return uc.applicationRepo.GetPendingApplications(ctx)
}

// CreateAdoption creates an adoption from an approved application
func (uc *AdoptionUseCase) CreateAdoption(ctx context.Context, req *CreateAdoptionRequest, creatorID primitive.ObjectID) (*entities.Adoption, error) {
	// Parse application ID
	applicationID, err := primitive.ObjectIDFromHex(req.ApplicationID)
	if err != nil {
		return nil, errors.NewBadRequest("invalid application ID")
	}

	// Get and verify application
	application, err := uc.applicationRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if !application.IsApproved() {
		return nil, errors.NewBadRequest("application must be approved before creating adoption")
	}

	// Check if adoption already exists for this application
	existingAdoption, err := uc.adoptionRepo.GetByApplicationID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	if existingAdoption != nil {
		return nil, errors.NewBadRequest("adoption already exists for this application")
	}

	// Create adoption
	adoption := entities.NewAdoption(
		applicationID,
		application.AnimalID,
		creatorID, // In real scenario, this should be the adopter's user ID
		req.AdoptionFee,
		creatorID,
	)

	// Set payment info
	paymentStatus := entities.PaymentStatusPending
	if req.PaymentStatus != "" {
		paymentStatus = entities.PaymentStatus(req.PaymentStatus)
	}
	adoption.PaymentStatus = paymentStatus
	adoption.AmountPaid = req.AmountPaid
	adoption.PaymentMethod = req.PaymentMethod

	// Set trial period
	adoption.TrialPeriod = req.TrialPeriod
	if req.TrialPeriod && req.TrialPeriodDays > 0 {
		trialEnd := time.Now().AddDate(0, 0, req.TrialPeriodDays)
		adoption.TrialEndDate = &trialEnd
	}

	// Set contract
	if req.ContractURL != "" {
		adoption.Contract.ContractURL = req.ContractURL
	}

	// Schedule follow-ups
	if req.ScheduleFollowUps && len(req.FollowUpIntervals) > 0 {
		for _, days := range req.FollowUpIntervals {
			followUpDate := time.Now().AddDate(0, 0, days)
			adoption.AddFollowUp(entities.FollowUpSchedule{
				ScheduledDate: followUpDate,
				Type:          "visit",
			})
		}
	}

	if err := uc.adoptionRepo.Create(ctx, adoption); err != nil {
		return nil, err
	}

	// Update application status to completed
	completedStatus := entities.ApplicationStatusCompleted
	application.Status = completedStatus
	_ = uc.applicationRepo.Update(ctx, application)

	// Update animal status to adopted
	animal, _ := uc.animalRepo.FindByID(ctx, application.AnimalID)
	if animal != nil {
		animal.Status = entities.AnimalStatusAdopted
		_ = uc.animalRepo.Update(ctx, animal)
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "adoption", "", "").
		WithEntityID(adoption.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return adoption, nil
}

// GetAdoptionByID retrieves an adoption by ID
func (uc *AdoptionUseCase) GetAdoptionByID(ctx context.Context, id primitive.ObjectID) (*entities.Adoption, error) {
	return uc.adoptionRepo.FindByID(ctx, id)
}

// UpdateAdoption updates an adoption
func (uc *AdoptionUseCase) UpdateAdoption(ctx context.Context, id primitive.ObjectID, req *UpdateAdoptionRequest, updaterID primitive.ObjectID) (*entities.Adoption, error) {
	adoption, err := uc.adoptionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track changes
	changes := make(map[string]interface{})

	if req.Status != nil {
		changes["status"] = *req.Status
		adoption.Status = *req.Status

		// Handle animal status update if adoption is returned
		if *req.Status == entities.AdoptionStatusReturned {
			animal, _ := uc.animalRepo.FindByID(ctx, adoption.AnimalID)
			if animal != nil {
				animal.Status = entities.AnimalStatusAvailable
				_ = uc.animalRepo.Update(ctx, animal)
			}
		}
	}

	if req.PaymentStatus != nil {
		adoption.PaymentStatus = *req.PaymentStatus
	}

	if req.AmountPaid != nil {
		adoption.AmountPaid = *req.AmountPaid
	}

	if req.PaymentDate != nil {
		adoption.PaymentDate = req.PaymentDate
	}

	if req.PaymentMethod != nil {
		adoption.PaymentMethod = *req.PaymentMethod
	}

	if req.ContractSigned != nil {
		adoption.Contract.SignedDate = req.ContractSigned
	}

	if req.ReturnDate != nil {
		adoption.ReturnDate = req.ReturnDate
	}

	if req.ReturnReason != nil {
		adoption.ReturnReason = *req.ReturnReason
	}

	if req.ReturnNotes != nil {
		adoption.ReturnNotes = *req.ReturnNotes
	}

	if req.Notes != nil {
		adoption.Notes = *req.Notes
	}

	adoption.UpdatedBy = updaterID

	if err := uc.adoptionRepo.Update(ctx, adoption); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(updaterID, entities.ActionUpdate, "adoption", "", "").
		WithEntityID(id).
		WithChanges(changes)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return adoption, nil
}

// DeleteAdoption deletes an adoption
func (uc *AdoptionUseCase) DeleteAdoption(ctx context.Context, id primitive.ObjectID, deleterID primitive.ObjectID) error {
	if _, err := uc.adoptionRepo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := uc.adoptionRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "adoption", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListAdoptions lists adoptions with filters
func (uc *AdoptionUseCase) ListAdoptions(ctx context.Context, req *ListAdoptionsRequest) ([]*entities.Adoption, int64, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}

	filter := repositories.AdoptionFilter{
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		FromDate:      req.FromDate,
		ToDate:        req.ToDate,
		TrialPeriod:   req.TrialPeriod,
		Limit:         req.Limit,
		Offset:        req.Offset,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
	}

	if req.AnimalID != "" {
		animalID, err := primitive.ObjectIDFromHex(req.AnimalID)
		if err == nil {
			filter.AnimalID = &animalID
		}
	}

	if req.AdopterID != "" {
		adopterID, err := primitive.ObjectIDFromHex(req.AdopterID)
		if err == nil {
			filter.AdopterID = &adopterID
		}
	}

	if req.ApplicationID != "" {
		applicationID, err := primitive.ObjectIDFromHex(req.ApplicationID)
		if err == nil {
			filter.ApplicationID = &applicationID
		}
	}

	return uc.adoptionRepo.List(ctx, filter)
}

// GetAdoptionsByAnimalID retrieves adoption for an animal
func (uc *AdoptionUseCase) GetAdoptionByAnimalID(ctx context.Context, animalID primitive.ObjectID) (*entities.Adoption, error) {
	return uc.adoptionRepo.GetByAnimalID(ctx, animalID)
}

// GetPendingFollowUps retrieves adoptions with pending follow-ups
func (uc *AdoptionUseCase) GetPendingFollowUps(ctx context.Context, days int) ([]*entities.Adoption, error) {
	return uc.adoptionRepo.GetPendingFollowUps(ctx, days)
}

// GetAdoptionStatistics retrieves adoption statistics
func (uc *AdoptionUseCase) GetAdoptionStatistics(ctx context.Context) (*repositories.AdoptionStatistics, error) {
	return uc.adoptionRepo.GetAdoptionStatistics(ctx)
}
