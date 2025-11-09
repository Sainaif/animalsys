package transfer

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransferUseCase struct {
	transferRepo repositories.TransferRepository
	animalRepo   repositories.AnimalRepository
	partnerRepo  repositories.PartnerRepository
	auditLogRepo repositories.AuditLogRepository
}

func NewTransferUseCase(
	transferRepo repositories.TransferRepository,
	animalRepo repositories.AnimalRepository,
	partnerRepo repositories.PartnerRepository,
	auditLogRepo repositories.AuditLogRepository,
) *TransferUseCase {
	return &TransferUseCase{
		transferRepo: transferRepo,
		animalRepo:   animalRepo,
		partnerRepo:  partnerRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateTransfer creates a new transfer
func (uc *TransferUseCase) CreateTransfer(ctx context.Context, transfer *entities.Transfer, userID primitive.ObjectID) error {
	// Validate required fields
	if transfer.AnimalID.IsZero() {
		return errors.NewBadRequest("Animal ID is required")
	}

	if transfer.PartnerID.IsZero() {
		return errors.NewBadRequest("Partner ID is required")
	}

	if transfer.Direction == "" {
		return errors.NewBadRequest("Transfer direction is required")
	}

	if transfer.Reason == "" {
		return errors.NewBadRequest("Transfer reason is required")
	}

	// Validate direction
	if transfer.Direction != entities.TransferDirectionIncoming && transfer.Direction != entities.TransferDirectionOutgoing {
		return errors.NewBadRequest("Invalid transfer direction")
	}

	// Check if animal exists
	animal, err := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.NewBadRequest("Animal not found")
		}
		return err
	}

	// Check if partner exists
	partner, err := uc.partnerRepo.FindByID(ctx, transfer.PartnerID)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.NewBadRequest("Partner not found")
		}
		return err
	}

	// Check if partner is active
	if partner.Status != entities.PartnerStatusActive {
		return errors.NewBadRequest("Partner is not active")
	}

	// Check capacity for incoming transfers
	if transfer.Direction == entities.TransferDirectionIncoming && !partner.HasCapacity() {
		return errors.NewBadRequest("Partner has no capacity for incoming transfers")
	}

	// Set metadata
	transfer.RequestedBy = userID
	now := time.Now()
	transfer.RequestedDate = now
	transfer.CreatedAt = now
	transfer.UpdatedAt = now

	// Set default status
	if transfer.Status == "" {
		transfer.Status = entities.TransferStatusPending
	}

	// Initialize arrays if nil
	if transfer.Documents == nil {
		transfer.Documents = []string{}
	}

	if err := uc.transferRepo.Create(ctx, transfer); err != nil {
		return err
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "transfer", animal.Name.English+" to/from "+partner.Name, "").
			WithEntityID(transfer.ID))

	return nil
}

// GetTransferByID retrieves a transfer by ID
func (uc *TransferUseCase) GetTransferByID(ctx context.Context, id primitive.ObjectID) (*entities.Transfer, error) {
	transfer, err := uc.transferRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

// UpdateTransfer updates a transfer
func (uc *TransferUseCase) UpdateTransfer(ctx context.Context, transfer *entities.Transfer, userID primitive.ObjectID) error {
	// Check if transfer exists
	existing, err := uc.transferRepo.FindByID(ctx, transfer.ID)
	if err != nil {
		return err
	}

	// Don't allow updates to completed or cancelled transfers
	if existing.Status == entities.TransferStatusCompleted {
		return errors.NewBadRequest("Cannot update a completed transfer")
	}
	if existing.Status == entities.TransferStatusCancelled {
		return errors.NewBadRequest("Cannot update a cancelled transfer")
	}

	// Preserve creation info and status change tracking
	transfer.RequestedBy = existing.RequestedBy
	transfer.RequestedDate = existing.RequestedDate
	transfer.ApprovedBy = existing.ApprovedBy
	transfer.ApprovedDate = existing.ApprovedDate
	transfer.CompletedBy = existing.CompletedBy
	transfer.CompletedDate = existing.CompletedDate
	transfer.CancelledBy = existing.CancelledBy
	transfer.CancelledDate = existing.CancelledDate
	transfer.CreatedAt = existing.CreatedAt
	transfer.UpdatedAt = time.Now()

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "").
			WithEntityID(transfer.ID))

	return nil
}

// DeleteTransfer deletes a transfer
func (uc *TransferUseCase) DeleteTransfer(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	// Check if transfer exists
	transfer, err := uc.transferRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deletion of pending transfers
	if transfer.Status != entities.TransferStatusPending {
		return errors.NewBadRequest("Only pending transfers can be deleted")
	}

	if err := uc.transferRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "transfer", animalName+" transfer", "").
			WithEntityID(id))

	return nil
}

// ListTransfers lists transfers with filtering and pagination
func (uc *TransferUseCase) ListTransfers(ctx context.Context, filter *repositories.TransferFilter) ([]*entities.Transfer, int64, error) {
	return uc.transferRepo.List(ctx, filter)
}

// GetTransfersByAnimal gets transfers for a specific animal
func (uc *TransferUseCase) GetTransfersByAnimal(ctx context.Context, animalID primitive.ObjectID) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetByAnimal(ctx, animalID)
}

