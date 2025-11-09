package communication

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommunicationUseCase struct {
	communicationRepo repositories.CommunicationRepository
	templateRepo      repositories.CommunicationTemplateRepository
	auditLogRepo      repositories.AuditLogRepository
}

func NewCommunicationUseCase(
	communicationRepo repositories.CommunicationRepository,
	templateRepo repositories.CommunicationTemplateRepository,
	auditLogRepo repositories.AuditLogRepository,
) *CommunicationUseCase {
	return &CommunicationUseCase{
		communicationRepo: communicationRepo,
		templateRepo:      templateRepo,
		auditLogRepo:      auditLogRepo,
	}
}

// CreateCommunication creates a new communication record
func (uc *CommunicationUseCase) CreateCommunication(ctx context.Context, communication *entities.Communication, userID primitive.ObjectID) error {
	// Validate template if provided
	if communication.TemplateID != nil {
		template, err := uc.templateRepo.FindByID(ctx, *communication.TemplateID)
		if err != nil {
			return errors.NewBadRequest("Invalid template ID")
		}

		// Increment template usage
		_ = uc.templateRepo.IncrementUsage(ctx, *communication.TemplateID)

		// If no subject/body provided, use template
		if communication.Subject == "" {
			communication.Subject = template.Subject
		}
		if communication.Body == "" {
			communication.Body = template.Body
		}
	}

	// Validate required fields
	if communication.RecipientEmail == "" && communication.RecipientPhone == "" {
		return errors.NewBadRequest("Either email or phone number is required")
	}

	if communication.Type == entities.TemplateTypeEmail && communication.RecipientEmail == "" {
		return errors.NewBadRequest("Email address is required for email communications")
	}

	if communication.Type == entities.TemplateTypeSMS && communication.RecipientPhone == "" {
		return errors.NewBadRequest("Phone number is required for SMS communications")
	}

	if err := uc.communicationRepo.Create(ctx, communication); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "communication", "", "").
			WithEntityID(communication.ID))

	return nil
}

// SendCommunicationFromTemplate sends a communication using a template
func (uc *CommunicationUseCase) SendCommunicationFromTemplate(
	ctx context.Context,
	templateID primitive.ObjectID,
	recipientType entities.RecipientType,
	recipientID primitive.ObjectID,
	recipientEmail string,
	recipientPhone string,
	variables map[string]string,
	userID primitive.ObjectID,
) (*entities.Communication, error) {
	// Get template
	template, err := uc.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, errors.NewBadRequest("Template not found")
	}

	if !template.Active {
		return nil, errors.NewBadRequest("Template is not active")
	}

	// Render template with variables
	subject := template.RenderSubject(variables)
	body := template.RenderBody(variables)

	// Create communication
	communication := &entities.Communication{
		Type:           template.Type,
		Category:       template.Category,
		Status:         entities.CommunicationStatusPending,
		RecipientType:  recipientType,
		RecipientID:    &recipientID,
		RecipientEmail: recipientEmail,
		RecipientPhone: recipientPhone,
		SenderID:       userID,
		TemplateID:     &templateID,
		Subject:        subject,
		Body:           body,
		MaxRetries:     3,
	}

	if err := uc.communicationRepo.Create(ctx, communication); err != nil {
		return nil, err
	}

	// Increment template usage
	_ = uc.templateRepo.IncrementUsage(ctx, templateID)

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "communication", "", "").
			WithEntityID(communication.ID))

	return communication, nil
}

// GetCommunicationByID retrieves a communication by ID
func (uc *CommunicationUseCase) GetCommunicationByID(ctx context.Context, id primitive.ObjectID) (*entities.Communication, error) {
	return uc.communicationRepo.FindByID(ctx, id)
}

// ListCommunications lists communications with filtering
func (uc *CommunicationUseCase) ListCommunications(ctx context.Context, filter *repositories.CommunicationFilter) ([]*entities.Communication, int64, error) {
	return uc.communicationRepo.List(ctx, filter)
}

// UpdateCommunicationStatus updates the status of a communication
func (uc *CommunicationUseCase) UpdateCommunicationStatus(ctx context.Context, id primitive.ObjectID, status entities.CommunicationStatus, userID primitive.ObjectID) error {
	communication, err := uc.communicationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.communicationRepo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "communication", "", "").
			WithEntityID(communication.ID).
			WithChanges(map[string]interface{}{
				"status": status,
			}))

	return nil
}

// MarkCommunicationAsOpened marks a communication as opened
func (uc *CommunicationUseCase) MarkCommunicationAsOpened(ctx context.Context, id primitive.ObjectID) error {
	return uc.communicationRepo.MarkAsOpened(ctx, id)
}

// MarkCommunicationAsClicked marks a communication as clicked
func (uc *CommunicationUseCase) MarkCommunicationAsClicked(ctx context.Context, id primitive.ObjectID) error {
	return uc.communicationRepo.MarkAsClicked(ctx, id)
}

