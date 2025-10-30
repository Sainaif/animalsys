package interfaces

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, id string, user *entities.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *entities.UserFilter) ([]*entities.User, int64, error)
	UpdateLastLogin(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, hashedPassword string) error
}

// AnimalRepository defines the interface for animal data operations
type AnimalRepository interface {
	Create(ctx context.Context, animal *entities.Animal) error
	GetByID(ctx context.Context, id string) (*entities.Animal, error)
	Update(ctx context.Context, id string, animal *entities.Animal) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *entities.AnimalFilter) ([]*entities.Animal, int64, error)
	AddMedicalRecord(ctx context.Context, id string, record entities.MedicalRecord) error
	AddPhoto(ctx context.Context, id string, photoURL string) error
	UpdateStatus(ctx context.Context, id string, status entities.AnimalStatus) error
	GetAvailableForAdoption(ctx context.Context, limit, offset int) ([]*entities.Animal, int64, error)
	GetByStatus(ctx context.Context, status entities.AnimalStatus, limit, offset int) ([]*entities.Animal, int64, error)
}

// AdoptionRepository defines the interface for adoption data operations
type AdoptionRepository interface {
	Create(ctx context.Context, adoption *entities.Adoption) error
	GetByID(ctx context.Context, id string) (*entities.Adoption, error)
	Update(ctx context.Context, id string, adoption *entities.Adoption) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *entities.AdoptionFilter) ([]*entities.Adoption, int64, error)
	GetByAnimalID(ctx context.Context, animalID string) ([]*entities.Adoption, error)
	GetByApplicantID(ctx context.Context, applicantID string) ([]*entities.Adoption, error)
	UpdateStatus(ctx context.Context, id string, status entities.AdoptionStatus) error
	GetPendingApplications(ctx context.Context, limit, offset int) ([]*entities.Adoption, int64, error)
}

// VolunteerRepository defines the interface for volunteer data operations
type VolunteerRepository interface {
	Create(ctx context.Context, volunteer *entities.Volunteer) error
	GetByID(ctx context.Context, id string) (*entities.Volunteer, error)
	GetByUserID(ctx context.Context, userID string) (*entities.Volunteer, error)
	Update(ctx context.Context, id string, volunteer *entities.Volunteer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status entities.VolunteerStatus, limit, offset int) ([]*entities.Volunteer, int64, error)
	AddTraining(ctx context.Context, id string, training entities.Training) error
	LogHours(ctx context.Context, id string, hours float64) error
	GetActiveVolunteers(ctx context.Context) ([]*entities.Volunteer, error)
	ResetMonthlyHours(ctx context.Context) error
}

// VolunteerHourRepository defines the interface for volunteer hour tracking
type VolunteerHourRepository interface {
	Create(ctx context.Context, entry *entities.VolunteerHourEntry) error
	GetByID(ctx context.Context, id string) (*entities.VolunteerHourEntry, error)
	GetByVolunteerID(ctx context.Context, volunteerID string, limit, offset int) ([]*entities.VolunteerHourEntry, int64, error)
	Update(ctx context.Context, id string, entry *entities.VolunteerHourEntry) error
	Delete(ctx context.Context, id string) error
	GetTotalHours(ctx context.Context, volunteerID string) (float64, error)
	GetHoursByDateRange(ctx context.Context, volunteerID string, startDate, endDate string) (float64, error)
}

// ScheduleRepository defines the interface for schedule data operations
type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entities.Schedule) error
	GetByID(ctx context.Context, id string) (*entities.Schedule, error)
	Update(ctx context.Context, id string, schedule *entities.Schedule) error
	Delete(ctx context.Context, id string) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.Schedule, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Schedule, error)
	GetByStatus(ctx context.Context, status entities.ShiftStatus) ([]*entities.Schedule, error)
	UpdateStatus(ctx context.Context, id string, status entities.ShiftStatus) error
	GetUpcomingShifts(ctx context.Context, userID string, limit int) ([]*entities.Schedule, error)
}

