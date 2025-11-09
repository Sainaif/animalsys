package user

import (
	"context"
	"time"

	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	"github.com/sainaif/animalsys/backend/internal/domain/repositories"
	"github.com/sainaif/animalsys/backend/pkg/errors"
	"github.com/sainaif/animalsys/backend/pkg/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserUseCase handles user management business logic
type UserUseCase struct {
	userRepo        repositories.UserRepository
	auditLogRepo    repositories.AuditLogRepository
	passwordService *security.PasswordService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(
	userRepo repositories.UserRepository,
	auditLogRepo repositories.AuditLogRepository,
	passwordService *security.PasswordService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		auditLogRepo:    auditLogRepo,
		passwordService: passwordService,
	}
}

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	Email     string             `json:"email" validate:"required,email"`
	Password  string             `json:"password" validate:"required,min=8"`
	FirstName string             `json:"first_name" validate:"required"`
	LastName  string             `json:"last_name" validate:"required"`
	Role      entities.UserRole  `json:"role" validate:"required"`
	Status    entities.UserStatus `json:"status"`
	Phone     string             `json:"phone,omitempty"`
	Language  string             `json:"language" validate:"required,oneof=en pl"`
	Theme     string             `json:"theme" validate:"required,oneof=light dark"`
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req *CreateUserRequest, creatorID primitive.ObjectID) (*entities.User, error) {
	// Validate password strength
	if err := uc.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Validate role
	if !entities.IsValidRole(req.Role) {
		return nil, errors.NewBadRequest("invalid role")
	}

	// Hash password
	passwordHash, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = entities.StatusActive
	}

	// Create user
	user := &entities.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Status:       status,
		Phone:        req.Phone,
		Language:     req.Language,
		Theme:        req.Theme,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(creatorID, entities.ActionCreate, "user", "", "").
		WithEntityID(user.ID)
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	// Remove sensitive data
	user.PasswordHash = ""

	return user, nil
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	FirstName *string             `json:"first_name,omitempty"`
	LastName  *string             `json:"last_name,omitempty"`
	Role      *entities.UserRole  `json:"role,omitempty"`
	Status    *entities.UserStatus `json:"status,omitempty"`
	Phone     *string             `json:"phone,omitempty"`
	Avatar    *string             `json:"avatar,omitempty"`
	Language  *string             `json:"language,omitempty" validate:"omitempty,oneof=en pl"`
	Theme     *string             `json:"theme,omitempty" validate:"omitempty,oneof=light dark"`
}

// UpdateUser updates a user
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID primitive.ObjectID, req *UpdateUserRequest, updaterID primitive.ObjectID) (*entities.User, error) {
	// Get existing user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Track changes for audit log
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.FirstName != nil && *req.FirstName != user.FirstName {
		changes["first_name"] = map[string]string{"old": user.FirstName, "new": *req.FirstName}
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil && *req.LastName != user.LastName {
		changes["last_name"] = map[string]string{"old": user.LastName, "new": *req.LastName}
		user.LastName = *req.LastName
	}

	if req.Role != nil && *req.Role != user.Role {
		if !entities.IsValidRole(*req.Role) {
			return nil, errors.NewBadRequest("invalid role")
		}
		changes["role"] = map[string]string{"old": string(user.Role), "new": string(*req.Role)}
		user.Role = *req.Role
	}

	if req.Status != nil && *req.Status != user.Status {
		if !entities.IsValidStatus(*req.Status) {
			return nil, errors.NewBadRequest("invalid status")
		}
		changes["status"] = map[string]string{"old": string(user.Status), "new": string(*req.Status)}
		user.Status = *req.Status
	}

	if req.Phone != nil && *req.Phone != user.Phone {
		changes["phone"] = map[string]string{"old": user.Phone, "new": *req.Phone}
		user.Phone = *req.Phone
	}

	if req.Avatar != nil && *req.Avatar != user.Avatar {
		user.Avatar = *req.Avatar
	}

	if req.Language != nil && *req.Language != user.Language {
		changes["language"] = map[string]string{"old": user.Language, "new": *req.Language}
		user.Language = *req.Language
	}

	if req.Theme != nil && *req.Theme != user.Theme {
		changes["theme"] = map[string]string{"old": user.Theme, "new": *req.Theme}
		user.Theme = *req.Theme
	}

	user.UpdatedAt = time.Now()

	// Update in database
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Create audit log if there were changes
	if len(changes) > 0 {
		auditLog := entities.NewAuditLog(updaterID, entities.ActionUpdate, "user", "", "").
			WithEntityID(user.ID).
			WithChanges(changes)
		_ = uc.auditLogRepo.Create(ctx, auditLog)
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.RefreshToken = ""

	return user, nil
}

// GetUserByID gets a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*entities.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.RefreshToken = ""

	return user, nil
}

// ListUsersRequest represents a user list request
type ListUsersRequest struct {
	Role   string `form:"role"`
	Status string `form:"status"`
	Search string `form:"search"`
	Limit  int64  `form:"limit"`
	Offset int64  `form:"offset"`
}

// ListUsersResponse represents a user list response
type ListUsersResponse struct {
	Users []*entities.User `json:"users"`
	Total int64            `json:"total"`
	Limit int64            `json:"limit"`
	Offset int64           `json:"offset"`
}

// ListUsers lists users with filters and pagination
func (uc *UserUseCase) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	filter := repositories.UserFilter{
		Role:   req.Role,
		Status: req.Status,
		Search: req.Search,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Remove sensitive data from all users
	for _, user := range users {
		user.PasswordHash = ""
		user.RefreshToken = ""
	}

	return &ListUsersResponse{
		Users:  users,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(ctx context.Context, userID primitive.ObjectID, deleterID primitive.ObjectID) error {
	// Prevent self-deletion
	if userID == deleterID {
		return errors.NewBadRequest("cannot delete your own account")
	}

	// Check if user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Delete user
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// Create audit log
	auditLog := entities.NewAuditLog(deleterID, entities.ActionDelete, "user", "", "").
		WithEntityID(userID).
		WithChanges(map[string]interface{}{
			"email": user.Email,
			"name":  user.FullName(),
		})
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResetPassword resets a user's password (admin only)
func (uc *UserUseCase) ResetPassword(ctx context.Context, userID primitive.ObjectID, req *ResetPasswordRequest, adminID primitive.ObjectID) error {
	// Validate password strength
	if err := uc.passwordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// Get user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := uc.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Clear refresh token to force re-login
	_ = uc.userRepo.UpdateRefreshToken(ctx, userID, "")

	// Create audit log
	auditLog := entities.NewAuditLog(adminID, entities.ActionUpdate, "user", "", "").
		WithEntityID(userID).
		WithChanges(map[string]interface{}{"password": "reset by admin"})
	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return nil
}
