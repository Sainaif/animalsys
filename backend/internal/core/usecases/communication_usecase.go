package usecases

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type CommunicationUseCase struct {
	communicationRepo interfaces.CommunicationRepository
	templateRepo      interfaces.CommunicationTemplateRepository
	auditRepo         interfaces.AuditLogRepository
}

func NewCommunicationUseCase(
	communicationRepo interfaces.CommunicationRepository,
	templateRepo interfaces.CommunicationTemplateRepository,
	auditRepo interfaces.AuditLogRepository,
) *CommunicationUseCase {
	return &CommunicationUseCase{
		communicationRepo: communicationRepo,
		templateRepo:      templateRepo,
		auditRepo:         auditRepo,
	}
}

// Communication Management

func (uc *CommunicationUseCase) CreateCommunication(ctx context.Context, req *entities.CommunicationCreateRequest, createdBy string) (*entities.Communication, error) {
	communication := entities.NewCommunication(
		req.Type,
		req.Subject,
		req.Content,
		req.Recipients,
	)

	communication.ScheduledFor = req.ScheduledFor
	communication.Priority = req.Priority
	communication.Tags = req.Tags
	communication.CreatedBy = createdBy

	// If scheduled for future, set status to scheduled
	if req.ScheduledFor != nil && req.ScheduledFor.After(time.Now()) {
		communication.Status = entities.CommunicationStatusScheduled
	} else {
		communication.Status = entities.CommunicationStatusDraft
	}

	if err := uc.communicationRepo.Create(ctx, communication); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "communication", communication.ID.Hex(), "Communication created")
	uc.auditRepo.Create(ctx, auditLog)

	return communication, nil
}

func (uc *CommunicationUseCase) GetCommunicationByID(ctx context.Context, id string) (*entities.Communication, error) {
	return uc.communicationRepo.GetByID(ctx, id)
}

func (uc *CommunicationUseCase) ListCommunications(ctx context.Context, status entities.CommunicationStatus, limit, offset int) ([]*entities.Communication, int64, error) {
	return uc.communicationRepo.List(ctx, status, limit, offset)
}

func (uc *CommunicationUseCase) UpdateCommunication(ctx context.Context, id string, req *entities.CommunicationUpdateRequest, updatedBy string) (*entities.Communication, error) {
	communication, err := uc.communicationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Subject != "" {
		communication.Subject = req.Subject
	}
	if req.Content != "" {
		communication.Content = req.Content
	}
	if req.Recipients != nil {
		communication.Recipients = req.Recipients
	}
	if req.ScheduledFor != nil {
		communication.ScheduledFor = req.ScheduledFor
		// Update status if scheduled for future
		if req.ScheduledFor.After(time.Now()) {
			communication.Status = entities.CommunicationStatusScheduled
		}
	}
	if req.Priority != "" {
		communication.Priority = req.Priority
	}
	if req.Tags != nil {
		communication.Tags = req.Tags
	}
	communication.UpdatedBy = updatedBy

	if err := uc.communicationRepo.Update(ctx, id, communication); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "communication", id, "Communication updated")
	uc.auditRepo.Create(ctx, auditLog)

	return communication, nil
}

func (uc *CommunicationUseCase) DeleteCommunication(ctx context.Context, id string, deletedBy string) error {
	if err := uc.communicationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "communication", id, "Communication deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *CommunicationUseCase) SendCommunication(ctx context.Context, id string, sentBy string) error {
	communication, err := uc.communicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// In real implementation, this would integrate with email/SMS service
	// For now, we just mark it as sent
	sentCount := len(communication.Recipients)
	failedCount := 0

	if err := uc.communicationRepo.MarkAsSent(ctx, id, sentCount, failedCount); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(sentBy, "", "", entities.ActionUpdate, "communication", id, "Communication sent")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *CommunicationUseCase) GetScheduledCommunications(ctx context.Context) ([]*entities.Communication, error) {
	return uc.communicationRepo.GetScheduled(ctx)
}

