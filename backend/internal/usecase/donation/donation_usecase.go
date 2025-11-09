package donation

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonationUseCase handles donation business logic
type DonationUseCase struct {
	donationRepo repositories.DonationRepository
	donorRepo    repositories.DonorRepository
	campaignRepo repositories.CampaignRepository
	auditLogRepo repositories.AuditLogRepository
}

// NewDonationUseCase creates a new donation use case
func NewDonationUseCase(
	donationRepo repositories.DonationRepository,
	donorRepo repositories.DonorRepository,
	campaignRepo repositories.CampaignRepository,
	auditLogRepo repositories.AuditLogRepository,
) *DonationUseCase {
	return &DonationUseCase{
		donationRepo: donationRepo,
		donorRepo:    donorRepo,
		campaignRepo: campaignRepo,
		auditLogRepo: auditLogRepo,
	}
}

// CreateDonation creates a new donation
func (uc *DonationUseCase) CreateDonation(ctx context.Context, donation *entities.Donation, userID primitive.ObjectID) error {
	// Validate donation
	if err := uc.validateDonation(donation); err != nil {
		return err
	}

	// Verify donor exists
	donor, err := uc.donorRepo.FindByID(ctx, donation.DonorID)
	if err != nil {
		return errors.NewBadRequest("Invalid donor ID")
	}

	// Cache donor information
	donation.DonorName = donor.GetFullName()
	donation.DonorEmail = donor.Contact.Email

	// Verify campaign if provided
	if donation.CampaignID != nil && !donation.CampaignID.IsZero() {
		campaign, err := uc.campaignRepo.FindByID(ctx, *donation.CampaignID)
		if err != nil {
			return errors.NewBadRequest("Invalid campaign ID")
		}
		donation.CampaignName = campaign.Name.English // Use English name for caching
	}

	// Calculate net amount
	donation.CalculateNetAmount()

	// Set metadata
	donation.ProcessedBy = userID
	donation.CreatedBy = userID
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Create(ctx, donation); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionCreate, "donation", "", "").
		WithEntityID(donation.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ProcessDonation processes a donation (marks as completed and updates statistics)
func (uc *DonationUseCase) ProcessDonation(ctx context.Context, donationID primitive.ObjectID, userID primitive.ObjectID) error {
	donation, err := uc.donationRepo.FindByID(ctx, donationID)
	if err != nil {
		return err
	}

	if donation.Status != entities.DonationStatusPending {
		return errors.NewBadRequest("Only pending donations can be processed")
	}

	// Mark as completed
	donation.Status = entities.DonationStatusCompleted
	now := time.Now()
	donation.PaymentDate = &now
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Update(ctx, donation); err != nil {
		return err
	}

	// Update donor statistics
	donor, err := uc.donorRepo.FindByID(ctx, donation.DonorID)
	if err == nil {
		isNewDonor := donor.DonationCount == 0
		donor.UpdateDonationStats(donation.Amount)
		donor.UpdatedBy = userID
		_ = uc.donorRepo.Update(ctx, donor)

		// Update campaign statistics if applicable
		if donation.CampaignID != nil && !donation.CampaignID.IsZero() {
			_ = uc.campaignRepo.UpdateCampaignStats(ctx, *donation.CampaignID, donation.Amount, isNewDonor)
		}
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donation", "", "").
		WithEntityID(donationID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// RefundDonation refunds a donation
func (uc *DonationUseCase) RefundDonation(ctx context.Context, donationID primitive.ObjectID, userID primitive.ObjectID) error {
	donation, err := uc.donationRepo.FindByID(ctx, donationID)
	if err != nil {
		return err
	}

	if donation.Status != entities.DonationStatusCompleted {
		return errors.NewBadRequest("Only completed donations can be refunded")
	}

	// Mark as refunded
	donation.Status = entities.DonationStatusRefunded
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Update(ctx, donation); err != nil {
		return err
	}

	// Update donor statistics (subtract the donation)
	donor, err := uc.donorRepo.FindByID(ctx, donation.DonorID)
	if err == nil {
		donor.TotalDonated -= donation.Amount
		donor.DonationCount--
		if donor.DonationCount > 0 {
			donor.AverageDonation = donor.TotalDonated / float64(donor.DonationCount)
		} else {
			donor.AverageDonation = 0
		}
		donor.UpdatedBy = userID
		_ = uc.donorRepo.Update(ctx, donor)

		// Update campaign statistics if applicable
		if donation.CampaignID != nil && !donation.CampaignID.IsZero() {
			campaign, err := uc.campaignRepo.FindByID(ctx, *donation.CampaignID)
			if err == nil {
				campaign.CurrentAmount -= donation.Amount
				campaign.DonationCount--
				if campaign.DonationCount > 0 {
					campaign.AverageDonation = campaign.CurrentAmount / float64(campaign.DonationCount)
				} else {
					campaign.AverageDonation = 0
				}
				_ = uc.campaignRepo.Update(ctx, campaign)
			}
		}
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donation", "", "").
		WithEntityID(donationID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetDonation gets a donation by ID
func (uc *DonationUseCase) GetDonation(ctx context.Context, id primitive.ObjectID) (*entities.Donation, error) {
	return uc.donationRepo.FindByID(ctx, id)
}

// UpdateDonation updates a donation
func (uc *DonationUseCase) UpdateDonation(ctx context.Context, donation *entities.Donation, userID primitive.ObjectID) error {
	// Validate donation
	if err := uc.validateDonation(donation); err != nil {
		return err
	}

	// Check if donation exists
	_, err := uc.donationRepo.FindByID(ctx, donation.ID)
	if err != nil {
		return err
	}

	// Recalculate net amount
	donation.CalculateNetAmount()
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Update(ctx, donation); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donation", "", "").
		WithEntityID(donation.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// DeleteDonation deletes a donation
func (uc *DonationUseCase) DeleteDonation(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	donation, err := uc.donationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deletion of pending or failed donations
	if donation.Status == entities.DonationStatusCompleted || donation.Status == entities.DonationStatusRefunded {
		return errors.NewBadRequest("Cannot delete completed or refunded donations. Use refund instead.")
	}

	if err := uc.donationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionDelete, "donation", "", "").
		WithEntityID(id)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ListDonations lists donations with filters
func (uc *DonationUseCase) ListDonations(ctx context.Context, filter *repositories.DonationFilter) ([]*entities.Donation, int64, error) {
	return uc.donationRepo.List(ctx, filter)
}

// GetDonationsByDonor gets all donations for a specific donor
func (uc *DonationUseCase) GetDonationsByDonor(ctx context.Context, donorID primitive.ObjectID) ([]*entities.Donation, error) {
	return uc.donationRepo.GetByDonorID(ctx, donorID)
}

// GetDonationsByCampaign gets all donations for a specific campaign
func (uc *DonationUseCase) GetDonationsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Donation, error) {
	return uc.donationRepo.GetByCampaignID(ctx, campaignID)
}

// GetRecurringDonations gets all active recurring donations
func (uc *DonationUseCase) GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error) {
	return uc.donationRepo.GetRecurringDonations(ctx)
}

// GetPendingThankYous gets donations that need thank you notes
func (uc *DonationUseCase) GetPendingThankYous(ctx context.Context) ([]*entities.Donation, error) {
	return uc.donationRepo.GetPendingThankYous(ctx)
}

// SendThankYou marks a donation as having thank you sent
func (uc *DonationUseCase) SendThankYou(ctx context.Context, donationID primitive.ObjectID, userID primitive.ObjectID) error {
	donation, err := uc.donationRepo.FindByID(ctx, donationID)
	if err != nil {
		return err
	}

	if donation.ThankYouSent {
		return errors.NewBadRequest("Thank you already sent for this donation")
	}

	donation.ThankYouSent = true
	now := time.Now()
	donation.ThankYouSentDate = &now
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Update(ctx, donation); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donation", "", "").
		WithEntityID(donationID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GenerateTaxReceipt generates a tax receipt for a donation
func (uc *DonationUseCase) GenerateTaxReceipt(ctx context.Context, donationID primitive.ObjectID, receiptURL string, userID primitive.ObjectID) error {
	donation, err := uc.donationRepo.FindByID(ctx, donationID)
	if err != nil {
		return err
	}

	if !donation.TaxDeductible {
		return errors.NewBadRequest("This donation is not tax deductible")
	}

	if donation.Status != entities.DonationStatusCompleted {
		return errors.NewBadRequest("Only completed donations can have tax receipts")
	}

	if donation.TaxReceipt.ReceiptNumber != "" {
		return errors.NewBadRequest("Tax receipt already generated for this donation")
	}

	// Generate receipt number
	donation.TaxReceipt.ReceiptNumber = donation.GenerateReceiptNumber()
	donation.TaxReceipt.ReceiptURL = receiptURL
	now := time.Now()
	donation.TaxReceipt.SentDate = &now
	donation.TaxReceipt.SentMethod = "email"
	donation.UpdatedBy = userID

	if err := uc.donationRepo.Update(ctx, donation); err != nil {
		return err
	}

	// Audit log
	auditLog := entities.NewAuditLog(userID, entities.ActionUpdate, "donation", "", "").
		WithEntityID(donationID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// GetDonationStatistics gets donation statistics
func (uc *DonationUseCase) GetDonationStatistics(ctx context.Context) (*repositories.DonationStatistics, error) {
	return uc.donationRepo.GetDonationStatistics(ctx)
}

// GetDonationsByDateRange gets donations within a date range
func (uc *DonationUseCase) GetDonationsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entities.Donation, error) {
	return uc.donationRepo.GetDonationsByDateRange(ctx, startDate, endDate)
}

// validateDonation validates donation data
func (uc *DonationUseCase) validateDonation(donation *entities.Donation) error {
	if donation.DonorID.IsZero() {
		return errors.NewBadRequest("Donor ID is required")
	}

	if donation.Type == "" {
		return errors.NewBadRequest("Donation type is required")
	}

	if donation.Amount <= 0 {
		return errors.NewBadRequest("Donation amount must be greater than zero")
	}

	if donation.Currency == "" {
		donation.Currency = "USD"
	}

	if donation.Payment.Method == "" {
		return errors.NewBadRequest("Payment method is required")
	}

	// Validate in-kind donations
	if donation.Type == entities.DonationTypeInKind {
		if len(donation.InKindItems) == 0 {
			return errors.NewBadRequest("In-kind donations must have at least one item")
		}
	}

	// Validate recurring donations
	if donation.IsRecurring {
		if donation.RecurringInfo == nil {
			return errors.NewBadRequest("Recurring donations must have recurring info")
		}
		if donation.RecurringInfo.Frequency == "" {
			return errors.NewBadRequest("Recurring frequency is required")
		}
	}

	return nil
}
