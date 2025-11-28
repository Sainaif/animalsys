package dto

import "github.com/sainaif/animalsys/backend/internal/domain/entities"

// TopDonorResponse represents the response for a top donor, including their total donated amount.
type TopDonorResponse struct {
	entities.Donor
	TotalDonated float64 `json:"total_donated"`
}

// RecurringDonorResponse represents the response for a recurring donor, including their recurring donation count and total amount.
type RecurringDonorResponse struct {
	entities.Donor
	RecurringCount  int     `json:"recurring_count"`
	RecurringAmount float64 `json:"recurring_amount"`
}
