package usecases

import (
	"context"
	"errors"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type ScheduleUseCase struct {
	scheduleRepo interfaces.ScheduleRepository
	userRepo     interfaces.UserRepository
	auditRepo    interfaces.AuditLogRepository
}

func NewScheduleUseCase(
	scheduleRepo interfaces.ScheduleRepository,
	userRepo interfaces.UserRepository,
	auditRepo interfaces.AuditLogRepository,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		auditRepo:    auditRepo,
	}
}

func (uc *ScheduleUseCase) Create(ctx context.Context, req *entities.ScheduleCreateRequest, createdBy string) (*entities.Schedule, error) {
	// Get user info
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	schedule := entities.NewSchedule(
		req.UserID,
		user.FirstName+" "+user.LastName,
		req.ShiftType,
		req.StartTime,
		req.EndTime,
	)
	schedule.Location = req.Location
	schedule.Notes = req.Notes
	schedule.CreatedBy = createdBy

	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "schedule", schedule.ID.Hex(), "Shift scheduled")
	uc.auditRepo.Create(ctx, auditLog)

	return schedule, nil
}

func (uc *ScheduleUseCase) GetByID(ctx context.Context, id string) (*entities.Schedule, error) {
	return uc.scheduleRepo.GetByID(ctx, id)
}

func (uc *ScheduleUseCase) List(ctx context.Context, filter *entities.ScheduleFilter) ([]*entities.Schedule, int64, error) {
	return uc.scheduleRepo.List(ctx, filter)
}

func (uc *ScheduleUseCase) Update(ctx context.Context, id string, req *entities.ScheduleUpdateRequest, updatedBy string) (*entities.Schedule, error) {
	schedule, err := uc.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.ShiftType != "" {
		schedule.ShiftType = req.ShiftType
	}
	if !req.StartTime.IsZero() {
		schedule.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		schedule.EndTime = req.EndTime
	}
	if req.Location != "" {
		schedule.Location = req.Location
	}
	if req.Status != "" {
		schedule.Status = req.Status
	}
	if req.Notes != "" {
		schedule.Notes = req.Notes
	}
	schedule.UpdatedBy = updatedBy

	if err := uc.scheduleRepo.Update(ctx, id, schedule); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "schedule", id, "Shift updated")
	uc.auditRepo.Create(ctx, auditLog)

	return schedule, nil
}

func (uc *ScheduleUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.scheduleRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "schedule", id, "Shift deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *ScheduleUseCase) GetUserSchedule(ctx context.Context, userID string, startDate, endDate string, limit, offset int) ([]*entities.Schedule, int64, error) {
	return uc.scheduleRepo.GetByUserID(ctx, userID, startDate, endDate, limit, offset)
}

func (uc *ScheduleUseCase) GetByDateRange(ctx context.Context, startDate, endDate string, limit, offset int) ([]*entities.Schedule, int64, error) {
	return uc.scheduleRepo.GetByDateRange(ctx, startDate, endDate, limit, offset)
}

func (uc *ScheduleUseCase) RequestSwap(ctx context.Context, scheduleID, requestingUserID, targetUserID, reason string) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	// Verify requesting user owns the shift
	if schedule.UserID != requestingUserID {
		return errors.New("only the assigned user can request a swap")
	}

	// Verify target user exists
	targetUser, err := uc.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	schedule.RequestSwap(targetUserID, targetUser.FirstName+" "+targetUser.LastName, reason)

	if err := uc.scheduleRepo.Update(ctx, scheduleID, schedule); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(requestingUserID, "", "", entities.ActionUpdate, "schedule", scheduleID, "Shift swap requested")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *ScheduleUseCase) ApproveSwap(ctx context.Context, scheduleID, targetUserID, approvedBy string) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.SwapRequest == nil {
		return errors.New("no swap request found")
	}

	if schedule.SwapRequest.TargetUserID != targetUserID {
		return errors.New("target user mismatch")
	}

	// Get target user info
	targetUser, err := uc.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	schedule.ApproveSwap(targetUserID, targetUser.FirstName+" "+targetUser.LastName)

	if err := uc.scheduleRepo.Update(ctx, scheduleID, schedule); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(approvedBy, "", "", entities.ActionUpdate, "schedule", scheduleID, "Shift swap approved")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *ScheduleUseCase) RejectSwap(ctx context.Context, scheduleID, rejectedBy string) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.SwapRequest == nil {
		return errors.New("no swap request found")
	}

	schedule.RejectSwap()

	if err := uc.scheduleRepo.Update(ctx, scheduleID, schedule); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(rejectedBy, "", "", entities.ActionUpdate, "schedule", scheduleID, "Shift swap rejected")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *ScheduleUseCase) MarkComplete(ctx context.Context, scheduleID, completedBy string) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	schedule.MarkComplete()

	if err := uc.scheduleRepo.Update(ctx, scheduleID, schedule); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(completedBy, "", "", entities.ActionUpdate, "schedule", scheduleID, "Shift marked as complete")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *ScheduleUseCase) MarkAbsent(ctx context.Context, scheduleID, reason, markedBy string) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	schedule.MarkAbsent(reason)

	if err := uc.scheduleRepo.Update(ctx, scheduleID, schedule); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(markedBy, "", "", entities.ActionUpdate, "schedule", scheduleID, "Shift marked as absent")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}