func (uc *CommunicationUseCase) ProcessScheduledCommunications(ctx context.Context) error {
	// Get all scheduled communications that are due
	scheduled, err := uc.communicationRepo.GetScheduled(ctx)
	if err != nil {
		return err
	}

	// Process each scheduled communication
	for _, comm := range scheduled {
		// In real implementation, this would send the communication
		// For now, we just mark it as sent
		sentCount := len(comm.Recipients)
		failedCount := 0

		if err := uc.communicationRepo.MarkAsSent(ctx, comm.ID.Hex(), sentCount, failedCount); err != nil {
			// Log error but continue processing others
			continue
		}

		// Audit
		auditLog := entities.NewAuditLog("system", "", "", entities.ActionUpdate, "communication", comm.ID.Hex(), "Scheduled communication sent")
		uc.auditRepo.Create(ctx, auditLog)
	}

	return nil
}

// Template Management

func (uc *CommunicationUseCase) CreateTemplate(ctx context.Context, req *entities.CommunicationTemplateCreateRequest, createdBy string) (*entities.CommunicationTemplate, error) {
	template := entities.NewCommunicationTemplate(
		req.Name,
		req.Type,
		req.Subject,
		req.Content,
	)

	template.Description = req.Description
	template.Variables = req.Variables
	template.Tags = req.Tags
	template.CreatedBy = createdBy

	if err := uc.templateRepo.Create(ctx, template); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "communication_template", template.ID.Hex(), "Template created: "+req.Name)
	uc.auditRepo.Create(ctx, auditLog)

	return template, nil
}

func (uc *CommunicationUseCase) GetTemplateByID(ctx context.Context, id string) (*entities.CommunicationTemplate, error) {
	return uc.templateRepo.GetByID(ctx, id)
}

func (uc *CommunicationUseCase) ListTemplates(ctx context.Context, commType entities.CommunicationType) ([]*entities.CommunicationTemplate, error) {
	return uc.templateRepo.List(ctx, commType)
}

func (uc *CommunicationUseCase) UpdateTemplate(ctx context.Context, id string, req *entities.CommunicationTemplateUpdateRequest, updatedBy string) (*entities.CommunicationTemplate, error) {
	template, err := uc.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Subject != "" {
		template.Subject = req.Subject
	}
	if req.Content != "" {
		template.Content = req.Content
	}
	if req.Variables != nil {
		template.Variables = req.Variables
	}
	if req.Tags != nil {
		template.Tags = req.Tags
	}
	template.UpdatedBy = updatedBy

	if err := uc.templateRepo.Update(ctx, id, template); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "communication_template", id, "Template updated")
	uc.auditRepo.Create(ctx, auditLog)

	return template, nil
}

func (uc *CommunicationUseCase) DeleteTemplate(ctx context.Context, id string, deletedBy string) error {
	if err := uc.templateRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "communication_template", id, "Template deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *CommunicationUseCase) CreateFromTemplate(ctx context.Context, templateID string, recipients []string, variables map[string]string, createdBy string) (*entities.Communication, error) {
	template, err := uc.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// Replace variables in content
	content := template.Content
	for key, value := range variables {
		// Simple variable replacement (in real implementation, use proper templating engine)
		// Variables should be in format {{variable_name}}
		content = replaceVariable(content, key, value)
	}

	// Create communication from template
	communication := entities.NewCommunication(
		template.Type,
		template.Subject,
		content,
		recipients,
	)
	communication.CreatedBy = createdBy
	communication.Status = entities.CommunicationStatusDraft

	if err := uc.communicationRepo.Create(ctx, communication); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "communication", communication.ID.Hex(), "Communication created from template")
	uc.auditRepo.Create(ctx, auditLog)

	return communication, nil
}

func (uc *CommunicationUseCase) GetCommunicationStatistics(ctx context.Context) (map[string]interface{}, error) {
	// Get all communications
	all, total, err := uc.communicationRepo.List(ctx, "", 0, 0)
	if err != nil {
		return nil, err
	}

	byStatus := make(map[string]int)
	byType := make(map[string]int)
	totalSent := 0
	totalFailed := 0

	for _, comm := range all {
		byStatus[string(comm.Status)]++
		byType[string(comm.Type)]++
		totalSent += comm.SentCount
		totalFailed += comm.FailedCount
	}

	stats := map[string]interface{}{
		"total_communications": total,
		"by_status":            byStatus,
		"by_type":              byType,
		"total_sent":           totalSent,
		"total_failed":         totalFailed,
	}

	return stats, nil
}

// Helper function for variable replacement
func replaceVariable(content, key, value string) string {
	// Simple replace - in production use proper templating engine
	placeholder := "{{" + key + "}}"
	return replaceAll(content, placeholder, value)
}

func replaceAll(s, old, new string) string {
	// Simple string replace implementation
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
