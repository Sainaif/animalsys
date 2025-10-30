package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type UserUseCase struct {
	userRepo  interfaces.UserRepository
	auditRepo interfaces.AuditLogRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository, auditRepo interfaces.AuditLogRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*entities.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) List(ctx context.Context, filter *entities.UserFilter) ([]*entities.User, int64, error) {
	return uc.userRepo.List(ctx, filter)
}

func (uc *UserUseCase) Update(ctx context.Context, id string, req *entities.UserUpdateRequest, updatedBy string) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Active != nil {
		user.Active = *req.Active
	}
	user.UpdatedBy = updatedBy

	if err := uc.userRepo.Update(ctx, id, user); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "user", id, "User updated")
	uc.auditRepo.Create(ctx, auditLog)

	return user, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "user", id, "User deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}