// GetTransfersByPartner gets transfers for a specific partner
func (uc *TransferUseCase) GetTransfersByPartner(ctx context.Context, partnerID primitive.ObjectID) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetByPartner(ctx, partnerID)
}

// GetTransfersByStatus gets transfers by status
func (uc *TransferUseCase) GetTransfersByStatus(ctx context.Context, status entities.TransferStatus) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetByStatus(ctx, status)
}

// GetPendingTransfers gets all pending transfers
func (uc *TransferUseCase) GetPendingTransfers(ctx context.Context) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetPendingTransfers(ctx)
}

// GetUpcomingTransfers gets transfers scheduled within the next N days
func (uc *TransferUseCase) GetUpcomingTransfers(ctx context.Context, days int) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetUpcomingTransfers(ctx, days)
}

// GetOverdueTransfers gets transfers that are overdue
func (uc *TransferUseCase) GetOverdueTransfers(ctx context.Context) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetOverdueTransfers(ctx)
}

// GetRequiringFollowUp gets transfers requiring follow-up
func (uc *TransferUseCase) GetRequiringFollowUp(ctx context.Context) ([]*entities.Transfer, error) {
	return uc.transferRepo.GetRequiringFollowUp(ctx)
}

// GetTransferStatistics gets transfer statistics
func (uc *TransferUseCase) GetTransferStatistics(ctx context.Context) (*repositories.TransferStatistics, error) {
	return uc.transferRepo.GetTransferStatistics(ctx)
}

// ApproveTransfer approves a transfer
func (uc *TransferUseCase) ApproveTransfer(ctx context.Context, transferID, userID primitive.ObjectID) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != entities.TransferStatusPending {
		return errors.NewBadRequest("Only pending transfers can be approved")
	}

	transfer.Approve(userID)

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "approved transfer").
			WithEntityID(transferID))

	return nil
}

// RejectTransfer rejects a transfer
func (uc *TransferUseCase) RejectTransfer(ctx context.Context, transferID, userID primitive.ObjectID, reason string) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != entities.TransferStatusPending {
		return errors.NewBadRequest("Only pending transfers can be rejected")
	}

	if reason == "" {
		return errors.NewBadRequest("Rejection reason is required")
	}

	transfer.Reject(reason)

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "rejected transfer: "+reason).
			WithEntityID(transferID))

	return nil
}

// StartTransit marks the transfer as in transit
func (uc *TransferUseCase) StartTransit(ctx context.Context, transferID, userID primitive.ObjectID) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != entities.TransferStatusApproved {
		return errors.NewBadRequest("Only approved transfers can be started")
	}

	transfer.StartTransit()

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "started transit").
			WithEntityID(transferID))

	return nil
}

// CompleteTransfer completes a transfer
func (uc *TransferUseCase) CompleteTransfer(ctx context.Context, transferID, userID primitive.ObjectID) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status == entities.TransferStatusCompleted {
		return errors.NewBadRequest("Transfer is already completed")
	}

	if transfer.Status == entities.TransferStatusCancelled || transfer.Status == entities.TransferStatusRejected {
		return errors.NewBadRequest("Cannot complete a cancelled or rejected transfer")
	}

	transfer.Complete(userID)

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Update partner statistics
	partner, err := uc.partnerRepo.FindByID(ctx, transfer.PartnerID)
	if err == nil {
		if transfer.Direction == entities.TransferDirectionIncoming {
			partner.IncrementTransfersIn()
		} else {
			partner.IncrementTransfersOut()
		}
		_ = uc.partnerRepo.Update(ctx, partner)
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "completed transfer").
			WithEntityID(transferID))

	return nil
}

// CancelTransfer cancels a transfer
func (uc *TransferUseCase) CancelTransfer(ctx context.Context, transferID, userID primitive.ObjectID, reason string) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status == entities.TransferStatusCompleted {
		return errors.NewBadRequest("Cannot cancel a completed transfer")
	}

	if transfer.Status == entities.TransferStatusCancelled {
		return errors.NewBadRequest("Transfer is already cancelled")
	}

	if reason == "" {
		return errors.NewBadRequest("Cancellation reason is required")
	}

	transfer.Cancel(userID, reason)

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "cancelled transfer: "+reason).
			WithEntityID(transferID))

	return nil
}

// ScheduleTransfer schedules a transfer for a specific date
func (uc *TransferUseCase) ScheduleTransfer(ctx context.Context, transferID, userID primitive.ObjectID, scheduledDate time.Time) error {
	transfer, err := uc.transferRepo.FindByID(ctx, transferID)
	if err != nil {
		return err
	}

	if transfer.Status != entities.TransferStatusApproved && transfer.Status != entities.TransferStatusPending {
		return errors.NewBadRequest("Only pending or approved transfers can be scheduled")
	}

	// Validate scheduled date is in the future
	if scheduledDate.Before(time.Now()) {
		return errors.NewBadRequest("Scheduled date must be in the future")
	}

	transfer.Schedule(scheduledDate)

	if err := uc.transferRepo.Update(ctx, transfer); err != nil {
		return err
	}

	// Get animal name for audit log
	animal, _ := uc.animalRepo.FindByID(ctx, transfer.AnimalID)
	animalName := "Animal"
	if animal != nil {
		animalName = animal.Name.English
	}

	// Create audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "transfer", animalName+" transfer", "scheduled transfer").
			WithEntityID(transferID))

	return nil
}