// GetPendingCommunications retrieves pending communications ready to be sent
func (uc *CommunicationUseCase) GetPendingCommunications(ctx context.Context) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetPending(ctx)
}

// GetCommunicationsForRetry retrieves failed communications ready for retry
func (uc *CommunicationUseCase) GetCommunicationsForRetry(ctx context.Context) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetForRetry(ctx)
}

// GetCommunicationsByRecipient retrieves communications for a specific recipient
func (uc *CommunicationUseCase) GetCommunicationsByRecipient(ctx context.Context, recipientType entities.RecipientType, recipientID primitive.ObjectID) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetByRecipient(ctx, recipientType, recipientID)
}

// GetCommunicationsByCampaign retrieves communications for a campaign
func (uc *CommunicationUseCase) GetCommunicationsByCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetByCampaign(ctx, campaignID)
}

// GetCommunicationsByBatch retrieves communications by batch ID
func (uc *CommunicationUseCase) GetCommunicationsByBatch(ctx context.Context, batchID string) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetByBatch(ctx, batchID)
}

// GetCommunicationStatistics retrieves communication statistics
func (uc *CommunicationUseCase) GetCommunicationStatistics(ctx context.Context, startDate, endDate time.Time) (*repositories.CommunicationStatistics, error) {
	return uc.communicationRepo.GetStatistics(ctx, startDate, endDate)
}

// DeleteCommunication deletes a communication
func (uc *CommunicationUseCase) DeleteCommunication(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	communication, err := uc.communicationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deletion of draft or failed communications
	if communication.Status != entities.CommunicationStatusPending && communication.Status != entities.CommunicationStatusFailed {
		return errors.NewBadRequest("Only pending or failed communications can be deleted")
	}

	if err := uc.communicationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "communication", "", "").
			WithEntityID(communication.ID))

	return nil
}

// CreateTemplate creates a new communication template
func (uc *CommunicationUseCase) CreateTemplate(ctx context.Context, template *entities.CommunicationTemplate, userID primitive.ObjectID) error {
	// Validate required fields
	if template.Name == "" {
		return errors.NewBadRequest("Template name is required")
	}

	if template.Subject == "" && template.Type == entities.TemplateTypeEmail {
		return errors.NewBadRequest("Subject is required for email templates")
	}

	if template.Body == "" {
		return errors.NewBadRequest("Template body is required")
	}

	if err := uc.templateRepo.Create(ctx, template); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionCreate, "communication_template", "", "").
			WithEntityID(template.ID))

	return nil
}

// UpdateTemplate updates an existing template
func (uc *CommunicationUseCase) UpdateTemplate(ctx context.Context, template *entities.CommunicationTemplate, userID primitive.ObjectID) error {
	// Validate template exists
	if _, err := uc.templateRepo.FindByID(ctx, template.ID); err != nil {
		return err
	}

	if err := uc.templateRepo.Update(ctx, template); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionUpdate, "communication_template", "", "").
			WithEntityID(template.ID))

	return nil
}

// GetTemplateByID retrieves a template by ID
func (uc *CommunicationUseCase) GetTemplateByID(ctx context.Context, id primitive.ObjectID) (*entities.CommunicationTemplate, error) {
	return uc.templateRepo.FindByID(ctx, id)
}

// ListTemplates lists templates with filtering
func (uc *CommunicationUseCase) ListTemplates(ctx context.Context, filter *repositories.CommunicationTemplateFilter) ([]*entities.CommunicationTemplate, int64, error) {
	return uc.templateRepo.List(ctx, filter)
}

// GetTemplatesByCategory retrieves templates by category
func (uc *CommunicationUseCase) GetTemplatesByCategory(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) ([]*entities.CommunicationTemplate, error) {
	return uc.templateRepo.GetByCategory(ctx, category, templateType)
}

// GetDefaultTemplate retrieves the default template for a category and type
func (uc *CommunicationUseCase) GetDefaultTemplate(ctx context.Context, category entities.TemplateCategory, templateType entities.TemplateType) (*entities.CommunicationTemplate, error) {
	return uc.templateRepo.GetDefault(ctx, category, templateType)
}

// GetActiveTemplates retrieves all active templates
func (uc *CommunicationUseCase) GetActiveTemplates(ctx context.Context) ([]*entities.CommunicationTemplate, error) {
	return uc.templateRepo.GetActiveTemplates(ctx)
}

// DeleteTemplate deletes a template
func (uc *CommunicationUseCase) DeleteTemplate(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	template, err := uc.templateRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.templateRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	_ = uc.auditLogRepo.Create(ctx,
		entities.NewAuditLog(userID, entities.ActionDelete, "communication_template", "", "").
			WithEntityID(template.ID))

	return nil
}
