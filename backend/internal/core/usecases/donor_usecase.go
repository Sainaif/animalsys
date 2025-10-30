package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type DonorUseCase struct {
	donorRepo    interfaces.DonorRepository
	donationRepo interfaces.DonationRepository
	auditRepo    interfaces.AuditLogRepository
}

func NewDonorUseCase(
	donorRepo interfaces.DonorRepository,
	donationRepo interfaces.DonationRepository,
	auditRepo interfaces.AuditLogRepository,
) *DonorUseCase {
	return &DonorUseCase{
		donorRepo:    donorRepo,
		donationRepo: donationRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *DonorUseCase) CreateDonor(ctx context.Context, req *entities.DonorCreateRequest, createdBy string) (*entities.Donor, error) {
	donor := entities.NewDonor(
		req.Name,
		req.Email,
		req.Phone,
		req.Type,
	)

	donor.Address = req.Address
	donor.TaxID = req.TaxID
	donor.PreferredContact = req.PreferredContact
	donor.Notes = req.Notes
	donor.CreatedBy = createdBy

	if err := uc.donorRepo.Create(ctx, donor); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "donor", donor.ID.Hex(), "Donor created")
	uc.auditRepo.Create(ctx, auditLog)

	return donor, nil
}

func (uc *DonorUseCase) GetDonorByID(ctx context.Context, id string) (*entities.Donor, error) {
	return uc.donorRepo.GetByID(ctx, id)
}

func (uc *DonorUseCase) ListDonors(ctx context.Context, filter *entities.DonorFilter) ([]*entities.Donor, int64, error) {
	return uc.donorRepo.List(ctx, filter)
}

func (uc *DonorUseCase) UpdateDonor(ctx context.Context, id string, req *entities.DonorUpdateRequest, updatedBy string) (*entities.Donor, error) {
	donor, err := uc.donorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		donor.Name = req.Name
	}
	if req.Email != "" {
		donor.Email = req.Email
	}
	if req.Phone != "" {
		donor.Phone = req.Phone
	}
	if req.Address != "" {
		donor.Address = req.Address
	}
	if req.TaxID != "" {
		donor.TaxID = req.TaxID
	}
	if req.PreferredContact != "" {
		donor.PreferredContact = req.PreferredContact
	}
	if req.Notes != "" {
		donor.Notes = req.Notes
	}
	donor.UpdatedBy = updatedBy

	if err := uc.donorRepo.Update(ctx, id, donor); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "donor", id, "Donor updated")
	uc.auditRepo.Create(ctx, auditLog)

	return donor, nil
}

func (uc *DonorUseCase) DeleteDonor(ctx context.Context, id string, deletedBy string) error {
	if err := uc.donorRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "donor", id, "Donor deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *DonorUseCase) RecordDonation(ctx context.Context, req *entities.DonationCreateRequest, recordedBy string) (*entities.Donation, error) {
	// Verify donor exists
	donor, err := uc.donorRepo.GetByID(ctx, req.DonorID)
	if err != nil {
		return nil, err
	}

	donation := entities.NewDonation(
		req.DonorID,
		donor.Name,
		req.Amount,
		req.Type,
		req.Date,
	)

	donation.Method = req.Method
	donation.Reference = req.Reference
	donation.Campaign = req.Campaign
	donation.IsRecurring = req.IsRecurring
	donation.RecurringFrequency = req.RecurringFrequency
	donation.TaxDeductible = req.TaxDeductible
	donation.Anonymous = req.Anonymous
	donation.Notes = req.Notes
	donation.RecordedBy = recordedBy

	if err := uc.donationRepo.Create(ctx, donation); err != nil {
		return nil, err
	}

	// Update donor statistics
	uc.donorRepo.UpdateTotalDonated(ctx, req.DonorID, req.Amount)
	uc.donorRepo.UpdateLastDonationDate(ctx, req.DonorID, req.Date)

	// Audit
	auditLog := entities.NewAuditLog(recordedBy, "", "", entities.ActionCreate, "donation", donation.ID.Hex(), "Donation recorded")
	uc.auditRepo.Create(ctx, auditLog)

	return donation, nil
}

