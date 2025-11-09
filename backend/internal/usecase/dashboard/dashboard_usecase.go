package dashboard

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
)

type DashboardUseCase struct {
	animalRepo     repositories.AnimalRepository
	adoptionRepo   repositories.AdoptionRepository
	donationRepo   repositories.DonationRepository
	volunteerRepo  repositories.VolunteerRepository
}

func NewDashboardUseCase(
	animalRepo repositories.AnimalRepository,
	adoptionRepo repositories.AdoptionRepository,
	donationRepo repositories.DonationRepository,
	volunteerRepo repositories.VolunteerRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		animalRepo:     animalRepo,
		adoptionRepo:   adoptionRepo,
		donationRepo:   donationRepo,
		volunteerRepo:  volunteerRepo,
	}
}

// GetDashboardMetrics retrieves all dashboard metrics
func (uc *DashboardUseCase) GetDashboardMetrics(ctx context.Context) (*entities.DashboardMetrics, error) {
	metrics := &entities.DashboardMetrics{
		GeneratedAt: time.Now(),
	}

	// Get animal statistics
	if animalStats, err := uc.animalRepo.GetStatistics(ctx); err == nil {
		metrics.Animals = entities.AnimalMetrics{
			Total:                animalStats.TotalAnimals,
			InShelter:            animalStats.AvailableForAdoption,
			Adopted:              animalStats.AdoptedThisYear,
			InFoster:             0, // Not tracked in current stats
			BySpecies:            animalStats.BySpecies,
			ByStatus:             animalStats.ByStatus,
			AverageDaysInShelter: 0, // Not tracked in current stats
		}

		// Update overview
		metrics.Overview.TotalAnimals = animalStats.TotalAnimals
		metrics.Overview.AnimalsInShelter = animalStats.AvailableForAdoption
		metrics.Overview.AnimalsAdopted = animalStats.AdoptedThisYear
		metrics.Overview.AnimalsInFoster = 0
	}

	// Get adoption statistics
	if adoptionStats, err := uc.adoptionRepo.GetAdoptionStatistics(ctx); err == nil {
		metrics.Adoptions = entities.AdoptionMetrics{
			TotalAdoptions:       adoptionStats.TotalAdoptions,
			PendingApplications:  adoptionStats.PendingAdoptions,
			ApprovedApplications: adoptionStats.CompletedAdoptions,
		}

		metrics.Overview.TotalAdoptions = adoptionStats.TotalAdoptions
	}

	// Get donation statistics
	if donationStats, err := uc.donationRepo.GetDonationStatistics(ctx); err == nil {
		metrics.Donations = entities.DonationMetrics{
			TotalAmount:      donationStats.TotalAmount,
			ThisWeek:         donationStats.AmountThisWeek,
			ThisMonth:        donationStats.AmountThisMonth,
			ThisYear:         donationStats.AmountThisYear,
			TotalDonors:      0, // Need to get from donor stats
			RecurringDonors:  donationStats.RecurringDonations,
			MajorDonors:      0, // Need to get from donor stats
			AverageDonation:  donationStats.AverageDonation,
			LargestDonation:  donationStats.LargestDonation,
		}

		metrics.Overview.TotalDonations = donationStats.TotalAmount
		metrics.Overview.DonationsThisMonth = donationStats.AmountThisMonth
	}

	// Get volunteer statistics
	if volunteerStats, err := uc.volunteerRepo.GetVolunteerStatistics(ctx); err == nil {
		metrics.Volunteers = entities.VolunteerMetrics{
			TotalVolunteers:  volunteerStats.TotalVolunteers,
			ActiveVolunteers: volunteerStats.ActiveVolunteers,
			TotalHours:       volunteerStats.TotalHours,
			AverageRating:    volunteerStats.AverageRating,
		}

		metrics.Overview.ActiveVolunteers = volunteerStats.ActiveVolunteers
		metrics.Overview.TotalVolunteerHours = volunteerStats.TotalHours
	}

	return metrics, nil
}

// GetOverviewMetrics retrieves overview-only metrics for quick display
func (uc *DashboardUseCase) GetOverviewMetrics(ctx context.Context) (*entities.OverviewMetrics, error) {
	overview := &entities.OverviewMetrics{}

	// Get key animal counts
	if animalStats, err := uc.animalRepo.GetStatistics(ctx); err == nil {
		overview.TotalAnimals = animalStats.TotalAnimals
		overview.AnimalsInShelter = animalStats.AvailableForAdoption
		overview.AnimalsAdopted = animalStats.AdoptedThisYear
		overview.AnimalsInFoster = 0
	}

	// Get donation totals
	if donationStats, err := uc.donationRepo.GetDonationStatistics(ctx); err == nil {
		overview.TotalDonations = donationStats.TotalAmount
		overview.DonationsThisMonth = donationStats.AmountThisMonth
	}

	// Get volunteer counts
	if volunteerStats, err := uc.volunteerRepo.GetVolunteerStatistics(ctx); err == nil {
		overview.ActiveVolunteers = volunteerStats.ActiveVolunteers
		overview.TotalVolunteerHours = volunteerStats.TotalHours
	}

	return overview, nil
}

// Helper function to convert map[string]int64 from stats
func convertToInt64Map(m map[string]int64) map[string]int64 {
	result := make(map[string]int64)
	for k, v := range m {
		result[k] = v
	}
	return result
}