// DocumentRepository defines the interface for document data operations
type DocumentRepository interface {
	Create(ctx context.Context, document *entities.Document) error
	GetByID(ctx context.Context, id string) (*entities.Document, error)
	Update(ctx context.Context, id string, document *entities.Document) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *entities.DocumentFilter) ([]*entities.Document, int64, error)
	GetByRelatedEntity(ctx context.Context, entityType, entityID string) ([]*entities.Document, error)
	GetExpiringSoon(ctx context.Context, days int) ([]*entities.Document, error)
}

// FinanceRepository defines the interface for finance data operations
type FinanceRepository interface {
	Create(ctx context.Context, finance *entities.Finance) error
	GetByID(ctx context.Context, id string) (*entities.Finance, error)
	Update(ctx context.Context, id string, finance *entities.Finance) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *entities.FinanceFilter) ([]*entities.Finance, int64, error)
	GetTotalByType(ctx context.Context, financeType entities.FinanceType, startDate, endDate string) (float64, error)
	GetByCategory(ctx context.Context, category string, startDate, endDate string) ([]*entities.Finance, error)
	GetSummary(ctx context.Context, startDate, endDate string) (map[string]float64, error)
	GetByFiscalYear(ctx context.Context, year int) ([]*entities.Finance, error)
}

// DonorRepository defines the interface for donor data operations
type DonorRepository interface {
	Create(ctx context.Context, donor *entities.Donor) error
	GetByID(ctx context.Context, id string) (*entities.Donor, error)
	GetByEmail(ctx context.Context, email string) (*entities.Donor, error)
	Update(ctx context.Context, id string, donor *entities.Donor) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, donorType entities.DonorType, limit, offset int) ([]*entities.Donor, int64, error)
	GetMajorDonors(ctx context.Context) ([]*entities.Donor, error)
	UpdateTotals(ctx context.Context, id string, amount float64) error
}

// DonationRepository defines the interface for donation data operations
type DonationRepository interface {
	Create(ctx context.Context, donation *entities.Donation) error
	GetByID(ctx context.Context, id string) (*entities.Donation, error)
	Update(ctx context.Context, id string, donation *entities.Donation) error
	Delete(ctx context.Context, id string) error
	GetByDonorID(ctx context.Context, donorID string, limit, offset int) ([]*entities.Donation, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.Donation, error)
	GetTotalByDonor(ctx context.Context, donorID string) (float64, error)
	GetRecurringDonations(ctx context.Context) ([]*entities.Donation, error)
}

// InventoryRepository defines the interface for inventory data operations
type InventoryRepository interface {
	Create(ctx context.Context, item *entities.InventoryItem) error
	GetByID(ctx context.Context, id string) (*entities.InventoryItem, error)
	Update(ctx context.Context, id string, item *entities.InventoryItem) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, itemType entities.InventoryItemType, limit, offset int) ([]*entities.InventoryItem, int64, error)
	GetLowStock(ctx context.Context) ([]*entities.InventoryItem, error)
	GetOutOfStock(ctx context.Context) ([]*entities.InventoryItem, error)
	GetExpiringSoon(ctx context.Context, days int) ([]*entities.InventoryItem, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

// StockMovementRepository defines the interface for stock movement data operations
type StockMovementRepository interface {
	Create(ctx context.Context, movement *entities.StockMovement) error
	GetByID(ctx context.Context, id string) (*entities.StockMovement, error)
	GetByItemID(ctx context.Context, itemID string, limit, offset int) ([]*entities.StockMovement, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.StockMovement, error)
}

// VeterinaryVisitRepository defines the interface for veterinary visit data operations
type VeterinaryVisitRepository interface {
	Create(ctx context.Context, visit *entities.VeterinaryVisit) error
	GetByID(ctx context.Context, id string) (*entities.VeterinaryVisit, error)
	Update(ctx context.Context, id string, visit *entities.VeterinaryVisit) error
	Delete(ctx context.Context, id string) error
	GetByAnimalID(ctx context.Context, animalID string, limit, offset int) ([]*entities.VeterinaryVisit, int64, error)
	GetUpcomingVisits(ctx context.Context, days int) ([]*entities.VeterinaryVisit, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]*entities.VeterinaryVisit, error)
}