func (uc *DonorUseCase) GetDonationByID(ctx context.Context, id string) (*entities.Donation, error) {
	return uc.donationRepo.GetByID(ctx, id)
}

func (uc *DonorUseCase) ListDonations(ctx context.Context, filter *entities.DonationFilter) ([]*entities.Donation, int64, error) {
	return uc.donationRepo.List(ctx, filter)
}

func (uc *DonorUseCase) UpdateDonation(ctx context.Context, id string, req *entities.DonationUpdateRequest, updatedBy string) (*entities.Donation, error) {
	donation, err := uc.donationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldAmount := donation.Amount

	// Update fields
	if req.Amount > 0 {
		donation.Amount = req.Amount
	}
	if req.Type != "" {
		donation.Type = req.Type
	}
	if req.Method != "" {
		donation.Method = req.Method
	}
	if req.Reference != "" {
		donation.Reference = req.Reference
	}
	if req.Campaign != "" {
		donation.Campaign = req.Campaign
	}
	if req.IsRecurring != nil {
		donation.IsRecurring = *req.IsRecurring
	}
	if req.RecurringFrequency != "" {
		donation.RecurringFrequency = req.RecurringFrequency
	}
	if req.TaxDeductible != nil {
		donation.TaxDeductible = *req.TaxDeductible
	}
	if req.Anonymous != nil {
		donation.Anonymous = *req.Anonymous
	}
	if req.Notes != "" {
		donation.Notes = req.Notes
	}
	donation.UpdatedBy = updatedBy

	if err := uc.donationRepo.Update(ctx, id, donation); err != nil {
		return nil, err
	}

	// Update donor statistics if amount changed
	if oldAmount != donation.Amount {
		amountDiff := donation.Amount - oldAmount
		uc.donorRepo.UpdateTotalDonated(ctx, donation.DonorID, amountDiff)
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "donation", id, "Donation updated")
	uc.auditRepo.Create(ctx, auditLog)

	return donation, nil
}

func (uc *DonorUseCase) DeleteDonation(ctx context.Context, id string, deletedBy string) error {
	// Get donation to update donor statistics
	donation, err := uc.donationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.donationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Update donor statistics
	uc.donorRepo.UpdateTotalDonated(ctx, donation.DonorID, -donation.Amount)

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "donation", id, "Donation deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *DonorUseCase) GetDonorDonations(ctx context.Context, donorID string, limit, offset int) ([]*entities.Donation, int64, error) {
	return uc.donationRepo.GetByDonorID(ctx, donorID, limit, offset)
}

func (uc *DonorUseCase) GetDonationsByDateRange(ctx context.Context, startDate, endDate string, limit, offset int) ([]*entities.Donation, int64, error) {
	return uc.donationRepo.GetByDateRange(ctx, startDate, endDate, limit, offset)
}

func (uc *DonorUseCase) GetTopDonors(ctx context.Context, limit int) ([]*entities.Donor, error) {
	return uc.donorRepo.GetTopDonors(ctx, limit)
}

func (uc *DonorUseCase) GetDonationStatistics(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	donations, _, err := uc.donationRepo.GetByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return nil, err
	}

	totalAmount := 0.0
	byType := make(map[string]float64)
	byMethod := make(map[string]float64)
	recurringCount := 0
	anonymousCount := 0

	for _, d := range donations {
		totalAmount += d.Amount
		byType[string(d.Type)] += d.Amount
		byMethod[d.Method] += d.Amount
		if d.IsRecurring {
			recurringCount++
		}
		if d.Anonymous {
			anonymousCount++
		}
	}

	stats := map[string]interface{}{
		"period": map[string]string{
			"start": startDate,
			"end":   endDate,
		},
		"total_amount":     totalAmount,
		"total_count":      len(donations),
		"by_type":          byType,
		"by_method":        byMethod,
		"recurring_count":  recurringCount,
		"anonymous_count":  anonymousCount,
		"average_donation": totalAmount / float64(len(donations)),
	}

	return stats, nil
}