// VaccinationRepository defines the interface for vaccination data operations
type VaccinationRepository interface {
	Create(ctx context.Context, vaccination *entities.Vaccination) error
	GetByID(ctx context.Context, id string) (*entities.Vaccination, error)
	Update(ctx context.Context, id string, vaccination *entities.Vaccination) error
	Delete(ctx context.Context, id string) error
	GetByAnimalID(ctx context.Context, animalID string) ([]*entities.Vaccination, error)
	GetDueSoon(ctx context.Context, days int) ([]*entities.Vaccination, error)
}

// CampaignRepository defines the interface for campaign data operations
type CampaignRepository interface {
	Create(ctx context.Context, campaign *entities.Campaign) error
	GetByID(ctx context.Context, id string) (*entities.Campaign, error)
	Update(ctx context.Context, id string, campaign *entities.Campaign) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status entities.CampaignStatus, limit, offset int) ([]*entities.Campaign, int64, error)
	GetActiveCampaigns(ctx context.Context) ([]*entities.Campaign, error)
	UpdateProgress(ctx context.Context, id string, progress float64) error
	AddExpense(ctx context.Context, id string, amount float64) error
}

// PartnerRepository defines the interface for partner data operations
type PartnerRepository interface {
	Create(ctx context.Context, partner *entities.Partner) error
	GetByID(ctx context.Context, id string) (*entities.Partner, error)
	Update(ctx context.Context, id string, partner *entities.Partner) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, partnerType entities.PartnerType, limit, offset int) ([]*entities.Partner, int64, error)
	GetActivePartners(ctx context.Context) ([]*entities.Partner, error)
	AddCollaboration(ctx context.Context, id string, collaboration entities.Collaboration) error
}

// CommunicationRepository defines the interface for communication data operations
type CommunicationRepository interface {
	Create(ctx context.Context, communication *entities.Communication) error
	GetByID(ctx context.Context, id string) (*entities.Communication, error)
	Update(ctx context.Context, id string, communication *entities.Communication) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status entities.CommunicationStatus, limit, offset int) ([]*entities.Communication, int64, error)
	GetScheduled(ctx context.Context) ([]*entities.Communication, error)
	MarkAsSent(ctx context.Context, id string, sentCount, failedCount int) error
}

// CommunicationTemplateRepository defines the interface for communication template data operations
type CommunicationTemplateRepository interface {
	Create(ctx context.Context, template *entities.CommunicationTemplate) error
	GetByID(ctx context.Context, id string) (*entities.CommunicationTemplate, error)
	Update(ctx context.Context, id string, template *entities.CommunicationTemplate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, commType entities.CommunicationType) ([]*entities.CommunicationTemplate, error)
}

// ReportRepository defines the interface for report data operations
type ReportRepository interface {
	Create(ctx context.Context, report *entities.Report) error
	GetByID(ctx context.Context, id string) (*entities.Report, error)
	Update(ctx context.Context, id string, report *entities.Report) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, reportType entities.ReportType, limit, offset int) ([]*entities.Report, int64, error)
	GetByCreatedBy(ctx context.Context, userID string, limit, offset int) ([]*entities.Report, int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
	MarkAsCompleted(ctx context.Context, id, fileURL string, fileSize int64) error
	MarkAsFailed(ctx context.Context, id, errorMsg string) error
}

// AuditLogRepository defines the interface for audit log data operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *entities.AuditLog) error
	GetByID(ctx context.Context, id string) (*entities.AuditLog, error)
	List(ctx context.Context, filter *entities.AuditLogFilter) ([]*entities.AuditLog, int64, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entities.AuditLog, int64, error)
	GetByResourceID(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entities.AuditLog, int64, error)
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}
